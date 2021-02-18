package model

import (
	"github.com/ggoop/mdf/framework/md"
)

// 注册数据模型,提供数据层，按模块注册数据模型
func Register() {
	//sys
	md.MDSv().Migrate(&Log{}, &Client{}, &Token{}, &AuthFree{}, &CodeRule{}, &CodeValue{}, &Blacklist{})
	md.MDSv().Migrate(&Profile{})
	//product
	md.MDSv().Migrate(&Product{}, &ProductModule{}, &ProductHost{}, &ProductService{})
	//user
	md.MDSv().Migrate(&User{}, &UserFavorite{})
	//ent
	md.MDSv().Migrate(&Ent{}, &EntUser{})
	//role
	md.MDSv().Migrate(&Role{}, &RoleUser{}, &RoleService{}, &RoleData{})
	//cron
	md.MDSv().Migrate(&Cron{}, &CronLog{})
	//oss
	md.MDSv().Migrate(&Oss{}, &OssObject{})
	//dti
	md.MDSv().Migrate(&DtiLocal{}, &DtiLocalParam{}, &DtiNode{}, &DtiParam{}, &DtiRemote{})

}
