package services

import (
	"github.com/ggoop/mdf/db"
	"sort"
	"sync"

	"github.com/ggoop/mdf/bootstrap/errors"
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/framework/files"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/utils"
)

type IProductSv interface {
}

type productSvImpl struct {
	*sync.Mutex
}

func ProductSv() IProductSv {
	return productSv
}

var productSv IProductSv = newProductSv()

/**
* 创建服务实例
 */
func newProductSv() *productSvImpl {
	return &productSvImpl{Mutex: &sync.Mutex{}}
}

func (s *productSvImpl) BatchImport(entID string, datas []files.ImportData) error {
	nameList := make(map[string]int)
	nameList["product"] = 1
	nameList["package"] = 2
	nameList["host"] = 3
	nameList["service"] = 4
	sort.Slice(datas, func(i, j int) bool { return nameList[datas[i].SheetName] < nameList[datas[j].SheetName] })
	for i, item := range datas {
		if item.SheetName == "product" {
			if _, err := s.importProduct(entID, datas[i]); err != nil {
				return err
			}
		}
		if item.SheetName == "host" {
			if _, err := s.importHost(entID, datas[i]); err != nil {
				return err
			}
		}
		if item.SheetName == "service" {
			if _, err := s.importService(entID, datas[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *productSvImpl) importProduct(entID string, datas files.ImportData) (int, error) {
	if len(datas.Data) == 0 {
		return 0, nil
	}
	for i, row := range datas.Data {
		item := &model.Product{}
		if cValue := files.GetCellValue("Code", row); cValue != "" {
			item.Code = cValue
		}
		if cValue := files.GetCellValue("Name", row); cValue != "" {
			item.Name = cValue
		}
		if item.Code == "" {
			glog.Error("产品编码为空", glog.Int("Line", i))
			continue
		}
		item.Icon = files.GetCellValue("Icon", row)
		item.EntID = entID
		s.SaveProduct(item)
	}
	return 0, nil
}
func (s *productSvImpl) importHost(entID string, datas files.ImportData) (int, error) {
	if len(datas.Data) == 0 {
		return 0, nil
	}
	for i, row := range datas.Data {
		item := &model.ProductHost{}
		if cValue := files.GetCellValue("Code", row); cValue != "" {
			item.Code = cValue
		}
		if cValue := files.GetCellValue("Name", row); cValue != "" {
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
func (s *productSvImpl) importService(entID string, datas files.ImportData) (int, error) {
	s.Lock()
	defer s.Unlock()
	if len(datas.Data) == 0 {
		return 0, nil
	}
	for i, row := range datas.Data {
		//product
		product := &model.Product{}
		if cValue := files.GetCellValue("ProductCode", row); cValue != "" {
			product, _ = s.GetProductByCode(entID, cValue)
		}
		if product == nil || product.Code == "" {
			glog.Error("产品编码为空", glog.Int("Line", i))
			continue
		}
		//host
		phost := &model.ProductHost{}
		if cValue := files.GetCellValue("HostCode", row); cValue != "" {
			phost, _ = s.GetHostByCode(entID, cValue)
		}
		//model
		pmodule := &model.ProductModule{ProductID: product.ID}
		if cValue := files.GetCellValue("ModuleCode", row); cValue != "" {
			pmodule.Code = cValue
		}
		if cValue := files.GetCellValue("ModuleName", row); cValue != "" {
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
		if cValue := files.GetCellValue("ServiceCode", row); cValue != "" {
			pService.Code = cValue
		}
		if cValue := files.GetCellValue("ServiceName", row); cValue != "" {
			pService.Name = cValue
		}
		if pService.Code == "" {
			glog.Error("服务编码为空", glog.Int("Line", i))
			continue
		}
		if phost.ID != "" {
			pService.HostID = phost.ID
		}

		pService.AppUri = files.GetCellValue("AppUri", row)
		pService.Uri = files.GetCellValue("Uri", row)
		pService.InApp = utils.ToSBool(files.GetCellValue("InApp", row))
		pService.Schema = files.GetCellValue("Schema", row)

		pService.Tags = files.GetCellValue("Tags", row)
		pService.Memo = files.GetCellValue("Memo", row)
		pService.BizType = files.GetCellValue("BizType", row)
		pService.Icon = files.GetCellValue("Icon", row)

		pService.InWeb = utils.ToSBool(files.GetCellValue("InWeb", row))
		pService.IsMaster = utils.ToSBool(files.GetCellValue("IsMaster", row))
		pService.IsSlave = utils.ToSBool(files.GetCellValue("IsSlave", row))
		pService.IsDefault = utils.ToSBool(files.GetCellValue("IsDefault", row))

		pService.Sequence = utils.ToInt(files.GetCellValue("Sequence", row))
		if pService.Sequence <= 0 {
			pService.Sequence = i + 1
		}
		pService.EntID = entID
		pService, _ = s.SaveService(pService)
	}
	return 0, nil
}

/*============product===========*/
func (s *productSvImpl) SaveProduct(item *model.Product) (*model.Product, error) {
	if item.EntID == "" || item.Code == "" {
		return nil, errors.ParamsRequired("entid or code")
	}
	oldItem := model.Product{}
	if item.ID != "" {
		db.Default().Model(item).Where("id=?", item.ID).Take(&oldItem)
	}
	if oldItem.ID == "" && item.Code != "" {
		db.Default().Model(item).Where("code=? and ent_id=?", item.Code, item.EntID).Take(&oldItem)
	}
	if oldItem.ID != "" {
		updates := make(map[string]interface{})
		if oldItem.Name != item.Name && item.Name != "" {
			updates["Name"] = item.Name
		}
		if len(updates) > 0 {
			db.Default().Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
		}
		item.ID = oldItem.ID
	} else {
		if item.ID == "" {
			item.ID = utils.GUID()
		}
		item.CreatedAt = utils.TimeNow()
		db.Default().Create(item)
	}
	return item, nil
}
func (s *productSvImpl) GetProductByCode(entID, idOrCode string) (*model.Product, error) {
	old := model.Product{}
	if err := db.Default().Where("ent_id=? and (code=? or id=?)", entID, idOrCode, idOrCode).Take(&old).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
func (s *productSvImpl) DeleteProducts(entID string, ids []string) error {
	systemRoles := 0
	db.Default().Model(model.Product{}).Where("ent_id=? and id in (?) and `system`=1", entID, ids).Count(&systemRoles)
	if systemRoles > 0 {
		return utils.ToError("系统预制产品不能删除!")
	}
	if err := db.Default().Delete(model.Product{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}

/*==================host============*/
func (s *productSvImpl) GetHostByCode(entID, idOrCode string) (*model.ProductHost, error) {
	old := model.ProductHost{}
	if err := db.Default().Where("ent_id=? and (code=? or id=?)", entID, idOrCode, idOrCode).Take(&old).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
func (s *productSvImpl) SaveHosts(item *model.ProductHost) (*model.ProductHost, error) {
	if item.EntID == "" || item.Code == "" {
		return nil, errors.ParamsRequired("entid or code")
	}
	oldItem := model.ProductHost{}
	if item.ID != "" {
		db.Default().Model(item).Where("id=?", item.ID).Take(&oldItem)
	}
	if oldItem.ID == "" && item.Code != "" {
		db.Default().Model(item).Where("code=? and ent_id=?", item.Code, item.EntID).Take(&oldItem)
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
			db.Default().Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
		}
		item.ID = oldItem.ID

	} else {
		if item.ID == "" {
			item.ID = utils.GUID()
		}
		if item.Name == "" {
			item.Name = item.Code
		}
		item.CreatedAt = utils.TimeNow()
		db.Default().Create(item)
	}
	return item, nil
}
func (s *productSvImpl) DeleteHosts(entID string, ids []string) error {
	systemDatas := 0
	db.Default().Model(model.ProductHost{}).Where("ent_id=? and id in (?) and `system`=1", entID, ids).Count(&systemDatas)
	if systemDatas > 0 {
		return utils.ToError("系统预制不能删除!")
	}
	if err := db.Default().Delete(model.ProductHost{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}

/*==================modules============*/
func (s *productSvImpl) GetModuleByCode(entID, idOrCode string) (*model.ProductModule, error) {
	old := model.ProductModule{}
	if err := db.Default().Where("ent_id=? and (code=? or id=?)", entID, idOrCode, idOrCode).Take(&old).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
func (s *productSvImpl) SaveModules(item *model.ProductModule) (*model.ProductModule, error) {
	if item.EntID == "" || item.Code == "" {
		return nil, errors.ParamsRequired("entid or code")
	}
	oldItem := model.ProductModule{}
	if item.ID != "" {
		db.Default().Model(item).Where("id=?", item.ID).Take(&oldItem)
	}
	if oldItem.ID == "" && item.Code != "" {
		db.Default().Model(item).Where("code=? and ent_id=?", item.Code, item.EntID).Take(&oldItem)
	}
	if oldItem.ID != "" {
		updates := make(map[string]interface{})
		if oldItem.Name != item.Name && item.Name != "" {
			updates["Name"] = item.Name
		}
		if len(updates) > 0 {
			db.Default().Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
		}
		item.ID = oldItem.ID

	} else {
		if item.ID == "" {
			item.ID = utils.GUID()
		}
		item.CreatedAt = utils.TimeNow()
		db.Default().Create(item)
	}
	return item, nil
}
func (s *productSvImpl) DeleteModules(entID string, ids []string) error {
	systemDatas := 0
	db.Default().Model(model.ProductModule{}).Where("ent_id=? and id in (?) and `system`=1", entID, ids).Count(&systemDatas)
	if systemDatas > 0 {
		return utils.ToError("系统预制不能删除!")
	}
	if err := db.Default().Delete(model.ProductModule{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}

//service
func (s *productSvImpl) SaveService(item *model.ProductService) (*model.ProductService, error) {
	if item.EntID == "" || item.Code == "" || item.ProductID == "" {
		return nil, errors.ParamsRequired("entid or code or ProductID")
	}
	oldItem := model.ProductService{}
	if item.ID != "" {
		db.Default().Model(item).Where("id=?", item.ID).Take(&oldItem)
	}
	if oldItem.ID == "" && item.Code != "" {
		db.Default().Model(item).Where("code=? and ent_id=?", item.Code, item.EntID).Take(&oldItem)
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
			db.Default().Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
		}
		item.ID = oldItem.ID

	} else {
		item.ID = utils.GUID()
		item.CreatedAt = utils.TimeNow()
		db.Default().Create(item)
	}
	return item, nil
}
func (s *productSvImpl) GetServiceByCode(entID, idOrCode string) (*model.ProductService, error) {
	old := model.ProductService{}
	if err := db.Default().Where("ent_id=? and (code=? or id=?)", entID, idOrCode, idOrCode).Take(&old).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
func (s *productSvImpl) DeleteServices(entID string, ids []string) error {
	systemRoles := 0
	db.Default().Model(model.ProductService{}).Where("ent_id=? and id in (?) and `system`=1", entID, ids).Count(&systemRoles)
	if systemRoles > 0 {
		return utils.ToError("系统预制不能删除!")
	}
	if err := db.Default().Delete(model.ProductService{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}
