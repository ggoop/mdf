package glog

import (
	"go.uber.org/zap"
	"time"
)

type Field = zap.Field

func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}
func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}
func Float32(key string, val float32) Field {
	return zap.Float32(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Int32(key string, val int32) Field {
	return zap.Int32(key, val)
}

func Int16(key string, val int16) Field {
	return zap.Int16(key, val)
}

func Int8(key string, val int8) Field {
	return zap.Int8(key, val)
}

func String(key string, val string) Field {
	return zap.String(key, val)
}

func Uint(key string, val uint) Field {
	return zap.Uint(key, val)
}

func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}

func Uint32(key string, val uint32) Field {
	return zap.Uint32(key, val)
}

func Uint16(key string, val uint16) Field {
	return zap.Uint16(key, val)
}

func Uint8(key string, val uint8) Field {
	return zap.Uint8(key, val)
}

func Uintptr(key string, val uintptr) Field {
	return zap.Uintptr(key, val)
}

func Reflect(key string, val interface{}) Field {
	return zap.Reflect(key, val)
}

func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}
func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}
