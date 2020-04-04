package table

import (
	"fmt"
)

func (table *Table) FindAll(params map[string][]string) ([]map[string]interface{}, error) {
	rawQuery := "SELECT * FROM "
	rawQuery = rawQuery + table.Name

	filter := table.getFilter(params)
	where, values := table.createWhereStringFromFilter(filter)
	pagination := table.createPaginationString(params)

	rawQuery = rawQuery + where + pagination

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
