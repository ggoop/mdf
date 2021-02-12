package utils

import (
	"database/sql/driver"
	"encoding/json"
)

type BizObject struct {
	data map[string]interface{}
}

func NewBizObject() *BizObject {
	return &BizObject{data: make(map[string]interface{})}
}
func (s *BizObject) Copy() *BizObject {
	c := NewBizObject()
	if s.data != nil {
		for k, v := range s.data {
			c.data[k] = v
		}
	}
	return c
}

func (t BizObject) Data() map[string]interface{} {
	return t.data
}
func (s BizObject) Reset() {
	s.data = make(map[string]interface{})
}
func (s *BizObject) Set(name string, value interface{}) *BizObject {
	if s.data == nil {
		s.data = make(map[string]interface{})
	}
	name = SnakeString(name)
	s.data[name] = value
	return s
}

func (s BizObject) Get(name string) (value interface{}, exists bool) {
	if s.data == nil {
		s.data = make(map[string]interface{})
	}
	value, exists = s.data[name]
	if !exists {
		name = SnakeString(name)
		value, exists = s.data[name]
	}
	return
}
func (c BizObject) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}
func (c BizObject) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c BizObject) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *BizObject) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" || string(bytes) == "" {
		return nil
	}
	data := make(map[string]interface{})
	json.Unmarshal(bytes, &data)
	nc := BizObject{data: data}
	*d = nc
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (d BizObject) MarshalJSON() ([]byte, error) {
	if d.data == nil || len(d.data) == 0 {
		return []byte(""), nil
	}
	bytes, err := json.Marshal(d.data)
	return bytes, err
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *BizObject) Scan(value interface{}) error {
	if value == nil {
		d = NewBizObject()
		return nil
	}
	switch v := value.(type) {
	case string:
		nd := NewBizObject()
		json.Unmarshal([]byte(v), &nd)
		d = nd
		return nil
	}
	d = NewBizObject()
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (d BizObject) Value() (driver.Value, error) {
	data, err := d.MarshalJSON()
	return string(data), err
}
func (t BizObject) IsValid() bool {
	return t.data != nil && len(t.data) > 0
}
