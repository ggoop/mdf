package utils

type Pager struct {
	Value    interface{}
	Error    error
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	LastPage int `json:"last_page"`
	Items    int `json:"items"` //当前条数
	Total    int `json:"total"` //总记录数
}
type PagerRes struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	LastPage int `json:"last_page"`
	Items    int `json:"items"` //当前条数
	Total    int `json:"total"` //总记录数
}

func (c *PagerRes) ToResource(obj *Pager) *Map {
	rtn := Map{"data": nil}
	if obj == nil {
		return &rtn
	}
	rtn["pager"] = ToPagerItem(obj)
	rtn["data"] = obj.Value
	return &rtn
}

func ToPagerRes(obj *Pager) *Map {
	return ToPagerItem(obj).ToResource(obj)
}
func ToPagerItem(obj *Pager) *PagerRes {
	return &PagerRes{
		Page:     obj.Page,
		PageSize: obj.PageSize,
		Total:    obj.Total,
		Items:    obj.Items,
		LastPage: obj.LastPage,
	}
}
