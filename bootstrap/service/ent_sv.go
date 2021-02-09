package service

import (
	"github.com/ggoop/mdf/bootstrap/errors"
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/utils"
)

type EntSv struct {
	repo *repositories.MysqlRepo
}

/**
* 创建服务实例
 */
func NewEntSv(repo *repositories.MysqlRepo) *EntSv {
	return &EntSv{repo: repo}
}

func (s *EntSv) GetEntBy(id string) (*model.Ent, error) {
	item := model.Ent{}
	if err := s.repo.Model(&model.Ent{}).Where("id = ?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *EntSv) GetUserDefaultEnt(userID string, entID string) *model.Ent {
	item := model.Ent{}
	if entID != "" {
		s.repo.Model(&model.Ent{}).Where("id in (?)", s.repo.Model(&model.EntUser{}).Select("ent_id").Where("user_id = ?", userID).QueryExpr()).Where("id=?", entID).Take(&item)
	}
	if item.ID != "" {
		return &item
	}
	var ents []string
	s.repo.Model(&model.EntUser{}).Select("ent_id").Where("user_id = ?", userID).Order("is_default desc,created_at desc,id desc").Pluck("ent_id", &ents)

	if len(ents) > 0 {
		s.repo.Model(&model.Ent{}).Where("id in (?)", ents[0]).Take(&item)
	}

	if item.ID != "" {
		return &item
	}
	return nil
}
func (s *EntSv) GetByOpenid(id string) (*model.Ent, error) {
	item := model.Ent{}
	if err := s.repo.Model(&model.Ent{}).Where("openid = ?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *EntSv) GetEntUserBy(entID, userID string) (*model.EntUser, error) {
	item := model.EntUser{}
	if err := s.repo.Model(item).Where("ent_id = ? and user_id=?", entID, userID).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *EntSv) GetEntByUser(userID string) ([]model.Ent, error) {
	items := make([]model.Ent, 0)
	if err := s.repo.Model(&model.Ent{}).Where("id in( ?)", s.repo.Model(model.EntUser{}).Select("ent_id").Where("user_id=?", userID).SubQuery()).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

/**
创建创业
*/
func (s *EntSv) IssueEnt(ent *model.Ent) (*model.Ent, error) {
	if ent.ID == "" {
		ent.ID = utils.GUID()
	}
	if ent.Openid == "" {
		ent.Openid = utils.GUID()
	}
	oldItem := model.Ent{}
	s.repo.Model(&model.Ent{}).Where("openid=?", ent.Openid).Take(&oldItem)
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
			s.repo.Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
			s.repo.Where("id=?", oldItem.ID).Take(&oldItem)
		}
		ent.ID = oldItem.ID
	} else {
		if err := s.repo.Create(ent).Error; err != nil {
			return nil, err
		}
	}
	return s.GetEntBy(ent.ID)
}

func (s *EntSv) UpdateEnt(id string, ent model.Ent) (*model.Ent, error) {
	oldItem := model.Ent{}
	s.repo.Model(&model.Ent{}).Where("id = ?", id).Take(&oldItem)
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
			s.repo.Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
			s.repo.Where("id=?", oldItem.ID).Take(&oldItem)
		}
	}
	return s.GetEntBy(id)
}
func (s *EntSv) DestroyEnt(id string) error {
	if err := s.repo.Where("id = ?", id).Delete(model.Ent{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *EntSv) AddMember(item *model.EntUser) (*model.EntUser, error) {
	if item.EntID == "" || item.UserID == "" {
		return nil, errors.ParamsRequired("entid or userid")
	}
	oldItem := model.EntUser{}
	s.repo.Model(oldItem).Take(&oldItem, "ent_id=? and user_id=?", item.EntID, item.UserID)
	if oldItem.ID != "" {
		updates := make(map[string]interface{})
		if oldItem.TypeID != item.TypeID && item.TypeID != "" && oldItem.TypeID == "" {
			updates["TypeID"] = item.TypeID
		}
		if len(updates) > 0 {
			s.repo.Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
			s.repo.Where("id=?", oldItem.ID).Take(oldItem)
		}
		item = &oldItem
	} else {
		item.ID = utils.GUID()
		if err := s.repo.Create(item).Error; err != nil {
			return nil, err
		}
		s.repo.Where("id=?", item.ID).Take(item)
	}
	return item, nil
}
func (s *EntSv) RemoveMember(item *model.EntUser) error {
	if item.EntID == "" || item.UserID == "" {
		return errors.ParamsRequired("entid or userid")
	}
	s.repo.Delete(item, "ent_id=? and user_id=?", item.EntID, item.UserID)
	return nil
}
