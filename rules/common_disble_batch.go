package rules

import (
	"fmt"
	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/md"
	"github.com/ggoop/mdf/mof"
	"github.com/ggoop/mdf/repositories"
)

type CommonDisableBatch struct {
	repo *repositories.MysqlRepo
}

func NewCommonDisableBatch(repo *repositories.MysqlRepo) *CommonDisableBatch {
	return &CommonDisableBatch{repo}
}
func (s CommonDisableBatch) Exec(req *mof.ReqContext, res *mof.ResContext) error {
	if len(req.IDS) == 0 {
		return glog.Error("缺少 ID 参数！")
	}
	//查找实体信息
	entity := md.GetEntity(req.MainEntity)
	if entity == nil {
		return glog.Error("找不到实体！")
	}
	if err := s.repo.Exec(fmt.Sprintf("update %s set enabled=0 where id in (?)", entity.TableName), req.IDS).Error; err != nil {
		return err
	}
	return nil
}
func (s CommonDisableBatch) GetRule() mof.RuleRegister {
	return mof.RuleRegister{Code: "disable_batch", Domain: "common"}
}
