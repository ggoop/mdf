package utils

import (
	"errors"
	"fmt"
)

// err : error ,
// code :x xx xxx xxxx,level-product-app-number 10bit
type GError struct {
	Err    error
	Code   int
	Format string
}

func (e GError) Error() string {
	if e.Format != "" {
		return fmt.Sprintf(e.Format, e.Err.Error(), e.Code)
	} else {
		return fmt.Sprintf("%s", e.Err.Error())
	}
}

// err : error ,
// code :x xx xxx xxxx,level-product-app-number 10bit
func ToError(err interface{}, code ...int) GError {
	var e error
	if v, ok := err.(GError); ok {
		e = v
	} else if v, ok := err.(error); ok {
		e = v
	} else if v, ok := err.(string); ok {
		e = errors.New(v)
	} else {
		e = errors.New(fmt.Sprint(err))
	}
	r := GError{Err: e}
	if code != nil && len(code) > 0 {
		r.Code = code[0]
	}
	return r
}
