package route

import (
	"github.com/kataras/iris"
	"net/http/httputil"
	"net/url"

	"github.com/ggoop/mdf/framework/http/middleware"
	"github.com/ggoop/mdf/framework/http/results"
	"github.com/ggoop/mdf/framework/reg"
)

func RegisterProxy(app *iris.Application, contextMid *middleware.Context) {
	/**
	  /proxy/endpoint/{endpointName}
	  执行代理服务接口,
	*/
	app.Any("/proxy/{node:string}/{uri:path}", func(ctx iris.Context) {
		uri := ctx.Params().Get("uri")
		//可以通过Header统一设置node
		nodeName := ctx.GetHeader("Node")
		if nodeName != "" && len(nodeName) > 2 {
			uri = ctx.Params().Get("node") + "/" + ctx.Params().Get("uri")
		} else {
			nodeName = ctx.Params().Get("node")
		}
		addr, err := reg.GetServerAddr(nodeName)
		if err != nil {
			results.ToError(err).Dispatch(ctx)
			return
		}
		if addr == "" {
			results.ToError("找不到服务地址！").Dispatch(ctx)
			return
		}
		remote, err := url.Parse(addr)
		if err != nil {
			results.ToError(err).Dispatch(ctx)
			return
		}
		r := ctx.Request()
		r.URL.Path = uri
		proxy := httputil.NewSingleHostReverseProxy(remote)
		w := ctx.ResponseWriter()
		proxy.ServeHTTP(w, r)
	})
}
