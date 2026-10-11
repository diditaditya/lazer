package table

import (
	"fmt"
	"strings"
)

func (table *Table) Create(data map[string]any) map[string]any {
	// prepare the query
	keys := []string{}
	vals := []any{}

	for key, val := range data {
		keys = append(keys, key)
		vals = append(vals, val)
	}

	fields := strings.Join(keys[:], ",")
	marks := []string{}
	for i := 0; i < len(vals); i++ {
		marks = append(marks, "?")
	}
	values := strings.Join(marks[:], ",")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table.Name, fields, values)

	// insert the data
	raw := table.Conn.Raw(query, vals...).Row()

	// create slice of interfaces to temporarily store the row
	row := make([]any, 0, len(table.RawColumns))
	for i := 0; i < len(table.RawColumns); i++ {
		var container any
		container = struct{}{}
		row = append(row, &container)
	}
	raw.Scan(row...)

	entry := table.transformRow(row, table.Name, table.ColumnNames)

	return entry
}
