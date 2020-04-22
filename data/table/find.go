package table

import (
	"fmt"
	"math"

	"lazer/laze"
	exception "lazer/error"
)

func (table *Table) count(where string, values []interface{}) (int, laze.Exception) {
	query := fmt.Sprintf("SELECT COUNT(%s) FROM %s %s", table.Pk, table.Name, where)
	counter, err := table.Conn.Raw(query, values...).Rows()
	if err != nil {
		ex := exception.FromError(err, exception.INTERNALERROR)
		return 0, ex
	}
	count := 0
	if counter.Next() {
		counter.Scan(&count)
	}
	return count, nil
}

func (table *Table) FindAll(params map[string][]string, include map[string]interface{}) ([]map[string]interface{}, map[string]interface{}, laze.Exception) {
	values := []interface{}{}

	fieldMarks, fields := table.getFields(include)
	
	if len(fields) == 0 {
		fmt.Println("fields length = 0")
		for i, field := range table.ColumnNames {
			fields = append(fields, field)
			fieldMarks = fieldMarks + field
			if i < len(table.ColumnNames) - 1 {
				fieldMarks = fieldMarks + ", "
			}
		}
	}

	rawQuery := "SELECT " + fieldMarks + " FROM "
	rawQuery = rawQuery + table.Name

	joinMarks, joined := table.getJoined(include)

	if len(joined) > 0 {
		rawQuery = rawQuery + " " + joinMarks
	}

	filter := table.getFilter(params)
	whereMarks, wheres := table.createWhereStringFromFilter(filter)
	values = append(values, wheres...)

	total, countErr := table.count(whereMarks, wheres)
	if countErr != nil {
		return nil, nil, countErr
	}
	pagination, page, limit, _ := table.getPagination(params)

	rawQuery = rawQuery + whereMarks + pagination

	pages := math.Ceil(float64(total) / float64(limit))

	rows, err := table.Conn.Raw(rawQuery, values...).Rows()

	defer rows.Close()

	if err != nil {
		ex := exception.FromError(err, exception.INTERNALERROR)
		return nil, nil, ex
	}

	meta := map[string]interface{}{
		"page": page,
		"pageSize": limit,
		"pages": pages,
		"total": total,
	}

	data := table.transform(rows, fields)
	return data, meta, nil
}

func (table *Table) FindByPk(value string) map[string]interface{} {
	condition := fmt.Sprintf("%s = ?", table.Pk)
	rows, err := table.Conn.Table(table.Name).Where(condition, value).Rows()

	defer rows.Close()

	if err != nil {
		fmt.Println("[DB] error fetching ", table.Name)
		fmt.Println(err)
	}

	data := table.transform(rows, table.ColumnNames)

	if len(data) > 0 {
		return data[0]
	}

	return map[string]interface{}{}
}
