package services

import (
	"encoding/json"
	"fmt"
	"github.com/ggoop/mdf/bootstrap/errors"
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/db"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/reg"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/ggoop/mdf/framework/md"

	"github.com/ggoop/mdf/utils"
)

//interface
type IDtiSv interface {
	GetLocalBy(code string) (*model.DtiLocal, error)
}

func DtiSv() IDtiSv {
	return dtiSv
}

var dtiSv IDtiSv = newDtiSvImlp()

//impl
type dtiSvImpl struct {
}

func newDtiSvImlp() *dtiSvImpl {
	return &dtiSvImpl{}
}
func (s *dtiSvImpl) GetLocalBy(code string) (*model.DtiLocal, error) {
	item := model.DtiLocal{}
	if err := db.Default().Where("id=? or code=? ", code, code).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *dtiSvImpl) UpdateOrCreateLocal(item model.DtiLocal) (*model.DtiLocal, error) {
	old := model.DtiLocal{}
	if !utils.StringIsCode(item.Code) {
		return nil, errors.CodeError("Code")
	}
	db.Default().Where("code=?", item.Code).Take(&old)
	if old.ID != "" {
		item.ID = old.ID
		updatas := make(map[string]interface{})
		if old.ProductCode != item.ProductCode {
			updatas["ProductCode"] = item.ProductCode
		}
		if old.Name != item.Name {
			updatas["Name"] = item.Name
		}
		if old.Memo != item.Memo {
			updatas["Memo"] = item.Memo
		}
		if old.Host != item.Host {
			updatas["Host"] = item.Host
		}
		if old.Path != item.Path {
			updatas["Path"] = item.Path
		}
		if old.Tags != item.Tags {
			updatas["Tags"] = item.Tags
		}
		if item.System.Valid() {
			updatas["System"] = item.System
		}
		if item.Enabled.Valid() {
			updatas["Enabled"] = item.Enabled
		}
		if len(updatas) > 0 {
			db.Default().Model(&old).Where("id=?", old.ID).Updates(updatas)
		}
	} else {
		if item.ID == "" {
			item.ID = utils.GUID()
		}
		db.Default().Create(&item)
	}
	s.UpdateOrCreateLocalParams(item.ID, item.Params)

	return s.GetLocalBy(item.ID)
}
func (s *dtiSvImpl) UpdateOrCreateLocalParams(localID string, params []model.DtiLocalParam) error {
	if len(params) == 0 {
		return nil
	}
	codes := make([]string, 0)
	for i, _ := range params {
		p := params[i]
		p.LocalID = localID
		if p.Sequence == 0 {
			p.Sequence = i
		}
		old := model.DtiLocalParam{}
		db.Default().Model(old).Take(&old, "local_id=? and code=?", localID, p.Code)
		if old.ID != "" {
			updatas := make(map[string]interface{})
			if old.Name != p.Name {
				updatas["Name"] = p.Name
			}
			if old.Memo != p.Memo {
				updatas["Memo"] = p.Memo
			}
			if old.TypeID != p.TypeID {
				updatas["TypeID"] = p.TypeID
			}
			if old.Value != p.Value {
				updatas["Value"] = p.Value
			}
			if old.ValueDef != p.ValueDef {
				updatas["ValueDef"] = p.ValueDef
			}
			if old.Sequence != p.Sequence {
				updatas["Sequence"] = p.Sequence
			}
			if p.Hidden.Valid() {
				updatas["Hidden"] = p.Hidden
			}
			if p.Required.Valid() {
				updatas["Required"] = p.Required
			}
			if len(updatas) > 0 {
				db.Default().Model(&old).Where("id=?", old.ID).Updates(updatas)
			}
		} else {
			db.Default().Create(&p)
		}
		codes = append(codes, p.Code)
	}
	if len(codes) > 0 {
		db.Default().Delete(model.DtiLocalParam{}, "local_id=? and code not in (?)", localID, codes)
	}
	return nil
}

//node
func (s *dtiSvImpl) GetNodeBy(entID, id string) (*model.DtiNode, error) {
	item := model.DtiNode{}
	if err := db.Default().Where("ent_id=?", entID).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *dtiSvImpl) UpdateOrCreateNode(entID string, item model.DtiNode) (*model.DtiNode, error) {
	if !utils.StringIsCode(item.Code) {
		return nil, errors.CodeError("Code")
	}
	if strings.Contains(item.Host, ".") {
		if url, err := url.Parse(item.Host); err != nil {
			return nil, err
		} else if url.Scheme == "" {
			url.Scheme = "http"
			item.Host = url.String()
		}
	}
	old := model.DtiNode{}
	if item.ID != "" {
		db.Default().Where("ent_id=?", entID).Where("id=?", item.ID).Take(&old)
	}
	if old.ID != "" {
		item.ID = old.ID
		updatas := make(map[string]interface{})
		if old.Name != item.Name && item.Name != "" {
			updatas["Name"] = item.Name
		}
		if old.Memo != item.Memo && item.Memo != "" {
			updatas["Memo"] = item.Memo
		}
		if item.System.Valid() && item.System.NotEqual(old.System) {
			updatas["System"] = item.System
		}
		if item.Public.Valid() && item.Public.NotEqual(old.Public) {
			updatas["Public"] = item.Public
		}
		if old.Host != item.Host && item.Host != "" {
			updatas["Host"] = item.Host
		}
		if len(updatas) > 0 {
			db.Default().Model(&old).Updates(updatas)
		}

	} else {
		item.EntID = entID
		if item.ID == "" {
			item.ID = utils.GUID()
		}
		db.Default().Create(&item)
		s.CompileTemplate(entID, item)
	}
	return s.GetNodeBy(entID, item.ID)
}
func (s *dtiSvImpl) DeleteNode(entID string, ids []string) error {
	if _, names := md.MDSv().QuotedBy(&model.DtiNode{}, ids); names != nil {
		return errors.IsQuoted(names...)
	}
	db.Default().Where("ent_id = ? and id in(?)", entID, ids).Delete(model.DtiNode{})
	return nil
}

func (s *dtiSvImpl) CompileTemplate(entID string, node model.DtiNode) error {
	if node.TemplateID == "" {
		return nil
	}
	dtiRemotes := make([]model.DtiRemote, 0)
	db.Default().Where("node_id=?", node.TemplateID).Find(&dtiRemotes)

	newItems := make([]interface{}, 0)
	count := 0
	for _, d := range dtiRemotes {
		if db.Default().Model(model.DtiRemote{}).Where("ent_id=?", entID).Where("code=?", d.Code).Count(&count); count > 0 {
			continue
		}
		item := model.DtiRemote{LocalID: d.LocalID, MethodID: d.MethodID, Code: d.Code, Name: d.Name, Enabled: utils.SBool_True}
		item.ID = utils.GUID()
		item.Path = d.Path
		item.Body = d.Body
		item.Header = d.Header
		item.Memo = d.Memo
		item.Query = d.Query
		item.CreatedAt = utils.TimeNow()
		item.EntID = entID
		item.NodeID = node.ID
		newItems = append(newItems, item)
	}
	if len(newItems) > 0 {
		db.Default().BatchInsert(newItems)
	}
	return nil
}

//params
func (s *dtiSvImpl) GetParamBy(entID, id string) (*model.DtiParam, error) {
	item := model.DtiParam{}
	if err := db.Default().Where("ent_id=? and id=?", entID, id).Preload("Node").Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *dtiSvImpl) UpdateOrCreateParam(entID string, item model.DtiParam) (*model.DtiParam, error) {
	old := model.DtiParam{}
	if !utils.StringIsCode(item.Code) {
		return nil, errors.CodeError("Code")
	}
	if item.ID != "" {
		db.Default().Where("ent_id=? and id=?", entID, item.ID).Take(&old)
	}
	if old.ID != "" {
		item.ID = old.ID
		updatas := make(map[string]interface{})
		if old.Name != item.Name && item.Name != "" {
			updatas["Name"] = item.Name
		}
		if old.Value != item.Value {
			updatas["Value"] = item.Value
		}
		if old.TypeID != item.TypeID && item.TypeID != "" {
			updatas["TypeID"] = item.TypeID
		}
		if old.NodeID != item.NodeID && item.NodeID != "" {
			updatas["NodeID"] = item.NodeID
		}
		if len(updatas) > 0 {
			db.Default().Model(&item).Where("id=?", old.ID).Updates(updatas)
		}
	} else {
		item.EntID = entID
		item.ID = utils.GUID()
		db.Default().Create(&item)
	}
	return s.GetParamBy(entID, item.ID)
}
func (s *dtiSvImpl) DeleteParam(entID string, ids []string) error {
	if _, names := md.MDSv().QuotedBy(&model.DtiParam{}, ids); names != nil {
		return errors.IsQuoted(names...)
	}
	db.Default().Where("ent_id = ? and id in(?)", entID, ids).Delete(model.DtiParam{})
	return nil
}

//remotes
func (s *dtiSvImpl) GetRemoteBy(entID, id string) (*model.DtiRemote, error) {
	item := model.DtiRemote{}
	if err := db.Default().Where("ent_id=?", entID).Preload("Node").Preload("Local").Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *dtiSvImpl) UpdateOrCreateRemote(entID string, item model.DtiRemote) (*model.DtiRemote, error) {
	old := model.DtiRemote{}
	if !utils.StringIsCode(item.Code) {
		return nil, errors.CodeError("Code")
	}

	if item.ID != "" {
		db.Default().Where("ent_id=?", entID).Where("id=?", item.ID).Take(&old)
	}
	if old.ID != "" {
		item.ID = old.ID
		updatas := make(map[string]interface{})
		if item.Enabled.Valid() && !old.Enabled.Equal(item.Enabled) {
			updatas["Enabled"] = item.Enabled
		}
		if old.Name != item.Name && item.Name != "" {
			updatas["Name"] = item.Name
		}
		if len(updatas) > 0 {
			db.Default().Model(&item).Updates(updatas)
		}
	} else {
		item.EntID = entID
		item.ID = utils.GUID()
		db.Default().Create(&item)
	}
	return s.GetRemoteBy(entID, item.ID)
}
func (s *dtiSvImpl) DeleteRemote(entID string, ids []string) error {
	if _, names := md.MDSv().QuotedBy(&model.DtiRemote{}, ids); names != nil {
		return errors.IsQuoted(names...)
	}
	db.Default().Where("ent_id = ? and id in(?)", entID, ids).Delete(model.DtiRemote{})
	return nil
}
func (s *dtiSvImpl) GetRemotes(entID string, dtiIDs []string) []model.DtiRemote {
	dtiRemotes := make([]model.DtiRemote, 0)
	query := db.Default().Where("ent_id=?", entID).Preload("Node").Preload("Local").Order("`sequence`")
	if len(dtiIDs) > 0 {
		query = query.Where("id in (?) or code in (?)", dtiIDs, dtiIDs)
	}
	query.Find(&dtiRemotes)
	return dtiRemotes
}
func (s *dtiSvImpl) RunRemotes(entID string, dtiIDs []string, params map[string]string, ctx *utils.TokenContext) error {
	dtiRemotes := make([]model.DtiRemote, 0)
	if err := db.Default().Where("ent_id=? and (id in (?) or code in (?))", entID, dtiIDs, dtiIDs).Preload("Node").Preload("Local").Order("`sequence`").Find(&dtiRemotes).Error; err != nil {
		return err
	}
	//分组
	groups := make(map[int][]model.DtiRemote)
	keys := make([]int, 0)
	for i, item := range dtiRemotes {
		if m, ok := groups[item.Sequence]; ok {
			m = append(m, dtiRemotes[i])
			groups[item.Sequence] = m
		} else {
			keys = append(keys, item.Sequence)
			m = make([]model.DtiRemote, 0)
			m = append(m, dtiRemotes[i])
			groups[item.Sequence] = m
		}
	}
	//并行执行
	for _, k := range keys {
		if m, ok := groups[k]; ok && len(m) > 0 {
			var n sync.WaitGroup
			for i, _ := range m {
				n.Add(1)
				go func(entID string, dtiRemote model.DtiRemote, params map[string]string, ctx *utils.TokenContext) {
					s.runRemote(entID, dtiRemote, params, ctx)
					n.Done()
				}(entID, m[i], params, ctx)
			}
			n.Wait()
		}
	}
	//按分组执行
	return nil
}
func (s *dtiSvImpl) runRemote(entID string, dtiRemote model.DtiRemote, params map[string]string, ctx *utils.TokenContext) error {
	defer func() {
		if r := recover(); r != nil {
			s.run_log_failed(entID, dtiRemote, fmt.Errorf("未知异常"))
		}
	}()
	// 开始
	s.run_log_begin(entID, dtiRemote, fmt.Sprintf("正在开始执行接口：%s", dtiRemote.Name))
	if dtiRemote.Node == nil || dtiRemote.Local == nil {
		return s.run_log_failed(entID, dtiRemote, fmt.Errorf("接口分类或者接口为空"))
	}
	// 准备请求参数
	if dtiRemote.MethodID == "" {
		return s.run_log_failed(entID, dtiRemote, fmt.Errorf("请求方式为空"))
	}
	// 准备请求参数
	if dtiRemote.Node.Host == "" {
		return s.run_log_failed(entID, dtiRemote, fmt.Errorf("接口主机地址为空"))
	}
	// 获取参数
	dtiParams := make([]model.DtiParam, 0)
	db.Default().Where("ent_id=? and (node_id=? or node_id='' or node_id is null)", entID, dtiRemote.NodeID).Find(&dtiParams)
	if len(params) > 0 {
		for k, v := range params {
			dtiParams = append(dtiParams, model.DtiParam{Code: k, Value: v})
		}
	}
	var remoteUrl *url.URL
	if dtiRemote.Node.Host != "" && strings.Index(strings.ToLower(dtiRemote.Node.Host), "http") == 0 {
		remoteUrl, _ = url.Parse(dtiRemote.Node.Host)
	} else if dtiRemote.Node.Host != "" {
		if u, _ := reg.FindServerByCode(dtiRemote.Node.Host); u != nil && u.Address != "" {
			remoteUrl, _ = url.Parse(u.Address)
		}
	}
	if dtiRemote.Path != "" {
		if strings.Index(dtiRemote.Path, "/") == 0 {
			remoteUrl.Path = remoteUrl.Path + dtiRemote.Path
		} else {
			remoteUrl.Path = remoteUrl.Path + "/" + dtiRemote.Path
		}
	}
	remoteUrl, _ = url.Parse(remoteUrl.String())
	remoteQuery := remoteUrl.Query()
	remoteQuery.Set("EntID", entID)
	//参数转为URL参数
	if len(params) > 0 {
		for k, v := range params {
			remoteQuery.Set(k, v)
		}
	}
	remoteUrl.RawQuery = remoteQuery.Encode()

	remoteUrlString := s.run_parseParams(remoteUrl.String(), dtiParams)
	s.run_log(entID, dtiRemote, fmt.Sprintf("远程节点路径为:%v", remoteUrlString))
	// 解析body
	dtiRemote.Body = s.run_parseParams(dtiRemote.Body, dtiParams)
	dtiRemote.MethodID = strings.ToUpper(dtiRemote.MethodID)
	if dtiRemote.MethodID == "" {
		dtiRemote.MethodID = "POST"
	}
	req, err := http.NewRequest(dtiRemote.MethodID, remoteUrlString, strings.NewReader(dtiRemote.Body))
	if err != nil {
		return s.run_log_failed(entID, dtiRemote, err)
	}
	// 解析header
	dtiRemote.Header = s.run_parseParams(dtiRemote.Header, dtiParams)
	if dtiRemote.Header != "" {
		header := make(map[string]string)
		json.Unmarshal([]byte(dtiRemote.Header), &header)
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("REMOTE_ID", dtiRemote.ID)
	req.Header.Set("REMOTE_CODE", dtiRemote.Code)
	req.Header.Set("USER_ID", ctx.UserID())
	req.Header.Set("ENT_ID", ctx.EntID())

	s.run_log(entID, dtiRemote, fmt.Sprintf("请求方式:%v，请求Body:%v", dtiRemote.MethodID, dtiRemote.Body))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return s.run_log_failed(entID, dtiRemote, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return s.run_log_failed(entID, dtiRemote, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return s.run_log_failed(entID, dtiRemote, fmt.Errorf("远程请求出错了, status:%v ,:%s", resp.StatusCode, string(body)))
	}
	s.run_log(entID, dtiRemote, fmt.Sprintf("接收到远程接口数据 %v 字节,ContentLength:%d,开始投递", len(body), resp.ContentLength))
	if err := s.run_handResult(entID, dtiRemote, body, params, ctx); err != nil {
		return s.run_log_failed(entID, dtiRemote, err)
	}
	s.run_log_succeed(entID, dtiRemote, "接口执行成功")

	return nil
}

func (s *dtiSvImpl) run_handResult(entID string, dtiRemote model.DtiRemote, body []byte, params map[string]string, ctx *utils.TokenContext) error {
	var remoteUrl *url.URL
	if dtiRemote.Local.Host != "" && strings.Index(strings.ToLower(dtiRemote.Local.Host), "http") == 0 {
		remoteUrl, _ = url.Parse(dtiRemote.Local.Host)
	} else if dtiRemote.Local.Host != "" {
		if u, _ := reg.FindServerByCode(dtiRemote.Local.Host); u != nil && u.Address != "" {
			remoteUrl, _ = url.Parse(u.Address)
		}
	}
	if remoteUrl.Host == "" {
		return s.run_log_failed(entID, dtiRemote, glog.Errorf("找不到%s地址", dtiRemote.Local.Code))
	}
	if dtiRemote.Local.Path != "" {
		if strings.Index(dtiRemote.Local.Path, "/") == 0 {
			remoteUrl.Path = remoteUrl.Path + dtiRemote.Local.Path
		} else {
			remoteUrl.Path = remoteUrl.Path + "/" + dtiRemote.Local.Path
		}
	}
	remoteUrl, _ = url.Parse(remoteUrl.String())
	remoteQuery := remoteUrl.Query()
	//参数转为URL参数
	if len(params) > 0 {
		for k, v := range params {
			remoteQuery.Set(k, v)
		}
	}
	remoteUrl.RawQuery = remoteQuery.Encode()

	s.run_log(entID, dtiRemote, fmt.Sprintf("处理结果数据地址:%v", remoteUrl.String()))

	req, err := http.NewRequest("POST", remoteUrl.String(), strings.NewReader(string(body)))
	if err != nil {
		return s.run_log_failed(entID, dtiRemote, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("REMOTE_ID", dtiRemote.ID)
	req.Header.Set("REMOTE_CODE", dtiRemote.Code)
	req.Header.Set("USER_ID", ctx.UserID())
	req.Header.Set("ENT_ID", ctx.EntID())

	req.Header.Set("Authorization", ctx.ToTokenString())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return s.run_log_failed(entID, dtiRemote, err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return s.run_log_failed(entID, dtiRemote, err)
	}
	if resp.StatusCode != 200 {
		return s.run_log_failed(entID, dtiRemote, fmt.Errorf("处理结果请求出错了, status:%v ,:%s", resp.StatusCode, string(resBody)))
	}
	return nil
}
func (s *dtiSvImpl) run_parseParams(oldString string, params []model.DtiParam) string {
	if oldString == "" {
		return oldString
	}
	for _, p := range params {
		oldString = strings.Replace(oldString, "#"+p.Code+"#", p.Value, -1)
		oldString = strings.Replace(oldString, "{"+p.Code+"}", p.Value, -1)
	}
	return oldString
}
func (s *dtiSvImpl) run_getParamValue(code string, params []model.DtiParam) string {
	code = strings.ToLower(code)
	for _, p := range params {
		if strings.ToLower(p.Code) == code {
			return p.Value
		}
	}
	return ""
}
func (s *dtiSvImpl) run_log_begin(entID string, dti model.DtiRemote, msg string) {
	db.Default().Model(model.DtiRemote{}).Where("ent_id=? and id=?", entID, dti.ID).Updates(map[string]interface{}{"StatusID": "running", "FmDate": utils.TimeNowPtr(), "Msg": msg})
	s.run_log(entID, dti, msg)
}
func (s *dtiSvImpl) run_log_succeed(entID string, dti model.DtiRemote, msg string) {
	db.Default().Model(model.DtiRemote{}).Where("ent_id=? and id=?", entID, dti.ID).Updates(map[string]interface{}{"StatusID": "succeed", "ToDate": utils.TimeNowPtr(), "Msg": msg})
	s.run_log(entID, dti, msg)
}
func (s *dtiSvImpl) run_log_failed(entID string, dti model.DtiRemote, err error) error {
	db.Default().Model(model.DtiRemote{}).Where("ent_id=? and id=?", entID, dti.ID).Updates(map[string]interface{}{"StatusID": "failed", "ToDate": utils.TimeNowPtr(), "Msg": err.Error()})
	s.run_log(entID, dti, err.Error())
	return err
}
func (s *dtiSvImpl) run_log(entID string, dti model.DtiRemote, msg string) {
	db.Default().Model(model.DtiRemote{}).Where("ent_id=? and id=?", entID, dti.ID).Updates(map[string]interface{}{"Msg": msg})
	LogSv().Create(model.Log{NodeID: dti.Code, NodeType: "dti", EntID: entID, Level: model.LOG_LEVEL_ERROR, Msg: msg})
}
