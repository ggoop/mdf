package bootstrap

import (
	"fmt"
	"github.com/ggoop/mdf/bootstrap/routes"
	"github.com/ggoop/mdf/framework/reg"
	"github.com/ggoop/mdf/gin"
	"github.com/ggoop/mdf/middleware/token"
	"io"
	"os"

	"github.com/ggoop/mdf/utils"
)

func Start() {
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
	engine := gin.New()

	engine.Use(token.Default())

	//注册路由
	routes.Register(engine)

	//启动注册服务
	reg.StartServer()

	engine.Run(fmt.Sprintf(":%s", utils.DefaultConfig.App.Port))
}
