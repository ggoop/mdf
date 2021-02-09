package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ggoop/mdf/bootstrap/model"
	resources "github.com/ggoop/mdf/bootstrap/resource"
	"github.com/ggoop/mdf/bootstrap/service"
	"github.com/ggoop/mdf/configs"
	"github.com/ggoop/mdf/utils"

	"github.com/kataras/iris"

	"github.com/ggoop/mdf/context"
	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/http/results"
	"github.com/ggoop/mdf/repositories"
)

type OAuthController struct {
	Ctx       iris.Context
	ProductSv *service.ProductSv
	Repo      *repositories.MysqlRepo
	UserSv    *service.UserSv
	EntSv     *service.EntSv
	TokenSv   *service.TokenSv
}

func SetAuthSession(ctx iris.Context, token *context.Context) bool {
	if token != nil {
		ctx.SetCookieKV(context.AuthSessionKey, token.ToTokenString())
	} else {
		ctx.RemoveCookie(context.AuthSessionKey)
	}
	return true
}

/**
免登接口
*/
func (c *OAuthController) GetFreeSign() results.Result {
	code := c.Ctx.URLParam("code")
	appType := c.Ctx.URLParam("type")

	authFree := model.AuthFree{}
	var resBodyObj struct {
		AppID        string `json:"app_id"`
		Openid       string `json:"openid"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		ErrMsg       string `json:"errmsg"`
		ErrCode      int    `json:"errcode"`
	}
	if code != "" && appType == "wx" {
		appid := configs.Default.GetValue("wx.appid")
		client := &http.Client{}
		remoteUrl, _ := url.Parse("https://api.weixin.qq.com")
		remoteUrl.Path = "/sns/oauth2/access_token"
		queryParams := remoteUrl.Query()
		queryParams.Set("appid", appid)
		queryParams.Set("secret", configs.Default.GetValue("wx.secret"))
		queryParams.Set("code", code)
		queryParams.Set("grant_type", "authorization_code")
		remoteUrl.RawQuery = queryParams.Encode()

		req, err := http.NewRequest("GET", remoteUrl.String(), nil)
		if err != nil {
			return results.ToError(err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return results.ToError(err)
		}
		defer resp.Body.Close()

		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return results.ToError(err)
		}

		if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
			return results.ToError(err)
		}
		if resp.StatusCode != 200 || resBodyObj.ErrMsg != "" {
			return results.ToError(resBodyObj.ErrMsg)
		}

		resBodyObj.AppID = appid

	} else {
		resBodyObj.AppID = c.Ctx.URLParam("appid")
		resBodyObj.Openid = c.Ctx.URLParam("openid")
	}
	if resBodyObj.Openid == "" {
		return results.ToError("找不到对应用户")
	}
	if resBodyObj.AppID == "" {
		return results.ToError("找不到应用")
	}
	c.Repo.Model(authFree).Where("src_openid=?", resBodyObj.Openid).Where("src_app=? and src_type=?", resBodyObj.AppID, appType).Take(&authFree)
	authFree.SrcAccessToken = resBodyObj.AccessToken
	authFree.SrcRefreshToken = resBodyObj.RefreshToken
	authFree.SrcOpenid = resBodyObj.Openid
	authFree.SrcApp = resBodyObj.AppID
	authFree.SrcType = appType

	if authFree.ID == "" {
		authFree.Qty = 1
		c.Repo.Create(&authFree)
	} else {
		updates := make(map[string]interface{})
		updates["Qty"] = authFree.Qty + 1
		updates["SrcAccessToken"] = resBodyObj.AccessToken
		updates["SrcRefreshToken"] = resBodyObj.RefreshToken
		c.Repo.Model(authFree).Where("id=?", authFree.ID).Updates(updates)
	}
	token := make(map[string]interface{})
	token["id"] = authFree.ID
	token["account"] = authFree.Account
	if authFree.UserID != "" {
		if u, _ := c.UserSv.GetUserBy(authFree.UserID); u != nil {
			token["user"] = (resources.UserRes{}).ToItem(u)
		}
		token["user_token"] = authFree.UserToken
	}
	if authFree.EntID != "" {
		if u, _ := c.EntSv.GetEntBy(authFree.EntID); u != nil {
			token["ent"] = (resources.EntRes{}).ToItem(u)
		}
		token["ent_token"] = authFree.EntToken
	}
	return results.ToJson(token)
}

func (c *OAuthController) PostFreeLogin() results.Result {
	var credit struct {
		Type      string      `json:"type"`
		ID        string      `json:"id"`
		Account   string      `json:"account"`
		UserID    string      `json:"user_id"`
		UserToken utils.SJson `json:"user_token"`
		EntID     string      `json:"ent_id"`
		EntToken  utils.SJson `json:"ent_token"`
	}
	if err := c.Ctx.ReadJSON(&credit); err != nil {
		return results.ToError(err)
	}
	if credit.ID == "" {
		return results.ToError("缺少参数")
	}
	updates := make(map[string]interface{})
	if strings.Contains(credit.Type, "ent") {
		updates["EntToken"] = credit.EntToken
		updates["EntID"] = credit.EntID
	}
	if strings.Contains(credit.Type, "user") {
		updates["UserID"] = credit.UserID
		updates["UserToken"] = credit.UserToken
	}
	if credit.Account != "" {
		updates["Account"] = credit.Account
	}
	if len(updates) > 0 {
		authFree := model.AuthFree{}
		if err := c.Repo.Model(authFree).Take(&authFree, "id=?", credit.ID).Error; err != nil {
			return results.ToError(err)
		} else {
			updates["Qty"] = authFree.Qty + 1
			c.Repo.Model(model.AuthFree{}).Where("id=?", authFree.ID).Updates(updates)
			return results.ToSingle(true)
		}
	}
	return results.ToError("免登失败")
}
func (c *OAuthController) GetFreeLogout() results.Result {
	id := c.Ctx.URLParam("id")
	authType := c.Ctx.URLParam("type")
	if id != "" && authType != "" {
		updates := make(map[string]interface{})
		if strings.Contains(authType, "ent") {
			updates["EntToken"] = ""
			updates["EntID"] = ""
		}
		if strings.Contains(authType, "user") {
			updates["UserID"] = ""
			updates["UserToken"] = ""
		}
		if len(updates) > 0 {
			c.Repo.Model(model.AuthFree{}).Where("id=?", id).Updates(updates)
			return results.ToSingle(true)
		}
	}
	return results.ToError("登出失败")
}

// 依据token获取token信息
func (c *OAuthController) GetTokenBy(tokenCode string) results.Result {
	if t, finder := c.TokenSv.Get(tokenCode); finder {
		token := context.Context{}
		token.SetUserID(t.UserID)
		token.SetClientID(t.ClientID)
		token.SetEntID(t.EntID)
		tokenStr := strings.Split(token.ToTokenString(), " ")
		if len(tokenStr) < 2 {
			return results.ToError("create token error!")
		}
		tData := utils.Map{}
		if t.UserID != "" {
			tData["user_id"] = t.UserID
		}
		if t.UserName != "" {
			tData["user_name"] = t.UserName
		}
		if t.EntID != "" {
			tData["ent_id"] = t.EntID
		}
		rtnData := utils.Map{"data": tData}
		rtnData["token"] = results.Map{
			"access_token": tokenStr[1],
			"type":         tokenStr[0],
			"expires_in":   token.ExpiresAt(),
			"token":        tokenCode,
		}
		return results.ToJson(rtnData)
	} else {
		return results.ToError("token无效")
	}
}

// 获取当前登录信息
func (c *OAuthController) GetCurrent() results.Result {
	ctx := c.Ctx.Values().Get(context.DefaultContextKey).(*context.Context)
	rtnData := utils.Map{}
	if vid := ctx.UserID(); vid != "" {
		if user, _ := c.UserSv.GetUserBy(vid); user != nil {
			rtnData["user"] = (&resources.UserRes{}).ToItem(user)
		}
	}
	if vid := ctx.EntID(); vid != "" {
		if ent, _ := c.EntSv.GetEntBy(vid); ent != nil {
			rtnData["ent"] = (&resources.EntRes{}).ToItem(ent)
		}
	}
	return results.ToJson(rtnData)
}

// 登录某个企业
func (c *OAuthController) PostEnt() results.Result {
	ctx := c.Ctx.Values().Get(context.DefaultContextKey).(*context.Context)
	if err := ctx.VerifyUser(); err != nil {
		return results.ToError(err)
	}
	var credit struct {
		Type  string `json:"type"`
		EndID string `json:"ent_id"`
	}
	if err := c.Ctx.ReadJSON(&credit); err != nil {
		return results.ToError(err)
	}
	ent, err := c.EntSv.GetEntBy(credit.EndID)
	if err != nil {
		return results.ToError(err)
	}
	token := ctx.Copy()
	token.SetEntID(ent.ID)
	tokenStr := strings.Split(token.ToTokenString(), " ")
	if len(tokenStr) < 2 {
		return results.ToError("create token error!")
	}
	SetAuthSession(c.Ctx, &token)

	rtnData := utils.Map{}
	rtnData["token"] = results.Map{"access_token": tokenStr[1], "type": tokenStr[0], "expires_in": token.ExpiresAt()}
	return results.ToJson(rtnData)
}

func (c *OAuthController) PostLogout() results.Result {
	SetAuthSession(c.Ctx, nil)
	return results.ToSingle(true)
}

// 用户登录接口
func (c *OAuthController) PostLogin() results.Result {
	ctx := c.Ctx.Values().Get(context.DefaultContextKey).(*context.Context)
	var credit struct {
		Type      string `json:"type"` //登录类型,account,mobile
		Account   string `json:"account"`
		Password  string `json:"password"`
		Mobile    string `json:"mobile"`
		Captcha   string `json:"captcha"`
		EntOpenid string `json:"ent_openid"`
	}
	if err := c.Ctx.ReadJSON(&credit); err != nil {
		return results.ToError(err)
	}
	user, err := c.UserSv.GetUserByAccount(credit.Account)
	if err != nil {
		return results.ToError("账号无效")
	}
	if !user.Enabled.IsTrue() {
		return results.ToError("账号已被锁定，请与管理员联系！")
	}
	if credit.Password == "" {
		return results.ToError("密码不能为空~")
	}
	credited := false
	if credit.Password != "" {
		credit.Password, err = utils.AesCFBEncrypt(credit.Password, user.Openid)
		if user.Password != credit.Password {
			return results.ToError("账号或者密码错误！")
		}
		credited = true
	}
	if !credited {
		err := fmt.Errorf("登录验证失败~")
		return results.ToError(err)
	}
	rtnData := utils.Map{}
	token := context.Context{}
	defaultEntID := ""
	//查询用户可登录的企业
	if credit.EntOpenid != "" {
		if ent, err := c.EntSv.GetByOpenid(credit.EntOpenid); err != nil {
			return results.ToError(err)
		} else {
			defaultEntID = ent.ID
		}
	}
	if ent := c.EntSv.GetUserDefaultEnt(user.ID, defaultEntID); ent != nil {
		token.SetEntID(ent.ID)
		rtnData["ent"] = (&resources.EntRes{}).ToItem(ent)
	}
	token.SetUserID(user.ID)
	token.SetClientID(ctx.ClientID())
	tokenStr := strings.Split(token.ToTokenString(), " ")
	if len(tokenStr) < 2 {
		return results.ToError("create token error!")
	}
	SetAuthSession(c.Ctx, &token)
	rtnData["user"] = utils.Map{"id": user.ID, "name": user.Name, "account": user.Account, "mobile": user.Mobile, "email": user.Email, "avatar": user.AvatarUrl}
	rtnData["token"] = results.Map{"access_token": tokenStr[1], "type": tokenStr[0], "expires_in": token.ExpiresAt()}
	return results.ToJson(rtnData)
}

// 获取注册token
func (c *OAuthController) PostToken() results.Result {
	var tokenInput struct {
		GrantType string `json:"grant_type"` //client_credentials,token,uc
		//client_credentials
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		//uc
		UserAccount string `json:"user_account"`
		UserID      string `json:"user_id"`
		EntOpenid   string `json:"ent_openid"`
		EntID       string `json:"ent_id"`
		//token
		Token string `json:"token"`
	}
	grantTypeClientCredentials := "client_credentials"
	grantTypeUc := "uc"
	grantTypeToken := "token"

	var err error
	if err = c.Ctx.ReadJSON(&tokenInput); err != nil {
		return results.ToError(err)
	}
	// 目前只支持client_credentials,uc,token
	grantTypes := []string{grantTypeClientCredentials, grantTypeUc, grantTypeToken}
	if utils.StringsContains(grantTypes, tokenInput.GrantType) >= 0 {
		glog.Infof("begin issue token")
	} else {
		return results.ToError(fmt.Errorf("not imp"))
	}
	if tokenInput.ClientID == "" {
		return results.ParamsRequired("client_id")
	}
	tokenContext := context.Context{}
	client := model.Client{}
	if err := c.Repo.Model(model.Client{}).Take(&client, "id=?", tokenInput.ClientID).Error; err != nil {
		return results.ToError(err)
	}
	if tokenInput.GrantType == grantTypeClientCredentials {
		if tokenInput.ClientSecret == "" {
			return results.ParamsRequired("client_secret")
		}
		if client.Secret != tokenInput.ClientSecret {
			return results.ToError("client_secret is error!")
		}
	}
	tokenContext.SetClientID(client.ID)

	creditUserID := ""
	if tokenInput.UserID != "" {
		if u, err := c.UserSv.GetUserBy(tokenInput.UserID); err != nil || u == nil || u.ID == "" {
			return results.ToError("用户不存在!")
		} else {
			creditUserID = u.ID
		}
	}
	if creditUserID == "" && tokenInput.UserAccount != "" {
		if u, err := c.UserSv.GetUserByAccount(tokenInput.UserAccount); err != nil || u == nil || u.ID == "" {
			return results.ToError("用户不存在!")
		} else {
			creditUserID = u.ID
		}
	}

	creditEntID := ""
	if tokenInput.EntID != "" {
		if u, err := c.EntSv.GetEntBy(tokenInput.EntID); err != nil || u == nil || u.ID == "" {
			return results.ToError("企业不存在!")
		} else {
			creditEntID = u.ID
		}
	}
	if creditEntID == "" && tokenInput.EntOpenid != "" {
		if u, err := c.EntSv.GetByOpenid(tokenInput.EntOpenid); err != nil || u == nil || u.ID == "" {
			return results.ToError("企业不存在!")
		} else {
			creditEntID = u.ID
		}
	}

	//如果是用户中心登录，则需要校验用户
	if tokenInput.GrantType == grantTypeUc {
		if creditUserID == "" || creditEntID == "" {
			return results.ToError("账号或者企业为空")
		}
		tokenContext.SetUserID(creditUserID)
		tokenContext.SetEntID(creditEntID)
	}
	//token 切换
	if tokenInput.GrantType == grantTypeToken {
		if tokenInput.Token == "" {
			return results.ToError("token为空")
		}
		if t, finder := c.TokenSv.Get(tokenInput.Token); finder {
			tokenContext.SetEntID(t.EntID)
			tokenContext.SetUserID(t.UserID)
		} else {
			return results.ToError("token不存在")
		}
		if creditUserID != "" && tokenContext.UserID() == "" {
			tokenContext.SetUserID(creditUserID)
		}
		if creditEntID != "" && tokenContext.EntID() == "" {
			tokenContext.SetEntID(creditEntID)
		}
	}
	tokenStr := strings.Split(tokenContext.ToTokenString(), " ")
	if len(tokenStr) < 2 {
		return results.ToError("create token error!")
	}

	SetAuthSession(c.Ctx, &tokenContext)

	t := &model.Token{}
	t.ClientID = tokenContext.ClientID()
	t.EntID = tokenContext.EntID()
	t.UserID = tokenContext.UserID()
	t.Type = tokenInput.GrantType
	t = c.TokenSv.Create(t)
	token := map[string]interface{}{
		"access_token": tokenStr[1],
		"type":         tokenStr[0],
		"expires_in":   tokenContext.ExpiresAt(),
		"token":        t.Token,
	}
	return results.ToJson(results.Map{"token": token})
}
