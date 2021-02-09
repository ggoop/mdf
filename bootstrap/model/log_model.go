package model

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

const LOG_LEVEL_ERROR = "error"
const LOG_LEVEL_WARN = "warn"
const LOG_LEVEL_INFO = "info"
const LOG_LEVEL_DEBUG = "debug"

type Log struct {
	md.Model
	EntID      string `gorm:"size:50;index:idx_ent_type"`
	Type       string `gorm:"not null;size:10;index:idx_ent_type;name:日志类型"` //login,op
	UserID     string `gorm:"size:50"`
	UserName   string `gorm:"size:50"`
	NodeType   string `gorm:"size:50;index:idx_ent_type" json:"node_type"`
	NodeID     string `gorm:"size:120;index:idx_ent_type" json:"node_id"`
	NodeAction string `gorm:"size:50;name:动作" json:"node_action"`
	DataID     string `gorm:"size:50" json:"data_id"`
	DataType   string `gorm:"size:50" json:"data_type"`
	Title      string `gorm:"size:50"`
	Msg        string `gorm:"type:text" json:"msg"`
	Level      string `gorm:"size:50" json:"level"` //error,warn,info,debug
	Status     string `gorm:"size:36;name:状态" json:"status"`
	ReqIP      string `gorm:"size:50"`
	ReqClient  string `gorm:"size:50;name:设备" json:"req_client"`
	ReqAgent   string `gorm:"size:500"`
}

func (s Log) Clone() Log {
	return Log{
		Type:       s.Type,
		EntID:      s.EntID,
		UserID:     s.UserID,
		UserName:   s.UserName,
		Title:      s.Title,
		NodeID:     s.NodeID,
		NodeAction: s.NodeAction,
		NodeType:   s.NodeType,
		DataID:     s.DataID,
		DataType:   s.DataType,
		Msg:        s.Msg,
		Level:      s.Level,
		Status:     s.Status,
		ReqIP:      s.ReqIP,
		ReqAgent:   s.ReqAgent,
	}
}
func (s Log) SetMsg(msg interface{}) Log {
	if msg == nil {
		return s
	} else if m, ok := msg.(string); ok {
		s.Msg = m
	} else if m, ok := msg.(error); ok {
		s.Msg = m.Error()
	} else if m, ok := msg.(utils.GError); ok {
		s.Msg = m.Error()
	}
	return s
}
func (s *Log) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".log", Domain: MD_DOMAIN, Name: "日志"}
}
func (s Log) TableName() string {
	return "sys_logs"
}
