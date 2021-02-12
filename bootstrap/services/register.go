package services

import (
	"github.com/ggoop/mdf/framework/di"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
)

func Register() {
	items := []interface{}{
		//sys
		NewCacheSv,
		NewProductSv, NewEntSv, NewTokenSv, NewUserSv, NewProfileSv, NewCronSv,
		NewOssSv, NewLogSv, NewMdSv,
		NewRoleSv,
		md.NewEntitySv, NewQuerySv,
		NewDtiSv,
	}
	for _, v := range items {
		if err := di.Global.Provide(v); err != nil {
			glog.Errorf("di Provide error:%s", err)
		}
	}
}
