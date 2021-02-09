package utils

import (
	"reflect"
	"strconv"
)

func ToFloat64ByReflect(v reflect.Value) float64 {
	switch v.Kind() {
	case reflect.String:
		val, _ := strconv.ParseFloat(v.String(), 64)
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
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	}
	return 0
}

func ToFloat64(v interface{}) float64 {
	return ToFloat64ByReflect(reflect.ValueOf(v))
}

func ToFloat32(v interface{}) float32 {
	return float32(ToInt64ByReflect(reflect.ValueOf(v)))
}

func ToFloat(v interface{}) float64 {
	return ToFloat64ByReflect(reflect.ValueOf(v))
}
