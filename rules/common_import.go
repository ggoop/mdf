package rules

import (
	"fmt"
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/bootstrap/service"
	"github.com/ggoop/mdf/files"
	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/md"
	"github.com/ggoop/mdf/mof"
	"github.com/ggoop/mdf/repositories"
	"github.com/ggoop/mdf/utils"
	"strings"
)

type CommonImport struct {
	repo *repositories.MysqlRepo
}

func NewCommonImport(repo *repositories.MysqlRepo) *CommonImport {
	return &CommonImport{repo}
}
func (s *CommonImport) Exec(req *mof.ReqContext, res *mof.ResContext) error {
	log := service.NewLogSv(s.repo)
	logData := model.Log{EntID: req.EntID, UserID: req.UserID, NodeType: req.Command, NodeID: req.PageID, DataID: req.MainEntity}
	log.CreateLog(logData.Clone().SetMsg("导入开始======begin======"))

	defer func() {
		log.CreateLog(logData.Clone().SetMsg("导入结束"))
	}()

	if req.Data == nil {
		err := glog.Error("没有要导入的数据")
		log.CreateLog(logData.Clone().SetMsg(err))
		return err
	}
	data := req.Data
	mainEntity := md.GetEntity(req.MainEntity)
	if mainEntity == nil {
		err := glog.Error("找不到主实体!")
		log.CreateLog(logData.Clone().SetMsg(err))
		return err
	}
	if items, ok := data.([]files.ImportData); ok {
		return s.doMultiple(req, res, mainEntity, items)
	} else if items, ok := data.(files.ImportData); ok {
		return s.importMapData(req, res, mainEntity, items.Datas)
	} else if items, ok := data.([]map[string]interface{}); ok {
		return s.importMapData(req, res, mainEntity, items)
	}
	return nil
}
func (s *CommonImport) doMultiple(req *mof.ReqContext, res *mof.ResContext, entity *md.MDEntity, datas []files.ImportData) error {
	for _, data := range datas {
		newEntity := entity
		if data.Entity.Code != "" {
			newEntity = md.GetEntity(data.Entity.Code)
		}
		if err := s.importMapData(req, res, newEntity, data.Datas); err != nil {
			return err
		}
	}
	return nil
}
func (s *CommonImport) importMapData(req *mof.ReqContext, res *mof.ResContext, entity *md.MDEntity, datas []map[string]interface{}) error {
	log := service.NewLogSv(s.repo)
	logData := model.Log{EntID: req.EntID, UserID: req.UserID, NodeType: req.Command, NodeID: req.PageID, DataID: req.MainEntity}
	log.CreateLog(logData.Clone().SetMsg(fmt.Sprintf("接收到需要导入的数据-%s：%v条", req.MainEntity, len(datas))))

	dbDatas := make([]map[string]interface{}, 0)
	mdsv := service.NewMdSv(s.repo)

	quotedMap := make(map[string]string)

	for _, item := range datas {
		dbItem := make(map[string]interface{})
		if v, ok := item[md.STATE_FIELD]; ok && (v == md.STATE_TEMP || v == md.STATE_NORMAL) {
			continue
		}
		for kk, kv := range item {
			field := entity.GetField(kk)
			if field == nil || kv == nil {
				continue
			}
			fieldName := ""
			if field.TypeType == md.TYPE_ENTITY {
				fieldName = field.DbName + "_id"
				if obj, is := kv.(map[string]interface{}); is && obj != nil && obj["id"] != nil {
					dbItem[fieldName] = obj["id"]
					quotedMap[fieldName] = fieldName
				} else if obj, is := kv.(string); is && obj != "" {
					qreq := mof.ReqContext{MainEntity: field.TypeID, Q: obj, EntID: req.EntID, UserID: req.UserID, Data: item}
					if obj, err := mdsv.TakeDataByQ(qreq); err != nil {
						log.CreateLog(logData.Clone().SetMsg(fmt.Sprintf("数据[%s]=[%s],查询失败：%v", qreq.MainEntity, qreq.Q, err.Error())))
					} else if len(obj) > 0 && obj["id"] != nil {
						dbItem[fieldName] = obj["id"]
						quotedMap[fieldName] = fieldName
					} else if len(obj) > 0 && obj["ID"] != nil {
						dbItem[fieldName] = obj["ID"]
						quotedMap[fieldName] = fieldName
					} else {
						log.CreateLog(logData.Clone().SetMsg(fmt.Sprintf("关联对象[%s],找不到[%s]对应数据!", qreq.MainEntity, qreq.Q)))
					}
				}
			} else if field.TypeType == md.TYPE_ENUM {
				fieldName = field.DbName + "_id"
				if obj, is := kv.(map[string]interface{}); is && obj != nil && obj["id"] != nil {
					dbItem[fieldName] = obj["id"]
					quotedMap[fieldName] = fieldName
				} else if obj, is := kv.(string); is && obj != "" {
					if vv := md.GetEnum(field.Limit, obj); vv != nil {
						dbItem[fieldName] = vv.ID
						quotedMap[fieldName] = fieldName
					} else {
						log.CreateLog(logData.Clone().SetMsg(fmt.Sprintf("关联枚举[%s],找不到[%s]对应数据!", field.Limit, obj)))
					}
				}
			} else if field.TypeType == md.TYPE_SIMPLE {
				fieldName = field.DbName
				if field.TypeID == md.FIELD_TYPE_BOOL {
					dbItem[fieldName] = files.GetMapSBoolValue(kk, item)
					quotedMap[fieldName] = fieldName
				} else if field.TypeID == md.FIELD_TYPE_DATETIME || field.TypeID == md.FIELD_TYPE_DATE {
					dbItem[fieldName] = files.GetMapTimeValue(kk, item)
					quotedMap[fieldName] = fieldName
				} else if field.TypeID == md.FIELD_TYPE_DECIMAL || field.TypeID == md.FIELD_TYPE_INT {
					dbItem[fieldName] = files.GetMapDecimalValue(kk, item)
					quotedMap[fieldName] = fieldName
				} else {
					dbItem[fieldName] = kv
					quotedMap[fieldName] = fieldName
				}
			}
		}
		if field := entity.GetField("ID"); field != nil {
			fieldName := field.DbName
			if _, ok := dbItem[fieldName]; !ok {
				dbItem[fieldName] = utils.GUID()
			}
			quotedMap[fieldName] = fieldName
		}
		if field := entity.GetField("EntID"); field != nil && field.DbName != "" {
			fieldName := field.DbName
			if _, ok := dbItem[fieldName]; !ok {
				dbItem[fieldName] = req.EntID
			}
			quotedMap[fieldName] = fieldName
		}
		if field := entity.GetField("CreatedBy"); field != nil && field.DbName != "" {
			fieldName := field.DbName
			if _, ok := dbItem[fieldName]; !ok {
				dbItem[fieldName] = req.UserID
			}
			quotedMap[fieldName] = fieldName
		}
		if field := entity.GetField("CreatedAt"); field != nil && field.DbName != "" {
			fieldName := field.DbName
			dbItem[fieldName] = utils.NewTime()
			quotedMap[fieldName] = fieldName
		}
		if field := entity.GetField("UpdatedAt"); field != nil && field.DbName != "" {
			fieldName := field.DbName
			dbItem[fieldName] = utils.NewTime()
			quotedMap[fieldName] = fieldName
		}
		if len(dbItem) > 0 {
			dbDatas = append(dbDatas, dbItem)
		}
	}
	quoted := make([]string, 0, len(quotedMap))

	for fk, _ := range quotedMap {
		quoted = append(quoted, fk)
	}

	placeholdersArr := make([]string, 0, len(quotedMap))
	valueVars := make([]interface{}, 0)
	var itemCount uint = 0
	var MaxBatchs uint = 100

	for _, data := range dbDatas {
		itemCount = itemCount + 1
		placeholders := make([]string, 0, len(quoted))
		for _, f := range quoted {
			placeholders = append(placeholders, "?")
			valueVars = append(valueVars, data[f])
		}
		placeholdersArr = append(placeholdersArr, "("+strings.Join(placeholders, ", ")+")")

		if itemCount >= MaxBatchs {
			if err := s.batchInsertSave(entity, quoted, placeholdersArr, valueVars...); err != nil {
				log.CreateLog(logData.Clone().SetMsg(fmt.Sprintf("数据库存储[%v]条记录出错了:%s!", itemCount, err.Error())))
				return err
			}
			itemCount = 0
			placeholdersArr = make([]string, 0, len(quotedMap))
			valueVars = make([]interface{}, 0)
		}
	}
	if itemCount > 0 {
		if err := s.batchInsertSave(entity, quoted, placeholdersArr, valueVars...); err != nil {
			log.CreateLog(logData.Clone().SetMsg(fmt.Sprintf("数据库存储[%v]条记录出错了:%s!", itemCount, err.Error())))
			return err
		}
	}
	return nil
}

func (s *CommonImport) batchInsertSave(entity *md.MDEntity, quoted []string, placeholders []string, valueVars ...interface{}) error {
	var sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", s.repo.Dialect().Quote(entity.TableName), strings.Join(quoted, ", "), strings.Join(placeholders, ", "))

	if err := s.repo.Exec(sql, valueVars...).Error; err != nil {
		return err
	}
	return nil
}

func (s *CommonImport) GetRule() mof.RuleRegister {
	return mof.RuleRegister{Code: "import", Domain: "common"}
}