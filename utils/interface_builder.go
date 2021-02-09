package utils

import "reflect"

type InterfaceBuilder struct {
	interfaces []interface{}
}

func NewInterfaceBuilder() *InterfaceBuilder {
	var interfaces InterfaceBuilder
	return &interfaces
}

func (builder *InterfaceBuilder) Append(arg interface{}) *InterfaceBuilder {
	builder.interfaces = append(builder.interfaces, arg)
	return builder
}

func (builder *InterfaceBuilder) Appends(args ...interface{}) *InterfaceBuilder {
	for i := range args {
		builder.interfaces = append(builder.interfaces, args[i])
	}
	return builder
}

func (builder *InterfaceBuilder) Clear() *InterfaceBuilder {
	var interfaces []interface{}
	builder.interfaces = interfaces
	return builder
}

func (builder *InterfaceBuilder) ToInterfaces() []interface{} {
	return builder.interfaces
}
func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}
