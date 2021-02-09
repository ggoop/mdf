package rules

import (
	"fmt"
	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/md"
	"github.com/ggoop/mdf/mof"
	"github.com/ggoop/mdf/repositories"
	"github.com/ggoop/mdf/utils"
)

type CommonDelete struct {
	repo *repositories.MysqlRepo
}

func NewCommonDelete(repo *repositories.MysqlRepo) *CommonDelete {
	return &CommonDelete{repo}
}
func (s CommonDelete) Exec(req *mof.ReqContext, res *mof.ResContext) error {
	if req.ID == "" {
		return glog.Error("缺少 ID 参数！")
	}
	//查找实体信息
	entity := md.GetEntity(req.MainEntity)
	if entity == nil {
		return glog.Error("找不到实体！")
	}
	if field := entity.GetField("System"); field != nil && field.DbName != "" {
		count := 0
		s.repo.Table(entity.TableName).Where(fmt.Sprintf("id=? and %s = 1", s.repo.Dialect().Quote("system")), req.ID).Count(&count)
		if count > 0 {
			return glog.Error("系统预制数据不可删除！")
		}
	}
	if df := entity.GetField("DeletedAt"); df != nil {
		if err := s.repo.Exec(fmt.Sprintf("update %s set %s=? where id=?", entity.TableName, df.DbName), utils.NewTime(), req.ID).Error; err != nil {
			return err
		}
	} else {
		if err := s.repo.Exec(fmt.Sprintf("delete from %s where id=?", entity.TableName), req.ID).Error; err != nil {
			return err
		}
	}
	return nil
}
func (s CommonDelete) GetRule() mof.RuleRegister {
	return mof.RuleRegister{Code: "delete", Domain: "common"}
}