package rules

import (
	"fmt"
	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/md"
	"github.com/ggoop/mdf/mof"
	"github.com/ggoop/mdf/repositories"
)

type CommonEnableBatch struct {
	repo *repositories.MysqlRepo
}

func NewCommonEnableBatch(repo *repositories.MysqlRepo) *CommonEnableBatch {
	return &CommonEnableBatch{repo}
}
func (s CommonEnableBatch) Exec(req *mof.ReqContext, res *mof.ResContext) error {
	if len(req.IDS) == 0 {
		return glog.Error("缺少 ID 参数！")
	}
	//查找实体信息
	entity := md.GetEntity(req.MainEntity)
	if entity == nil {
		return glog.Error("找不到实体！")
	}
	if err := s.repo.Exec(fmt.Sprintf("update %s set enabled=1 where id in (?)", entity.TableName), req.IDS).Error; err != nil {
		return err
	}
	return nil
}
func (s CommonEnableBatch) GetRule() mof.RuleRegister {
	return mof.RuleRegister{Code: "enable_batch", Domain: "common"}
}
