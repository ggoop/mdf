package routes

import (
	"github.com/ggoop/mdf/gin"
	"net/http"
)

func apiMd(group *gin.RouterGroup) {
	r := group.Group("md")
	{
		r.GET("test", func(c *gin.Context) {
			c.String(http.StatusOK, "dddd")
		})
	}
}
