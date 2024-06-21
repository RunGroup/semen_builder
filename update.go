package semen_builder

import (
	"reflect"
	"strings"
)

type SetValue struct {
	Field string
	Value any
}

type Increment struct {
	Value any
}

type Decrement struct {
	Value any
}

func (v SetValue) toSql() (string, any) {
	var arg any
	sql := ""

	if reflect.TypeOf(v.Value).String() == reflect.TypeOf(Increment{}).String() {
		sql += " " + v.Field + " = " + v.Field + " + ?,"
		arg = v.Value.(Increment).Value
	} else if reflect.TypeOf(v.Value).String() == reflect.TypeOf(Decrement{}).String() {
		sql += " " + v.Field + " = " + v.Field + " - ?,"
		arg = v.Value.(Decrement).Value
	} else {
		sql += " " + v.Field + " = ?,"
		arg = v.Value
	}

	return sql, arg
}

type Update struct {
	Table  string
	Values []SetValue
	Where  *Condition
}

func (u Update) ToSql() (sql string, args []any) {
	sql = "UPDATE " + u.Table + " SET"
	args = []any{}

	for _, value := range u.Values {
		valueSql, arg := value.toSql()
		sql += valueSql
		args = append(args, arg)
	}

	sql = strings.TrimRight(sql, ",")

	if u.Where != nil {
		conditionsSql, conditionsArgs := u.Where.toSql()
		args = append(args, conditionsArgs...)
		sql += " WHERE " + conditionsSql
	}

	return sql, args
}
