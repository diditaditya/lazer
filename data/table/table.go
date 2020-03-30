package table

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"lazer/laze"
	exception "lazer/error"
)

type Table struct {
	Name string
	Conn *gorm.DB
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

func (table *Table) GetPkColumn() {
	for k,v := range table.RawColumns {
		if v.Key == "PRI" {
			table.Pk = k
			break
		}
	}
}

func (table *Table) FindAll(params map[string][]string) ([]map[string]interface{}, error) {

	rawQuery := "SELECT * FROM "
	rawQuery = rawQuery + table.Name

	filter := table.getFilter(params)
	where, values := table.createWhereStringFromFilter(filter)

	rawQuery = rawQuery + where

	rows, err := table.Conn.Raw(rawQuery, values...).Rows()

	defer rows.Close()

	if err != nil {
		fmt.Println("[table] error fetching ", table.Name)
		fmt.Println(err)
		return nil, err
	}

	data := table.transform(rows)
	return data, nil
}

func (table *Table) FindByPk(value string) map[string]interface{} {
	condition := fmt.Sprintf("%s = ?", table.Pk)
	rows, err := table.Conn.Table(table.Name).Where(condition, value).Rows()

	defer rows.Close()

	if err != nil {
		fmt.Println("[DB] error fetching ", table.Name)
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

func (table *Table) Delete(params map[string][]string) laze.Exception {
	// safe guard by requiring params
	if len(params) == 0 {
		ex := exception.New(exception.BADREQUEST, "parameters are required")
		return ex
	}

	rawQuery := "DELETE FROM "
	rawQuery = rawQuery + table.Name

	filter := table.getFilter(params)
	where, values := table.createWhereStringFromFilter(filter)

	rawQuery = rawQuery + where

	rows, err := table.Conn.Raw(rawQuery, values...).Rows()

	defer rows.Close()

	if err != nil {
		fmt.Println("[table] error deleting from", table.Name)
		fmt.Println(err)
		ex := exception.FromError(err, exception.INTERNALERROR)
		return ex
	}

	return nil
}