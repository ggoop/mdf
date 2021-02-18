package routes

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/framework/reg"
	"github.com/ggoop/mdf/gin"
	"net/http/httputil"
	"net/url"
)

func routeProxy(engine *gin.Engine) {
	handle := func(c *gin.Context) {
		uri := c.Param("uri")
		//可以通过Header统一设置node
		nodeName := c.GetHeader("Node")
		if nodeName != "" && len(nodeName) > 2 {
			uri = c.Param("node") + "/" + c.Param("uri")
		} else {
			nodeName = c.Param("node")
		}
		addr, err := reg.FindServerByCode(nodeName)
		if err != nil {
			md.NewResContext().SetError(err).Bind(c)
			return
		}
		if addr == nil || addr.Address == "" {
			md.NewResContext().SetError("找不到服务地址").Bind(c)
			return
		}
		remote, err := url.Parse(addr.Address)
		if err != nil {
			md.NewResContext().SetError(err).Bind(c)
			return
		}
		r := c.Request
		r.URL.Path = uri
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.ServeHTTP(c.Writer, r)
	}
	engine.GET("/proxy/:node/*uri", handle)
	engine.POST("/proxy/:node/*uri", handle)
}
