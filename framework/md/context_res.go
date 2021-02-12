package md

import "github.com/ggoop/mdf/utils"

type ResContext struct {
	data  utils.Map
	error error
}

func (s ResContext) New() ResContext {
	return ResContext{}
}

func (s ResContext) Data() utils.Map {
	if s.data == nil {
		s.data = utils.Map{}
	}
	return s.data
}
func (s *ResContext) Set(name string, value interface{}) *ResContext {
	if s.data == nil {
		s.data = utils.Map{}
	}
	s.data[name] = value
	return s
}

func (s ResContext) SetError(err error) {
	s.error = err
}
func (s ResContext) Error() error {
	return s.error
}
