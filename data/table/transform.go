package table

import (
	"bytes"
	// "fmt"
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

func castToInterface(raw interface{}) interface{} {
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
			return str
		} else {
			return casted
		}
	default:
		return casted
	}
}

func (table *Table) transformRow(row []interface{}, tableName string, fields []string) map[string]interface{} {
	// the stored row has must be mapped to the column names
	// this is rather unassuring, will the indices always be correct?
	if len(fields) == 0 { fields = table.ColumnNames }
	mapped := make(map[string]interface{})
	for i := 0; i < len(fields); i++ {
		name := fields[i]
		mapped[name] = castToInterface(row[i])
	}
	return mapped
}

func (table *Table) transform(rows *sql.Rows, fields []string, include trait.Joined) []map[string]interface{} {
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
		entry := table.transformRow(row, table.Name, fields)

		data = append(data, entry)
	}

	result := transformIncludes(data, table.Name, include)

	return result
}

func clearCurrentTableName(row map[string]interface{}, tableName string) (map[string]interface{}, map[string]interface{}) {
	data := map[string]interface{}{}
	joined := map[string]interface{}{}
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
		var nulled map[string]interface{}
		joined = nulled
	}
	return data, joined
}

type collection struct {
	data		map[string]interface{}
	joined	[]map[string]interface{}
}

func transformIncludes(raw []map[string]interface{}, tableName string, include trait.Joined) (result []map[string]interface{}) {
	incJoined := include.GetJoined()
	
	if include == nil || len(incJoined) == 0 {
		for _, row := range raw {
			data, _ := clearCurrentTableName(row, tableName)
			result = append(result, data)
		}
		return result
	}

	collected := map[string]*collection{}
	pk := include.GetTablePk()

	// collect the duplicates which indicates joined tables
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
							data: base,
							joined: []map[string]interface{}{},
						}
						if joined != nil {
							collected[pkVal].joined = append(collected[pkVal].joined, joined)
						}
					}
					break
				}
			}
		}
	}

	// check the collected rows, recursively process if joined
	for _, coll := range collected {
		for _, joined := range incJoined {
			joinedTableName := joined.GetTableName()
			if len(coll.joined) > 0 {
				rawJoined := coll.joined
				refType := joined.GetReferenceType()
				if refType == "hasMany" { rawJoined = []map[string]interface{}{ coll.joined[0] } }
				coll.data[joinedTableName] = transformIncludes(rawJoined, joinedTableName, joined)
			} else {
				var nulled map[string]interface{}
				coll.data[joinedTableName] = nulled
			}
		}
		result = append(result, coll.data)
	}

	return result
}
