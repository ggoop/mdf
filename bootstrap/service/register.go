package service

import (
	"github.com/kataras/iris/hero"

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
	afterRegister()
}
func afterRegister() {
	if err := di.Global.Invoke(func(s1 *ProductSv, s2 *EntSv, s3 *TokenSv, s4 *UserSv, s5 *ProfileSv, s6 *CronSv) {
		hero.Register(s1, s2, s3, s4, s5, s6)
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	if err := di.Global.Invoke(func(s1 *md.EntitySv, s2 *OssSv, s3 *LogSv, s4 *MdSv) {
		hero.Register(s1, s2, s3, s4)
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	if err := di.Global.Invoke(func(s1 *RoleSv) {
		hero.Register(s1)
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	if err := di.Global.Invoke(func(s1 *QuerySv) {
		hero.Register(s1)
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	if err := di.Global.Invoke(func(s2 *DtiSv) {
		hero.Register(s2)
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
}
