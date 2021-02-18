package actions

import (
	"github.com/ggoop/mdf/framework/md"
)

func Register() {
	//注册到mof框架
	md.ActionSv().RegisterAction(
		newCommonDelete(),
		newCommonDisable(),
		newCommonEnable(),
		newCommonImport(),
		newCommonQuery(),
		newCommonSave(),
	)
}
