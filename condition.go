package semen_builder

import (
	"fmt"
	"strings"
)

type Condition struct {
	Field      string
	Operator   string
	Args       []any
	Arg        any
	mode       string
	conditions []Condition
}

func (c *Condition) And(condition Condition) *Condition {
	condition.mode = "AND"
	c.conditions = append(c.conditions, condition)

	return c
}

func (c *Condition) Or(condition Condition) *Condition {
	condition.mode = "OR"
	c.conditions = append(c.conditions, condition)

	return c
}

func (c *Condition) toSql() (sql string, args []any) {
	sql = ""
	args = []any{}

	if c.Field != "" {
		if len(c.Args) > 0 {
			argsSql := strings.TrimRight(strings.Repeat("?,", len(c.Args)), ",")
			sql += fmt.Sprintf("%s %s (%s)", c.Field, c.Operator, argsSql)

			args = append(args, c.Args...)
		} else if c.Arg != nil {
			sql += fmt.Sprintf("%s %s ?", c.Field, c.Operator)
			args = append(args, c.Arg)
		} else {
			sql += fmt.Sprintf("%s %s", c.Field, c.Operator)
		}
	}

	if len(c.conditions) > 0 {
		sql += " " + c.conditions[0].mode + " "
	}

	for i, condition := range c.conditions {
		if len(condition.conditions) > 0 {
			sql += "("
		}
		nestedSql, nestedArgs := condition.toSql()

		if i != 0 {
			sql += " " + condition.mode + " "
		}

		args = append(args, nestedArgs...)
		sql += "" + nestedSql

		if len(condition.conditions) > 0 {
			sql += ")"
		}
	}

	return sql, args
}

func AppendOr(where *Condition, condition Condition) *Condition {
	if where == nil {
		return &condition
	} else {
		where.Or(condition)
		return where
	}
}

func AppendAnd(where *Condition, condition Condition) *Condition {
	if where == nil {
		return &condition
	} else {
		where.And(condition)
		return where
	}
}
