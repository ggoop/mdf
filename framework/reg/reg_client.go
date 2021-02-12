package reg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ggoop/mdf/utils"
	"github.com/robfig/cron"
	"io"
	"strings"

	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/ggoop/mdf/framework/glog"
)

//code 到 RegObject 的缓存
var clientRegMap map[string]*RegObject = make(map[string]*RegObject)

func setRegObjectCache(code string, item *RegObject) {
	clientRegMap[strings.ToLower(code)] = item
}
func getRegObjectCache(code string) *RegObject {
	if obj, ok := clientRegMap[strings.ToLower(code)]; ok {
		return obj
	}
	return nil
}

/**
获取注册中心地址
*/
func getRegistryHost() string {
	registry := utils.DefaultConfig.App.Registry
	if registry == "" {
		registry = fmt.Sprintf("http://127.0.0.1:%s", utils.DefaultConfig.App.Port)
	}
	return registry
}

/**
通过token获取上下文
*/
func GetTokenContext(tokenCode string) (*utils.TokenContext, error) {
	//1、权限注册中心、2、应用注册中心，3、本地
	authAddr := ""
	if ser, err := FindServerByCode(utils.DefaultConfig.Auth.Code); ser != nil {
		authAddr = ser.Address
	} else {
		glog.Error(err)
	}
	if authAddr == "" {
		authAddr = fmt.Sprintf("http://127.0.0.1:%s", utils.DefaultConfig.App.Port)
	}
	client := &http.Client{}
	client.Timeout = 2 * time.Second
	remoteUrl, _ := url.Parse(authAddr)
	remoteUrl.Path = fmt.Sprintf("/api/oauth/token/%s", tokenCode)
	req, err := http.NewRequest("GET", remoteUrl.String(), nil)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	var resBodyObj struct {
		Msg   string `json:"msg"`
		Token struct {
			AccessToken string `json:"access_token"`
			Type        string `json:"type"`
		} `json:"token"`
	}
	if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
		glog.Error(err)
		return nil, err
	}
	if resp.StatusCode != 200 || resBodyObj.Msg != "" {
		glog.Error(resBodyObj.Msg)
		return nil, err
	}
	token := utils.NewTokenContext()
	token, _ = token.FromTokenString(fmt.Sprintf("%s %s", resBodyObj.Token.Type, resBodyObj.Token.AccessToken))
	return token, nil
}

var m_cronCache *cron.Cron

/**
由配置文件信息，注册
*/
func StartClient() {
	if m_cronCache == nil {
		m_cronCache := cron.New()
		m_cronCache.AddFunc("@every 120s", registerDefault)
		m_cronCache.Start()
	}
	registerDefault()
}
func registerDefault() {
	address := utils.DefaultConfig.App.Address
	if address == "" {
		address = fmt.Sprintf("http://127.0.0.1:%s", utils.DefaultConfig.App.Port)
	}
	Register(RegObject{
		Code:    utils.DefaultConfig.App.Code,
		Name:    utils.DefaultConfig.App.Name,
		Address: address,
		Configs: utils.DefaultConfig,
	})
}
func Register(item RegObject) error {
	client := &http.Client{}
	client.Timeout = 3 * time.Second
	postBody, err := json.Marshal(item)
	if err != nil {
		glog.Error(err)
		return err
	}
	regHost := getRegistryHost()
	remoteUrl, _ := url.Parse(regHost)
	remoteUrl.Path = "/api/regs/register"
	req, err := http.NewRequest("POST", remoteUrl.String(), bytes.NewBuffer([]byte(postBody)))
	if err != nil {
		glog.Error(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		glog.Error(err)
		return err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return err
	}
	var resBodyObj struct {
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
	if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
		glog.Error(err)
		return err
	}
	if resp.StatusCode != 200 || resBodyObj.Msg != "" {
		glog.Error(resBodyObj.Msg)
		return err
	}
	glog.Error("成功注册：", glog.Any("Item", item), glog.Any("RegHost", regHost))
	return nil
}
func DoHttpRequest(serverCode, method, path string, body io.Reader) ([]byte, error) {
	regs, err := FindServerByCode(serverCode)
	if err != nil {
		return nil, err
	}
	if regs == nil || regs.Address == "" {
		return nil, glog.Error("找不到服务,", glog.String("serverCode", serverCode))
	}
	serverUrl := regs.Address
	client := &http.Client{}
	remoteUrl, err := url.Parse(serverUrl)
	if err != nil {
		return nil, err
	}
	remoteUrl.Path = path
	req, err := http.NewRequest(method, remoteUrl.String(), body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		var resBodyObj struct {
			Msg string `json:"msg"`
		}
		if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
			return nil, err
		}
		return nil, glog.Error(resBodyObj.Msg)
	}
	return resBody, nil
}

/**
通过编码找到注册对象
*/
func FindServerByCode(serverCode string) (*RegObject, error) {
	if serverCode == "" {
		return nil, nil
	}
	//优先从缓存里取
	if cv := getRegObjectCache(serverCode); cv != nil {
		return cv, nil
	}
	client := &http.Client{}
	client.Timeout = 2 * time.Second
	remoteUrl, _ := url.Parse(getRegistryHost())
	remoteUrl.Path = fmt.Sprintf("/api/regs/%s", serverCode)
	req, err := http.NewRequest("GET", remoteUrl.String(), nil)

	if err != nil {
		glog.Error(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	var resBodyObj struct {
		Msg  string     `json:"msg"`
		Data *RegObject `json:"data"`
	}
	glog.Error(string(resBody))
	if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
		glog.Error(err)
		return nil, err
	}
	if resp.StatusCode != 200 || resBodyObj.Msg != "" {
		glog.Error(resBodyObj.Msg)
		return nil, err
	}
	//设置缓存
	setRegObjectCache(resBodyObj.Data.Code, resBodyObj.Data)

	return resBodyObj.Data, nil
}
