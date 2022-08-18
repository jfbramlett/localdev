package etl

import (
	"fmt"
	"github.com/splice/platform/localdev/v2/localdev/internal/etl/datbase"
)

type Reference struct {
	Name         string `json:"name"`
	ReferenceCol string `json:"reference_col"`
	FkCol        string `json:"fk_col"`
}

type Query struct {
	Name      string     `json:"name"`
	Table     string     `json:"table"`
	Join      []string   `json:"join"`
	Where     []string   `json:"where"`
	Limit     int        `json:"limit"`
	Reference *Reference `json:"reference"`
}

func (q *Query) ToFKSQL(fkCol string) string {
	return q.toSQL(fkCol, true)
}

func (q *Query) ToRefSQL(refCol string, args []interface{}) string {
	return fmt.Sprintf("SELECT p.* FROM %s p WHERE p.%s in (%s)", q.Table, refCol, datbase.ValuesToParameters(args))
}

func (q *Query) ToSQL() string {
	return q.toSQL("*", true)
}

func (q *Query) toSQL(col string, includeLimit bool) string {
	stmt := fmt.Sprintf("SELECT p.%s FROM %s p", col, q.Table)
	if len(q.Join) > 0 {
		for _, j := range q.Join {
			stmt = fmt.Sprintf("%s JOIN %s", stmt, j)
		}
	}
	if len(q.Where) > 0 {
		stmt = fmt.Sprintf("%s WHERE ", stmt)
		for idx, w := range q.Where {
			if idx > 0 {
				stmt = fmt.Sprintf("%s AND ", stmt)
			}
			stmt = fmt.Sprintf("%s %s", stmt, w)
		}
	}
	if includeLimit && q.Limit > 0 {
		stmt = fmt.Sprintf("%s LIMIT %d", stmt, q.Limit)
	}

	return stmt
}

type Queries struct {
	Queries []*Query `json:"queries"`
}
