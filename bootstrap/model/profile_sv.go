package model

import (
	"github.com/ggoop/mdf/bootstrap/errors"
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/utils"
)

type ProfileSv struct {
	repo *repositories.MysqlRepo
}

/**
* 创建服务实例
 */
func NewProfileSv(repo *repositories.MysqlRepo) *ProfileSv {
	return &ProfileSv{repo: repo}
}

func (s *ProfileSv) SaveProfiles(item *Profile) (*Profile, error) {
	if item.EntID == "" || item.Code == "" {
		return nil, errors.ParamsRequired("entid or code")
	}
	oldItem := Profile{}
	if item.ID != "" {
		s.repo.Model(item).Where("id=? and ent_id=?", item.ID, item.EntID).Take(&oldItem)
	}
	if oldItem.ID == "" && item.Code != "" {
		s.repo.Model(item).Where("code=? and ent_id=?", item.Code, item.EntID).Take(&oldItem)
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
			s.repo.Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
		}
		item.ID = oldItem.ID

	} else {
		item.ID = utils.GUID()
		if item.Name == "" {
			item.Name = item.Code
		}
		item.CreatedAt = utils.NewTime()
		if err := s.repo.Create(item).Error; err != nil {
			return nil, err
		}
	}
	return item, nil
}
func (s *ProfileSv) DeleteProfiles(entID string, ids []string) error {
	if err := s.repo.Delete(Profile{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}
func (s *ProfileSv) GetValue(entID, code string) (string, error) {
	item := Profile{}
	if err := s.repo.Model(item).First(&item, "code=? and ent_id=?", code, entID).Error; err != nil {
		return "", err
	}
	return item.Value, nil
}
func (s *ProfileSv) SetValue(entID, code, value string) error {
	item := Profile{}
	s.repo.Model(item).First(&item, "code=? and ent_id=?", code, entID)
	if item.ID == "" {
		item.Code = code
		item.ID = utils.GUID()
		item.DefaultValue = value
		item.Value = value
		if err := s.repo.Create(&item).Error; err != nil {
			return err
		}
	} else {
		updates := utils.Map{}
		updates["Value"] = value
		s.repo.Model(item).Where("id=?", item.ID).Update(updates)
	}
	return nil
}
