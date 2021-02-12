package rules

import (
	"fmt"

	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type CommonQuery struct {
	repo *repositories.MysqlRepo
}

func NewCommonQuery(repo *repositories.MysqlRepo) *CommonQuery {
	return &CommonQuery{repo}
}
func (s *CommonQuery) GetRule() md.RuleRegister {
	return md.RuleRegister{Code: "query", OwnerType: md.RuleType_Widget, OwnerID: "common"}
}

func (s CommonQuery) Exec(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext) error {
	if req.Entity == "" {
		return glog.Error("缺少 MainEntity 参数！")
	}
	//查找实体信息
	entity := md.GetEntity(req.Entity)
	if entity == nil {
		return glog.Error("找不到实体！")
	}
	exector := md.NewExector(entity.TableName)
	for _, f := range entity.Fields {
		if f.TypeType == md.TYPE_SIMPLE {
			exector.Select(fmt.Sprintf("$$%s as \"%s\"", f.Code, f.DbName))
			if f.TypeID != "" && f.DbName != "" {
				exector.SetFieldDataType(f.DbName, f.TypeID)
			}
		}
	}
	if req.ID != "" {
		exector.Where("id=?", req.ID)
	}
	if sysField := entity.GetField("EntID"); sysField != nil && token.EntID() != "" {
		exector.Where(sysField.Code+" = ?", token.EntID())
	}
	exector.Page(req.Page, req.PageSize)

	count, err := exector.Count(s.repo)
	if err != nil {
		return err
	}
	if datas, err := exector.Query(s.repo); err != nil {
		return err
	} else if len(datas) > 0 {
		s.loadEnums(datas, entity)
		s.loadEntities(datas, entity)
		res.Set("data", datas)
		res.Set("pager", utils.Pager{Total: count, PageSize: req.PageSize, Page: req.Page})
	}
	return nil
}
func (s CommonQuery) loadEnums(datas []map[string]interface{}, entity *md.MDEntity) error {
	for _, f := range entity.Fields {
		if f.TypeType == md.TYPE_ENUM {
			for ri, data := range datas {
				if fv, ok := data[f.DbName+"_id"]; ok && fv != nil && fv.(string) != "" {
					datas[ri][f.DbName] = md.GetEnum(f.Limit, fv.(string))
				}
			}
		}
	}
	return nil
}
func (s CommonQuery) loadEntities(datas []map[string]interface{}, entity *md.MDEntity) error {
	for _, f := range entity.Fields {
		if f.TypeType == md.TYPE_ENTITY && f.TypeID != "" && (f.Kind == md.KIND_TYPE_BELONGS_TO || f.Kind == md.KIND_TYPE_HAS_ONE) {
			ids := make([]interface{}, 0)
			for _, data := range datas {
				if fv, ok := data[f.DbName+"_id"]; ok && fv != nil && fv.(string) != "" {
					ids = append(ids, fv)
				}
			}
			if len(ids) > 0 {
				refEntity := md.GetEntity(f.TypeID)
				if refEntity != nil {
					exector := md.NewExector(refEntity.TableName)
					for _, f := range refEntity.Fields {
						if f.TypeType == md.TYPE_SIMPLE {
							exector.Select(fmt.Sprintf("$$%s as \"%s\"", f.Code, f.DbName))
						}
					}
					exector.Where(fmt.Sprintf("%s in ( ? )", f.AssociationKey), ids)
					if refDatas, err := exector.Query(s.repo); err != nil {
						glog.Error(err)
					} else if len(refDatas) > 0 {
						dataMap := make(map[string]interface{})
						for i, _ := range refDatas {
							d := refDatas[i]
							dataMap[d["id"].(string)] = d
						}
						for i, data := range datas {
							if fv, ok := data[f.DbName+"_id"]; ok && fv != nil && fv.(string) != "" {
								datas[i][f.DbName] = dataMap[fv.(string)]
							}
						}
					}
				}
			}
		}
	}
	return nil
}