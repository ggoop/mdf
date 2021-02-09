package controller

import (
	"github.com/ggoop/mdf/bootstrap/service"
	"github.com/ggoop/mdf/http/results"
	"github.com/ggoop/mdf/repositories"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

type ViewController struct {
	Ctx       iris.Context
	Repo      *repositories.MysqlRepo
	Session   *sessions.Session
	UserSv    *service.UserSv
	EntSv     *service.EntSv
	ProfileSv *service.ProfileSv
}

func (c *ViewController) AnyTest() results.Result {
	c.Ctx.Header("Access-Control-Allow-Origin", "*")
	c.Ctx.Header("Access-Control-Allow-Headers", "*")
	c.Ctx.Header("Access-Control-Allow-Credentials", "true")
	return results.ToJson(results.Map{"data": true})
}

/**
登录
*/
func (c *ViewController) GetLogin() {
	c.Ctx.View("index.html")
}

/**
登出
*/
func (c *ViewController) GetInit() {
	c.Ctx.View("index.html")
}

/**
注册
*/
func (c *ViewController) GetRegister() {
	c.Ctx.View("index.html")
}
