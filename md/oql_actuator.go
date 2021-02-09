package md

import (
	"github.com/ggoop/mdf/utils"
)

type OQLActuator interface {
	GetName() string
	Count(oql *OQL, value interface{}) *OQL
	Pluck(oql *OQL, column string, value interface{}) *OQL
	Take(oql *OQL, out interface{}) *OQL
	Find(oql *OQL, out interface{}) *OQL
	Create(oql *OQL, data interface{}) *OQL
	Update(oql *OQL, data interface{}) *OQL
	Delete(oql *OQL) *OQL
}

var oqlActuatorMap = map[string]OQLActuator{}

func RegisterOQLActuator(name string, query OQLActuator) {
	oqlActuatorMap[name] = query
}
func GetOQLActuator(names ...string) OQLActuator {
	if names == nil || len(names) == 0 {
		return oqlActuatorMap[utils.DefaultConfig.Db.Driver]
	}
	return oqlActuatorMap[names[0]]
}
func init() {
	aa := &commonActuator{}
	RegisterOQLActuator(aa.GetName(), aa)
}

//公共查询
type commonActuator struct {
}

func (commonActuator) GetName() string {
	return "common"
}
func (s *commonActuator) Count(oql *OQL, value interface{}) *OQL {
	oql.Select(false).Select("count(*)")
	return oql
}
func (s *commonActuator) Pluck(oql *OQL, column string, value interface{}) *OQL {
	return oql
}
func (s *commonActuator) Take(oql *OQL, out interface{}) *OQL {
	return oql
}
func (s *commonActuator) Find(oql *OQL, out interface{}) *OQL {
	return oql
}
func (s *commonActuator) Create(oql *OQL, data interface{}) *OQL {
	return oql
}
func (s *commonActuator) Update(oql *OQL, data interface{}) *OQL {
	return oql
}
func (s *commonActuator) Delete(oql *OQL) *OQL {
	return oql
}
