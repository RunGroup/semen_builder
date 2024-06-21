package semen_builder

import "fmt"

type On struct {
	Fields   []string
	FieldArg []any
	mode     string
	on       []On
}

func (o *On) And(on On) *On {
	on.mode = "AND"
	o.on = append(o.on, on)

	return o
}

func (o *On) Or(on On) *On {
	on.mode = "OR"
	o.on = append(o.on, on)

	return o
}

func (o *On) toSql() (sql string, args []any) {
	sql = ""
	args = []any{}

	if o.mode != "" {
		sql += " " + o.mode + " "
	}

	if len(o.Fields) > 0 {
		sql += fmt.Sprintf("%s = %s", o.Fields[0], o.Fields[1])
	} else if len(o.FieldArg) > 0 {
		sql += fmt.Sprintf("%s = ?", o.FieldArg[0])
		args = append(args, o.FieldArg[1])
	}

	for i, on := range o.on {
		if len(on.on) > 0 {
			sql += "("
		}
		nestedSql, nestedArgs := on.toSql()
		args = append(args, nestedArgs...)

		if i != 0 {
			sql += " " + on.mode + " "
		}

		sql += "" + nestedSql

		if len(on.on) > 0 {
			sql += ")"
		}
	}

	return sql, args
}

type Join struct {
	Table string
	As    string
	On    *On
	Mode  string
}

func (j Join) toSql() (sql string, args []any) {
	if j.Mode == "" {
		j.Mode = Inner
	}
	onSql, args := j.On.toSql()
	sql = fmt.Sprintf("%s JOIN %s as %s ON %s", j.Mode, j.Table, j.As, onSql)

	return sql, args
}
