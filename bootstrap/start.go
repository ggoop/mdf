package bootstrap

import (
	"fmt"
	"github.com/ggoop/mdf/bootstrap/actions"
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/bootstrap/routes"
	"github.com/ggoop/mdf/bootstrap/rules"
	"github.com/ggoop/mdf/bootstrap/seeds"
	"github.com/ggoop/mdf/db"
	"github.com/ggoop/mdf/framework/reg"
	"github.com/ggoop/mdf/gin"
	"github.com/ggoop/mdf/middleware/token"
	"github.com/ggoop/mdf/utils"
	"os"
)

func Start() {

	engine := gin.New()

	initContext(engine)
	runArg := ""
	if os.Args != nil && len(os.Args) > 0 {
		if len(os.Args) > 1 {
			runArg = os.Args[1]
		}
	}
	if runArg == "upgrade" || runArg == "init" || runArg == "debug" {
		model.Register()
		seeds.Register()
	}
	//动作注册
	actions.Register()
	//规则注册
	rules.Register()

	engine.Use(token.Default())

	utils.CreatePath("dist")
	engine.LoadHTMLGlob("dist/*.html")
	//注册路由
	routes.Register(engine)

	//启动注册服务
	reg.StartServer()

	engine.Run(fmt.Sprintf(":%s", utils.DefaultConfig.App.Port))
}
func initContext(engine *gin.Engine) {
	db.CreateDB(utils.DefaultConfig.Db.Database)
}
