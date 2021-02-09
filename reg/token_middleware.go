package reg

import (
	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/utils"
	"strings"

	"github.com/ggoop/mdf/context"

	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

type TokenMiddlewareConfig struct {
	Sessions *sessions.Sessions
}

type tokenMiddlewareHandle struct {
	Sessions *sessions.Sessions
}

func NewTokenMiddleware(cfg ...TokenMiddlewareConfig) func(ctx iris.Context) {
	item := tokenMiddlewareHandle{}
	if len(cfg) > 0 {
		item.Sessions = cfg[0].Sessions
	}
	return item.Handle
}
func (m *tokenMiddlewareHandle) Handle(ctx iris.Context) {
	var (
		uc  *context.Context
		err error
	)
	if ctxValue := ctx.Values().Get(context.DefaultContextKey); ctxValue != nil {
		if old := ctxValue.(*context.Context); old != nil && old.ID() != "" {
			uc = old
		}
	}
	if tokenCode := ctx.URLParam("token"); tokenCode != "" && uc == nil && strings.ToUpper(ctx.Method()) == "GET" {
		if uc, err = GetTokenContext(tokenCode); err != nil {
			glog.Error(err)
		} else if uc != nil {
			uc.SetID(utils.GUID())
			ctx.Values().Set(context.DefaultContextKey, uc)
			ctx.SetCookieKV(context.AuthSessionKey, uc.ToTokenString())
		}
	}
	if tokenCode := ctx.GetHeader("token"); tokenCode != "" && uc == nil {
		if uc, err = GetTokenContext(tokenCode); err != nil {
			glog.Error(err)
		} else if uc != nil {
			uc.SetID(utils.GUID())
			ctx.Values().Set(context.DefaultContextKey, uc)
			ctx.SetCookieKV(context.AuthSessionKey, uc.ToTokenString())
		}
	}
	if tokenCode := ctx.GetCookie("token"); tokenCode != "" && uc == nil {
		if uc, err = GetTokenContext(tokenCode); err != nil {
			glog.Error(err)
		} else if uc != nil {
			uc.SetID(utils.GUID())
			ctx.Values().Set(context.DefaultContextKey, uc)
			ctx.SetCookieKV(context.AuthSessionKey, uc.ToTokenString())
		}
	}
	ctx.Next()

}
