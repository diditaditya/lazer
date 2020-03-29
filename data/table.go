package data

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Table struct {
	name string
	conn *gorm.DB
	ColumnNames []string
	RawColumns map[string]RawColumn
	Pk string
}

type RawColumn struct {
	Field string
	Type string
	Null string
	Key string
	Default string
	Extra string
}

func (table *Table) getPkColumn() {
	for k,v := range table.RawColumns {
		if v.Key == "PRI" {
			table.Pk = k
			break
		}
	}
}

func (table *Table) transformRow(row []interface{}) map[string]interface{} {
	// the stored row has must be mapped to the column names
	// this is rather unassuring, will the indices always be correct?
	mapped := make(map[string]interface{})
	for i := 0; i < len(table.ColumnNames); i++ {
		name := table.ColumnNames[i]
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

func (table *Table) transform(rows *sql.Rows) []map[string]interface{} {
	// create slice of map to hold the data
	data := []map[string]interface{}{}

	for rows.Next() {
		// create slice of interfaces to temporarily store the row
		row := make([]interface{}, 0, len(table.RawColumns))
		for i := 0; i < len(table.RawColumns); i++ {
			var container interface{}
			container = struct{}{}
			row = append(row, &container)
		}
		rows.Scan(row...)
		entry := table.transformRow(row)

		data = append(data, entry)
	}

	return data
}

func (table *Table) getFilter(params map[string][]string) map[string][]string {
	filter := make(map[string][]string)
	if len(params) > 0 {
		for key, val := range params {
			if _, ok := table.RawColumns[key]; ok {
				filter[key] = val
			}
		}
	}
	return filter
}

func isStringDataType(dataType string) bool {
	isChar := strings.Contains(dataType, "char")
	isTest := strings.Contains(dataType, "text")
	if isChar || isTest {
		return true
	}
	return false
}

func (table *Table) createWhereStringFromFilter(filter map[string][]string) (string, []interface{}) {
	if len(filter) == 0 {
		return "", []interface{}{}
	}

	where := " WHERE "
	values := []interface{}{}
	counter := 0
	for key, vals := range filter {
		for idx, val := range vals {
			where = where + key

			equator := " = "
			dataType := table.RawColumns[key].Type
			isString := isStringDataType(dataType)
			if (isString) {
				equator = " LIKE "
			}
			where = where + equator

			where = where + "?"
			if idx < (len(vals) - 1) {
				where = " " + where + " OR "
			}
			value := val
			if (isString) {
				value = fmt.Sprintf("%%%s%%", val)
			}
			values = append(values, value)
		}
		
		if counter < (len(filter) - 1) {
			where = where + " AND "
			counter = counter + 1
		}
	}
	return where, values
}

func (table *Table) FindAll(params map[string][]string) ([]map[string]interface{}, error) {

	rawQuery := "SELECT * FROM "
	rawQuery = rawQuery + table.name

	filter := table.getFilter(params)
	where, values := table.createWhereStringFromFilter(filter)

	rawQuery = rawQuery + where

	rows, err := table.conn.Raw(rawQuery, values...).Rows()

	defer rows.Close()

	if err != nil {
		fmt.Println("[table] error fetching ", table.name)
		fmt.Println(err)
		return nil, err
	}

	data := table.transform(rows)
	return data, nil
}

func (table *Table) FindByPk(value string) map[string]interface{} {
	condition := fmt.Sprintf("%s = ?", table.Pk)
	rows, err := table.conn.Table(table.name).Where(condition, value).Rows()

	defer rows.Close()

	if err != nil {
		fmt.Println("[DB] error fetching ", table.name)
		fmt.Println(err)
	}

	data := table.transform(rows)

	if len(data) > 0 {
		return data[0]
	}

	return map[string]interface{}{}
}

func (table *Table) Create(data map[string]interface{}) map[string]interface{} {
	// prepare the query
	keys := []string{}
	vals := []interface{}{}

	for key, val := range data {
		keys = append(keys, key)
		// value, ok := val.(string)
		// if ok {
		// 	vals = append(vals, value)
		// }
		vals = append(vals, val)
	}

	fields := strings.Join(keys[:], ",")
	marks := []string{}
	for i := 0; i < len(vals); i++ {
		marks = append(marks, "?")
	}
	values := strings.Join(marks[:], ",")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table.name, fields, values)

	// insert the data
	raw := table.conn.Raw(query, vals...).Row()

	// create slice of interfaces to temporarily store the row
	row := make([]interface{}, 0, len(table.RawColumns))
	for i := 0; i < len(table.RawColumns); i++ {
		var container interface{}
		container = struct{}{}
		row = append(row, &container)
	}
	raw.Scan(row...)

	entry := table.transformRow(row)

	return entry
}