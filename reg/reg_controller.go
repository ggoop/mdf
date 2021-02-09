package reg

import (
	"github.com/ggoop/mdf/http/results"
	"github.com/kataras/iris"
)

type RegController struct {
	Ctx iris.Context
}

func (c *RegController) PostRegister() results.Result {
	item := RegObject{}
	if err := c.Ctx.ReadJSON(&item); err != nil {
		return results.ToError(err)
	}
	reg_store.Add(item)
	return results.ToJson(results.Map{"data": true})
}
func (c *RegController) GetBy(code string) results.Result {
	return results.ToJson(results.Map{"data": reg_store.Get(RegObject{Code: code})})
}
func (c *RegController) Get() results.Result {
	return results.ToJson(results.Map{"data": reg_store.GetAll()})
}
