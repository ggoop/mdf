package rules

import (
	"fmt"

	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
)

type CommonEnable struct {
	repo *repositories.MysqlRepo
}

func NewCommonEnable(repo *repositories.MysqlRepo) *CommonEnable {
	return &CommonEnable{repo}
}
func (s CommonEnable) Exec(req *md.ReqContext, res *md.ResContext) error {
	if req.ID == "" {
		return glog.Error("缺少 ID 参数！")
	}
	//查找实体信息
	entity := md.GetEntity(req.Entity)
	if entity == nil {
		return glog.Error("找不到实体！")
	}
	if err := s.repo.Exec(fmt.Sprintf("update %s set enabled=1 where id =?", entity.TableName), req.ID).Error; err != nil {
		return err
	}
	return nil
}
func (s CommonEnable) GetRule() md.RuleRegister {
	return md.RuleRegister{Code: "enable", OwnerType: md.RuleType_Widget, OwnerID: "common"}
}
