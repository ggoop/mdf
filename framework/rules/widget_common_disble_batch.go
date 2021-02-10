package rules

import (
	"fmt"

	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
)

type CommonDisableBatch struct {
	repo *repositories.MysqlRepo
}

func NewCommonDisableBatch(repo *repositories.MysqlRepo) *CommonDisableBatch {
	return &CommonDisableBatch{repo}
}
func (s CommonDisableBatch) Exec(req *md.ReqContext, res *md.ResContext) error {
	if len(req.IDS) == 0 {
		return glog.Error("缺少 ID 参数！")
	}
	//查找实体信息
	entity := md.GetEntity(req.Entity)
	if entity == nil {
		return glog.Error("找不到实体！")
	}
	if err := s.repo.Exec(fmt.Sprintf("update %s set enabled=0 where id in (?)", entity.TableName), req.IDS).Error; err != nil {
		return err
	}
	return nil
}
func (s CommonDisableBatch) GetRule() md.RuleRegister {
	return md.RuleRegister{Code: "disable_batch", Owner: "common"}
}
