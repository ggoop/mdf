package rules

import (
	"fmt"
	"github.com/ggoop/mdf/db"
	"strings"

	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type commonSave struct {
}

func newCommonSave() *commonSave {
	return &commonSave{}
}
func (s *commonSave) Exec(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext) {
	reqData := make(map[string]interface{})
	if data, ok := req.Data.(map[string]interface{}); !ok {
		res.SetError("data数据格式不正确")
		return
	} else {
		reqData = data
	}
	if reqData == nil {
		res.SetError("没有要保存的数据！")
		return
	}
	//查找实体信息
	entity := md.MDSv().GetEntity(req.Entity)
	if entity == nil {
		res.SetError("找不到实体！")
		return
	}
	state := ""
	if s, ok := reqData[utils.STATE_FIELD]; ok && s != nil {
		state = s.(string)
	}
	if state == utils.STATE_UPDATED && req.ID == "" {

	}
	//如果有ID，则为修改保存
	if req.ID != "" {
		oldData := make(map[string]interface{})
		exector := md.NewOQL().From(entity.TableName)
		for _, f := range entity.Fields {
			if f.TypeType == utils.TYPE_SIMPLE {
				exector.Select(f.Code)
			}
		}
		exector.Where("id=?", req.ID)
		datas := make([]map[string]interface{}, 0)
		if err := exector.Find(&datas).Error(); err != nil {
			res.SetError(err)
			return
		} else if len(datas) > 0 {
			oldData = datas[0]
		}
		if len(oldData) == 0 {
			res.SetError("找不到要修改的数据！")
			return
		}
		s.doActionUpdate(token, req, res, entity, reqData, oldData)
	} else {
		s.doActionCreate(token, req, res, entity, reqData)
	}
}

func (s *commonSave) fillEntityDefaultValue(entity *md.MDEntity, data map[string]interface{}) map[string]interface{} {
	for _, field := range entity.Fields {
		//如果字段设置了默认值，且没有传入字段值时，取默认值
		if field.DbName != "" && field.DefaultValue != "" {
			if dv, ok := data[field.DbName]; !ok || dv == nil {
				data[field.DbName] = field.CompileValue(field.DefaultValue)
			}
		}
	}
	return data
}
func (s *commonSave) dataToEntityData(entity *md.MDEntity, data map[string]interface{}) map[string]interface{} {
	changeData := make(map[string]interface{})
	for di, dv := range data {
		field := entity.GetField(di)
		if field == nil || field.TypeType != utils.TYPE_SIMPLE {
			continue
		}
		dbFieldName := field.DbName
		if obj, is := dv.(map[string]interface{}); is && obj != nil && obj["code"] != nil {
			changeData[dbFieldName] = obj["code"]
		} else {
			changeData[dbFieldName] = field.CompileValue(dv)
		}
	}
	// 处理枚举和实体
	for di, dv := range data {
		field := entity.GetField(di)
		if field == nil {
			continue
		}
		if field.TypeType == utils.TYPE_ENTITY || field.TypeType == utils.TYPE_ENUM {
			dbField := entity.GetField(field.Code + "ID")
			if dbField == nil || dbField.TypeType != utils.TYPE_SIMPLE || dbField.DbName == "" {
				continue
			}
			if obj, is := dv.(map[string]interface{}); is && obj != nil && obj["id"] != nil {
				changeData[dbField.DbName] = obj["id"]
			} else {
				changeData[dbField.DbName] = ""
			}
			continue
		}
	}
	return changeData
}
func (s *commonSave) doActionCreate(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext, entity *md.MDEntity, reqData map[string]interface{}) error {
	reqData["id"] = utils.GUID()
	if sysField := entity.GetField("EntID"); sysField != nil && token.EntID() != "" {
		reqData[sysField.DbName] = token.EntID()
	}
	if sysField := entity.GetField("CreatedBy"); sysField != nil && token.UserID() != "" {
		reqData[sysField.DbName] = token.UserID
	}
	fieldCreated := entity.GetField("CreatedAt")
	if fieldCreated != nil && fieldCreated.DbName != "" {
		reqData[fieldCreated.DbName] = utils.TimeNow()
	}
	fieldUpdatedAt := entity.GetField("UpdatedAt")
	if fieldUpdatedAt != nil && fieldUpdatedAt.DbName != "" {
		reqData[fieldUpdatedAt.DbName] = utils.TimeNow()
	}
	//取传入的值
	changeData := s.dataToEntityData(entity, reqData)
	//配置默认值
	changeData = s.fillEntityDefaultValue(entity, changeData)
	if len(changeData) == 0 {
		return glog.Error("没有要保存的数据！")
	}
	//数据校验
	if err := s.dataCheck(req, res, entity, changeData); err != nil {
		return err
	}
	//树规则
	isTree := false
	fieldParent := entity.GetField("ParentID")
	if fieldParent != nil && fieldParent.DbName != "" {
		isTree = true
	}
	isLeafField := entity.GetField("IsLeaf")
	if isTree {
		if changeData["id"] != "" && changeData["id"] == changeData[fieldParent.DbName] {
			return glog.Error("树结构，父节点不能等于当前节点!")
		}
		if isLeafField != nil && isLeafField.DbName != "" {
			changeData[isLeafField.DbName] = 1
		}
	}

	//开始保存数据
	fields := make([]string, 0)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	for f, v := range changeData {
		if vv, is := v.(utils.SBool); is && !vv.Valid() {
			continue
		}
		if vv, is := v.(utils.SJson); is && !vv.Valid() {
			continue
		}
		fields = append(fields, db.Default().Dialect().Quote(f))
		placeholders = append(placeholders, "?")
		values = append(values, v)
	}
	sql := fmt.Sprintf("insert into %s (%s) values (%s)", db.Default().Dialect().Quote(entity.TableName), strings.Join(fields, ","), strings.Join(placeholders, ","))
	if err := db.Default().Table(entity.TableName).Exec(sql, values...).Error; err != nil {
		return err
	}
	//处理树节点标识
	if isTree {
		//更新父节点标识
		if parentID, ok := changeData[fieldParent.DbName].(string); ok && parentID != "" && isLeafField != nil {
			updates := make(map[string]interface{})
			updates[isLeafField.DbName] = 0
			if fieldUpdatedAt != nil && fieldUpdatedAt.DbName != "" {
				updates[fieldUpdatedAt.DbName] = utils.TimeNow()
			}
			if err := db.Default().Table(entity.TableName).Where("id=?", parentID).Updates(updates).Error; err != nil {
				return err
			}
		}
	}
	//保存关联实体
	if err := s.saveRelationData(token, req, res, entity, reqData); err != nil {
		return err
	}
	res.Set("data", changeData)
	return nil
}
func (s *commonSave) doActionUpdate(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext, entity *md.MDEntity, reqData map[string]interface{}, oldData map[string]interface{}) error {
	fieldUpdatedAt := entity.GetField("UpdatedAt")
	if fieldUpdatedAt != nil && fieldUpdatedAt.DbName != "" {
		reqData[fieldUpdatedAt.DbName] = utils.TimeNow()
	}
	if sysField := entity.GetField("ID"); sysField != nil && req.ID != "" {
		reqData[sysField.DbName] = req.ID
	}
	data := s.dataToEntityData(entity, reqData)

	changeData := make(map[string]interface{})
	for nk, nv := range data {
		if nk == "id" {
			continue
		}
		isChanged := true
		oldValue := oldData[nk]
		field := entity.GetField(nk)
		if field == nil {
			continue
		}
		fieldType := strings.ToLower(field.TypeID)
		//布尔类型判断
		if fieldType == "bool" || fieldType == "boolean" {
			newV := utils.ToSBool(nv)
			oldV := utils.ToSBool(oldValue)
			if !newV.Valid() || newV.Equal(oldV) {
				isChanged = false
			}
		} else {
			if nv == oldValue {
				isChanged = false
			}
		}
		if isChanged {
			changeData[nk] = nv
		}
	}
	//树规则
	isTree := false
	fieldParent := entity.GetField("ParentID")
	if fieldParent != nil && fieldParent.DbName != "" {
		isTree = true
	}
	isLeafField := entity.GetField("IsLeaf")
	if isTree {
		if req.ID == changeData[fieldParent.DbName] {
			return glog.Error("树结构，父节点不能等于当前节点!")
		}
	}
	//数据校验
	if err := s.dataCheck(req, res, entity, changeData); err != nil {
		return err
	}
	if len(changeData) > 0 {
		//开始保存数据
		if err := db.Default().Table(entity.TableName).Where("id=?", req.ID).Updates(changeData).Error; err != nil {
			return err
		}
	}
	//保存关联实体
	if err := s.saveRelationData(token, req, res, entity, reqData); err != nil {
		return err
	}

	if len(changeData) > 0 {
		if isTree && isLeafField != nil {
			oldParentID := ""
			if mv, ok := oldData[fieldParent.DbName]; ok {
				oldParentID = mv.(string)
			}
			//如果修改了父节点
			if newParentID, ok := changeData[fieldParent.DbName]; ok {
				if newParentID != "" { //如果设置父节点不为空，则更新父节点为非叶子节点
					updates := make(map[string]interface{})
					updates[isLeafField.DbName] = utils.SBool_False
					if fieldUpdatedAt != nil && fieldUpdatedAt.DbName != "" {
						updates[fieldUpdatedAt.DbName] = utils.TimeNow()
					}
					if err := db.Default().Table(entity.TableName).Where("id=?", newParentID).Updates(updates).Error; err != nil {
						return err
					}
				}
				if oldParentID != "" { //如果设置父节点为空，则更新父节点叶子节点状态
					count := 0
					updates := make(map[string]interface{})
					if fieldUpdatedAt != nil && fieldUpdatedAt.DbName != "" {
						updates[fieldUpdatedAt.DbName] = utils.TimeNow()
					}
					if db.Default().Table(entity.TableName).Where(fmt.Sprintf("%s=?", fieldParent.DbName), oldParentID).Count(&count); count == 0 {
						updates[isLeafField.DbName] = 1
					} else {
						updates[isLeafField.DbName] = 0
					}
					if err := db.Default().Table(entity.TableName).Where("id=?", oldParentID).Updates(updates).Error; err != nil {
						return err
					}
				}
			}
		}
		res.Set("data", changeData)
	}
	return nil
}

func (s *commonSave) saveRelationData(token *utils.TokenContext, req *md.ReqContext, res *md.ResContext, entity *md.MDEntity, reqData map[string]interface{}) error {
	for _, nv := range entity.Fields {
		if nv.Kind == md.KIND_TYPE_HAS_MANT {
			if do, ok := reqData[nv.DbName].([]interface{}); ok && len(do) > 0 {
				for _, dr := range do {
					if ds, ok := dr.(map[string]interface{}); ok {
						state := ""

						if s, ok := ds[utils.STATE_FIELD]; ok && s != nil {
							state = s.(string)
						}
						if state == "" {
							glog.Error("实体对应状态为空，跳过更新！", glog.String("state", state))
							continue
						}

						newReq := md.NewReqContext()
						newReq.OwnerType = req.OwnerType
						newReq.OwnerID = req.OwnerID
						refEntity := md.MDSv().GetEntity(nv.TypeID)
						if f := refEntity.GetField(nv.ForeignKey); f != nil {
							ds[f.DbName] = reqData["id"]
						}
						ruleID := ""
						if state == utils.STATE_CREATED || state == utils.STATE_UPDATED {
							ruleID = "save"
						}
						if state == utils.STATE_DELETED {
							ruleID = "delete"
						}
						if ruleID == "" {
							glog.Error("该状态找不到对应规则", glog.String("state", state))
							continue
						}
						if state == utils.STATE_UPDATED || state == utils.STATE_DELETED {
							if id, ok := ds["id"].(string); ok && id != "" {
								newReq.ID = id
							}
						}
						newReq.Data = ds
						newReq.Entity = refEntity.ID
						newReq.Rule = ruleID

						if rtn := md.ActionSv().DoAction(token, &newReq); rtn.Error() != nil {
							return rtn.Error()
						}
					}
				}

			}
		}
	}
	return nil
}
func (s *commonSave) dataCheck(req *md.ReqContext, res *md.ResContext, entity *md.MDEntity, data map[string]interface{}) error {
	return nil
}
func (s *commonSave) Register() md.RuleRegister {
	return md.RuleRegister{Code: "save", OwnerType: md.RuleType_Widget, OwnerID: "common"}
}
