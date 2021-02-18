package services

import (
	"github.com/ggoop/mdf/bootstrap/errors"
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/db"
	"github.com/ggoop/mdf/utils"
)

type IProfileSv interface {
}
type profileSvImpl struct {
}

var profileSv IProfileSv = newProfileSvImpl()

func ProfileSv() IProfileSv {
	return profileSv
}

/**
* 创建服务实例
 */
func newProfileSvImpl() *profileSvImpl {
	return &profileSvImpl{}
}

func (s *profileSvImpl) SaveProfiles(item *model.Profile) (*model.Profile, error) {
	if item.EntID == "" || item.Code == "" {
		return nil, errors.ParamsRequired("entid or code")
	}
	oldItem := model.Profile{}
	if item.ID != "" {
		db.Default().Model(item).Where("id=? and ent_id=?", item.ID, item.EntID).Take(&oldItem)
	}
	if oldItem.ID == "" && item.Code != "" {
		db.Default().Model(item).Where("code=? and ent_id=?", item.Code, item.EntID).Take(&oldItem)
	}
	if oldItem.ID != "" {
		updates := make(map[string]interface{})
		if oldItem.Name != item.Name && item.Name != "" {
			updates["Name"] = item.Name
		}
		if oldItem.Value != item.Value && item.Value != "" {
			updates["Value"] = item.Value
		}
		if oldItem.Memo != item.Memo && item.Memo != "" {
			updates["Memo"] = item.Memo
		}
		if oldItem.DefaultValue != item.DefaultValue && item.DefaultValue != "" {
			updates["DefaultValue"] = item.DefaultValue
		}
		if item.System.Valid() {
			updates["System"] = item.System
		}

		if len(updates) > 0 {
			db.Default().Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
		}
		item.ID = oldItem.ID

	} else {
		item.ID = utils.GUID()
		if item.Name == "" {
			item.Name = item.Code
		}
		item.CreatedAt = utils.TimeNow()
		if err := db.Default().Create(item).Error; err != nil {
			return nil, err
		}
	}
	return item, nil
}
func (s *profileSvImpl) DeleteProfiles(entID string, ids []string) error {
	if err := db.Default().Delete(model.Profile{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}
func (s *profileSvImpl) GetValue(entID, code string) (string, error) {
	item := model.Profile{}
	if err := db.Default().Model(item).First(&item, "code=? and ent_id=?", code, entID).Error; err != nil {
		return "", err
	}
	return item.Value, nil
}
func (s *profileSvImpl) SetValue(entID, code, value string) error {
	item := model.Profile{}
	db.Default().Model(item).First(&item, "code=? and ent_id=?", code, entID)
	if item.ID == "" {
		item.Code = code
		item.ID = utils.GUID()
		item.DefaultValue = value
		item.Value = value
		if err := db.Default().Create(&item).Error; err != nil {
			return err
		}
	} else {
		updates := utils.Map{}
		updates["Value"] = value
		db.Default().Model(item).Where("id=?", item.ID).Update(updates)
	}
	return nil
}
