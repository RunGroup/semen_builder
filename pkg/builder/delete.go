package builder

import (
	"fmt"
)

type Delete struct {
	From  string
	As    string
	Joins []Join
	Where *Condition
}

func (d Delete) ToSql() (sql string, args []any) {
	sql = fmt.Sprintf("DELETE FROM %s", d.From)
	args = []any{}

	if len(d.Joins) > 0 {
		for _, join := range d.Joins {
			joinSql, joinArgs := join.toSql()
			sql += fmt.Sprintf(" %s", joinSql)
			args = append(args, joinArgs...)
		}
	}

	conditionsSql, conditionsArgs := d.Where.toSql()
	args = append(args, conditionsArgs...)
	if conditionsSql != "" {
		sql += " WHERE " + conditionsSql
	}

	return sql, args
}
