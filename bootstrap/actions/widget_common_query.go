package actions

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type commonQuery struct {
}

func newCommonQuery() *commonQuery {
	return &commonQuery{}
}
func (s *commonQuery) Register() md.RuleRegister {
	return md.RuleRegister{Code: "query", OwnerType: md.RuleType_Widget, OwnerID: "common"}
}

func (s commonQuery) Exec(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext) {

}
