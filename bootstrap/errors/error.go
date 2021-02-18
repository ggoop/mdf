package errors

import (
	"fmt"
	"github.com/ggoop/mdf/utils"
)

//xx xx xxx
func New(err interface{}, code int) utils.GError {
	if v, ok := err.(utils.GError); ok {
		return v
	} else if v, ok := err.(error); ok {
		return utils.ToError(v, code)
	} else {
		return utils.ToError(v, code)
	}
}
func ParamsRequired(params ...string) utils.GError {
	return New(fmt.Errorf("参数 %s 不能为空!", params), 10100100)
}
func ParamsFailed(params ...string) utils.GError {
	return New(fmt.Errorf("参数 %s 不正确!", params), 10100101)
}
func CodeError(params ...string) utils.GError {
	return New(fmt.Errorf("%s 格式不正确!", params), 10100102)
}
func ExistError(params ...string) utils.GError {
	return New(fmt.Errorf("%s 已存在!", params), 10100102)
}
func NotExistError(params ...string) utils.GError {
	return New(fmt.Errorf("%s 不已存!", params), 10100102)
}
func IsDeleted(params ...string) utils.GError {
	return New(fmt.Errorf("%s 已经被删除!", params), 10100102)
}
func IsQuoted(params ...string) utils.GError {
	return New(fmt.Errorf("已经被 %s 引用!", params), 10100102)
}
