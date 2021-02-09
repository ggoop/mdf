package model

import (
	"sync"

	"github.com/ggoop/mdf/bootstrap/errors"
	"github.com/ggoop/mdf/utils"

	"github.com/ggoop/mdf/framework/db/repositories"
)

type ProductSv struct {
	*sync.Mutex
	repo *repositories.MysqlRepo
}

/**
* 创建服务实例
 */
func NewProductSv(repo *repositories.MysqlRepo) *ProductSv {
	return &ProductSv{repo: repo, Mutex: &sync.Mutex{}}
}

/*============product===========*/
func (s *ProductSv) SaveProduct(item *Product) (*Product, error) {
	if item.EntID == "" || item.Code == "" {
		return nil, errors.ParamsRequired("entid or code")
	}
	oldItem := Product{}
	if item.ID != "" {
		s.repo.Model(item).Where("id=?", item.ID).Take(&oldItem)
	}
	if oldItem.ID == "" && item.Code != "" {
		s.repo.Model(item).Where("code=? and ent_id=?", item.Code, item.EntID).Take(&oldItem)
	}
	if oldItem.ID != "" {
		updates := make(map[string]interface{})
		if oldItem.Name != item.Name && item.Name != "" {
			updates["Name"] = item.Name
		}
		if len(updates) > 0 {
			s.repo.Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
		}
		item.ID = oldItem.ID
	} else {
		if item.ID == "" {
			item.ID = utils.GUID()
		}
		item.CreatedAt = utils.NewTime()
		s.repo.Create(item)
	}
	return item, nil
}
func (s *ProductSv) GetProductByCode(entID, idOrCode string) (*Product, error) {
	old := Product{}
	if err := s.repo.Where("ent_id=? and (code=? or id=?)", entID, idOrCode, idOrCode).Take(&old).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
func (s *ProductSv) DeleteProducts(entID string, ids []string) error {
	systemRoles := 0
	s.repo.Model(Product{}).Where("ent_id=? and id in (?) and `system`=1", entID, ids).Count(&systemRoles)
	if systemRoles > 0 {
		return utils.NewError("系统预制产品不能删除!")
	}
	if err := s.repo.Delete(Product{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}

/*==================host============*/
func (s *ProductSv) GetHostByCode(entID, idOrCode string) (*ProductHost, error) {
	old := ProductHost{}
	if err := s.repo.Where("ent_id=? and (code=? or id=?)", entID, idOrCode, idOrCode).Take(&old).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
func (s *ProductSv) SaveHosts(item *ProductHost) (*ProductHost, error) {
	if item.EntID == "" || item.Code == "" {
		return nil, errors.ParamsRequired("entid or code")
	}
	oldItem := ProductHost{}
	if item.ID != "" {
		s.repo.Model(item).Where("id=?", item.ID).Take(&oldItem)
	}
	if oldItem.ID == "" && item.Code != "" {
		s.repo.Model(item).Where("code=? and ent_id=?", item.Code, item.EntID).Take(&oldItem)
	}
	if oldItem.ID != "" {
		updates := make(map[string]interface{})
		if oldItem.Name != item.Name && item.Name != "" {
			updates["Name"] = item.Name
		}
		if item.System.Valid() {
			updates["System"] = item.System
		}
		updates["DevHost"] = item.DevHost
		updates["TestHost"] = item.TestHost
		updates["PreHost"] = item.PreHost
		updates["ProdHost"] = item.ProdHost

		if len(updates) > 0 {
			s.repo.Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
		}
		item.ID = oldItem.ID

	} else {
		if item.ID == "" {
			item.ID = utils.GUID()
		}
		if item.Name == "" {
			item.Name = item.Code
		}
		item.CreatedAt = utils.NewTime()
		s.repo.Create(item)
	}
	return item, nil
}
func (s *ProductSv) DeleteHosts(entID string, ids []string) error {
	systemDatas := 0
	s.repo.Model(ProductHost{}).Where("ent_id=? and id in (?) and `system`=1", entID, ids).Count(&systemDatas)
	if systemDatas > 0 {
		return utils.NewError("系统预制不能删除!")
	}
	if err := s.repo.Delete(ProductHost{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}

/*==================modules============*/
func (s *ProductSv) GetModuleByCode(entID, idOrCode string) (*ProductModule, error) {
	old := ProductModule{}
	if err := s.repo.Where("ent_id=? and (code=? or id=?)", entID, idOrCode, idOrCode).Take(&old).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
func (s *ProductSv) SaveModules(item *ProductModule) (*ProductModule, error) {
	if item.EntID == "" || item.Code == "" {
		return nil, errors.ParamsRequired("entid or code")
	}
	oldItem := ProductModule{}
	if item.ID != "" {
		s.repo.Model(item).Where("id=?", item.ID).Take(&oldItem)
	}
	if oldItem.ID == "" && item.Code != "" {
		s.repo.Model(item).Where("code=? and ent_id=?", item.Code, item.EntID).Take(&oldItem)
	}
	if oldItem.ID != "" {
		updates := make(map[string]interface{})
		if oldItem.Name != item.Name && item.Name != "" {
			updates["Name"] = item.Name
		}
		if len(updates) > 0 {
			s.repo.Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
		}
		item.ID = oldItem.ID

	} else {
		if item.ID == "" {
			item.ID = utils.GUID()
		}
		item.CreatedAt = utils.NewTime()
		s.repo.Create(item)
	}
	return item, nil
}
func (s *ProductSv) DeleteModules(entID string, ids []string) error {
	systemDatas := 0
	s.repo.Model(ProductModule{}).Where("ent_id=? and id in (?) and `system`=1", entID, ids).Count(&systemDatas)
	if systemDatas > 0 {
		return utils.NewError("系统预制不能删除!")
	}
	if err := s.repo.Delete(ProductModule{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}

//service
func (s *ProductSv) SaveService(item *ProductService) (*ProductService, error) {
	if item.EntID == "" || item.Code == "" || item.ProductID == "" {
		return nil, errors.ParamsRequired("entid or code or ProductID")
	}
	oldItem := ProductService{}
	if item.ID != "" {
		s.repo.Model(item).Where("id=?", item.ID).Take(&oldItem)
	}
	if oldItem.ID == "" && item.Code != "" {
		s.repo.Model(item).Where("code=? and ent_id=?", item.Code, item.EntID).Take(&oldItem)
	}
	if oldItem.ID != "" {
		updates := make(map[string]interface{})

		if oldItem.Code != item.Code && item.Code != "" {
			updates["Code"] = item.Code
		}
		if oldItem.Uri != item.Uri && item.Uri != "" {
			updates["Uri"] = item.Uri
		}
		if oldItem.AppUri != item.AppUri && item.AppUri != "" {
			updates["AppUri"] = item.AppUri
		}
		if oldItem.Params != item.Params && item.Params != "" {
			updates["Params"] = item.Params
		}
		if oldItem.ProductID != item.ProductID && item.ProductID != "" {
			updates["ProductID"] = item.ProductID
		}
		if oldItem.Name != item.Name && item.Name != "" {
			updates["Name"] = item.Name
		}
		if oldItem.Icon != item.Icon && item.Icon != "" {
			updates["Icon"] = item.Icon
		}
		if oldItem.HostID != item.HostID && item.HostID != "" {
			updates["HostID"] = item.HostID
		}
		if oldItem.BizType != item.BizType && item.BizType != "" {
			updates["BizType"] = item.BizType
		}
		if oldItem.Sequence != item.Sequence && item.Sequence > 0 {
			updates["Sequence"] = item.Sequence
		}
		if oldItem.Memo != item.Memo && item.Memo != "" {
			updates["Memo"] = item.Memo
		}
		if oldItem.Tags != item.Tags && item.Tags != "" {
			updates["Tags"] = item.Tags
		}

		if item.InApp.Valid() && oldItem.InApp.NotEqual(item.InApp) {
			updates["InApp"] = item.InApp
		}
		if oldItem.Schema != item.Schema && item.Schema != "" {
			updates["Schema"] = item.Schema
		}
		if item.InWeb.Valid() && oldItem.InWeb.NotEqual(item.InWeb) {
			updates["InWeb"] = item.InWeb
		}
		if item.IsMaster.Valid() && oldItem.IsMaster.NotEqual(item.IsMaster) {
			updates["IsMaster"] = item.IsMaster
		}
		if item.IsSlave.Valid() && oldItem.IsSlave.NotEqual(item.IsSlave) {
			updates["IsSlave"] = item.IsSlave
		}
		if item.IsDefault.Valid() && oldItem.IsDefault.NotEqual(item.IsDefault) {
			updates["IsDefault"] = item.IsDefault
		}
		if len(updates) > 0 {
			s.repo.Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
		}
		item.ID = oldItem.ID

	} else {
		item.ID = utils.GUID()
		item.CreatedAt = utils.NewTime()
		s.repo.Create(item)
	}
	return item, nil
}
func (s *ProductSv) GetServiceByCode(entID, idOrCode string) (*ProductService, error) {
	old := ProductService{}
	if err := s.repo.Where("ent_id=? and (code=? or id=?)", entID, idOrCode, idOrCode).Take(&old).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
func (s *ProductSv) DeleteServices(entID string, ids []string) error {
	systemRoles := 0
	s.repo.Model(ProductService{}).Where("ent_id=? and id in (?) and `system`=1", entID, ids).Count(&systemRoles)
	if systemRoles > 0 {
		return utils.NewError("系统预制不能删除!")
	}
	if err := s.repo.Delete(ProductService{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}
