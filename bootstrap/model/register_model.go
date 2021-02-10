package model

import (
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/di"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
)

// 注册数据模型,提供数据层，按模块注册数据模型
func RegisterMD() {
	if err := di.Global.Invoke(func(db *repositories.MysqlRepo) {
		//sys
		md.Migrate(db, &Log{}, &Client{}, &Token{}, &AuthFree{}, &CodeRule{}, &CodeValue{}, &Blacklist{})
		md.Migrate(db, &Profile{})
		//product
		md.Migrate(db, &Product{}, &ProductModule{}, &ProductHost{}, &ProductService{})
		//user
		md.Migrate(db, &User{}, &UserFavorite{})
		//ent
		md.Migrate(db, &Ent{}, &EntUser{})
		//role
		md.Migrate(db, &Role{}, &RoleUser{}, &RoleService{}, &RoleData{})
		//cron
		md.Migrate(db, &Cron{}, &CronLog{})
		//oss
		md.Migrate(db, &Oss{}, &OssObject{})
		//dti
		md.Migrate(db, &DtiLocal{}, &DtiLocalParam{}, &DtiNode{}, &DtiParam{}, &DtiRemote{})
		//query
		md.Migrate(db, &md.Query{}, &md.QueryFilter{}, &md.QueryColumn{}, &md.QueryOrder{}, &md.QueryWhere{}, &md.QueryCase{})

	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
}
