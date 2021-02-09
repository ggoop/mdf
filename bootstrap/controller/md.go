package controller

import (
	"fmt"
	"github.com/kataras/iris"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/bootstrap/resource"
	"github.com/ggoop/mdf/bootstrap/service"
	"github.com/ggoop/mdf/framework/context"
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/files"
	"github.com/ggoop/mdf/framework/http/results"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/framework/mof"
	"github.com/ggoop/mdf/utils"
)

type MdController struct {
	Ctx   iris.Context
	MdSv  *service.MdSv
	Repo  *repositories.MysqlRepo
	OssSv *service.OssSv
}

func (c *MdController) AnyInit() results.Result {
	ctx := c.Ctx.Values().Get(context.DefaultContextKey).(*context.Context)
	if ctx != nil {

	}
	types := make([]string, 0)
	if q := c.Ctx.URLParam("q"); q != "" {
		types = strings.Split(c.Ctx.URLParam("q"), ",")
	}
	if err := c.MdSv.InitData(types); err != nil {
		return results.ToError(err)
	}
	return results.ToJson(true)
}

// 获取页面
func (c *MdController) PostPageFetch() results.Result {
	ctx := c.Ctx.Values().Get(context.DefaultContextKey).(*context.Context)
	var postInput mof.ReqContext
	if err := c.Ctx.ReadJSON(&postInput); err != nil {
		return results.ToError(err)
	}
	postInput.EntID = ctx.EntID()
	postInput.UserID = ctx.UserID()
	if rtn, err := c.MdSv.GetPageInfo(postInput); err != nil {
		return results.ToError(err)
	} else {
		return results.ToJson(rtn)
	}
}
func (c *MdController) PostActionDo() results.Result {
	ctx := c.Ctx.Values().Get(context.DefaultContextKey).(*context.Context)
	var postInput mof.ReqContext
	if err := c.Ctx.ReadJSON(&postInput); err != nil {
		return results.ToError(err)
	}
	if postInput.EntID == "" {
		postInput.EntID = ctx.EntID()
	}
	postInput.UserID = ctx.UserID()
	if rtn, err := c.MdSv.DoAction(postInput); err != nil {
		return results.ToError(err)
	} else {
		return results.ToJson(rtn)
	}
}

/**
导入,支持文件导入，数据导入
*/
func (c *MdController) PostImport() results.Result {
	ctx := c.Ctx.Values().Get(context.DefaultContextKey).(*context.Context)
	if err := ctx.Valid(true, false); err != nil {
		return results.ToError(err)
	}
	var postInput mof.ReqContext
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
	postInput.Command = "import"
	postInput.UserID = ctx.UserID()
	if rtn, err := c.MdSv.DoAction(postInput); err != nil {
		return results.ToError(err)
	} else {
		return results.ToJson(rtn)
	}
}

/**
下载模板
*/
func (c *MdController) GetByTemplate(pageID string) results.Result {
	if pageID == "" || pageID == "0" {
		return results.ParamsRequired("id")
	}
	ctx := c.Ctx.Values().Get(context.DefaultContextKey).(*context.Context)
	fileItem, _ := c.OssSv.GetObjectBy(pageID)
	if fileItem == nil {
		fileItem, _ = c.OssSv.GetObjectByCode(ctx.EntID(), pageID)
	}
	if fileItem == nil {
		fileItem, _ = c.OssSv.GetObjectByCode("", pageID)
	}
	if fileItem == nil {
		if mdPage, err := c.MdSv.GetPage(pageID); err != nil {
			return results.ToError(err)
		} else {
			fileItem = &model.OssObject{Name: mdPage.Name + ".xlsx", Path: pageID + ".xlsx", OssBucket: "storage/templates"}
		}
	}
	if fileItem == nil {
		return results.NotFound("file")
	}
	filePath := utils.JoinCurrentPath(filepath.Join(fileItem.OssBucket, fileItem.Path))
	if !utils.PathExists(filePath) {
		fileItem.Path = pageID + ".template.xlsx"
		filePath = utils.JoinCurrentPath(filepath.Join(fileItem.OssBucket, fileItem.Path))
		if !utils.PathExists(filePath) {
			return results.NotFound("file")
		}
	}
	c.Ctx.Header("Content-Type", fileItem.MimeType)
	fileName := fileItem.Name
	if fileItem.OriginalName != "" {
		fileName = fileItem.OriginalName
	}
	c.Ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s;filename*=%s", fileName, url.QueryEscape(fileName)))
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return results.ToError(err)
	}
	c.Ctx.Write(file)
	return nil
}

func (c *MdController) GetEnumsBy(groupID string) results.Result {
	enumSv := md.NewEntitySv(c.Repo)
	items, err := enumSv.GetEnumBy(groupID)
	if err != nil {
		return results.ToError(err)
	}
	return results.ToJson((&resources.EnumRes{}).ToResource(&items))
}
func (c *MdController) GetEntities() results.Result {
	var postInput struct {
		Page     int    `form:"page"`
		PageSize int    `form:"page_size"`
		Q        string `form:"q"`
	}
	if err := c.Ctx.ReadForm(&postInput); err != nil {
		return results.ToError(err)
	}
	type objectItem struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	}
	if postInput.Q != "" {
		postInput.Q = "%" + postInput.Q + "%"
	}
	items := make([]objectItem, 0)
	query := c.Repo.Table(fmt.Sprintf("%v as d", c.Repo.NewScope(md.MDEntity{}).TableName()))
	query = query.Select([]string{"d.id", "d.code", "d.name"})
	query = query.Where("d.domain = ?", "suite").Order("d.code")
	if postInput.Q != "" {
		query = query.Where("d.code like ? or d.name like ? or d.full_name like ?", postInput.Q, postInput.Q, postInput.Q)
	}
	pager, err := service.PaginateScan(query, &items, postInput.Page, postInput.PageSize)
	if err != nil {
		return results.ToError(err)
	}
	return results.ToJson((&resources.PagerRes{}).ToResource(pager))
}
