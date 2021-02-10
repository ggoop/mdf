package rules

import (
	"fmt"

	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
)

type CommonEnableBatch struct {
	repo *repositories.MysqlRepo
}

func NewCommonEnableBatch(repo *repositories.MysqlRepo) *CommonEnableBatch {
	return &CommonEnableBatch{repo}
}
func (s CommonEnableBatch) Exec(req *md.ReqContext, res *md.ResContext) error {
	if len(req.IDS) == 0 {
		return glog.Error("缺少 ID 参数！")
	}
	//查找实体信息
	entity := md.GetEntity(req.Entity)
	if entity == nil {
		return glog.Error("找不到实体！")
	}
	if err := s.repo.Exec(fmt.Sprintf("update %s set enabled=1 where id in (?)", entity.TableName), req.IDS).Error; err != nil {
		return err
	}
	return nil
}
func (s CommonEnableBatch) GetRule() md.RuleRegister {
	return md.RuleRegister{Code: "enable_batch", Owner: "common"}
}
