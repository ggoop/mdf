package actions

import (
	"github.com/ggoop/mdf/utils"

	"github.com/ggoop/mdf/framework/md"
)

type commonEnable struct {
}

func newCommonEnable() *commonEnable {
	return &commonEnable{}
}

func (s commonEnable) Register() md.RuleRegister {
	return md.RuleRegister{Code: "enable", OwnerType: md.RuleType_Widget, OwnerID: "common"}
}

func (s commonEnable) Exec(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext) {

}
