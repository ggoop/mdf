package rules

import (
	"github.com/ggoop/mdf/di"
	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/mof"
	"github.com/ggoop/mdf/repositories"
)

func Register() {

	//注册到mof框架
	if err := di.Global.Invoke(func(repo *repositories.MysqlRepo) {
		mof.RegisterActionRule(
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
