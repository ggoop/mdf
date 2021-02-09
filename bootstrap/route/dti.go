package route

import (
	"github.com/kataras/iris"

	"github.com/ggoop/mdf/bootstrap/dti"
	"github.com/ggoop/mdf/framework/http/middleware"
)

func RegisterDti(app *iris.Application, contextMid *middleware.Context) {
	/**
	/dti/dbName/serverName
	如果 group=localhost 则，直接使用本地数据库
	执行代理服务接口,
	*/
	app.Any("/dti/{group:string}/{name:path}", func(ctx iris.Context) {
		hand := &dti.DtiHandProc{Group: ctx.Params().Get("group"), Ctx: ctx, Path: ctx.Params().Get("name")}
		if hand.Group != "" {
			hand.Do()
		}
	})
	app.Any("/dti/{name:path}", func(ctx iris.Context) {
		hand := &dti.DtiHandProc{Group: "dti", Ctx: ctx, Path: ctx.Params().Get("name")}
		if hand.Group != "" {
			hand.Do()
		}
	})
}
