package md

import (
	"github.com/ggoop/mdf/utils"
)

func init() {
	aa := &mysqlActuator{}
	RegisterOQLActuator(aa.GetName(), aa)
}

//公共查询
type mysqlActuator struct {
	from   string
	offset int
	limit  int
}

func (mysqlActuator) GetName() string {
	return utils.ORM_DRIVER_MYSQL
}
func (s *mysqlActuator) Count(oql *OQL, value interface{}) *OQL {
	return oql
}
func (s *mysqlActuator) Pluck(oql *OQL, column string, value interface{}) *OQL {
	return oql
}
func (s *mysqlActuator) Take(oql *OQL, out interface{}) *OQL {
	return oql
}
func (s *mysqlActuator) Find(oql *OQL, out interface{}) *OQL {
	return oql
}
func (s *mysqlActuator) Create(oql *OQL, data interface{}) *OQL {
	return oql
}
func (s *mysqlActuator) Update(oql *OQL, data interface{}) *OQL {
	return oql
}
func (s *mysqlActuator) Delete(oql *OQL) *OQL {
	return oql
}
