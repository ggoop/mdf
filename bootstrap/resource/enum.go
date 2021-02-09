package resources

import (
	"github.com/ggoop/mdf/http/results"
	"github.com/ggoop/mdf/md"
	"github.com/ggoop/mdf/utils"
)

type EnumRes struct {
	Type         string `json:"type"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	IsEnumObject bool   `json:"_isEnumObject"`
}

func (c *EnumRes) ToItem(obj *md.MDEnum) (res EnumRes) {
	if obj == nil {
		return res
	}
	return EnumRes{Type: obj.EntityID, ID: obj.ID, Name: obj.Name, IsEnumObject: true}
}
func (c *EnumRes) ToItems(obj *[]md.MDEnum) []EnumRes {
	lists := make([]EnumRes, 0)
	if len(*obj) > 0 {
		for _, v := range *obj {
			lists = append(lists, c.ToItem(&v))
		}
	}
	return lists
}
func (c *EnumRes) ToResource(obj interface{}) *results.Map {
	rtn := results.Map{"data": nil}
	if utils.IsNil(obj) {
		return &rtn
	}
	// 单个实体
	if datas, is := obj.(*md.MDEnum); is {
		return &results.Map{"data": c.ToItem(datas)}
	}
	if datas, is := obj.(md.MDEnum); is {
		return &results.Map{"data": c.ToItem(&datas)}
	}
	//多个
	if datas, is := obj.(*[]md.MDEnum); is {
		return &results.Map{"data": c.ToItems(datas)}
	}
	if datas, is := obj.([]md.MDEnum); is {
		return &results.Map{"data": c.ToItems(&datas)}
	}

	//分页
	if datas, is := obj.(*utils.Pager); is {
		rtn["pager"] = ToPagerItem(datas)
		if datas, is := datas.Value.(*[]md.MDEnum); is {
			rtn["data"] = c.ToItems(datas)
		}
	}
	return &rtn
}
