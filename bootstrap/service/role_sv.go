package service

import (
	"fmt"
	"sync"

	"github.com/ggoop/mdf/bootstrap/errors"
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"

	"github.com/ggoop/mdf/framework/db/repositories"
)

type RoleSv struct {
	*sync.Mutex
	repo *repositories.MysqlRepo
}

/**
* 创建服务实例
 */
func NewRoleSv(repo *repositories.MysqlRepo) *RoleSv {
	return &RoleSv{repo: repo, Mutex: &sync.Mutex{}}
}
func (s *RoleSv) GetRoleBy(entID, idOrCode string) (*model.Role, error) {
	item := model.Role{}
	if err := s.repo.Model(&model.Role{}).Where("(id=? or code=? ) and ent_id=?", idOrCode, idOrCode, entID).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *RoleSv) RoleIsExists(entID, code string) bool {
	count := 0
	s.repo.Model(&model.Role{}).Where("code=? and ent_id=?", code, entID).Count(&count)
	return count > 0
}
func (s *RoleSv) SaveRole(entID string, item model.Role) (*model.Role, error) {
	if item.Code == "" {
		return nil, errors.ParamsRequired("编码")
	}
	if item.ID != "" {
		old := model.Role{}
		if err := s.repo.Model(&model.Role{}).Where("id=? and ent_id=?", item.ID, entID).Take(&old).Error; err != nil {
			return nil, err
		}
		//系统，只能修改备注，和状态
		if old.System.IsTrue() {
			s.repo.Model(&old).Updates(map[string]interface{}{"Memo": item.Memo, "Enabled": item.Enabled})
		} else {
			//如果修改编码
			if old.Code != item.Code {
				if s.RoleIsExists(entID, item.Code) {
					return nil, errors.ExistError("名称", item.Code)
				}
			}
			s.repo.Model(&old).Updates(map[string]interface{}{"Code": item.Code, "Name": item.Name, "Memo": item.Memo, "Enabled": item.Enabled})
		}
	} else {
		if s.RoleIsExists(entID, item.Code) {
			return nil, errors.ExistError("名称", item.Code)
		}
		item.ID = utils.GUID()
		item.EntID = entID
		s.repo.Create(&item)
	}
	return s.GetRoleBy(entID, item.ID)
}

//删除
func (s *RoleSv) DeleteRole(entID string, ids []string) error {
	systemDatas := 0
	s.repo.Model(model.Role{}).Where("id in (?) and `system`=1", ids).Count(&systemDatas)
	if systemDatas > 0 {
		return utils.NewError("系统预制不能删除!")
	}
	if err := s.repo.Delete(model.Role{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}

//应用&菜单授权
func (s *RoleSv) SaveService(entID string, item model.RoleService) error {
	old := model.RoleService{}
	s.repo.Model(old).Where("ent_id=? and role_id=? and service_id=?", entID, item.RoleID, item.ServiceID).Take(&old)
	if old.ID != "" {
		updates := make(map[string]interface{})
		if item.Readable != old.Readable {
			updates["Readable"] = item.Readable
		}
		if item.Deletable != old.Deletable {
			updates["Deletable"] = item.Deletable
		}
		if item.Writable != old.Writable {
			updates["Writable"] = item.Writable
		}
		if len(updates) > 0 {
			s.repo.Model(old).Where("id=?", old.ID).Updates(updates)
		}
	} else {
		item.CreatedAt = utils.NewTime()
		item.ID = utils.GUID()
		item.EntID = entID
		s.repo.Create(&item)
	}
	return nil
}

//删除菜单 & 应用授权
func (s *RoleSv) DeleteService(entID string, item model.RoleService) error {
	if err := s.repo.Delete(model.RoleService{}, "ent_id=? and role_id=? and  service_id =?", entID, item.RoleID, item.ServiceID).Error; err != nil {
		return err
	}
	return nil
}
func (s *RoleSv) GetServiceIdsByUser(entID, userID string) []string {
	ids := make([]string, 0)
	query := s.repo.Table(fmt.Sprintf("%v  ru", s.repo.NewScope(model.RoleUser{}).TableName()))
	query = query.Joins(fmt.Sprintf("inner join %v  rs on ru.role_id=rs.role_id", s.repo.NewScope(model.RoleService{}).TableName()))
	query = query.Select([]string{"rs.service_id"})
	query = query.Where("rs.ent_id=? and ru.user_id=?", entID, userID)
	query.Pluck("rs.service_id", &ids)
	return ids
}

//通过用户，获取实体数据权限
func (s *RoleSv) GetEntityQueryPermitByUser(entID, userID string, entities ...string) []model.RoleData {
	return s.GetEntityPermitByUser(entID, userID, "", entities...)
}
func (s *RoleSv) GetEntityPermitByUser(entID, userID, scope string, entities ...string) []model.RoleData {
	query := s.repo.Table(fmt.Sprintf("%v  d", s.repo.NewScope(model.RoleData{}).TableName()))
	query = query.Joins(fmt.Sprintf("inner join %v  ru on ru.role_id=d.role_id", s.repo.NewScope(model.RoleUser{}).TableName()))
	query = query.Joins(fmt.Sprintf("inner join %v  e on e.id=d.entity_id", s.repo.NewScope(md.MDEntity{}).TableName()))
	query = query.Select([]string{"d.id", "d.created_at", "d.scope", "e.id as entity_id", "d.role_id", "d.exp"})
	if scope != "" {
		query = query.Where("d.scope like ?", "%"+scope+"%")
	}
	query = query.Where("d.ent_id=? and ru.ent_id=? and ru.user_id=?", entID, entID, userID)
	if entities != nil && len(entities) > 0 {
		query = query.Where("e.code in (?) or e.table_name in (?) or e.id in (?)", entities, entities, entities)
	}
	items := make([]model.RoleData, 0)
	query.Find(&items)
	return items
}
