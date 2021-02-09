package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ggoop/mdf/bootstrap/errors"
	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/reg"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ggoop/mdf/utils"
	"github.com/robfig/cron"

	"github.com/ggoop/mdf/repositories"

	"github.com/ggoop/mdf/bootstrap/model"
)

var m_cronCache *cron.Cron

type CronSv struct {
	repo  *repositories.MysqlRepo
	logSv *LogSv
	*sync.Mutex
}

func NewCronSv(repo *repositories.MysqlRepo) *CronSv {
	return &CronSv{repo: repo, logSv: NewLogSv(repo), Mutex: &sync.Mutex{}}
}
func (s *CronSv) Start() {
	if m_cronCache == nil {
		m_cronCache := cron.New()
		m_cronCache.AddFunc("@every 20s", s.runJob)
		m_cronCache.Start()
	}
}

func (s *CronSv) runJob() {
	s.Lock()
	defer func() {
		s.Unlock()
	}()
	// 获取需要执行的任务,可用的，等待执行，且已到执行时间
	items := make([]model.Cron, 0)
	s.repo.Where("enabled=1 and status_id in (?) and ne_time<=?", []string{"waiting"}, utils.NewTime().Unix()).Preload("Endpoint").Find(&items)
	if len(items) == 0 {
		return
	}
	for i, _ := range items {
		s.jobHandle(items[i])
	}
	time.Sleep(time.Second)
}
func (s *CronSv) jobHandle(item model.Cron) {
	defer func() {
		if r := recover(); r != nil {
			s.runJobItemReset(item)
		}
	}()
	item.NumRun = item.NumRun + 1
	s.repo.Model(model.Cron{}).Where("id=?", item.ID).Updates(map[string]interface{}{"StatusID": "running", "NumRun": item.NumRun, "LastStatusID": "running", "LastTime": utils.NewTime()})

	s.logSv.Create(model.Log{NodeID: item.ID, NodeType: "Cron", Level: model.LOG_LEVEL_INFO, Msg: fmt.Sprintf("开始执行计划任务：%s", item.Name)}) //log

	s.createLog(&model.CronLog{EntID: item.EntID, CronID: item.ID, Title: "开始执行", StatusID: "running"})

	go s.runJobItem(item)
}
func (s *CronSv) runJobItem(item model.Cron) {
	defer func() {
		if r := recover(); r != nil {
			s.runJobItemReset(item)
		}
	}()
	tokenContext := item.Context
	if item.ClientID != "" {
		tokenContext.SetClientID(item.ClientID)
	}
	if item.UserID != "" {
		tokenContext.SetUserID(item.UserID)
	}
	if item.EntID != "" {
		tokenContext.SetEntID(item.EntID)
	}
	client := &http.Client{}

	if item.Endpoint == nil || item.Endpoint.Host == "" || item.Endpoint.Code == "" {
		s.runJobItemFailed(item, glog.Error("没有配置可执行接口"))
		return
	}

	postBody := item.Body
	postUri := ""
	if addr, err := reg.GetServerAddr(item.Endpoint.Host); err != nil {
		s.runJobItemFailed(item, err)
		return
	} else {
		remoteUrl, _ := url.Parse(addr)
		remoteUrl.Path = item.Endpoint.Path
		postUri = remoteUrl.String()
	}
	req, err := http.NewRequest("POST", postUri, bytes.NewBuffer([]byte(postBody)))
	if err != nil {
		s.runJobItemFailed(item, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenContext.ToTokenString())
	req.Header.Set("JOB_ID", item.ID)
	req.Header.Set("JOB_CODE", item.Code)
	req.Header.Set("USER_ID", tokenContext.UserID())
	req.Header.Set("ENT_ID", tokenContext.EntID())
	resp, err := client.Do(req)
	if err != nil {
		s.runJobItemFailed(item, err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.runJobItemFailed(item, err)
		return
	}
	resBodyObj := resBodyDTO{}
	if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
		glog.Error(utils.ToString(resBody))
		s.runJobItemFailed(item, err)
		return
	}
	if resp.StatusCode != 200 || resBodyObj.Msg != "" {
		s.runJobItemFailed(item, errors.New(resBodyObj.Msg, 3020))
		return
	}
	s.runJobItemSucceed(item)
}
func (s *CronSv) runJobItemFailed(item model.Cron, err error) {

	s.createLog(&model.CronLog{EntID: item.EntID, CronID: item.ID, Title: "执行计划出错", StatusID: "failed", Msg: err.Error()})

	s.logSv.Create(model.Log{NodeID: item.ID, NodeType: "Cron", Level: model.LOG_LEVEL_INFO, Msg: fmt.Sprintf("执行计划出错：%s", err.Error())}) //log
	updates := make(map[string]interface{})
	item.NumPeriod = item.NumPeriod + 1
	item.NumFailed = item.NumFailed + 1
	updates["NumFailed"] = item.NumFailed
	updates["NumPeriod"] = item.NumPeriod
	updates["LastStatusID"] = "failed"
	updates["LastMsg"] = err.Error()
	if item.UnitID == "once" {
		if item.NumPeriod > item.Retry { //大于重试次数，则停止
			updates["StatusID"] = "completed"
			s.logSv.Create(model.Log{NodeID: item.ID, NodeType: "Cron", Level: model.LOG_LEVEL_INFO, Msg: fmt.Sprintf("执行次数大于重试次数，设置状态为完成!")}) //log
		} else {
			updates["StatusID"] = "waiting"
		}
	} else {
		updates["StatusID"] = "waiting"
		if item.NumPeriod > item.Retry {
			s.logSv.Create(model.Log{NodeID: item.ID, NodeType: "Cron", Level: model.LOG_LEVEL_INFO, Msg: fmt.Sprintf("执行次数大于重试次数，重置下次执行时间!")}) //log
			currTime := time.Unix(item.NeTime, 0)
			if item.UnitID == "month" {
				currTime = currTime.AddDate(0, 1, 0)
			} else if item.UnitID == "week" {
				currTime = currTime.AddDate(0, 0, 7)
			} else if item.UnitID == "day" {
				currTime = currTime.AddDate(0, 0, 1)
			} else if item.UnitID == "hour" {
				currTime = currTime.Add(1 * time.Hour)
			} else {
				currTime = currTime.AddDate(1, 0, 0)
			}
			updates["NumPeriod"] = 0
			updates["NeTime"] = currTime.Unix()
		}
	}
	s.repo.Model(model.Cron{}).Where("id=?", item.ID).Updates(updates)
}
func (s *CronSv) runJobItemSucceed(item model.Cron) {

	s.logSv.Create(model.Log{NodeID: item.ID, NodeType: "Cron", Level: model.LOG_LEVEL_INFO, Msg: fmt.Sprintf("执行计划成功")}) //log
	s.createLog(&model.CronLog{EntID: item.EntID, CronID: item.ID, Title: "执行成功", StatusID: "succeed"})

	updates := make(map[string]interface{})
	item.NumPeriod = item.NumPeriod + 1
	item.NumSuccess = item.NumSuccess + 1

	updates["NumSuccess"] = item.NumSuccess
	updates["NumPeriod"] = item.NumPeriod
	updates["LastStatusID"] = "succeed"
	updates["LastMsg"] = ""
	if item.UnitID == "once" {
		updates["StatusID"] = "completed"
	} else {
		s.logSv.Create(model.Log{NodeID: item.ID, NodeType: "Cron", Level: model.LOG_LEVEL_INFO, Msg: fmt.Sprintf("重置下次执行时间!")}) //log
		updates["StatusID"] = "waiting"
		currTime := time.Unix(item.NeTime, 0)
		if item.UnitID == "month" {
			currTime = currTime.AddDate(0, 1, 0)
		} else if item.UnitID == "week" {
			currTime = currTime.AddDate(0, 0, 7)
		} else if item.UnitID == "day" {
			currTime = currTime.AddDate(0, 0, 1)
		} else if item.UnitID == "hour" {
			currTime = currTime.Add(1 * time.Hour)
		} else {
			currTime = currTime.AddDate(1, 0, 0)
		}
		updates["NumPeriod"] = 0
		updates["NeTime"] = currTime.Unix()
	}
	s.repo.Model(model.Cron{}).Where("id=?", item.ID).Updates(updates)
}
func (s *CronSv) runJobItemReset(item model.Cron) {
	s.repo.Model(model.Cron{}).Where("id=? and status_id=?", item.ID, "running").Updates(map[string]interface{}{"StatusID": "waiting", "LastMsg": "被异常中断"})
	s.logSv.Create(model.Log{NodeID: item.ID, NodeType: "Cron", Level: model.LOG_LEVEL_INFO, Msg: fmt.Sprintf("被异常中断")}) //log

	s.createLog(&model.CronLog{EntID: item.EntID, CronID: item.ID, Title: "被异常中断", StatusID: "completed"})
}
func (s *CronSv) createLog(item *model.CronLog) {
	if item.ID != "" {
		item.ID = utils.GUID()
	}
	s.repo.Create(item)
}

/**
创建一个调度任务
*/
func (s *CronSv) Create(entID string, item *model.Cron) (*model.Cron, error) {
	if item.UnitID == "" {
		return nil, errors.ParamsRequired("频率")
	}
	if item.Code == "" {
		item.Code = utils.GUID()
	}
	if item.Cycle < 1 {
		item.Cycle = 1
	}
	if item.FmTime.IsZero() {
		item.FmTime = utils.NewTime()
	}
	if item.FmTime.Unix() < utils.NewTimePtr().Unix() {
		item.FmTime = utils.NewTime()
	}
	//相同的 tag+url+UnitID 只能出现一次
	if item.Tag != "" {
		count := 0
		if s.repo.Model(model.Cron{}).Where("enabled=1 and status_id in (?) tag=? and unit_id=? and ent_id=?", []string{"waiting", "running"}, item.Tag, item.UnitID, item.EntID).Count(&count); count > 0 {
			return nil, errors.ExistError("频率相同的任务")
		}
	}
	if item.EndpointID == "" && item.Endpoint != nil {
		if item.Endpoint.ID != "" {
			item.EndpointID = item.Endpoint.ID
		} else if item.Endpoint.Code != "" {
			if e, err := NewDtiSv(s.repo).GetLocalBy(item.Endpoint.Code); err != nil {
				return nil, err
			} else {
				item.EndpointID = e.ID
			}
		}
	}
	if item.EndpointID == "" {
		return nil, glog.Error("endpoint为空")
	}
	//下次执行时间
	item.NeTime = item.FmTime.Unix()
	item.EntID = entID
	item.Enabled = utils.SBool_True
	item.StatusID = "waiting"
	item.EntID = entID
	item.ID = utils.GUID()
	if err := s.repo.Create(item).Error; err != nil {
		return nil, err
	}
	return s.GetBy(item.ID)
}
func (s *CronSv) GetBy(id string) (*model.Cron, error) {
	item := model.Cron{}
	if err := s.repo.Model(item).Where("id=?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *CronSv) Destroy(entID string, ids []string) error {
	systemDatas := 0
	s.repo.Model(model.Cron{}).Where("id in (?) and `system`=1", ids).Count(&systemDatas)
	if systemDatas > 0 {
		return utils.NewError("系统预制不能删除!")
	}
	if err := s.repo.Delete(model.Cron{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}
