package utils

import (
	"fmt"
)

// Model base model definition, including fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, which could be embedded in your models
//    type User struct {
//      gorm.Model
//    }
type ORMModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt Time
	UpdatedAt Time
}

const (
	//mssql
	ORM_DRIVER_MSSQL = "mssql"
	//mysql
	ORM_DRIVER_MYSQL = "mysql"
)

type SqlStruct struct {
	SelectSql string
	FromSql   string
	JoinSql   string
	GroupSql  string
	HavingSql string
	OrderSql  string
	LimitSql  string
	WhereSql  string
	LimitNum  int64
	OffsetNum int64
}

func (sql SqlStruct) CombinedConditionSql() string {
	return sql.JoinSql + sql.WhereSql + sql.GroupSql + sql.HavingSql + sql.OrderSql + sql.LimitSql
}

func (sql SqlStruct) AddWhere(where string) {
	if sql.WhereSql == "" {
		sql.WhereSql = "where"
	}
	sql.WhereSql = fmt.Sprintf("%v and %v", sql.WhereSql, where)
}
func (sql SqlStruct) CombinedSql() string {
	return fmt.Sprintf("SELECT %v FROM %v %v", sql.SelectSql, sql.FromSql, sql.CombinedConditionSql())
}
