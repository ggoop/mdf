package md

import (
	"github.com/ggoop/mdf/gin"
	"github.com/ggoop/mdf/utils"
	"net/http"
)

type ResContext struct {
	data  utils.Map
	error error
}

func NewResContext() *ResContext {
	return &ResContext{}
}

func (s ResContext) Data() utils.Map {
	if s.data == nil {
		s.data = utils.Map{}
	}
	return s.data
}
func (s ResContext) Has(name string) bool {
	if _, ok := s.data[name]; ok {
		return true
	}
	return false
}
func (s *ResContext) Set(name string, value interface{}) *ResContext {
	if s.data == nil {
		s.data = utils.Map{}
	}
	s.data[name] = value
	return s
}

func (s *ResContext) SetError(err interface{}) *ResContext {
	s.error = utils.ToError(err)
	return s
}
func (s ResContext) Error() error {
	return s.error
}

func (s ResContext) Bind(c *gin.Context) {
	if s.error != nil {
		if _, ok := s.data["code"]; !ok {
			s.Set("code", http.StatusBadRequest)
		}
		s.Set("msg", s.error.Error())

		c.JSON(http.StatusBadRequest, s.data)
	} else {
		if _, ok := s.data["code"]; !ok {
			s.Set("code", 200)
		}
		c.JSON(http.StatusOK, s.data)
	}
}
func (s *ResContext) Adjust(fn func(res *ResContext)) *ResContext {
	fn(s)
	return s
}
