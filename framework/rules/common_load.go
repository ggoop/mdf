package rules

import (
	"fmt"

	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/framework/mof"
	"github.com/ggoop/mdf/framework/query"
)

type CommonLoad struct {
	repo *repositories.MysqlRepo
}

func NewCommonLoad(repo *repositories.MysqlRepo) *CommonLoad {
	return &CommonLoad{repo}
}

func (s CommonLoad) GetRule() mof.RuleRegister {
	return mof.RuleRegister{Code: "load", Owner: "common"}
}

func (s CommonLoad) Exec(req *mof.ReqContext, res *mof.ResContext) error {
	if req.ID == "" {
		return glog.Error("缺少 ID 参数！")
	}
	//查找实体信息
	entity := md.GetEntity(req.Entity)
	if entity == nil {
		return glog.Error("找不到实体！")
	}
	exector := query.NewExector(entity.TableName)
	for _, f := range entity.Fields {
		if f.TypeType == md.TYPE_SIMPLE {
			exector.Select(fmt.Sprintf("$$%s as \"%s\"", f.Code, f.DbName))
			if f.TypeID != "" && f.DbName != "" {
				exector.SetFieldDataType(f.DbName, f.TypeID)
			}
		}
	}
	exector.Where("id=?", req.ID)
	if datas, err := exector.Query(s.repo); err != nil {
		return err
	} else if len(datas) > 0 {
		data := datas[0]
		s.loadEnums(data, entity)
		s.loadEntities(data, entity)
		res.SetData("data", data)
	}
	return nil
}
func (s CommonLoad) loadEnums(data map[string]interface{}, entity *md.MDEntity) error {
	for _, f := range entity.Fields {
		if f.TypeType == md.TYPE_ENUM {
			if fv, ok := data[f.DbName+"_id"]; ok && fv != nil && fv.(string) != "" {
				data[f.DbName] = md.GetEnum(f.Limit, fv.(string))
			}
		}
	}
	return nil
}
func (s CommonLoad) loadEntities(data map[string]interface{}, entity *md.MDEntity) error {
	for _, f := range entity.Fields {
		if f.TypeType == md.TYPE_ENTITY && f.TypeID != "" && (f.Kind == md.KIND_TYPE_BELONGS_TO || f.Kind == md.KIND_TYPE_HAS_ONE) {
			if fv, ok := data[f.DbName+"_id"]; ok && fv != nil && fv.(string) != "" {
				refEntity := md.GetEntity(f.TypeID)
				if refEntity != nil {
					exector := query.NewExector(refEntity.TableName)
					for _, f := range refEntity.Fields {
						if f.TypeType == md.TYPE_SIMPLE {
							exector.Select(fmt.Sprintf("$$%s as \"%s\"", f.Code, f.DbName))
						}
					}
					exector.Where(fmt.Sprintf("%s = ?", f.AssociationKey), fv)
					if refDatas, err := exector.Query(s.repo); err != nil {
						glog.Error(err)
					} else if len(refDatas) > 0 {
						data[f.DbName] = refDatas[0]
					}
				}
			}
		}
	}
	return nil
}
