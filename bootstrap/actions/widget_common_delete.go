package actions

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type CommonDelete struct {
}

func newCommonDelete() *CommonDelete {
	return &CommonDelete{}
}
func (s CommonDelete) Register() md.RuleRegister {
	return md.RuleRegister{Code: "delete", OwnerType: md.RuleType_Widget, OwnerID: "common"}
}

func (s CommonDelete) Exec(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext) {

}
