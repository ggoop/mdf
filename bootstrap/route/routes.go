package route

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"

	"github.com/ggoop/mdf/bootstrap/controller"
	"github.com/ggoop/mdf/bootstrap/service"
	"github.com/ggoop/mdf/framework/context"
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/di"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/http/middleware"
)

// 路由服务注册
func Register() {
	if err := di.Global.Invoke(func(app *iris.Application, contextMid *middleware.Context) {
		registerSys(app, contextMid)
		registerView(app, contextMid)
		RegisterProxy(app, contextMid)
		RegisterDti(app, contextMid)
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
}

func registerView(app *iris.Application, contextMid *middleware.Context) {
	m := mvc.New(app.Party("/", contextMid.Default))
	m.Handle(new(controller.ViewController))
	app.Get("/", contextMid.Default, func(ctx iris.Context) {
		token := ctx.Values().Get(context.DefaultContextKey).(*context.Context)
		if token.VerifyUser() != nil {
			ctx.Redirect("/login")
			return
		} else {
			if u, err := service.NewUserSv(repositories.NewMysqlRepo()).GetUserBy(token.UserID()); err != nil || u == nil {
				controller.SetAuthSession(ctx, nil)
				ctx.Redirect("/login")
				return
			}
		}
		ctx.View("index.html")
	})
	app.Get("/apps/{directory:path}", contextMid.Default, func(ctx iris.Context) {
		token := ctx.Values().Get(context.DefaultContextKey).(*context.Context)
		if token.VerifyUser() != nil {
			ctx.Redirect("/login")
			return
		}
		ctx.View("index.html")
	})
	app.Get("/md/{directory:path}", contextMid.Default, func(ctx iris.Context) {
		token := ctx.Values().Get(context.DefaultContextKey).(*context.Context)
		if token.VerifyUser() != nil {
			ctx.Redirect("/login")
			return
		}
		ctx.View("index.html")
	})
	app.Get("/mobile/{directory:path}", contextMid.Default, func(ctx iris.Context) {
		ctx.View("mobile/index.html")
	})

	app.Options("{directory:path}", func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.StatusCode(204)
	})

}
