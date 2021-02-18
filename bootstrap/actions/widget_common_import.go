package actions

import (
	"github.com/ggoop/mdf/framework/files"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type commonImport struct {
}

func newCommonImport() *commonImport {
	return &commonImport{}
}
func (s *commonImport) Register() md.RuleRegister {
	return md.RuleRegister{Code: "import", OwnerType: md.RuleType_Widget, OwnerID: "common"}
}
func (s *commonImport) Exec(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext) {
	if len(req.Files) > 0 {
		datas := make([]files.ImportData, 0)
		for _, file := range req.Files {
			if f, err := file.Open(); err != nil {
				res.SetError(err)
				return
			} else {
				if ds, err := files.NewExcelSv().GetExcelDatasByReader(f); err != nil {
					res.SetError(err)
					return
				} else if len(ds) > 0 {
					datas = append(datas, ds...)
				}
			}
		}
		req.Data = datas
	}
}
