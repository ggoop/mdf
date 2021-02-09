package service

import (
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/repositories"
	"github.com/ggoop/mdf/utils"
)

type LogSv struct {
	repo *repositories.MysqlRepo
}

/**
* 创建服务实例
 */
func NewLogSv(repo *repositories.MysqlRepo) *LogSv {
	return &LogSv{repo: repo}
}

func (s *LogSv) CreateOp(item model.Log) {
	item.Type = "op"
	s.Create(item)
}
func (s *LogSv) CreateLog(item model.Log) {
	item.Type = "log"
	s.Create(item)
}
func (s *LogSv) Create(item model.Log) {
	item.ID = utils.GUID()
	if item.Type == "" {
		item.Type = "log"
	}
	if err := s.repo.Create(&item).Error; err != nil {
		glog.Error(err)
	}
	glog.Errorf("%s-%s: %v", item.NodeType, item.NodeID, item.Msg)
}

func (s *LogSv) Log(item model.Log) {
	glog.Errorf("%s-%s: %v", item.NodeType, item.NodeID, item.Msg)
}
