package routes

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/gin"
	"github.com/ggoop/mdf/middleware/token"
)

func apiMd(r *gin.RouterGroup) {
	r.POST("import", func(c *gin.Context) {
		md.ActionSv().DoAction(
			token.Get(c),
			md.NewReqContext().Bind(c).Adjust(func(req *md.ReqContext) {
				req.Action = "import"
			}),
		).Bind(c)
	})
	r.POST("save", func(c *gin.Context) {
		md.ActionSv().DoAction(
			token.Get(c),
			md.NewReqContext().Bind(c).Adjust(func(req *md.ReqContext) {
				req.Action = "save"
			}),
		).Bind(c)
	})
	r.POST("delete", func(c *gin.Context) {
		md.ActionSv().DoAction(
			token.Get(c),
			md.NewReqContext().Bind(c).Adjust(func(req *md.ReqContext) {
				req.Action = "delete"
			}),
		).Bind(c)
	})
	r.POST("enable", func(c *gin.Context) {
		md.ActionSv().DoAction(
			token.Get(c),
			md.NewReqContext().Bind(c).Adjust(func(req *md.ReqContext) {
				req.Action = "enable"
			}),
		).Bind(c)
	})
	r.POST("disable", func(c *gin.Context) {
		md.ActionSv().DoAction(
			token.Get(c),
			md.NewReqContext().Bind(c).Adjust(func(req *md.ReqContext) {
				req.Action = "disable"
			}),
		).Bind(c)
	})
	r.POST("query", func(c *gin.Context) {
		md.ActionSv().DoAction(
			token.Get(c),
			md.NewReqContext().Bind(c).Adjust(func(req *md.ReqContext) {
				req.Action = "query"
			}),
		).Bind(c)
	})

}
