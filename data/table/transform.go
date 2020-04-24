package table

import (
	"bytes"
	"database/sql"
	"strings"
)

func unpack(raw string) (tableName string, field string) {
	sep := strings.Split(raw[:], ".")
	field = sep[0]
	if len(sep) > 1 {
		tableName = sep[0]
		field = sep[1]
	}
	return tableName, field
}

func (table *Table) transformRow(row []interface{}, fields []string) map[string]interface{} {
	// the stored row has must be mapped to the column names
	// this is rather unassuring, will the indices always be correct?
	if len(fields) == 0 { fields = table.ColumnNames }
	mapped := make(map[string]interface{})
	for i := 0; i < len(fields); i++ {
		name := fields[i]
		tableName, field := unpack(fields[i])
		if tableName == table.Name { name = field }
		mapped[name] = row[i]
	}

	// the mapped row still has array buffers in it, let's handle them
	entry := make(map[string]interface{})
	for key, raw := range mapped {
		// get the raw type
		casted := *(raw).(*interface{})

		// you can print out the type as follow
		// rawType := fmt.Sprintf("%T", casted)
		// fmt.Println(rawType)

		// handle the data based on its type
		switch casted.(type) {
		case []uint8:
			// for now, let's turn them all into string
			arByte, ok := casted.([]uint8)
			if ok {
				str := bytes.NewBuffer(arByte).String()
				entry[key] = str
			} else {
				entry[key] = casted
			}
		default:
			entry[key] = casted
		}
	}

	return entry
}

func (table *Table) transform(rows *sql.Rows, fields []string) []map[string]interface{} {
	// create slice of map to hold the data
	data := []map[string]interface{}{}
	if len(fields) == 0 { fields = table.ColumnNames }

	for rows.Next() {
		// create slice of interfaces to temporarily store the row
		row := make([]interface{}, 0, len(fields))
		for i := 0; i < len(fields); i++ {
			var container interface{}
			container = struct{}{}
			row = append(row, &container)
		}
		rows.Scan(row...)
		entry := table.transformRow(row, fields)

		data = append(data, entry)
	}

	return data
}
