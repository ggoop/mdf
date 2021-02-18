package utils

import "github.com/shopspring/decimal"

func ToDecimal(v interface{}) decimal.Decimal {
	if vv, ok := v.(decimal.Decimal); ok {
		return vv
	} else if vv, ok := v.(string); ok {
		rv, _ := decimal.NewFromString(vv)
		return rv
	} else if vv, ok := v.(float64); ok {
		return decimal.NewFromFloat(vv)
	} else if vv, ok := v.(float32); ok {
		return decimal.NewFromFloat32(vv)
	}
	return decimal.Zero
}
