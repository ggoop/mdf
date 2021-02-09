package model

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type Role struct {
	md.Model
	EntID   string      `gorm:"size:50"`
	Code    string      `gorm:"not null" json:"code"`
	Name    string      `gorm:"not null" json:"name"`
	Memo    string      `json:"memo"`
	Scope   string      `json:"scope"` //作用范围,如:用户:user,应用:app
	Enabled utils.SBool `gorm:"not null;default:1;name:启用"`
	System  utils.SBool `gorm:"not null;default:0;name:系统的"`
}

func (t Role) TableName() string {
	return "sys_roles"
}
func (s *Role) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".role", Domain: MD_DOMAIN, Name: "角色"}
}

type RoleUser struct {
	md.Model
	EntID  string `gorm:"size:50"`
	UserID string `gorm:"size:50" json:"user_id"`
	RoleID string `gorm:"size:50" json:"role_id"`
	Role   *Role  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
}

func (t RoleUser) TableName() string {
	return "sys_role_users"
}
func (s *RoleUser) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".role.user", Domain: MD_DOMAIN, Name: "用户角色"}
}

type RoleService struct {
	md.Model
	EntID     string          `gorm:"size:50"`
	ServiceID string          `gorm:"size:50" json:"service_id"`
	Service   *ProductService `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	RoleID    string          `gorm:"size:50" json:"role_id"`
	Role      *Role           `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Readable  utils.SBool     `gorm:"name:可读取" json:"readable"`
	Writable  utils.SBool     `gorm:"name:可写入" json:"writable"`
	Deletable utils.SBool     `gorm:"name:可删除" json:"deletable"`
}

func (t RoleService) TableName() string {
	return "sys_role_services"
}
func (s *RoleService) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".role.service", Domain: MD_DOMAIN, Name: "角色对应的服务"}
}

type RoleData struct {
	md.Model
	EntID    string       `gorm:"size:50"`
	RoleID   string       `gorm:"size:50" json:"role_id"`
	Role     *Role        `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	EntityID string       `gorm:"size:50" json:"entity_id"` //类型ID
	Entity   *md.MDEntity `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Scope    string       `gorm:"size:50"` //query
	Exp      string       `json:"exp"`
	Memo     string       `json:"memo"`
	IsRead   utils.SBool  `gorm:"name:读取权限" json:"is_read"`
	IsWrite  utils.SBool  `gorm:"name:写入权限" json:"is_write"`
	IsDelete utils.SBool  `gorm:"name:删除权限" json:"is_delete"`
}

func (t RoleData) TableName() string {
	return "sys_role_datas"
}
func (s *RoleData) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".role.data", Domain: MD_DOMAIN, Name: "数据权限"}
}
