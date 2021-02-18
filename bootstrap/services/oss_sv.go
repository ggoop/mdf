package services

import (
	"fmt"
	"github.com/ggoop/mdf/db"
	"strings"

	"github.com/ggoop/mdf/bootstrap/errors"
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/utils"
)

type IOssSv interface {
}

func OssSv() IOssSv {
	return ossSv
}

type ossSvImpl struct {
}

var ossSv IOssSv = newOssSvImpl()

/**
* 创建服务实例
 */
func newOssSvImpl() *ossSvImpl {
	return &ossSvImpl{}
}

func (s *ossSvImpl) GetObjectBy(id string) (*model.OssObject, error) {
	if id == "" {
		return nil, nil
	}
	item := model.OssObject{}
	if err := db.Default().Where("id=?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *ossSvImpl) GetObjectByCode(entID, code string) (*model.OssObject, error) {
	item := model.OssObject{}
	if err := db.Default().Where("ent_id=? and code=?", entID, code).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *ossSvImpl) ObjectNameIsExists(entID, directoryID, name string) bool {
	count := 0
	db.Default().Model(&model.OssObject{}).Where("name=? and directory_id=? and ent_id=?", name, directoryID, entID).Count(&count)
	return count > 0
}

//保存对象
func (s *ossSvImpl) SaveObject(item *model.OssObject) (*model.OssObject, error) {
	if item.ID == "" {
		item.ID = utils.GUID()
	}
	if item.Code == "" {
		item.Code = utils.GUID()
	}
	if item.Type == "" {
		item.Type = "obj"
	}
	if err := db.Default().Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// Directory
func (s *ossSvImpl) SaveDirectory(entID string, item model.OssObject) (*model.OssObject, error) {
	if item.Name == "" {
		return nil, errors.ParamsRequired("名称")
	}
	if item.ID != "" {
		old, err := s.GetObjectBy(item.ID)
		if err != nil {
			return nil, err
		}
		//如果修改名称，则需要校验名称是否已存在
		if old.Name != item.Name {
			if s.ObjectNameIsExists(entID, old.DirectoryID, item.Name) {
				return nil, errors.ExistError("名称", item.Name)
			}
		}
		db.Default().Model(&old).Updates(map[string]interface{}{"Name": item.Name, "Tag": item.Tag})
	} else {
		if s.ObjectNameIsExists(entID, item.DirectoryID, item.Name) {
			return nil, errors.ExistError("名称", item.Name)
		}
		item.Type = "dir"
		item.EntID = entID
		s.SaveObject(&item)
	}
	return s.GetObjectBy(item.ID)
}

func (s *ossSvImpl) GetConfig(id string) (*model.Oss, error) {
	old := model.Oss{}
	if err := db.Default().Take(&old, "id=?", id).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
func (s *ossSvImpl) GetEntConfig(entID string) (*model.Oss, error) {
	old := model.Oss{}
	if err := db.Default().Order("is_default desc").Take(&old, "ent_id=?", entID).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
func (s *ossSvImpl) SaveConfig(item model.Oss) (*model.Oss, error) {
	if ind := strings.Index(item.Endpoint, "//"); ind >= 0 {
		item.Endpoint = string(([]byte(item.Endpoint)[ind+2:]))
	}
	old, _ := s.GetConfig(item.ID)
	if old != nil && old.ID != "" {
		updates := make(map[string]interface{})
		count := 0
		if old.Type != item.Type && item.Type != "" {
			return nil, fmt.Errorf("不能修改存储空间类型")
		}
		if old.Bucket != item.Bucket && item.Bucket != "" {
			if db.Default().Model(model.OssObject{}).Where("oss_id in (?)", item.ID).Count(&count); count > 0 {
				return nil, fmt.Errorf("当前存储已被使用，或者存在文件，不能修改存储空间")
			}
			updates["Bucket"] = item.Bucket
		}
		if old.Endpoint != item.Endpoint && item.Endpoint != "" {
			if db.Default().Model(model.OssObject{}).Where("oss_id in (?)", item.ID).Count(&count); count > 0 {
				return nil, fmt.Errorf("当前存储已被使用，或者存在文件，不能修改存储空间地址")
			}
			updates["Endpoint"] = item.Endpoint
		}
		if old.AccessKeySecret != item.AccessKeySecret {
			updates["AccessKeySecret"] = item.AccessKeySecret
		}
		if old.AccessKeyID != item.AccessKeyID {
			updates["AccessKeyID"] = item.AccessKeyID
		}
		if old.Region != item.Region {
			updates["Region"] = item.Region
		}
		if len(updates) > 0 {
			db.Default().Model(model.Oss{}).Where("id=?", old.ID).Updates(updates)
		}
	} else {
		if item.ID == "" {
			item.ID = utils.GUID()
		}
		item.CreatedAt = utils.TimeNow()
		db.Default().Create(item)
	}
	return s.GetConfig(item.ID)
}
