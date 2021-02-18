package utils

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/ggoop/mdf/framework/glog"
)

// JSONTime format json time field by myself
type Time struct {
	time.Time
}

const (
	Layout_YYYYMM          = "2006-01"
	Layout_YYYYMMDD        = "2006-01-02"
	Layout_YYYYMMDDHHIISS  = "2006-01-02 15:04:05"
	Layout_YYYYMMDDHHIISST = "2006-01-02T15:04:05"
	Layout_YYYYMMDD2       = "20060102"
)

func TimeNow() Time {
	return ToTime(time.Now())
}
func TimeNowPtr() *Time {
	return ToTimePtr(time.Now())
}
func ToTime(value interface{}) Time {
	if value == "" || value == nil {
		return Time{}
	}
	if v, ok := value.(time.Time); ok {
		return Time{v}
	}
	if v, ok := value.(string); ok {
		v = strings.Replace(v, `"`, "", -1)
		layout := Layout_YYYYMMDDHHIISS
		data := []rune(v)
		if len(data) >= len(Layout_YYYYMMDDHHIISST) && strings.Contains(v, "T") {
			layout = Layout_YYYYMMDDHHIISST
			data = data[:len(Layout_YYYYMMDDHHIISST)]
		} else if len(data) >= len(Layout_YYYYMMDDHHIISS) {
			layout = Layout_YYYYMMDDHHIISS
			data = data[:len(Layout_YYYYMMDDHHIISS)]
		} else if len(data) >= len(Layout_YYYYMMDD) {
			layout = Layout_YYYYMMDD
			data = data[:len(Layout_YYYYMMDD)]
		} else if len(data) >= len(Layout_YYYYMM) {
			layout = Layout_YYYYMM
			data = data[:len(Layout_YYYYMM)]
		}
		now, err := time.ParseInLocation(layout, string(data), time.Local)
		if err != nil {
			glog.Error(err)
		}
		return Time{now}
	}
	return Time{time.Now()}
}
func ToTimePtr(value interface{}) *Time {
	t := ToTime(value)
	return &t
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t Time) MarshalJSON() ([]byte, error) {
	if t.UnixNano() <= 0 || t.Unix() <= 0 || t.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}
func (t *Time) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		*t = Time{}
		return nil
	}
	ds := strings.Replace(string(data), `"`, "", -1)
	layout := Layout_YYYYMMDDHHIISS
	data = []byte(ds)
	if len(data) >= len(Layout_YYYYMMDDHHIISST) && strings.Contains(ds, "T") {
		layout = Layout_YYYYMMDDHHIISST
		data = data[:len(Layout_YYYYMMDDHHIISST)]
	} else if len(data) >= len(Layout_YYYYMMDDHHIISS) {
		layout = Layout_YYYYMMDDHHIISS
		data = data[:len(Layout_YYYYMMDDHHIISS)]
	} else if len(data) >= len(Layout_YYYYMMDD) {
		layout = Layout_YYYYMMDD
		data = data[:len(Layout_YYYYMMDD)]
	} else if len(data) >= len(Layout_YYYYMM) {
		layout = Layout_YYYYMM
		data = data[:len(Layout_YYYYMM)]
	}
	now, _ := time.ParseInLocation(layout, string(data), time.Local)
	if now.UnixNano() < 0 || now.Unix() <= 0 {
		*t = Time{}
	} else {
		*t = Time{now}
	}
	return nil
}

// deserialization.
func (d *Time) UnmarshalText(text []byte) error {
	*d = ToTime(string(text))
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization.
func (d Time) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

// Value insert timestamp into mysql need this function.
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *Time) Scan(v interface{}) error {
	*t = ToTime(v)
	return nil
}
func (d Time) Valid() bool {
	return !d.IsZero()
}
