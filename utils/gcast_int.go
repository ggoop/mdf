package utils

import (
	"reflect"
	"strconv"
	"strings"
)

func ToInt64ByReflect(v reflect.Value) int64 {
	switch v.Kind() {
	case reflect.String:
		var val int64
		if strings.Contains(v.Interface().(string), ".") {
			fval, _ := strconv.ParseFloat(v.Interface().(string), 64)
			val = int64(fval)
		} else {
			val, _ = strconv.ParseInt(v.Interface().(string), 10, 64)
		}
		return val
	case reflect.Array, reflect.Map, reflect.Slice:
		return 0
	case reflect.Bool:
		if v.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return int64(v.Float())
	}
	return 0
}

func ToInt64(v interface{}) int64 {
	return ToInt64ByReflect(reflect.ValueOf(v))
}

func ToInt32(v interface{}) int32 {
	return int32(ToInt64ByReflect(reflect.ValueOf(v)))
}

func ToInt(v interface{}) int {
	return int(ToInt64ByReflect(reflect.ValueOf(v)))
}

/// MARK: Int64

func ToUint64ByReflect(v reflect.Value) uint64 {
	switch v.Kind() {
	case reflect.String:
		var val uint64
		if strings.Contains(v.Interface().(string), ".") {
			fval, _ := strconv.ParseFloat(v.Interface().(string), 64)
			val = uint64(fval)
		} else {
			val, _ = strconv.ParseUint(v.Interface().(string), 10, 64)
		}
		return val
	case reflect.Array, reflect.Map, reflect.Slice:
		return 0
	case reflect.Bool:
		if v.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uint64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return uint64(v.Float())
	}
	return 0
}

func ToUint64(v interface{}) uint64 {
	return ToUint64ByReflect(reflect.ValueOf(v))
}

func ToUint32(v interface{}) uint32 {
	return uint32(ToUint64ByReflect(reflect.ValueOf(v)))
}

func ToUint(v interface{}) uint {
	return uint(ToUint64ByReflect(reflect.ValueOf(v)))
}
