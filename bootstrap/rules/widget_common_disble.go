package rules

import (
	"fmt"
	"github.com/ggoop/mdf/utils"

	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
)

type CommonDisable struct {
	repo *repositories.MysqlRepo
}

func NewCommonDisable(repo *repositories.MysqlRepo) *CommonDisable {
	return &CommonDisable{repo}
}
func (s CommonDisable) Exec(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext) error {
	if req.ID == "" {
		return glog.Error("缺少 ID 参数！")
	}
	//查找实体信息
	entity := md.GetEntity(req.Entity)
	if entity == nil {
		return glog.Error("找不到实体！")
	}
	if err := s.repo.Exec(fmt.Sprintf("update %s set enabled=0 where id =?", entity.TableName), req.ID).Error; err != nil {
		return err
	}
	return nil
}
func (s CommonDisable) GetRule() md.RuleRegister {
	return md.RuleRegister{Code: "disable", OwnerType: md.RuleType_Widget, OwnerID: "common"}
}
