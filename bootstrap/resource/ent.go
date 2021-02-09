package resources

import (
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/framework/http/results"
	"github.com/ggoop/mdf/utils"
)

type EntRes struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	Memo      string `json:"memo"`
}

func (c EntRes) ToItem(obj *model.Ent) EntRes {
	item := EntRes{
		ID: obj.ID, CreatedAt: obj.CreatedAt.Format("2006-01-02 15:04:05"),
		Name: obj.Name, Memo: obj.Memo,
	}
	return item
}
func (c EntRes) ToItems(obj *[]model.Ent) []EntRes {
	lists := make([]EntRes, 0)
	if len(*obj) > 0 {
		for _, v := range *obj {
			lists = append(lists, c.ToItem(&v))
		}
	}
	return lists
}
func (c EntRes) ToResource(obj interface{}) *results.Map {
	rtn := results.Map{"data": nil}
	if utils.IsNil(obj) {
		return &rtn
	}
	// 单个实体
	if datas, is := obj.(*model.Ent); is {
		return &results.Map{"data": c.ToItem(datas)}
	}
	//多个
	if datas, is := obj.(*[]model.Ent); is {
		return &results.Map{"data": c.ToItems(datas)}
	}

	//分页
	if datas, is := obj.(*utils.Pager); is {
		rtn["pager"] = ToPagerItem(datas)
		if datas, is := datas.Value.(*[]model.Ent); is {
			rtn["data"] = c.ToItems(datas)
		}
	}
	return &rtn
}
