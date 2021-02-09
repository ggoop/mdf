package bootstrap

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/kataras/iris/sessions"
	"os"
	"time"

	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/bootstrap/route"
	"github.com/ggoop/mdf/bootstrap/seed"
	"github.com/ggoop/mdf/bootstrap/service"
	"github.com/ggoop/mdf/framework/configs"
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/di"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/http/middleware"
	"github.com/ggoop/mdf/framework/http/middleware/cors"
	"github.com/ggoop/mdf/framework/http/middleware/logger"
	grecover "github.com/ggoop/mdf/framework/http/middleware/recover"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/framework/reg"
	"github.com/ggoop/mdf/framework/rules"
	"github.com/ggoop/mdf/utils"
)

func StartHttp() error {
	defer func() {
		if r := recover(); r != nil {
			glog.Error(r)
			os.Exit(0)
		}
	}()
	app := iris.New()
	// 创建容器
	di.SetGlobal(di.New())
	app.Use(cors.AllowAll())

	app.Use(grecover.New())
	app.Use(logger.New())

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.JSON(iris.Map{"msg": "Not Found : " + ctx.Path()})
	})

	// 启用session
	sessionManager := sessions.New(sessions.Config{
		Cookie:       "site_session_id",
		Expires:      48 * time.Hour,
		AllowReclaim: true,
	})
	hero.Register(sessionManager.Start)
	if err := di.Global.Provide(func() *sessions.Sessions {
		return sessionManager
	}); err != nil {
		return glog.Errorf("注册缓存服务异常:%s", err)
	}
	//注册中间件
	if err := di.Global.Provide(func() *middleware.Context {
		return &middleware.Context{Sessions: sessionManager}
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	//添加tokencode 拦截器
	app.Use(reg.NewTokenMiddleware())

	// 基础组件注册
	comRegister(app)

	runArg := ""
	if os.Args != nil && len(os.Args) > 0 {
		if len(os.Args) > 1 {
			runArg = os.Args[1]
		}
	}
	if runArg == "upgrade" || runArg == "init" || runArg == "debug" {
		dateFm := utils.NewTime()
		glog.Errorf("数据升级开始:%v \r", dateFm)
		if err := di.Global.Invoke(func(repo *repositories.MysqlRepo) {
			//数据迁移
			model.RegisterMD()
			//数据填充
			seed.Register()

		}); err != nil {
			return glog.Errorf("数据库组件没有准备好:%s", err)
		}
		dateTo := utils.NewTime()
		glog.Errorf("数据升级完成:%v,%v \r", dateTo.Sub(dateFm.Time).String(), dateTo)
	}
	if runArg == "upgrade" || runArg == "init" {
		os.Exit(0)
		return nil
	}
	//初始缓存
	initCache()

	//启动 JOB
	go initJob()

	//服务注册
	go reg.Start()

	// 启动服务
	dc := iris.DefaultConfiguration()
	if utils.PathExists("env/iris.yaml") {
		dc = iris.YAML(utils.JoinCurrentPath("env/iris.yaml"))
	}
	dc.DisableBodyConsumptionOnUnmarshal = true
	dc.FireMethodNotAllowed = true

	if err := app.Run(iris.Addr(":"+configs.Default.App.Port), iris.WithConfiguration(dc)); err != nil {
		return glog.Errorf("Run service error:%s\n", err.Error())
	}
	return nil
}

func comRegister(app *iris.Application) {
	// 注册app
	if err := di.Global.Provide(func() *iris.Application {
		return app
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	mysqkRepo := repositories.NewMysqlRepo()
	mysqkRepo.SetLogger(glog.GetLogger(""))
	mysqkRepo.DB.DB().SetConnMaxLifetime(0)
	repositories.SetDefault(mysqkRepo)

	// 控制层，注入数据访问
	hero.Register(mysqkRepo)

	if err := di.Global.Provide(func() *repositories.MysqlRepo {
		return mysqkRepo
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}

	// 路由注册
	route.Register()

	//服务注册
	service.Register()

	//规则注册
	rules.Register()
}

// 枚举缓存
func initCache() {
	if err := di.Global.Invoke(func(s1 *md.EntitySv) {
		s1.InitCache()
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
}

//启动 JOB
func initJob() {
	time.Sleep(10 * time.Second)
	if err := di.Global.Invoke(func(s1 *service.CronSv) {
		s1.Start()
	}); err != nil {
		glog.Errorf("di Provide SysCronSv error:%s", err)
	}
}
