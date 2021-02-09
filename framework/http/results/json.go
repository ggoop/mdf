package results

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/ggoop/mdf/framework/configs"
	"github.com/ggoop/mdf/framework/glog"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"

	"github.com/ggoop/mdf/utils"
)

type (
	Result = mvc.Result
	Map    = iris.Map
)

func ToJson(data interface{}) Result {
	return mvc.Response{
		Object: data,
	}
}
func Unauthenticated() Result {
	return ToError(fmt.Errorf("Unauthenticated"), iris.StatusUnauthorized)
}
func ParamsRequired(params ...string) Result {
	return ToError(fmt.Errorf("参数 %s 不能为空!", params), iris.StatusBadRequest)
}
func ParamsFailed(params ...string) Result {
	return ToError(fmt.Errorf("参数 %s 不正确!", params), iris.StatusUnsupportedMediaType)
}
func NotFound(params ...string) Result {
	return ToError(fmt.Errorf("找不到 %s", strings.Join(params, " ")), iris.StatusNotFound)
}
func ToSingle(data interface{}) Result {
	return mvc.Response{
		Object: iris.Map{"data": data},
	}
}

// func(err,http code, err code)
func ToError(err interface{}, code ...int) Result {
	res := mvc.Response{}
	obj := iris.Map{}
	httpCode := iris.StatusBadRequest
	dataCode := iris.StatusBadRequest
	if ev, ok := err.(utils.GError); ok {
		obj["msg"] = ev.Error()
		if ev.Code > 0 {
			dataCode = ev.Code
			httpCode = ev.Code
		}
	} else if ev, ok := err.(error); ok {
		obj["msg"] = ev.Error()
	} else {
		obj["msg"] = err
	}
	//http 状态码
	if code != nil && len(code) > 0 {
		if code[0] >= 100 && code[0] < 1000 {
			httpCode = code[0]
			dataCode = code[0]
		} else {
			dataCode = code[0]
		}
	}
	//使用指定的异常代码
	if code != nil && len(code) > 1 {
		dataCode = code[1]
	}
	if httpCode <= 100 || httpCode > 1000 {
		httpCode = iris.StatusBadRequest
	}
	obj["code"] = dataCode
	res.Code = httpCode
	res.Object = obj
	if configs.Default.App.Debug {
		glog.Error(fmt.Sprintf("%s,at %s", obj["msg"], debug.Stack()))
	}
	return res
}
