package routes

import "github.com/ggoop/mdf/gin"

func routeApi(engine *gin.Engine) {
	group := engine.Group("api")
	apiMd(group.Group("md"))
	apiAuth(group.Group("auth"))
}
