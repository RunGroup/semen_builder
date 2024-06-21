package semen_builder

import (
	"fmt"
	"reflect"
	"time"
)

type Value struct {
	Field string
	Arg   any
	Fn    string
}

type Insert struct {
	Table        string
	Values       []Value
	StructValues any
	Timestamps   bool
}

func (i Insert) ToSql() (sql string, args []any) {
	sql = "INSERT INTO " + i.Table
	args = []any{}

	allValues := []Value{}

	if i.StructValues != nil {
		t := reflect.TypeOf(i.StructValues)
		fields := reflect.ValueOf(i.StructValues)
		fieldCount := t.NumField()

		for index := range fieldCount {
			if t.Field(index).Tag.Get("sb") == "skip" {
				continue
			}

			allValues = append(allValues, Value{
				Field: fmt.Sprintf("%s", t.Field(index).Tag.Get("db")),
				Arg:   fields.Field(index).Interface(),
			})
		}
	}

	allValues = append(allValues, i.Values...)

	if i.Timestamps {
		timestamp := time.Now().UTC()

		allValues = append(allValues, Value{Field: "created_at", Arg: timestamp})
		allValues = append(allValues, Value{Field: "updated_at", Arg: timestamp})
	}

	sql += " ("
	for index, value := range allValues {
		sql += value.Field
		if index != len(allValues)-1 {
			sql += ","
		}
	}
	sql += ") VALUES "

	sql += "("
	for index, value := range allValues {
		if value.Arg != nil {
			sql += "?"
			args = append(args, value.Arg)
		} else if value.Fn != "" {
			sql += value.Fn
		}

		if index != len(allValues)-1 {
			sql += ","
		}
	}
	sql += ")"

	return sql, args
}
