package actions

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type commonDisable struct {
}

func newCommonDisable() *commonDisable {
	return &commonDisable{}
}

func (s commonDisable) Register() md.RuleRegister {
	return md.RuleRegister{Code: "disable", OwnerType: md.RuleType_Widget, OwnerID: "common"}
}
func (s commonDisable) Exec(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext) {

}
