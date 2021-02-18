package rules

import (
	"fmt"
	"github.com/ggoop/mdf/db"
	"github.com/ggoop/mdf/utils"

	"github.com/ggoop/mdf/framework/md"
)

type commonDisable struct {
}

func newCommonDisable() *commonDisable {
	return &commonDisable{}
}
func (s commonDisable) Exec(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext) {
	if req.ID == "" {
		res.SetError("缺少 ID 参数！")
		return
	}
	//查找实体信息
	entity := md.MDSv().GetEntity(req.Entity)
	if entity == nil {
		res.SetError("找不到实体！")
		return
	}
	if err := db.Default().Exec(fmt.Sprintf("update %s set enabled=0 where id =?", entity.TableName), req.ID).Error; err != nil {
		res.SetError(err)
		return
	}
}
func (s commonDisable) Register() md.RuleRegister {
	return md.RuleRegister{Code: "disable", OwnerType: md.RuleType_Widget, OwnerID: "common"}
}
