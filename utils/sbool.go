package utils

import (
	"bytes"
	"database/sql/driver"
)

// JSONTime format json time field by myself
type SBool struct {
	value string
	valid bool // Valid is true if Bool is not NULL
}

var SBool_True = SBool{value: "1", valid: true}

var SBool_False = SBool{value: "0", valid: true}

var SBool_Null = SBool{value: "", valid: false}

func ToSBool(val interface{}) SBool {
	s := SBool{}
	s.Parse(val)
	return s
}
func (t *SBool) Parse(val interface{}) SBool {
	t.value = t._ToBoolValue(val)
	t.valid = t.value != ""
	return *t
}
func (t SBool) String() string {
	return t.value
}

func (t SBool) IsTrue() bool {
	return t.value == "1"
}
func (t SBool) IsFalse() bool {
	return t.value == "0"
}
func (t SBool) Valid() bool {
	return t.valid
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t SBool) MarshalJSON() ([]byte, error) {
	if !t.valid {
		return []byte("null"), nil
	}
	if t.IsTrue() {
		return []byte("true"), nil
	}
	if t.IsFalse() {
		return []byte("false"), nil
	}
	return []byte(t.String()), nil
}
func (t *SBool) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		t.value = ""
		return nil
	}
	v := SBool{value: t._ToBoolValue(string(bytes.Trim(data, "\"")))}
	v.valid = v.value != ""
	*t = v
	return nil
}

func (t SBool) Not() SBool {
	if t.IsTrue() {
		return SBool_False
	}
	if t.IsFalse() {
		return SBool_True
	}
	return t
}
func (t SBool) Equal(v SBool) bool {
	return t.value == v.value
}

func (t SBool) NotEqual(v SBool) bool {
	return t.value != v.value
}

// Value implements the driver Valuer interface.
func (t SBool) Value() (driver.Value, error) {
	if !t.valid {
		return nil, nil
	}
	return t.value, nil
}

// Scan implements the Scanner interface.
func (t *SBool) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	t.value = t._ToBoolValue(v)
	t.valid = t.value != ""
	return nil
}
func (t SBool) _ToBoolValue(value interface{}) string {
	if value == nil {
		return ""
	}
	bValue := ""

	sb := ToString(value)
	if sb == "true" || sb == "1" || sb == "t" {
		bValue = "1"
	} else if sb == "false" || sb == "0" || sb == "f" {
		bValue = "0"
	}
	return bValue
}
func (t SBool) OrmDataType(driver string) string {
	return "varchar(2)"
}
