package rules

import (
	"github.com/ggoop/mdf/framework/md"
)

func Register() {
	//注册到mof框架
	md.ActionSv().RegisterRule(
		newCommonQuery(),
		newCommonSave(),
		newCommonDelete(),

		newCommonEnable(),
		newCommonDisable(),

		newCommonImport(),

		newMdImportPre(),
	)
}
