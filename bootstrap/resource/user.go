package resources

import (
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/framework/http/results"
	"github.com/ggoop/mdf/utils"
)

type UserRes struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
	Memo      string `json:"memo"`
}

func (c UserRes) ToItem(obj *model.User) UserRes {
	return UserRes{
		ID:     obj.ID,
		Mobile: obj.Mobile, Email: obj.Email,
		Name: obj.Name, AvatarUrl: obj.AvatarUrl, Memo: obj.Memo,
	}
}
func (c UserRes) ToItems(obj *[]model.User) []UserRes {
	lists := make([]UserRes, 0)
	if len(*obj) > 0 {
		for _, v := range *obj {
			lists = append(lists, c.ToItem(&v))
		}
	}
	return lists
}
func (c UserRes) ToResource(obj interface{}) *results.Map {
	rtn := results.Map{"data": nil}
	if utils.IsNil(obj) {
		return &rtn
	}
	// 单个实体
	if datas, is := obj.(*model.User); is {
		return &results.Map{"data": c.ToItem(datas)}
	}
	//多个
	if datas, is := obj.(*[]model.User); is {
		return &results.Map{"data": c.ToItems(datas)}
	}

	//分页
	if datas, is := obj.(*utils.Pager); is {
		rtn["pager"] = ToPagerItem(datas)
		if datas, is := datas.Value.(*[]model.User); is {
			rtn["data"] = c.ToItems(datas)
		}
	}
	return &rtn
}
