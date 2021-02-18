package rules

import (
	"github.com/ggoop/mdf/db"
	"github.com/ggoop/mdf/framework/files"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type mdImportPre struct {
}

func newMdImportPre() *mdImportPre {
	return &mdImportPre{}
}
func (s *mdImportPre) Register() md.RuleRegister {
	return md.RuleRegister{Code: "md_import_pre", OwnerType: md.RuleType_Widget, OwnerID: "md"}
}
func (s *mdImportPre) Exec(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext) {
	if req.Data == nil {
		res.SetError("没有要导入的数据")
		return
	}
	if items, ok := req.Data.([]files.ImportData); !ok {
		res.SetError("导入的数据非法！")
		return
	} else {
		s.doProcess(token, req, res, items)
	}
}
func (s *mdImportPre) doProcess(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext, data []files.ImportData) {
	widgetCodes := make([]string, 0)
	for i, _ := range data {
		d := data[i]
		if d.EntityCode == "md.widget" {
			for _, r := range d.Data {
				if cv, co := r["Code"]; co && cv != "" {
					widgetCodes = append(widgetCodes, cv)
				}
			}
		}
	}
	if len(widgetCodes) == 0 {
		res.SetError("没有任何组件数据")
		return
	}
	//先按删除数据
	sql := "delete from md_widget_layouts where widget_id in (select id from md_widgets where code in (?))"
	if err := db.Default().Exec(sql, widgetCodes).Error; err != nil {
		glog.Error(err)
	}
}
