package routes

import "github.com/ggoop/mdf/gin"

func Register(engine *gin.Engine) {
	routeView(engine)
	routeApi(engine)
	routeProxy(engine)
	routeDti(engine)
}
