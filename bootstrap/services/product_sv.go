package services

import (
	"sort"
	"sync"

	"github.com/ggoop/mdf/bootstrap/errors"
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/framework/files"
	"github.com/ggoop/mdf/framework/glog"
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

func (s *ProductSv) BatchImport(entID string, datas []files.ImportData) error {
	nameList := make(map[string]int)
	nameList["product"] = 1
	nameList["package"] = 2
	nameList["host"] = 3
	nameList["service"] = 4
	sort.Slice(datas, func(i, j int) bool { return nameList[datas[i].Name] < nameList[datas[j].Name] })
	for i, item := range datas {
		if item.Name == "product" {
			if _, err := s.importProduct(entID, datas[i]); err != nil {
				return err
			}
		}
		if item.Name == "host" {
			if _, err := s.importHost(entID, datas[i]); err != nil {
				return err
			}
		}
		if item.Name == "service" {
			if _, err := s.importService(entID, datas[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *ProductSv) importProduct(entID string, datas files.ImportData) (int, error) {
	if len(datas.Datas) == 0 {
		return 0, nil
	}
	for i, row := range datas.Datas {
		item := &model.Product{}
		if cValue := files.GetMapStringValue("Code", row); cValue != "" {
			item.Code = cValue
		}
		if cValue := files.GetMapStringValue("Name", row); cValue != "" {
			item.Name = cValue
		}
		if item.Code == "" {
			glog.Error("产品编码为空", glog.Int("Line", i))
			continue
		}
		item.Icon = files.GetMapStringValue("Icon", row)
		item.EntID = entID
		s.SaveProduct(item)
	}
	return 0, nil
}
func (s *ProductSv) importHost(entID string, datas files.ImportData) (int, error) {
	if len(datas.Datas) == 0 {
		return 0, nil
	}
	for i, row := range datas.Datas {
		item := &model.ProductHost{}
		if cValue := files.GetMapStringValue("Code", row); cValue != "" {
			item.Code = cValue
		}
		if cValue := files.GetMapStringValue("Name", row); cValue != "" {
			item.Name = cValue
		}
		if item.Code == "" {
			glog.Error("编码为空", glog.Int("Line", i))
			continue
		}
		item.EntID = entID
		s.SaveHosts(item)
	}
	return 0, nil
}
func (s *ProductSv) importService(entID string, datas files.ImportData) (int, error) {
	s.Lock()
	defer s.Unlock()
	if len(datas.Datas) == 0 {
		return 0, nil
	}
	for i, row := range datas.Datas {
		//product
		product := &model.Product{}
		if cValue := files.GetMapStringValue("ProductCode", row); cValue != "" {
			product, _ = s.GetProductByCode(entID, cValue)
		}
		if product == nil || product.Code == "" {
			glog.Error("产品编码为空", glog.Int("Line", i))
			continue
		}
		//host
		phost := &model.ProductHost{}
		if cValue := files.GetMapStringValue("HostCode", row); cValue != "" {
			phost, _ = s.GetHostByCode(entID, cValue)
		}
		//model
		pmodule := &model.ProductModule{ProductID: product.ID}
		if cValue := files.GetMapStringValue("ModuleCode", row); cValue != "" {
			pmodule.Code = cValue
		}
		if cValue := files.GetMapStringValue("ModuleName", row); cValue != "" {
			pmodule.Name = cValue
		}
		if pmodule.Code == "" {
			glog.Error("模块编码为空", glog.Int("Line", i))
			continue
		}
		pmodule.EntID = entID
		pmodule, _ = s.SaveModules(pmodule)

		//service
		pService := &model.ProductService{ProductID: product.ID, ModuleID: pmodule.ID}
		if cValue := files.GetMapStringValue("ServiceCode", row); cValue != "" {
			pService.Code = cValue
		}
		if cValue := files.GetMapStringValue("ServiceName", row); cValue != "" {
			pService.Name = cValue
		}
		if pService.Code == "" {
			glog.Error("服务编码为空", glog.Int("Line", i))
			continue
		}
		if phost.ID != "" {
			pService.HostID = phost.ID
		}

		pService.AppUri = files.GetMapStringValue("AppUri", row)
		pService.Uri = files.GetMapStringValue("Uri", row)
		pService.InApp = files.GetMapSBoolValue("InApp", row)
		pService.Schema = files.GetMapStringValue("Schema", row)

		pService.Tags = files.GetMapStringValue("Tags", row)
		pService.Memo = files.GetMapStringValue("Memo", row)
		pService.BizType = files.GetMapStringValue("BizType", row)
		pService.Icon = files.GetMapStringValue("Icon", row)

		pService.InWeb = files.GetMapSBoolValue("InWeb", row)
		pService.IsMaster = utils.SBool_Parse(files.GetMapBoolValue("IsMaster", row, false))
		pService.IsSlave = utils.SBool_Parse(files.GetMapBoolValue("IsSlave", row, false))
		pService.IsDefault = utils.SBool_Parse(files.GetMapBoolValue("IsDefault", row, false))

		pService.Sequence = files.GetMapIntValue("Sequence", row)
		if pService.Sequence <= 0 {
			pService.Sequence = i + 1
		}
		pService.EntID = entID
		pService, _ = s.SaveService(pService)
	}
	return 0, nil
}

/*============product===========*/
func (s *ProductSv) SaveProduct(item *model.Product) (*model.Product, error) {
	if item.EntID == "" || item.Code == "" {
		return nil, errors.ParamsRequired("entid or code")
	}
	oldItem := model.Product{}
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
func (s *ProductSv) GetProductByCode(entID, idOrCode string) (*model.Product, error) {
	old := model.Product{}
	if err := s.repo.Where("ent_id=? and (code=? or id=?)", entID, idOrCode, idOrCode).Take(&old).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
func (s *ProductSv) DeleteProducts(entID string, ids []string) error {
	systemRoles := 0
	s.repo.Model(model.Product{}).Where("ent_id=? and id in (?) and `system`=1", entID, ids).Count(&systemRoles)
	if systemRoles > 0 {
		return utils.NewError("系统预制产品不能删除!")
	}
	if err := s.repo.Delete(model.Product{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}

/*==================host============*/
func (s *ProductSv) GetHostByCode(entID, idOrCode string) (*model.ProductHost, error) {
	old := model.ProductHost{}
	if err := s.repo.Where("ent_id=? and (code=? or id=?)", entID, idOrCode, idOrCode).Take(&old).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
func (s *ProductSv) SaveHosts(item *model.ProductHost) (*model.ProductHost, error) {
	if item.EntID == "" || item.Code == "" {
		return nil, errors.ParamsRequired("entid or code")
	}
	oldItem := model.ProductHost{}
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
	s.repo.Model(model.ProductHost{}).Where("ent_id=? and id in (?) and `system`=1", entID, ids).Count(&systemDatas)
	if systemDatas > 0 {
		return utils.NewError("系统预制不能删除!")
	}
	if err := s.repo.Delete(model.ProductHost{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}

/*==================modules============*/
func (s *ProductSv) GetModuleByCode(entID, idOrCode string) (*model.ProductModule, error) {
	old := model.ProductModule{}
	if err := s.repo.Where("ent_id=? and (code=? or id=?)", entID, idOrCode, idOrCode).Take(&old).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
func (s *ProductSv) SaveModules(item *model.ProductModule) (*model.ProductModule, error) {
	if item.EntID == "" || item.Code == "" {
		return nil, errors.ParamsRequired("entid or code")
	}
	oldItem := model.ProductModule{}
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
	s.repo.Model(model.ProductModule{}).Where("ent_id=? and id in (?) and `system`=1", entID, ids).Count(&systemDatas)
	if systemDatas > 0 {
		return utils.NewError("系统预制不能删除!")
	}
	if err := s.repo.Delete(model.ProductModule{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}

//service
func (s *ProductSv) SaveService(item *model.ProductService) (*model.ProductService, error) {
	if item.EntID == "" || item.Code == "" || item.ProductID == "" {
		return nil, errors.ParamsRequired("entid or code or ProductID")
	}
	oldItem := model.ProductService{}
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
func (s *ProductSv) GetServiceByCode(entID, idOrCode string) (*model.ProductService, error) {
	old := model.ProductService{}
	if err := s.repo.Where("ent_id=? and (code=? or id=?)", entID, idOrCode, idOrCode).Take(&old).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
func (s *ProductSv) DeleteServices(entID string, ids []string) error {
	systemRoles := 0
	s.repo.Model(model.ProductService{}).Where("ent_id=? and id in (?) and `system`=1", entID, ids).Count(&systemRoles)
	if systemRoles > 0 {
		return utils.NewError("系统预制不能删除!")
	}
	if err := s.repo.Delete(model.ProductService{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}
