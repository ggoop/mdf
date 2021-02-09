package rules

import (
	"fmt"

	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/framework/mof"
	"github.com/ggoop/mdf/utils"
)

type CommonDeleteBatch struct {
	repo *repositories.MysqlRepo
}

func NewCommonDeleteBatch(repo *repositories.MysqlRepo) *CommonDeleteBatch {
	return &CommonDeleteBatch{repo}
}
func (s CommonDeleteBatch) Exec(req *mof.ReqContext, res *mof.ResContext) error {
	if len(req.IDS) == 0 {
		return glog.Error("缺少 ID 参数！")
	}
	//查找实体信息
	entity := md.GetEntity(req.Entity)
	if entity == nil {
		return glog.Error("找不到实体！")
	}
	if df := entity.GetField("DeletedAt"); df != nil {
		if err := s.repo.Exec(fmt.Sprintf("update %s set %s=? where id in (?)", entity.TableName, df.DbName), utils.NewTime(), req.IDS).Error; err != nil {
			return err
		}
	} else {
		if err := s.repo.Exec(fmt.Sprintf("delete from %s where id in (?)", entity.TableName), req.IDS).Error; err != nil {
			return err
		}
	}
	return nil
}
func (s CommonDeleteBatch) GetRule() mof.RuleRegister {
	return mof.RuleRegister{Code: "delete_batch", Owner: "common"}
}
