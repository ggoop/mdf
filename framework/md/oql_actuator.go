package md

import (
	"github.com/ggoop/mdf/utils"
)

type OQLActuator interface {
	GetName() string
	Count(oql OQL, value interface{}) OQL
	Pluck(oql OQL, column string, value interface{}) OQL
	Take(oql OQL, out interface{}) OQL
	Find(oql OQL, out interface{}) OQL
	Create(oql OQL, data interface{}) OQL
	Update(oql OQL, data interface{}) OQL
	Delete(oql OQL) OQL
}

var oqlActuatorMap = make(map[string]OQLActuator)

func RegisterOQLActuator(query OQLActuator) {
	oqlActuatorMap[query.GetName()] = query
}
func GetOQLActuator(names ...string) OQLActuator {
	if names == nil || len(names) == 0 {
		return oqlActuatorMap[utils.DefaultConfig.Db.Driver]
	}
	return oqlActuatorMap[names[0]]
}
