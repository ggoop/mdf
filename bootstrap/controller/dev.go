package controller

import (
	"github.com/kataras/iris"

	"github.com/ggoop/mdf/bootstrap/service"
	"github.com/ggoop/mdf/framework/context"
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/files"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/http/results"
	"github.com/ggoop/mdf/framework/md"
)

type DevController struct {
	Ctx   iris.Context
	MdSv  *service.MdSv
	Repo  *repositories.MysqlRepo
	OssSv *service.OssSv
}

func (c *DevController) PostWidgetImport() results.Result {
	ctx := c.Ctx.Values().Get(context.DefaultContextKey).(*context.Context)
	if err := ctx.Valid(true, false); err != nil {
		glog.Error(err)
	}
	var postInput md.ReqContext
	if err := c.Ctx.ReadForm(&postInput); err != nil {
		c.Ctx.ReadJSON(&postInput)
	}
	file, info, err := c.Ctx.FormFile("file")
	if err != nil {
		return results.ToError(err)
	}
	if file != nil {
		postInput.Tag = info.Filename
		defer file.Close()
		if datas, err := files.NewExcelSv().GetExcelDatasByReader(file); err != nil {
			return results.ToError(err)
		} else {
			postInput.Data = datas
		}
	}
	if postInput.EntID == "" {
		postInput.EntID = ctx.EntID()
	}
	postInput.Action = "import"
	postInput.UserID = ctx.UserID()
	if rtn := c.MdSv.DoAction(postInput); rtn.Error != nil {
		return results.ToError(err)
	} else {
		return results.ToJson(rtn)
	}
}
