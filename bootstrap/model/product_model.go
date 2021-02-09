package model

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

/**
AppProduct
*/
type Product struct {
	md.Model
	EntID  string      `gorm:"size:50" json:"ent_id"`
	Code   string      `gorm:"size:50" json:"code"`
	Name   string      `gorm:"name:菜单名称" json:"name"`
	Icon   string      `gorm:"name:图标" json:"icon"`
	System utils.SBool `gorm:"not null;default:0;name:系统的" json:"system"`
}

func (t Product) TableName() string {
	return "sys_app_products"
}
func (s *Product) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".app.product", Domain: MD_DOMAIN, Name: "应用产品"}
}

/**
站点,主机，服务地址
*/
type ProductHost struct {
	md.Model
	EntID    string      `gorm:"size:50" json:"ent_id"`
	Code     string      `gorm:"size:50" json:"code"`
	Name     string      `gorm:"name:名称" json:"name"`
	DevHost  string      `gorm:"name:开发环境" json:"dev_host"`
	TestHost string      `gorm:"name:测试环境" json:"test_host"`
	PreHost  string      `gorm:"name:预发环境" json:"pre_host"`
	ProdHost string      `gorm:"name:生产环境" json:"prod_host"`
	System   utils.SBool `gorm:"not null;default:0;name:系统的" json:"system"`
}

func (t ProductHost) TableName() string {
	return "sys_app_hosts"
}
func (s *ProductHost) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".app.host", Domain: MD_DOMAIN, Name: "站点"}
}

/**
Modules
*/
type ProductModule struct {
	md.Model
	EntID     string      `gorm:"size:50" json:"ent_id"`
	Code      string      `gorm:"size:50" json:"code"`
	Name      string      `gorm:"name:菜单名称" json:"name"`
	ProductID string      `gorm:"size:50" json:"product_id" json:"product_id"`
	Product   *Product    `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false" json:"product"`
	System    utils.SBool `gorm:"not null;default:0;name:系统的" json:"system"`
}

func (t ProductModule) TableName() string {
	return "sys_app_modules"
}
func (s *ProductModule) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".app.module", Domain: MD_DOMAIN, Name: "应用模块"}
}

/**
AppService
*/
type ProductService struct {
	md.Model
	EntID     string         `gorm:"size:50" json:"ent_id"`
	ProductID string         `gorm:"size:50" json:"product_id"`
	Product   *Product       `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false" json:"product"`
	HostID    string         `gorm:"size:50" json:"host_id"`
	Host      *ProductHost   `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false" json:"host"`
	ModuleID  string         `gorm:"size:50" json:"module_id"`
	Module    *ProductModule `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false" json:"module"`
	Code      string         `gorm:"size:50" json:"code"`
	Name      string         `gorm:"name:名称" json:"name"`
	Params    string         `gorm:"name:json格式参数" json:"params"`
	Uri       string         `gorm:"name:导航用的uri" json:"uri"`
	AppUri    string         `gorm:"name:导航用的uri" json:"app_uri"`
	Icon      string         `gorm:"name:图标" json:"icon"`
	Sequence  int            `json:"sequence" json:"sequence"`

	InWeb  utils.SBool `gorm:"not null;default:0;name:是否WEB" json:"in_web"`
	InApp  utils.SBool `gorm:"not null;default:0;name:是否APP" json:"in_app"`
	Schema string      `gorm:"size:50" json:"schema"`

	Memo      string      `gorm:"name:备注" json:"memo"`
	Tags      string      `gorm:"name:备注" json:"tags"`
	BizType   string      `gorm:"size:10;name:服务类型" json:"biz_type"` //m管理/b业务/a全员
	IsMaster  utils.SBool `gorm:"not null;default:0;name:主节点" json:"is_master"`
	IsSlave   utils.SBool `gorm:"not null;default:0;name:主节点" json:"is_slave"`
	IsDefault utils.SBool `gorm:"not null;default:0;name:是否默认" json:"is_default"`
	System    utils.SBool `gorm:"not null;default:0;name:系统的" json:"system"`
}

func (t ProductService) TableName() string {
	return "sys_app_services"
}
func (s *ProductService) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".app.service", Domain: MD_DOMAIN, Name: "应用服务"}
}
