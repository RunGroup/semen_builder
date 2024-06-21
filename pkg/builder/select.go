package builder

import (
	"fmt"
	"reflect"
	"strings"
)

type Select struct {
	From         string
	As           string
	Fields       []any
	StructFields any
	Joins        []Join
	Where        *Condition
	MultiGroupBy []string
	GroupBy      string
	Order        Order
	Orders       []Order
	Limit        uint
	Offset       uint
}

func (s Select) ToSql() (sql string, args []any) {
	allFields := []string{}

	if s.StructFields != nil {
		fields := []string{}
		t := reflect.TypeOf(s.StructFields)
		fieldCount := t.NumField()

		for i := range fieldCount {
			if t.Field(i).Tag.Get("sb") == "skip" {
				continue
			}

			if s.As != "" {
				fields = append(fields, fmt.Sprintf("%s.%s", s.As, t.Field(i).Tag.Get("db")))
			} else {
				fields = append(fields, fmt.Sprintf("%s", t.Field(i).Tag.Get("db")))
			}
		}

		allFields = append(allFields, fields...)
	}

	for _, field := range s.Fields {
		fieldType := reflect.TypeOf(field).Name()

		if fieldType == reflect.String.String() {
			if s.As != "" {
				allFields = append(allFields, fmt.Sprintf("%s.%s", s.As, field.(string)))
			} else {
				allFields = append(allFields, field.(string))
			}
		} else if fieldType == reflect.TypeOf(Field{}).Name() {
			allFields = append(allFields, field.(Field).toSql(s.As))
		} else if fieldType == reflect.TypeOf(RawField{}).Name() {
			allFields = append(allFields, field.(RawField).Sql)
		}
	}

	fieldsSql := strings.Join(allFields, ", ")
	fieldsSql = strings.TrimRight(fieldsSql, ", ")

	sql = fmt.Sprintf("SELECT %s FROM %s", fieldsSql, s.From)

	if s.As != "" {
		sql = fmt.Sprintf("%s AS %s", sql, s.As)
	}
	args = []any{}

	if len(s.Joins) > 0 {
		for _, join := range s.Joins {
			joinSql, joinArgs := join.toSql()
			sql += fmt.Sprintf(" %s", joinSql)
			args = append(args, joinArgs...)
		}
	}

	if s.Where != nil {
		conditionsSql, conditionsArgs := s.Where.toSql()
		args = append(args, conditionsArgs...)
		if conditionsSql != "" {
			sql += " WHERE " + conditionsSql
		}
	}

	if s.GroupBy != "" {
		sql += " GROUP BY " + s.GroupBy
	} else if len(s.MultiGroupBy) > 0 {
		sql += " GROUP BY " + strings.Join(s.MultiGroupBy, ", ")
	}

	if s.Order.Direction != "" && s.Order.Field != "" {
		sql += fmt.Sprintf(" ORDER BY %s %s", s.Order.Field, s.Order.Direction)
	} else if len(s.Orders) > 0 {
		sql += " ORDER BY"
		for _, order := range s.Orders {
			sql += fmt.Sprintf(" %s %s", order.Field, order.Direction) + ","
		}
		sql = strings.TrimRight(sql, ",")
	}

	if s.Limit > 0 {
		sql += " LIMIT ?"
		args = append(args, s.Limit)
	}
	if s.Offset > 0 {
		sql += " OFFSET ?"
		args = append(args, s.Offset)
	}

	return sql, args
}
