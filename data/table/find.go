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
	rawQuery := "SELECT * FROM "
	rawQuery = rawQuery + table.Name

	filter := table.getFilter(params)
	where, values := table.createWhereStringFromFilter(filter)

	fieldMarks, fields := table.getFields(include)
	fmt.Printf("field marks: %v\n", fieldMarks)
	fmt.Printf("fields: %v\n", fields)

	joinMarks, joined := table.getJoined(include)
	fmt.Printf("join marks: %v\n", joinMarks)
	fmt.Printf("joined: %v\n", joined)

	total, countErr := table.count(where, values)
	if countErr != nil {
		return nil, nil, countErr
	}
	pagination, page, limit, _ := table.getPagination(params)

	rawQuery = rawQuery + where + pagination

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

	data := table.transform(rows)
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

	data := table.transform(rows)

	if len(data) > 0 {
		return data[0]
	}

	return map[string]interface{}{}
}
