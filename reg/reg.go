package reg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"io"
	"strings"

	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/ggoop/mdf/configs"
	"github.com/ggoop/mdf/context"
	"github.com/ggoop/mdf/glog"
)

//code 到 RegObject 的缓存
var _codeRegObjectMap map[string]*RegObject = make(map[string]*RegObject)

func setRegObjectCache(item *RegObject) {
	if item != nil {
		_codeRegObjectMap[strings.ToLower(item.Code)] = item
	}
}
func getRegObjectCache(code string) *RegObject {
	if obj, ok := _codeRegObjectMap[strings.ToLower(code)]; ok {
		return obj
	}
	return nil
}

/**
获取注册中心地址
*/
func GetRegistry() string {
	registry := configs.Default.App.Registry
	if registry == "" {
		registry = fmt.Sprintf("http://127.0.0.1:%s", configs.Default.App.Port)
	}
	return registry
}
func localRegistry() string {
	return fmt.Sprintf("http://127.0.0.1:%s", configs.Default.App.Port)
}

/**
通过token获取上下文
*/
func GetTokenContext(tokenCode string) (*context.Context, error) {
	//1、权限注册中心、2、应用注册中心，3、本地
	authAddr := ""
	if sers, err := FindByCode(configs.Default.Auth.Address); sers != nil {
		authAddr = sers.Address
	} else {
		glog.Error(err)
	}
	if authAddr == "" {
		authAddr = configs.Default.App.Registry
	}
	if authAddr == "" {
		authAddr = fmt.Sprintf("http://127.0.0.1:%s", configs.Default.App.Port)
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
	token := &context.Context{}
	token, _ = token.FromTokenString(fmt.Sprintf("%s %s", resBodyObj.Token.Type, resBodyObj.Token.AccessToken))
	return token, nil
}

var m_cronCache *cron.Cron

/**
由配置文件信息，注册
*/
func RegisterDefault() {
	if m_cronCache == nil {
		m_cronCache := cron.New()
		m_cronCache.AddFunc("@every 120s", registerDefault)
		m_cronCache.Start()
	}
	registerDefault()
}
func registerDefault() {
	address := configs.Default.App.Address
	if address == "" {
		address = fmt.Sprintf("http://127.0.0.1:%s", configs.Default.App.Port)
	}
	Register(RegObject{
		Code:          configs.Default.App.Code,
		Name:          configs.Default.App.Name,
		Address:       address,
		PublicAddress: configs.Default.App.PublicAddress,
		Configs:       configs.Default,
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
	regHost := GetRegistry()
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
func DoHttpRequest(serverName, method, path string, body io.Reader) ([]byte, error) {
	regs, err := FindByCode(serverName)
	if err != nil {
		return nil, err
	}
	if regs == nil || regs.Address == "" {
		return nil, glog.Error("找不到服务,", glog.String("serverName", serverName))
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

func GetServerAddr(code string) (string, error) {
	//如果是 http,或者 https，则直接返回
	if code != "" && strings.Index(strings.ToLower(code), "http") == 0 {
		return code, nil
	}
	if d, err := FindByCode(code); err != nil {
		return "", err
	} else if d != nil && d.Address != "" {
		return d.Address, nil
	} else if d != nil && d.PublicAddress != "" {
		return d.PublicAddress, nil
	}
	return "", nil
}

/**
通过编码找到注册对象
*/
func FindByCode(code string) (*RegObject, error) {
	if code == "" {
		return nil, nil
	}
	//优先从缓存里取
	if cv := getRegObjectCache(code); cv != nil {
		return cv, nil
	}
	client := &http.Client{}
	client.Timeout = 2 * time.Second
	remoteUrl, _ := url.Parse(GetRegistry())
	remoteUrl.Path = fmt.Sprintf("/api/regs/%s", code)
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
	setRegObjectCache(resBodyObj.Data)

	return resBodyObj.Data, nil
}
