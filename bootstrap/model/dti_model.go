package model

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

//本地接口
type DtiLocal struct {
	md.Model
	EntID       string          `gorm:"size:50"`
	ProductCode string          `gorm:"size:50" json:"product_code"`
	Code        string          `gorm:"size:50" json:"code"`
	Name        string          `json:"name"`
	Memo        string          `json:"memo"`
	Host        string          `json:"host"` //注册中心名称，或者主机地址（带有http就是主机地址）
	Path        string          `gorm:"size:50" json:"path"`
	Enabled     utils.SBool     `gorm:"default:true" json:"enabled"`
	Tags        string          `gorm:"size:100;name:标签" json:"tags"`
	System      utils.SBool     `gorm:"not null;default:0;name:系统的" json:"system"`
	Params      []DtiLocalParam `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:LocalID"`
}

func (s DtiLocal) TableName() string {
	return "sys_dti_locals"
}
func (s *DtiLocal) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".dti.local", Domain: MD_DOMAIN, Name: "本地接口"}
}

/**
接口参数
*/
type DtiLocalParam struct {
	md.Model
	LocalID  string      `gorm:"size:50" json:"local_id"`
	Local    *DtiLocal   `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Code     string      `json:"code"`
	Name     string      `json:"name"`
	Memo     string      `json:"memo"`
	TypeID   string      `json:"type_id"`
	Type     *md.MDEnum  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:sys.dti.data.type"`
	Value    string      `json:"value"`
	Required utils.SBool `gorm:"default:0" json:"required"`
	ValueDef string      `gorm:"name:值定义" json:"value_def"` // 枚举类型/实体全名
	Hidden   utils.SBool `gorm:"default:0" json:"hidden"`
	Sequence int         `gorm:"default:0" json:"sequence"`
}

func (s DtiLocalParam) TableName() string {
	return "sys_dti_local_params"
}
func (s *DtiLocalParam) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".dti.local.param", Domain: MD_DOMAIN, Name: "接口参数"}
}

//远程节点
type DtiNode struct {
	md.Model
	EntID      string      `gorm:"size:50"`
	TemplateID string      `gorm:"size:50" json:"template_id"`
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	Memo       string      `json:"memo"`
	Host       string      `json:"host"`
	Public     utils.SBool `gorm:"default:0" json:"public"`
	Enabled    utils.SBool `gorm:"default:true" json:"enabled"`
	System     utils.SBool `gorm:"not null;default:0;name:系统的"`
}

func (s DtiNode) TableName() string {
	return "sys_dti_nodes"
}
func (s *DtiNode) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".dti.node", Domain: MD_DOMAIN, Name: "远程节点"}
}

//参数
type DtiParam struct {
	md.Model
	EntID  string     `gorm:"size:50"`
	NodeID string     `gorm:"size:50" json:"node_id"`
	Node   *DtiNode   `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Code   string     `json:"code"`
	Name   string     `json:"name"`
	Memo   string     `json:"memo"`
	TypeID string     `json:"type_id"`
	Type   *md.MDEnum `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:sys.dti.param.type"`
	Value  string     `json:"value"`
}

func (s DtiParam) TableName() string {
	return "sys_dti_params"
}
func (s *DtiParam) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".dti.param", Domain: MD_DOMAIN, Name: "接口参数"}
}

//
type DtiRemote struct {
	md.Model
	EntID    string      `gorm:"size:50"`
	Code     string      `json:"code"`
	Name     string      `json:"name"`
	NodeID   string      `gorm:"size:50" json:"node_id"`
	Node     *DtiNode    `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	LocalID  string      `gorm:"size:50" json:"local_id"`
	Local    *DtiLocal   `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	MethodID string      `gorm:"size:50" json:"method_id"` //请求类型
	Method   *md.MDEnum  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:sys.dti.method"`
	Path     string      `gorm:"size:50" json:"path"`
	Header   string      `gorm:"size:50" json:"header"`
	Body     string      `json:"body"`
	Query    string      `json:"query"`
	Memo     string      `json:"memo"`
	Enabled  utils.SBool `gorm:"default:true" json:"enabled"`
	Sequence int         `json:"sequence"`
	FmDate   utils.Time  `json:"fm_date"`
	ToDate   utils.Time  `json:"to_date"`
	StatusID string      `json:"status_id"`
	Status   *md.MDEnum  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:sys.dti.status"`
	Msg      string      `gorm:"type:text" json:"msg"`
}

func (s DtiRemote) TableName() string {
	return "sys_dti_remotes"
}
func (s *DtiRemote) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".dti.remote", Domain: MD_DOMAIN, Name: "接口"}
}
