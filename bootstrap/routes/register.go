package routes

import "github.com/ggoop/mdf/gin"

func Register(engine *gin.Engine) {
	mdRoute := engine.Group("api")
	apiMd(mdRoute)
}
