package rules

import (
	"fmt"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type commonQuery struct {
}

func newCommonQuery() *commonQuery {
	return &commonQuery{}
}
func (s *commonQuery) Register() md.RuleRegister {
	return md.RuleRegister{Code: "query", OwnerType: md.RuleType_Widget, OwnerID: "common"}
}

func (s commonQuery) Exec(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext) {
	if req.Entity == "" {
		res.SetError("缺少 MainEntity 参数！")
		return
	}
	//查找实体信息
	entity := md.MDSv().GetEntity(req.Entity)
	if entity == nil {
		res.SetError("找不到实体！")
		return
	}
	exector := md.NewOQL().From(entity.TableName)
	for _, f := range entity.Fields {
		if f.TypeType == utils.TYPE_SIMPLE {
			exector.Select(fmt.Sprintf("$$%s as \"%s\"", f.Code, f.DbName))
		}
	}
	if req.ID != "" {
		exector.Where("id=?", req.ID)
	}
	if sysField := entity.GetField("EntID"); sysField != nil && token.EntID() != "" {
		exector.Where(sysField.Code+" = ?", token.EntID())
	}
	count := 0
	if err := exector.Count(&count).Error(); err != nil {
		res.SetError(err)
		return
	}
	datas := make([]map[string]interface{}, 0)
	if err := exector.Find(datas).Error(); err != nil {
		res.SetError(err)
		return
	} else if len(datas) > 0 {
		s.loadEnums(datas, entity)
		s.loadEntities(datas, entity)
		res.Set("data", datas)
		res.Set("pager", utils.Pager{Total: count, PageSize: req.PageSize, Page: req.Page})
	}
}
func (s commonQuery) loadEnums(datas []map[string]interface{}, entity *md.MDEntity) error {
	for _, f := range entity.Fields {
		if f.TypeType == utils.TYPE_ENUM {
			for ri, data := range datas {
				if fv, ok := data[f.DbName+"_id"]; ok && fv != nil && fv.(string) != "" {
					datas[ri][f.DbName] = md.MDSv().GetEnum(f.Limit, fv.(string))
				}
			}
		}
	}
	return nil
}
func (s commonQuery) loadEntities(datas []map[string]interface{}, entity *md.MDEntity) error {
	for _, f := range entity.Fields {
		if f.TypeType == utils.TYPE_ENTITY && f.TypeID != "" && (f.Kind == md.KIND_TYPE_BELONGS_TO || f.Kind == md.KIND_TYPE_HAS_ONE) {
			ids := make([]interface{}, 0)
			for _, data := range datas {
				if fv, ok := data[f.DbName+"_id"]; ok && fv != nil && fv.(string) != "" {
					ids = append(ids, fv)
				}
			}
			if len(ids) > 0 {
				refEntity := md.MDSv().GetEntity(f.TypeID)
				if refEntity != nil {
					exector := md.NewOQL().From(refEntity.TableName)
					for _, f := range refEntity.Fields {
						if f.TypeType == utils.TYPE_SIMPLE {
							exector.Select(fmt.Sprintf("$$%s as \"%s\"", f.Code, f.DbName))
						}
					}
					exector.Where(fmt.Sprintf("%s in ( ? )", f.AssociationKey), ids)
					refDatas := make([]map[string]interface{}, 0)
					if err := exector.Find(&refDatas).Error(); err != nil {
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
