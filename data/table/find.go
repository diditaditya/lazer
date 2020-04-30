package table

import (
	"fmt"
	"math"
	"strings"

	exception "lazer/error"
	"lazer/laze"
	"lazer/data/trait"
)

func (table *Table) count(where string, values []interface{}) (int, laze.Exception) {
	query := fmt.Sprintf("SELECT COUNT(DISTINCT(%s)) FROM %s %s", table.Pk, table.Name, where)
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

func (table *Table) getOrderBy(params map[string][]string) (orderBy string) {
	mapper := map[string]string{
		"a":    "ASC",
		"asc":  "ASC",
		"ASC":  "ASC",
		"d":    "DESC",
		"dsc":  "DESC",
		"desc": "DESC",
		"DESC": "DESC",
	}

	orders := []string{}
	if rawOrders, ok := params["sort"]; ok {
		fmt.Printf("rawOrders %v\n", rawOrders);
		for _, raw := range rawOrders {
			ordered := strings.Split(raw, ",")
			if table.isField(ordered[0]) {
				order := ordered[0]
				if len(ordered) > 1 {
					if sortType, found := mapper[ordered[1]]; found {
						order = order + " " + sortType
					}
				}
				orders = append(orders, order)
			}
		}
	}

	orderBy = strings.Join(orders[:], ", ")
	if len(orderBy) > 0 {
		orderBy = " ORDER BY " + orderBy
	}
	return orderBy
}

func (table *Table) FindAll(params map[string][]string, include trait.Joined) ([]map[string]interface{}, map[string]interface{}, laze.Exception) {
	fieldMarks, fields := table.getFields(include)
	orderBy := table.getOrderBy(params)

	rawQuery := "SELECT " + fieldMarks + " FROM "
	rawQuery = rawQuery + table.Name

	joinMarks, joined := table.getJoined(include)

	if len(joined) > 0 {
		rawQuery = rawQuery + " " + joinMarks
	}

	filter := table.getFilter(params)
	whereMarks, wheres := table.createWhereStringFromFilter(filter)

	total, countErr := table.count(whereMarks, wheres)
	if countErr != nil {
		return nil, nil, countErr
	}
	pagination, page, limit, _ := table.getPagination(params)

	rawQuery = rawQuery + whereMarks + orderBy + pagination

	pages := math.Ceil(float64(total) / float64(limit))

	rows, err := table.Conn.Raw(rawQuery, wheres...).Rows()

	defer rows.Close()

	if err != nil {
		ex := exception.FromError(err, exception.INTERNALERROR)
		return nil, nil, ex
	}

	meta := map[string]interface{}{
		"page":     page,
		"pageSize": limit,
		"pages":    pages,
		"total":    total,
	}

	data := table.transform(rows, fields, include)
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

	data := table.transform(rows, table.ColumnNames, nil)

	if len(data) > 0 {
		return data[0]
	}

	return map[string]interface{}{}
}
