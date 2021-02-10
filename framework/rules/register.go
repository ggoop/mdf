package rules

import (
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/di"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
)

func Register() {

	//注册到mof框架
	if err := di.Global.Invoke(func(repo *repositories.MysqlRepo) {
		md.RegisterActionRule(
			NewCommonLoad(repo),
			NewCommonQuery(repo),
			NewCommonSave(repo),
			NewCommonImport(repo),

			NewCommonDelete(repo),
			NewCommonDeleteBatch(repo),

			NewCommonEnable(repo),
			NewCommonEnableBatch(repo),

			NewCommonDisable(repo),
			NewCommonDisableBatch(repo),
		)
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
}
