package rules

import (
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/md"
)

func Register() {
	//注册到mof框架
	md.RegisterActionRule(
		NewCommonLoad(repositories.Default()),
		NewCommonQuery(repositories.Default()),
		NewCommonSave(repositories.Default()),
		NewCommonImport(repositories.Default()),

		NewCommonDelete(repositories.Default()),
		NewCommonDeleteBatch(repositories.Default()),

		NewCommonEnable(repositories.Default()),
		NewCommonEnableBatch(repositories.Default()),

		NewCommonDisable(repositories.Default()),
		NewCommonDisableBatch(repositories.Default()),
	)
}
