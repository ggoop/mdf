package rules

import (
	"fmt"
	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/md"
	"github.com/ggoop/mdf/mof"
	"github.com/ggoop/mdf/repositories"
)

type CommonDisable struct {
	repo *repositories.MysqlRepo
}

func NewCommonDisable(repo *repositories.MysqlRepo) *CommonDisable {
	return &CommonDisable{repo}
}
func (s CommonDisable) Exec(req *mof.ReqContext, res *mof.ResContext) error {
	if req.ID == "" {
		return glog.Error("缺少 ID 参数！")
	}
	//查找实体信息
	entity := md.GetEntity(req.MainEntity)
	if entity == nil {
		return glog.Error("找不到实体！")
	}
	if err := s.repo.Exec(fmt.Sprintf("update %s set enabled=0 where id =?", entity.TableName), req.ID).Error; err != nil {
		return err
	}
	return nil
}
func (s CommonDisable) GetRule() mof.RuleRegister {
	return mof.RuleRegister{Code: "disable", Domain: "common"}
}