package services

import (
	"fmt"
	"github.com/ggoop/mdf/db"
	"sync"

	"github.com/ggoop/mdf/bootstrap/errors"
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type IRoleSv interface {
}
type roleSvImpl struct {
	*sync.Mutex
}

var roleSv IRoleSv = newRoleSvImpl()

func RoleSv() IRoleSv {
	return roleSv
}

/**
* 创建服务实例
 */
func newRoleSvImpl() *roleSvImpl {
	return &roleSvImpl{Mutex: &sync.Mutex{}}
}
func (s *roleSvImpl) GetRoleBy(entID, idOrCode string) (*model.Role, error) {
	item := model.Role{}
	if err := db.Default().Model(&model.Role{}).Where("(id=? or code=? ) and ent_id=?", idOrCode, idOrCode, entID).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *roleSvImpl) RoleIsExists(entID, code string) bool {
	count := 0
	db.Default().Model(&model.Role{}).Where("code=? and ent_id=?", code, entID).Count(&count)
	return count > 0
}
func (s *roleSvImpl) SaveRole(entID string, item model.Role) (*model.Role, error) {
	if item.Code == "" {
		return nil, errors.ParamsRequired("编码")
	}
	if item.ID != "" {
		old := model.Role{}
		if err := db.Default().Model(&model.Role{}).Where("id=? and ent_id=?", item.ID, entID).Take(&old).Error; err != nil {
			return nil, err
		}
		//系统，只能修改备注，和状态
		if old.System.IsTrue() {
			db.Default().Model(&old).Updates(map[string]interface{}{"Memo": item.Memo, "Enabled": item.Enabled})
		} else {
			//如果修改编码
			if old.Code != item.Code {
				if s.RoleIsExists(entID, item.Code) {
					return nil, errors.ExistError("名称", item.Code)
				}
			}
			db.Default().Model(&old).Updates(map[string]interface{}{"Code": item.Code, "Name": item.Name, "Memo": item.Memo, "Enabled": item.Enabled})
		}
	} else {
		if s.RoleIsExists(entID, item.Code) {
			return nil, errors.ExistError("名称", item.Code)
		}
		item.ID = utils.GUID()
		item.EntID = entID
		db.Default().Create(&item)
	}
	return s.GetRoleBy(entID, item.ID)
}

//删除
func (s *roleSvImpl) DeleteRole(entID string, ids []string) error {
	systemDatas := 0
	db.Default().Model(model.Role{}).Where("id in (?) and `system`=1", ids).Count(&systemDatas)
	if systemDatas > 0 {
		return utils.ToError("系统预制不能删除!")
	}
	if err := db.Default().Delete(model.Role{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}

//应用&菜单授权
func (s *roleSvImpl) SaveService(entID string, item model.RoleService) error {
	old := model.RoleService{}
	db.Default().Model(old).Where("ent_id=? and role_id=? and service_id=?", entID, item.RoleID, item.ServiceID).Take(&old)
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
			db.Default().Model(old).Where("id=?", old.ID).Updates(updates)
		}
	} else {
		item.CreatedAt = utils.TimeNow()
		item.ID = utils.GUID()
		item.EntID = entID
		db.Default().Create(&item)
	}
	return nil
}

//删除菜单 & 应用授权
func (s *roleSvImpl) DeleteService(entID string, item model.RoleService) error {
	if err := db.Default().Delete(model.RoleService{}, "ent_id=? and role_id=? and  service_id =?", entID, item.RoleID, item.ServiceID).Error; err != nil {
		return err
	}
	return nil
}
func (s *roleSvImpl) GetServiceIdsByUser(entID, userID string) []string {
	ids := make([]string, 0)
	query := db.Default().Table(fmt.Sprintf("%v  ru", db.Default().NewScope(model.RoleUser{}).TableName()))
	query = query.Joins(fmt.Sprintf("inner join %v  rs on ru.role_id=rs.role_id", db.Default().NewScope(model.RoleService{}).TableName()))
	query = query.Select([]string{"rs.service_id"})
	query = query.Where("rs.ent_id=? and ru.user_id=?", entID, userID)
	query.Pluck("rs.service_id", &ids)
	return ids
}

//通过用户，获取实体数据权限
func (s *roleSvImpl) GetEntityQueryPermitByUser(entID, userID string, entities ...string) []model.RoleData {
	return s.GetEntityPermitByUser(entID, userID, "", entities...)
}
func (s *roleSvImpl) GetEntityPermitByUser(entID, userID, scope string, entities ...string) []model.RoleData {
	query := db.Default().Table(fmt.Sprintf("%v  d", db.Default().NewScope(model.RoleData{}).TableName()))
	query = query.Joins(fmt.Sprintf("inner join %v  ru on ru.role_id=d.role_id", db.Default().NewScope(model.RoleUser{}).TableName()))
	query = query.Joins(fmt.Sprintf("inner join %v  e on e.id=d.entity_id", db.Default().NewScope(md.MDEntity{}).TableName()))
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
