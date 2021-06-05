package table

import (
	"fmt"
	"strings"
)

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
	isText := strings.Contains(dataType, "text")
	if isChar || isText {
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
			where = where + table.Name + "." + key

			equator := " = "
			dataType := table.RawColumns[key].Type
			isString := isStringDataType(dataType)
			if isString {
				equator = " LIKE "
			}
			where = where + equator

			where = where + "?"
			if idx < (len(vals) - 1) {
				where = " " + where + " OR "
			}
			value := val
			if isString {
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
