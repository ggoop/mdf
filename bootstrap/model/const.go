package model

const MD_DOMAIN = "sys"

const (
	OWNER_TYPE_ENUM_CREATOR = "creator"
	OWNER_TYPE_ENUM_MANAGER = "manager"
	OWNER_TYPE_ENUM_MEMBER  = "member"

	STATUS_ENUM_LOCKED   = "locked"   //锁定
	STATUS_ENUM_CREATED  = "created"  //已创建
	STATUS_ENUM_VERIFIED = "verified" //已验证
	STATUS_ENUM_REVOKED  = "revoked"  //已注销
)
const (
	CONFIG_NAME_LOCAL = "local"
	SYS_ENT_ID        = "0"
)
