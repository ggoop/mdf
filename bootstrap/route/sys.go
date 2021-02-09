package route

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"

	"github.com/ggoop/mdf/bootstrap/controller"
	"github.com/ggoop/mdf/framework/http/middleware"
	"github.com/ggoop/mdf/framework/reg"
)

func registerSys(app *iris.Application, contextMid *middleware.Context) {
	//dev route
	{
		m := mvc.New(app.Party("/api/dev", contextMid.Default))
		m.Handle(new(controller.DevController))
	}
	//md route
	{
		m := mvc.New(app.Party("/api/md", contextMid.Default))
		m.Handle(new(controller.MdController))
	}
	//regs 注册中心
	{
		m := mvc.New(app.Party("/api/regs", contextMid.Default))
		m.Handle(new(reg.RegController))
	}
	//oauth route
	{
		m := mvc.New(app.Party("/api/oauth", contextMid.Default))
		m.Handle(new(controller.OAuthController))
	}
}
