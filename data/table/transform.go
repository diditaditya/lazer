package table

import (
	"database/sql"
	"strings"

	"lazer/data/trait"
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

func castToInterface(raw any) any {
	// get the raw type
	ptr := raw.(*any)
	casted := *ptr

	// you can print out the type as follow
	// rawType := fmt.Sprintf("%T", casted)
	// fmt.Println(rawType)

	// handle the data based on its type
	switch val := casted.(type) {
	case []byte:
		// for now, let's turn them all into string
		return string(val)
	default:
		return casted
	}
}

func (table *Table) transformRow(row []any, tableName string, fields []string) map[string]any {
	// the stored row has must be mapped to the column names
	// this is rather unassuring, will the indices always be correct?
	if len(fields) == 0 {
		fields = table.ColumnNames
	}
	mapped := make(map[string]any)
	for i := 0; i < len(fields); i++ {
		name := fields[i]
		mapped[name] = castToInterface(row[i])
	}
	return mapped
}

func (table *Table) transform(rows *sql.Rows, fields []string, include trait.Joined) []map[string]any {
	// create slice of map to hold the data
	data := []map[string]any{}
	if len(fields) == 0 {
		fields = table.ColumnNames
	}

	for rows.Next() {
		// create slice of interfaces to temporarily store the row
		row := make([]any, 0, len(fields))
		for i := 0; i < len(fields); i++ {
			row = append(row, new(any))
		}
		rows.Scan(row...)
		entry := table.transformRow(row, table.Name, fields)

		data = append(data, entry)
	}

	result := transformIncludes(data, table.Name, include)

	return result
}

func clearCurrentTableName(row map[string]any, tableName string) (map[string]any, map[string]any) {
	data := map[string]any{}
	joined := map[string]any{}
	for key, val := range row {
		table, col := unpack(key)
		field := key
		if table == tableName {
			field = col
			data[field] = val
		} else {
			joined[field] = val
		}
	}

	isAllNull := true
	for _, val := range joined {
		if val != nil {
			isAllNull = false
			break
		}
	}

	if isAllNull {
		var nulled map[string]any
		joined = nulled
	}
	return data, joined
}

type collection struct {
	order  int
	data   map[string]any
	joined []map[string]any
}

func transformIncludes(raw []map[string]any, tableName string, include trait.Joined) (result []map[string]any) {
	if include == nil {
		for _, row := range raw {
			data, _ := clearCurrentTableName(row, tableName)
			result = append(result, data)
		}
		return result
	}

	incJoined := include.GetJoined()

	if len(incJoined) == 0 {
		for _, row := range raw {
			data, _ := clearCurrentTableName(row, tableName)
			result = append(result, data)
		}
		return result
	}

	collected := map[string]*collection{}
	pk := include.GetTablePk()

	// collect the duplicates which indicates joined tables
	counter := 0
	for _, row := range raw {
		for key, val := range row {
			table, field := unpack(key)
			isPkField := table == tableName && field == pk
			if isPkField {
				pkVal, ok := val.(string)
				if ok {
					_, pkFound := collected[pkVal]
					if pkFound {
						_, joined := clearCurrentTableName(row, tableName)
						collected[pkVal].joined = append(collected[pkVal].joined, joined)
					} else {
						base, joined := clearCurrentTableName(row, tableName)
						collected[pkVal] = &collection{
							order:  counter,
							data:   base,
							joined: []map[string]any{},
						}
						counter = counter + 1
						if joined != nil {
							collected[pkVal].joined = append(collected[pkVal].joined, joined)
						}
					}
					break
				}
			}
		}
	}

	// reorder the collection
	ordered := make([]*collection, len(collected))
	for _, coll := range collected {
		ordered[coll.order] = coll
	}

	// check the collected rows, recursively process if joined
	for _, coll := range ordered {
		for _, joined := range incJoined {
			joinedTableName := joined.GetTableName()
			if len(coll.joined) > 0 {
				rawJoined := coll.joined
				refType := joined.GetReferenceType()
				if refType == "hasMany" {
					rawJoined = []map[string]any{coll.joined[0]}
				}
				coll.data[joinedTableName] = transformIncludes(rawJoined, joinedTableName, joined)
			} else {
				var nulled map[string]any
				coll.data[joinedTableName] = nulled
			}
		}
		result = append(result, coll.data)
	}

	return result
}
