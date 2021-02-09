package di

import (
	"go.uber.org/dig"
)

var Global = New()

func SetGlobal(container *dig.Container) {
	if Global == nil {
		Global = container
	}
}
func New(opts ...dig.Option) *dig.Container {
	container := dig.New(opts...)
	return container
}
