package services

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ggoop/mdf/bootstrap/errors"
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/files"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type QuerySv struct {
	repo *repositories.MysqlRepo
}

func NewQuerySv(repo *repositories.MysqlRepo) *QuerySv {
	return &QuerySv{repo: repo}
}
func (s *QuerySv) GetQuery(queryIDOrCode string) (*md.Query, error) {
	q := md.Query{}
	if err := s.repo.Preload("Columns").Preload("Orders").Preload("Wheres").Preload("Filters").Where("id=? or code=?", queryIDOrCode, queryIDOrCode).Take(&q).Error; err != nil {
		return nil, err
	}
	return &q, nil
}
func (s *QuerySv) AddQueries(item []md.Query) error {
	for _, item := range item {
		s.SaveQuery(model.SYS_ENT_ID, item, true)
	}
	return nil
}

func (s *QuerySv) SaveQuery(entID string, item md.Query, reInit bool) (*md.Query, error) {
	if !utils.StringIsCode(item.Code) {
		return nil, errors.CodeError(item.Code)
	}
	if item.Type == "" {
		item.Type = md.TYPE_ENTITY
	}
	old := md.Query{}
	s.repo.Where("id=? or code=?", item.ID, item.Code).Take(&old)
	if old.ID != "" {
		item.ID = old.ID
		updates := make(map[string]interface{})
		if old.Name != item.Name && item.Name != "" {
			updates["Name"] = item.Name
		}
		if old.Code != item.Code && item.Code != "" {
			updates["Code"] = item.Code
		}
		if old.Type != item.Type && item.Type != "" {
			updates["Type"] = item.Type
		}
		if old.PageSize != item.PageSize && item.PageSize > 0 {
			updates["PageSize"] = item.PageSize
		}
		if old.Entry != item.Entry && item.Entry != "" {
			updates["Entry"] = item.Entry
		}
		if old.ContextJson != item.ContextJson && item.ContextJson != "" {
			updates["ContextJson"] = item.ContextJson
		}
		if old.Condition != item.Condition && item.Condition != "" {
			updates["Condition"] = item.Condition
		}
		if len(updates) > 0 {
			s.repo.Model(old).Where("id=?", old.ID).Updates(updates)
		}
	} else {
		if item.ID == "" {
			item.ID = utils.GUID()
		}
		s.repo.Create(&item)
	}
	{
		items := make([]md.QueryFilter, 0)
		items = append(items, item.Filters...)
		for _, d := range item.Wheres {
			for _, f := range d.Filters {
				f.OwnerField = d.Field
				items = append(items, f)
			}
		}
		s.saveQueryFilters(item.ID, items, reInit)
	}
	s.saveQueryWheres(item.ID, item.ID, "query", item.Wheres, reInit)
	s.saveQueryColumns(item.ID, item.ID, "query", item.Columns, reInit)
	s.saveQueryOrders(item.ID, item.ID, "query", item.Orders, reInit)
	return nil, nil
}
func (s *QuerySv) saveQueryFilters(queryID string, items []md.QueryFilter, reInit bool) {
	itemNames := make([]string, 0)
	for i, ditem := range items {
		if ditem.Sequence == 0 {
			ditem.Sequence = i
		}
		itemNames = append(itemNames, fmt.Sprintf("%s:%s", ditem.OwnerField, ditem.Field))
		oldItem := md.QueryFilter{}
		if s.repo.Where("query_id=? and owner_field=? and field=?", queryID, ditem.OwnerField, ditem.Field).Take(&oldItem); oldItem.ID == "" {
			ditem.ID = utils.GUID()
			ditem.QueryID = queryID
			s.repo.Create(&ditem)
		} else {
			updates := make(map[string]interface{})
			if ditem.Operator != oldItem.Operator && ditem.Operator != "" {
				updates["Operator"] = ditem.Operator
			}
			if ditem.DataType != oldItem.DataType && ditem.DataType != "" {
				updates["DataType"] = ditem.DataType
			}
			if ditem.DataSource != oldItem.DataSource && ditem.DataSource != "" {
				updates["DataSource"] = ditem.DataSource
			}
			if ditem.Expr != oldItem.Expr && ditem.Expr != "" {
				updates["Expr"] = ditem.Expr
			}
			if ditem.Title != oldItem.Title && ditem.Title != "" {
				updates["Title"] = ditem.Title
			}
			if ditem.Value.Valid() && ditem.Value.NotEqual(oldItem.Value) {
				updates["Value"] = ditem.Value
			}
			if ditem.Enabled.Valid() && ditem.Enabled.NotEqual(oldItem.Enabled) {
				updates["Enabled"] = ditem.Enabled
			}
			if len(updates) > 0 {
				s.repo.Model(oldItem).Where("id=?", oldItem.ID).Updates(updates)
			}
		}
	}
	if reInit {
		if len(itemNames) > 0 {
			s.repo.Delete(md.QueryFilter{}, "query_id=? and concat(owner_field,':',field) not in (?)", queryID, itemNames)
		} else {
			s.repo.Delete(md.QueryFilter{}, "query_id=?", queryID)
		}
	}
}
func (s *QuerySv) saveQueryOrders(queryID, ownerID, ownerType string, items []md.QueryOrder, reInit bool) {
	itemNames := make([]string, 0)
	for i, d := range items {
		if d.Sequence == 0 {
			d.Sequence = i
		}
		itemNames = append(itemNames, d.Field)
		f := md.QueryOrder{}
		if s.repo.Where("owner_id=? and field=?", ownerID, d.Field).Take(&f); f.ID == "" {
			d.ID = utils.GUID()
			d.QueryID = queryID
			d.OwnerID = ownerID
			d.OwnerType = ownerType
			s.repo.Create(&d)
		} else {
			updates := make(map[string]interface{})
			if f.QueryID != queryID && queryID != "" {
				updates["QueryID"] = queryID
			}
			if f.OwnerType != ownerType && ownerType != "" {
				updates["OwnerType"] = ownerType
			}
			if f.Title != d.Title && d.Title != "" {
				updates["Title"] = d.Title
			}
			if f.Expr != d.Expr && d.Expr != "" {
				updates["Expr"] = d.Expr
			}
			if f.Sequence != d.Sequence && d.Sequence > 0 {
				updates["Sequence"] = d.Sequence
			}
			if f.Order != d.Order && d.Order != "" {
				updates["Order"] = d.Order
			}
			if f.Enabled.NotEqual(d.Enabled) && d.Enabled.Valid() {
				updates["Enabled"] = d.Enabled
			}
			if f.IsDefault.NotEqual(d.IsDefault) && d.IsDefault.Valid() {
				updates["IsDefault"] = d.IsDefault
			}
			if f.Fixed.NotEqual(d.Fixed) && d.Fixed.Valid() {
				updates["Fixed"] = d.Fixed
			}
			if f.Hidden.NotEqual(d.Hidden) && d.Hidden.Valid() {
				updates["Hidden"] = d.Hidden
			}
			if len(updates) > 0 {
				s.repo.Model(f).Where("id=?", f.ID).Updates(updates)
			}
		}
	}
	if reInit {
		if ownerType == "query" {
			if len(itemNames) > 0 {
				s.repo.Delete(md.QueryOrder{}, "query_id=? and field not in (?)", queryID, itemNames)
			} else {
				s.repo.Delete(md.QueryOrder{}, "query_id=?", queryID)
			}
		} else {
			if len(itemNames) > 0 {
				s.repo.Delete(md.QueryOrder{}, "owner_id=? and field not in (?)", ownerID, itemNames)
			} else {
				s.repo.Delete(md.QueryOrder{}, "owner_id=?", ownerID)
			}
		}
	}
}
func (s *QuerySv) saveQueryColumns(queryID, ownerID, ownerType string, items []md.QueryColumn, reInit bool) {
	itemNames := make([]string, 0)
	for i, d := range items {
		if d.Sequence == 0 {
			d.Sequence = i + 1
		}
		itemNames = append(itemNames, d.Field)
		f := md.QueryColumn{}
		if s.repo.Where("owner_id=? and field=?", ownerID, d.Field).Take(&f); f.ID == "" {
			d.ID = utils.GUID()
			d.OwnerID = ownerID
			d.OwnerType = ownerType
			d.QueryID = queryID

			s.repo.Create(&d)
		} else {
			updates := make(map[string]interface{})
			if f.QueryID != queryID && queryID != "" {
				updates["QueryID"] = queryID
			}
			if f.OwnerType != ownerType && ownerType != "" {
				updates["OwnerType"] = ownerType
			}
			if f.Title != d.Title && d.Title != "" {
				updates["Title"] = d.Title
			}
			if f.Name != d.Name && d.Name != "" {
				updates["Name"] = d.Name
			}
			if f.Sequence != d.Sequence && d.Sequence > 0 {
				updates["Sequence"] = d.Sequence
			}
			if f.Expr != d.Expr && d.Expr != "" {
				updates["Expr"] = d.Expr
			}
			if f.Sequence != d.Sequence && d.Sequence > 0 {
				updates["Sequence"] = d.Sequence
			}
			if f.Width != d.Width && d.Width != "" {
				updates["Width"] = d.Width
			}
			if f.Tags != d.Tags && d.Tags != "" {
				updates["Tags"] = d.Tags
			}
			if f.DataType != d.DataType && d.DataType != "" {
				updates["DataType"] = d.DataType
			}
			if f.KeyField.NotEqual(d.KeyField) && d.KeyField.Valid() {
				updates["KeyField"] = d.KeyField
			}
			if f.CodeField.NotEqual(d.CodeField) && d.CodeField.Valid() {
				updates["CodeField"] = d.CodeField
			}
			if f.NameField.NotEqual(d.NameField) && d.NameField.Valid() {
				updates["NameField"] = d.NameField
			}
			if f.Enabled.NotEqual(d.Enabled) && d.Enabled.Valid() {
				updates["Enabled"] = d.Enabled
			}
			if f.IsDefault.NotEqual(d.IsDefault) && d.IsDefault.Valid() {
				updates["IsDefault"] = d.IsDefault
			}
			if f.Fixed.NotEqual(d.Fixed) && d.Fixed.Valid() {
				updates["Fixed"] = d.Fixed
			}
			if f.Hidden.NotEqual(d.Hidden) && d.Hidden.Valid() {
				updates["Hidden"] = d.Hidden
			}
			if len(updates) > 0 {
				s.repo.Model(f).Where("id=?", f.ID).Updates(updates)
			}
		}
	}
	if reInit {
		if ownerType == "query" {
			if len(itemNames) > 0 {
				s.repo.Delete(md.QueryColumn{}, "query_id=? and field not in (?)", queryID, itemNames)
			} else {
				s.repo.Delete(md.QueryColumn{}, "query_id=?", queryID)
			}
		} else {
			if len(itemNames) > 0 {
				s.repo.Delete(md.QueryColumn{}, "owner_id=? and field not in (?)", ownerID, itemNames)
			} else {
				s.repo.Delete(md.QueryColumn{}, "owner_id=?", ownerID)
			}
		}
	}
}
func (s *QuerySv) saveQueryWheres(queryID, ownerID, ownerType string, items []md.QueryWhere, reInit bool) {
	itemNames := make([]string, 0)
	for i, d := range items {
		if d.Sequence == 0 {
			d.Sequence = i
		}
		itemNames = append(itemNames, d.Field)
		f := md.QueryWhere{}
		if s.repo.Where("owner_id=? and field=?", ownerID, d.Field).Take(&f); f.ID == "" {
			d.ID = utils.GUID()
			d.QueryID = queryID
			d.OwnerID = ownerID
			d.OwnerType = ownerType
			s.repo.Create(&d)
		} else {
			updates := make(map[string]interface{})
			if f.QueryID != queryID && queryID != "" {
				updates["QueryID"] = queryID
			}
			if f.OwnerType != ownerType && ownerType != "" {
				updates["OwnerType"] = ownerType
			}
			if f.DataType != d.DataType && d.DataType != "" {
				updates["DataType"] = d.DataType
			}
			if f.DataSource != d.DataSource && d.DataSource != "" {
				updates["DataSource"] = d.DataSource
			}
			if f.Title != d.Title && d.Title != "" {
				updates["Title"] = d.Title
			}
			if f.Expr != d.Expr && d.Expr != "" {
				updates["Expr"] = d.Expr
			}
			if f.Operator != d.Operator && d.Operator != "" {
				updates["Operator"] = d.Operator
			}
			if f.Value.NotEqual(d.Value) && d.Value.Valid() {
				updates["Value"] = d.Value
			}
			if f.Sequence != d.Sequence && d.Sequence > 0 {
				updates["Sequence"] = d.Sequence
			}
			if f.Enabled.NotEqual(d.Enabled) && d.Enabled.Valid() {
				updates["Enabled"] = d.Enabled
			}
			if f.IsDefault.NotEqual(d.IsDefault) && d.IsDefault.Valid() {
				updates["IsDefault"] = d.IsDefault
			}
			if f.Fixed.NotEqual(d.Fixed) && d.Fixed.Valid() {
				updates["Fixed"] = d.Fixed
			}
			if f.Hidden.NotEqual(d.Hidden) && d.Hidden.Valid() {
				updates["Hidden"] = d.Hidden
			}
			if d.IsBasic.Valid() && f.IsBasic.NotEqual(d.IsBasic) {
				updates["IsBasic"] = d.IsBasic
			}
			if len(updates) > 0 {
				s.repo.Model(f).Where("id=?", f.ID).Updates(updates)
			}
		}
	}
	if reInit {
		if ownerType == "query" {
			if len(itemNames) > 0 {
				if len(itemNames) > 0 {
					s.repo.Delete(md.QueryWhere{}, "query_id=? and field not in (?)", queryID, itemNames)
				} else {
					s.repo.Delete(md.QueryWhere{}, "query_id=?", queryID)
				}
			}
		} else {
			if len(itemNames) > 0 {
				s.repo.Delete(md.QueryWhere{}, "owner_id=? and field not in (?)", ownerID, itemNames)
			} else {
				s.repo.Delete(md.QueryWhere{}, "owner_id=?", ownerID)
			}
		}
	}
}
func (s *QuerySv) SaveCase(entID, queryID string, item md.QueryCase) (*md.QueryCase, error) {
	old := md.QueryCase{}
	if item.ID != "" {
		s.repo.Where("ent_id=?", entID).Where("id=?", item.ID).Take(&old)
	}
	if old.ID != "" {
		if old.UserID != item.UserID {
			return nil, glog.Error("没有权限修改")
		}
		updates := make(map[string]interface{})
		if old.Name != item.Name && item.Name != "" {
			updates["Name"] = item.Name
		}
		if old.ScopeType != item.ScopeType && item.ScopeType != "" {
			updates["ScopeType"] = item.ScopeType
		}
		if old.ScopeID != item.ScopeID && item.ScopeID != "" {
			updates["ScopeID"] = item.ScopeID
		}
		if old.PageSize != item.PageSize && item.PageSize > 0 {
			updates["PageSize"] = item.PageSize
		}
		if old.Memo != item.Memo && item.Memo != "" {
			updates["Memo"] = item.Memo
		}
		if old.IsDefault.NotEqual(item.IsDefault) && item.IsDefault.Valid() {
			updates["IsDefault"] = item.IsDefault
		}
		if len(updates) > 0 {
			s.repo.Model(old).Where("id=?", old.ID).Updates(updates)
		}
		item.ID = old.ID
	} else {
		item.ID = utils.GUID()
		item.QueryID = queryID
		item.EntID = entID
		//如果是用户保存的预制方案，则调整为用户
		if (item.ScopeType == "sys" && item.UserID != "") || item.ScopeType == "" {
			item.ScopeType = "user"
		}
		s.repo.Create(&item)
	}
	//同一企业和查询，只能有一个默认方案
	if item.IsDefault.Valid() && item.IsDefault.IsTrue() {
		updates := make(map[string]interface{})
		updates["IsDefault"] = utils.SBool_False
		if err := s.repo.Model(md.QueryCase{}).Where("ent_id=? and query_id=? and user_id=? and id!=?", entID, queryID, item.UserID, item.ID).Updates(updates).Error; err != nil {
			glog.Error(err)
		}
	}

	s.saveQueryWheres(item.QueryID, item.ID, "case", item.Wheres, true)
	s.saveQueryColumns(item.QueryID, item.ID, "case", item.Columns, true)
	s.saveQueryOrders(item.QueryID, item.ID, "case", item.Orders, true)
	return s.GetCase(entID, queryID, item.ID)
}
func (s *QuerySv) GetCases(entID, userID, queryID string, scopeTypes ...string) ([]md.QueryCase, error) {
	items := make([]md.QueryCase, 0)
	queryItem, err := s.GetQuery(queryID)
	if err != nil {
		return nil, err
	}
	if queryItem.ID != "" {
		query := s.repo.Where("ent_id=? and query_id=?", entID, queryItem.ID)
		query = query.Where("user_id=? or scope_type=?", userID, "public").Order("scope_type,id desc")
		query.Find(&items)
	}
	return items, nil
}
func (s *QuerySv) GetDefaultCase(entID, userID, queryID, caseID string) (*md.QueryCase, error) {
	queryItem, err := s.GetQuery(queryID)
	if err != nil {
		return nil, err
	}
	caseItem := md.QueryCase{}

	q := s.repo.Preload("Columns").Preload("Orders").Preload("Wheres").Preload("Filters")
	//依据查询方案精确查询
	if caseID != "" && caseID != "0" {
		q.Where("ent_id=? and id=?", entID, caseID).Take(&caseItem)
	}
	//查询用户默认方案
	if caseItem.ID == "" && queryItem.ID != "" && userID != "" && caseID != "0" {
		q.Where("ent_id=? and query_id=? and user_id=?", entID, queryItem.ID, userID).Order("is_default desc,id desc").Take(&caseItem)
	}
	//查询预制方案
	if caseItem.ID == "" && queryItem.ID != "" {
		caseItem = s.queryToCase(queryItem)
	}
	if caseItem.QueryID == "" {
		return nil, nil
	}
	return &caseItem, nil
}
func (s *QuerySv) queryToCase(q *md.Query) md.QueryCase {
	caseItem := md.QueryCase{}
	caseItem.ScopeType = "sys"
	caseItem.CreatedAt = q.CreatedAt
	caseItem.UpdatedAt = q.UpdatedAt
	caseItem.Name = "系统默认方案"
	for i, f := range q.Columns {
		if f.IsDefault.IsTrue() {
			caseItem.Columns = append(caseItem.Columns, q.Columns[i])
		}
	}
	for i, f := range q.Wheres {
		if f.IsDefault.IsTrue() {
			caseItem.Wheres = append(caseItem.Wheres, q.Wheres[i])
		}
	}
	for i, f := range q.Orders {
		if f.IsDefault.IsTrue() {
			caseItem.Orders = append(caseItem.Orders, q.Orders[i])
		}
	}
	caseItem.Filters = q.Filters
	caseItem.Memo = q.Memo
	caseItem.QueryID = q.ID
	caseItem.PageSize = q.PageSize

	return caseItem
}

//依据查询，获取查询方案
func (s *QuerySv) GetCase(entID, queryID, caseID string, preloads ...string) (*md.QueryCase, error) {
	queryItem, err := s.GetQuery(queryID)
	if err != nil {
		return nil, err
	}
	old := md.QueryCase{}
	if caseID != "0" && caseID != "" {
		if err := s.repo.Preload("Columns").Preload("Orders").Preload("Wheres").Preload("Filters").Where("ent_id=?", entID).Where("id=?", caseID).Take(&old).Error; err != nil {
			return nil, err
		}
	} else {
		old = s.queryToCase(queryItem)
	}

	return &old, nil
}
func (s *QuerySv) DeleteCases(entID, userID, queryID string, caseIDs []string) error {
	if err := s.repo.Where("ent_id=?", entID).Where("id in (?)", caseIDs).Delete(md.QueryCase{}).Error; err != nil {
		return err
	}
	return nil
}
func (s *QuerySv) BatchImport(entID string, datas []files.ImportData) error {
	if len(datas) <= 0 {
		return nil
	}
	nameList := make(map[string]int)
	nameList["queries"] = 1
	nameList["wheres"] = 2
	nameList["columns"] = 3
	nameList["orders"] = 4
	nameList["filters"] = 5
	sort.Slice(datas, func(i, j int) bool { return nameList[datas[i].Name] < nameList[datas[j].Name] })

	qitems := make([]md.Query, 0)
	qwheres := make([]md.QueryWhere, 0)
	qcolumns := make([]md.QueryColumn, 0)
	qorders := make([]md.QueryOrder, 0)
	qfilters := make([]md.QueryFilter, 0)
	for _, item := range datas {
		if strings.ToLower(item.Name) == "queries" {
			if d, err := s.importDataToQueries(entID, item); err != nil {
				return err
			} else if len(d) > 0 {
				qitems = append(qitems, d...)
			}
		}
		if strings.ToLower(item.Name) == "wheres" {
			if d, err := s.importDataToWheres(entID, item); err != nil {
				return err
			} else if len(d) > 0 {
				qwheres = append(qwheres, d...)
			}
		}
		if strings.ToLower(item.Name) == "columns" {
			if d, err := s.importDataToColumns(entID, item); err != nil {
				return err
			} else if len(d) > 0 {
				qcolumns = append(qcolumns, d...)
			}
		}
		if strings.ToLower(item.Name) == "orders" {
			if d, err := s.importDataToOrders(entID, item); err != nil {
				return err
			} else if len(d) > 0 {
				qorders = append(qorders, d...)
			}
		}
		if strings.ToLower(item.Name) == "filters" {
			if d, err := s.importDataToFilters(entID, item); err != nil {
				return err
			} else if len(d) > 0 {
				qfilters = append(qfilters, d...)
			}
		}
	}
	for i, q := range qitems {
		for di, d := range qwheres {
			if q.ID == d.OwnerID {
				qitems[i].Wheres = append(qitems[i].Wheres, qwheres[di])
			}
		}
		for di, d := range qcolumns {
			if q.ID == d.OwnerID {
				qitems[i].Columns = append(qitems[i].Columns, qcolumns[di])
			}
		}
		for di, d := range qorders {
			if q.ID == d.OwnerID {
				qitems[i].Orders = append(qitems[i].Orders, qorders[di])
			}
		}
		for di, d := range qfilters {
			if q.ID == d.QueryID {
				qitems[i].Filters = append(qitems[i].Filters, qfilters[di])
			}
		}
	}
	for _, q := range qitems {
		if _, err := s.SaveQuery(entID, q, true); err != nil {
			return err
		}
	}
	return nil
}
func (s *QuerySv) importDataToQueries(entID string, data files.ImportData) ([]md.Query, error) {
	if len(data.Datas) == 0 {
		return nil, nil
	}
	items := make([]md.Query, 0)
	for _, row := range data.Datas {
		item := md.Query{}
		if cValue := files.GetMapStringValue("ID", row); cValue != "" {
			item.ID = cValue
		} else {
			continue
		}
		if cValue := files.GetMapStringValue("Code", row); cValue != "" {
			item.Code = cValue
		}
		if cValue := files.GetMapStringValue("Name", row); cValue != "" {
			item.Name = cValue
		}
		if cValue := files.GetMapStringValue("Type", row); cValue != "" {
			item.Type = cValue
		}
		if cValue := files.GetMapStringValue("Entry", row); cValue != "" {
			item.Entry = cValue
		}
		if cValue := files.GetMapStringValue("Memo", row); cValue != "" {
			item.Memo = cValue
		}
		if cValue := files.GetMapStringValue("Condition", row); cValue != "" {
			item.Condition = cValue
		}
		item.PageSize = files.GetMapIntValue("PageSize", row)
		items = append(items, item)
	}
	return items, nil
}
func (s *QuerySv) importDataToWheres(entID string, data files.ImportData) ([]md.QueryWhere, error) {
	if len(data.Datas) == 0 {
		return nil, nil
	}
	items := make([]md.QueryWhere, 0)
	for _, row := range data.Datas {
		item := md.QueryWhere{}
		if cValue := files.GetMapStringValue("OwnerID", row); cValue != "" {
			item.OwnerID = cValue
		} else {
			continue
		}
		if cValue := files.GetMapStringValue("OwnerType", row); cValue != "" {
			item.OwnerType = cValue
		}
		if cValue := files.GetMapStringValue("Field", row); cValue != "" {
			item.Field = cValue
		}
		if cValue := files.GetMapStringValue("Expr", row); cValue != "" {
			item.Expr = cValue
		}
		if cValue := files.GetMapStringValue("Title", row); cValue != "" {
			item.Title = cValue
		}
		if cValue := files.GetMapStringValue("Operator", row); cValue != "" {
			item.Operator = cValue
		}
		if cValue := files.GetMapStringValue("DataType", row); cValue != "" {
			item.DataType = cValue
		}
		if cValue := files.GetMapStringValue("DataSource", row); cValue != "" {
			item.DataSource = cValue
		}
		if cValue := files.GetMapStringValue("Logical", row); cValue != "" {
			item.Logical = cValue
		}
		item.Sequence = files.GetMapIntValue("Sequence", row)
		item.Value = files.GetMapSJsonValue("Value", row)
		item.IsBasic = files.GetMapSBoolValue("IsBasic", row)

		item.Fixed = files.GetMapSBoolValue("Fixed", row)
		item.IsDefault = files.GetMapSBoolValue("IsDefault", row)
		item.Enabled = files.GetMapSBoolValue("Enabled", row)
		item.Hidden = files.GetMapSBoolValue("Hidden", row)

		items = append(items, item)
	}
	return items, nil
}
func (s *QuerySv) importDataToColumns(entID string, data files.ImportData) ([]md.QueryColumn, error) {
	if len(data.Datas) == 0 {
		return nil, nil
	}
	items := make([]md.QueryColumn, 0)
	for _, row := range data.Datas {
		item := md.QueryColumn{}
		if cValue := files.GetMapStringValue("OwnerID", row); cValue != "" {
			item.OwnerID = cValue
		} else {
			continue
		}
		if cValue := files.GetMapStringValue("OwnerType", row); cValue != "" {
			item.OwnerType = cValue
		}
		if cValue := files.GetMapStringValue("Field", row); cValue != "" {
			item.Field = cValue
		}
		if cValue := files.GetMapStringValue("Expr", row); cValue != "" {
			item.Expr = cValue
		}
		if cValue := files.GetMapStringValue("Title", row); cValue != "" {
			item.Title = cValue
		}
		if cValue := files.GetMapStringValue("Name", row); cValue != "" {
			item.Name = cValue
		}
		if cValue := files.GetMapStringValue("DataType", row); cValue != "" {
			item.DataType = cValue
		}
		if cValue := files.GetMapStringValue("Tags", row); cValue != "" {
			item.Tags = cValue
		}
		if cValue := files.GetMapStringValue("Width", row); cValue != "" {
			item.Width = cValue
		}

		item.Sequence = files.GetMapIntValue("Sequence", row)
		item.KeyField = files.GetMapSBoolValue("KeyField", row)
		item.CodeField = files.GetMapSBoolValue("CodeField", row)
		item.NameField = files.GetMapSBoolValue("NameField", row)

		item.Fixed = files.GetMapSBoolValue("Fixed", row)
		item.IsDefault = files.GetMapSBoolValue("IsDefault", row)
		item.Enabled = files.GetMapSBoolValue("Enabled", row)
		item.Hidden = files.GetMapSBoolValue("Hidden", row)

		items = append(items, item)
	}
	return items, nil
}

func (s *QuerySv) importDataToOrders(entID string, data files.ImportData) ([]md.QueryOrder, error) {
	if len(data.Datas) == 0 {
		return nil, nil
	}
	items := make([]md.QueryOrder, 0)
	for _, row := range data.Datas {
		item := md.QueryOrder{}
		if cValue := files.GetMapStringValue("OwnerID", row); cValue != "" {
			item.OwnerID = cValue
		} else {
			continue
		}
		if cValue := files.GetMapStringValue("OwnerType", row); cValue != "" {
			item.OwnerType = cValue
		}
		if cValue := files.GetMapStringValue("Field", row); cValue != "" {
			item.Field = cValue
		}
		if cValue := files.GetMapStringValue("Expr", row); cValue != "" {
			item.Expr = cValue
		}
		if cValue := files.GetMapStringValue("Title", row); cValue != "" {
			item.Title = cValue
		}
		if cValue := files.GetMapStringValue("Order", row); cValue != "" {
			item.Order = cValue
		}

		item.Sequence = files.GetMapIntValue("Sequence", row)

		item.Fixed = files.GetMapSBoolValue("Fixed", row)
		item.IsDefault = files.GetMapSBoolValue("IsDefault", row)
		item.Enabled = files.GetMapSBoolValue("Enabled", row)
		item.Hidden = files.GetMapSBoolValue("Hidden", row)

		items = append(items, item)
	}
	return items, nil
}

func (s *QuerySv) importDataToFilters(entID string, data files.ImportData) ([]md.QueryFilter, error) {
	if len(data.Datas) == 0 {
		return nil, nil
	}
	items := make([]md.QueryFilter, 0)
	for _, row := range data.Datas {
		item := md.QueryFilter{}
		if cValue := files.GetMapStringValue("QueryID", row); cValue != "" {
			item.QueryID = cValue
		} else {
			continue
		}
		if cValue := files.GetMapStringValue("OwnerField", row); cValue != "" {
			item.OwnerField = cValue
		}
		if cValue := files.GetMapStringValue("Field", row); cValue != "" {
			item.Field = cValue
		}
		if cValue := files.GetMapStringValue("Expr", row); cValue != "" {
			item.Expr = cValue
		}
		if cValue := files.GetMapStringValue("Title", row); cValue != "" {
			item.Title = cValue
		}
		if cValue := files.GetMapStringValue("Operator", row); cValue != "" {
			item.Operator = cValue
		}
		if cValue := files.GetMapStringValue("DataType", row); cValue != "" {
			item.DataType = cValue
		}
		if cValue := files.GetMapStringValue("DataSource", row); cValue != "" {
			item.DataSource = cValue
		}
		item.Sequence = files.GetMapIntValue("Sequence", row)
		item.Value = files.GetMapSJsonValue("Value", row)
		item.Enabled = files.GetMapSBoolValue("Enabled", row)
		items = append(items, item)
	}
	return items, nil
}
