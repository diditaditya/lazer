package table

import (
	"bytes"
	"fmt"
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

	return data
}

func transformIncludes(raw map[string]interface{}, tableName string, include trait.Joined) (result map[string]interface{}) {
	joinedTables := map[string]interface{}{}
	incJoined := include.GetJoined()
	for _, joinedTable := range incJoined {
		tableJoined := joinedTable.GetTableName()
		joinType := joinedTable.GetReferenceType()
		joinedTables[tableJoined] = joinType
	}

	// if the table is joined
	// 	go through each row, check for duplicates based on primary key
	//	collect the duplicates
	//	create new row
	//	set the fields of the current table as attributes of the new row
	//	check the joined tables reference type
	//	if the type is belongs to
	//		set the joined table name as attribute to new row, accepting slice of rows
	//		set the collected to the attribute, and run again from the beginning by using the collection

	fmt.Printf("joined: %v\n", joinedTables)
	return joinedTables
}
