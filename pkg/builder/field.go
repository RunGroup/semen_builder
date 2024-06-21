package builder

import "fmt"

type Field struct {
	Name       string
	TableAlias string
	As         string
}

func (f Field) toSql(tableAlias string) string {
	tAs := ""
	if f.TableAlias != "" {
		tAs = f.TableAlias
	} else if tableAlias != "" {
		tAs = tableAlias
	}

	sql := ""

	if tAs != "" {
		sql = fmt.Sprintf("%s.%s", tAs, f.Name)
	} else {
		sql = fmt.Sprintf("%s", f.Name)
	}

	if f.As != "" {
		sql += fmt.Sprintf(" AS %s", f.As)
	}

	return sql
}

type RawField struct {
	Sql string
}
