package utils

import (
	"reflect"
)

func ToBoolByReflect(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		val := v.String()
		return "1" == val || "true" == val
	case reflect.Array, reflect.Map, reflect.Slice:
		return 0 != v.Len()
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return 0 != v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return 0 != v.Uint()
	case reflect.Float32, reflect.Float64:
		return 0 != v.Float()
	}
	return false
}

func ToBool(v interface{}) bool {
	return ToBoolByReflect(reflect.ValueOf(v))
}
