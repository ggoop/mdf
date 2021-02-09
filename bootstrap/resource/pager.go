package resources

import (
	"github.com/ggoop/mdf/http/results"
	"github.com/ggoop/mdf/utils"
)

type PagerRes struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	LastPage int `json:"last_page"`
	Items    int `json:"items"` //当前条数
	Total    int `json:"total"` //总记录数
}

func (c *PagerRes) ToResource(obj *utils.Pager) *results.Map {
	rtn := results.Map{"data": nil}
	if utils.IsNil(obj) {
		return &rtn
	}
	rtn["pager"] = ToPagerItem(obj)
	rtn["data"] = obj.Value
	return &rtn
}

func ToPagerItem(obj *utils.Pager) *PagerRes {
	return &PagerRes{
		Page:     obj.Page,
		PageSize: obj.PageSize,
		Total:    obj.Total,
		Items:    obj.Items,
		LastPage: obj.LastPage,
	}
}
