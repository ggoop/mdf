package services

import (
	"github.com/ggoop/mdf/bootstrap/errors"
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/db"
	"github.com/ggoop/mdf/utils"
)

type IEntSv interface {
	GetEntBy(id string) (*model.Ent, error)
}

var entSvInstance IEntSv = newEntSv()

type entSvImpl struct {
}

func EntSv() IEntSv {
	return entSvInstance
}

/**
* 创建服务实例
 */
func newEntSv() *entSvImpl {
	return &entSvImpl{}
}

func (s *entSvImpl) GetEntBy(id string) (*model.Ent, error) {
	item := model.Ent{}
	if err := db.Default().Model(&model.Ent{}).Where("id = ?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *entSvImpl) GetByOpenid(id string) (*model.Ent, error) {
	item := model.Ent{}
	if err := db.Default().Model(&model.Ent{}).Where("openid = ?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *entSvImpl) GetEntUserBy(entID, userID string) (*model.EntUser, error) {
	item := model.EntUser{}
	if err := db.Default().Model(item).Where("ent_id = ? and user_id=?", entID, userID).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *entSvImpl) GetEntByUser(userID string) ([]model.Ent, error) {
	items := make([]model.Ent, 0)
	if err := db.Default().Model(&model.Ent{}).Where("id in( ?)", db.Default().Model(model.EntUser{}).Select("ent_id").Where("user_id=?", userID).SubQuery()).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

/**
创建创业
*/
func (s *entSvImpl) IssueEnt(ent *model.Ent) (*model.Ent, error) {
	if ent.ID == "" {
		ent.ID = utils.GUID()
	}
	if ent.Openid == "" {
		ent.Openid = utils.GUID()
	}
	oldItem := model.Ent{}
	db.Default().Model(&model.Ent{}).Where("openid=?", ent.Openid).Take(&oldItem)
	if oldItem.ID != "" {
		updates := make(map[string]interface{})
		if oldItem.Memo != ent.Memo && ent.Memo != "" {
			updates["Memo"] = ent.Memo
		}
		if oldItem.Name != ent.Name && ent.Name != "" {
			updates["Name"] = ent.Name
		}
		if oldItem.StatusID != ent.StatusID && ent.StatusID != "" {
			updates["StatusID"] = ent.StatusID
		}
		if oldItem.TypeID != ent.TypeID && ent.TypeID != "" {
			updates["TypeID"] = ent.TypeID
		}
		if len(updates) > 0 {
			db.Default().Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
			db.Default().Where("id=?", oldItem.ID).Take(&oldItem)
		}
		ent.ID = oldItem.ID
	} else {
		if err := db.Default().Create(ent).Error; err != nil {
			return nil, err
		}
	}
	return s.GetEntBy(ent.ID)
}

func (s *entSvImpl) UpdateEnt(id string, ent model.Ent) (*model.Ent, error) {
	oldItem := model.Ent{}
	db.Default().Model(&model.Ent{}).Where("id = ?", id).Take(&oldItem)
	if oldItem.ID != "" {
		updates := make(map[string]interface{})
		if oldItem.Memo != ent.Memo && ent.Memo != "" {
			updates["Memo"] = ent.Memo
		}
		if oldItem.Name != ent.Name && ent.Name != "" {
			updates["Name"] = ent.Name
		}
		if oldItem.Gateway != ent.Gateway && ent.Gateway != "" {
			updates["Gateway"] = ent.Gateway
		}
		if oldItem.StatusID != ent.StatusID && ent.StatusID != "" {
			updates["StatusID"] = ent.StatusID
		}
		if oldItem.TypeID != ent.TypeID && ent.TypeID != "" {
			updates["TypeID"] = ent.TypeID
		}
		if len(updates) > 0 {
			db.Default().Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
			db.Default().Where("id=?", oldItem.ID).Take(&oldItem)
		}
	}
	return s.GetEntBy(id)
}
func (s *entSvImpl) DestroyEnt(id string) error {
	if err := db.Default().Where("id = ?", id).Delete(model.Ent{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *entSvImpl) AddMember(item *model.EntUser) (*model.EntUser, error) {
	if item.EntID == "" || item.UserID == "" {
		return nil, errors.ParamsRequired("entid or userid")
	}
	oldItem := model.EntUser{}
	db.Default().Model(oldItem).Take(&oldItem, "ent_id=? and user_id=?", item.EntID, item.UserID)
	if oldItem.ID != "" {
		updates := make(map[string]interface{})
		if oldItem.TypeID != item.TypeID && item.TypeID != "" && oldItem.TypeID == "" {
			updates["TypeID"] = item.TypeID
		}
		if len(updates) > 0 {
			db.Default().Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
			db.Default().Where("id=?", oldItem.ID).Take(oldItem)
		}
		item = &oldItem
	} else {
		item.ID = utils.GUID()
		if err := db.Default().Create(item).Error; err != nil {
			return nil, err
		}
		db.Default().Where("id=?", item.ID).Take(item)
	}
	return item, nil
}
func (s *entSvImpl) RemoveMember(item *model.EntUser) error {
	if item.EntID == "" || item.UserID == "" {
		return errors.ParamsRequired("entid or userid")
	}
	db.Default().Delete(item, "ent_id=? and user_id=?", item.EntID, item.UserID)
	return nil
}
