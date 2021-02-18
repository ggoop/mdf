package rules

import (
	"fmt"
	"github.com/ggoop/mdf/db"

	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type CommonDelete struct {
}

func newCommonDelete() *CommonDelete {
	return &CommonDelete{}
}
func (s CommonDelete) Exec(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext) {
	if req.ID == "" {
		res.SetError("缺少 ID 参数！")
		return
	}
	//查找实体信息
	entity := md.MDSv().GetEntity(req.Entity)
	if entity == nil {
		res.SetError("找不到实体")
		return
	}
	if field := entity.GetField("System"); field != nil && field.DbName != "" {
		count := 0
		db.Default().Table(entity.TableName).Where(fmt.Sprintf("id=? and %s = 1", db.Default().Dialect().Quote("system")), req.ID).Count(&count)
		if count > 0 {
			res.SetError("系统预制数据不可删除")
			return
		}
	}
	if df := entity.GetField("DeletedAt"); df != nil {
		if err := db.Default().Exec(fmt.Sprintf("update %s set %s=? where id=?", entity.TableName, df.DbName), utils.TimeNow(), req.ID).Error; err != nil {
			res.SetError(err)
			return
		}
	} else {
		if err := db.Default().Exec(fmt.Sprintf("delete from %s where id=?", entity.TableName), req.ID).Error; err != nil {
			res.SetError(err)
			return
		}
	}
}
func (s CommonDelete) Register() md.RuleRegister {
	return md.RuleRegister{Code: "delete", OwnerType: md.RuleType_Widget, OwnerID: "common"}
}
