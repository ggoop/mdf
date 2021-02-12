package model

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

const (
	CRON_UNIT_ENUM_once  = "once"
	CRON_UNIT_ENUM_month = "month"
	CRON_UNIT_ENUM_week  = "week"
	CRON_UNIT_ENUM_day   = "day"
	CRON_UNIT_ENUM_hour  = "hour"
)

type Cron struct {
	md.Model
	ClientID string             `gorm:"size:50"`
	EntID    string             `gorm:"size:50"`
	UserID   string             `gorm:"size:50"`
	Tag      string             `json:"tag"`
	Code     string             `json:"code"`
	Name     string             `json:"name"`
	Memo     string             `json:"memo"`
	StatusID string             `gorm:"size:50" json:"status_id"` //waiting,running,completed
	Status   *md.MDEnum         `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:sys.corn.status"`
	Enabled  utils.SBool        `gorm:"default:true" json:"enabled"`
	Context  utils.TokenContext `gorm:"type:text;name:上下文" json:"context"` //上下文参数

	//调度任务周期
	NeTime int64      `json:"ne_time"` //下次执行时间 ,unix
	FmTime utils.Time `json:"fm_time"` //开始
	ToTime utils.Time `json:"to_time"` //结束
	UnitID string     `json:"unit_id"` //sys.cron.unit,频率 执行一次,每月,每周,每天,每小时
	Unit   *md.MDEnum `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:sys.cron.unit"`
	Cycle  int        `json:"cycle"` //周期,1天，2月，4小时
	Spec   string     `json:"spec"`
	Retry  int        `json:"retry"`

	EndpointID string    `gorm:"size:50" json:"endpoint_id"`
	Endpoint   *DtiLocal `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Header     string    `gorm:"size:50" json:"header"`
	Body       string    `json:"body"`
	Query      string    `json:"query"`

	//执行情况
	NumRun       int         //执行次数
	NumSuccess   int         //成功次数
	NumFailed    int         //失败次数
	NumPeriod    int         //当前频率执行次数
	LastMsg      string      `json:"last_msg"`
	LastStatusID string      `json:"last_status_id"` //running,succeed,failed
	LastStatus   *md.MDEnum  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:sys.corn.status"`
	LastTime     utils.Time  `json:"last_time"`
	System       utils.SBool `gorm:"not null;default:0;name:系统的" json:"system"`
}

func (s Cron) TableName() string {
	return "sys_crons"
}
func (s *Cron) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".cron", Domain: MD_DOMAIN, Name: "计划任务"}
}

/**
日志
*/
type CronLog struct {
	md.Model
	EntID    string     `gorm:"size:50;index:idx_ent_type"`
	CronID   string     `gorm:"size:50;index:idx_ent_type" json:"cron_id"`
	Cron     *Cron      `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Title    string     `json:"title"`
	Msg      string     `gorm:"type:text" json:"msg"`
	StatusID string     `json:"status_id"` //running,succeed,failed
	Status   *md.MDEnum `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:sys.corn.status"`
}

func (s CronLog) TableName() string {
	return "sys_cron_logs"
}
func (s *CronLog) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".cron.log", Domain: MD_DOMAIN, Name: "任务日志"}
}
