package table

import (
	"fmt"

	"lazer/laze"
	exception "lazer/error"
)

func (table *Table) recordExistsByPk(pkValue string) laze.Exception {
	condition := fmt.Sprintf("%s = ?", table.Pk)

	found, err := table.Conn.Table(table.Name).Where(condition, pkValue).Rows()

	if err != nil {
		ex := exception.FromError(err, exception.INTERNALERROR)
		return ex
	}

	if !found.Next() {
		ex := exception.New(exception.UNPROCESSABLE, "record not found")
		return ex
	}
	return nil
}

func (table *Table) createUpdateQuery(data map[string]interface{}, params map[string][]string) (string, []interface{}) {
	filter := table.getFilter(params)
	whereStr, filterValues := table.createWhereStringFromFilter(filter)

	values := []interface{}{}
	columnsStr := ""
	counter := 0
	for key, val := range data {
		if _, ok := table.RawColumns[key]; ok {
			if key != table.Pk {
				values = append(values, val)
				columnValue := fmt.Sprintf("%s = ?", key)
				columnsStr = columnsStr + columnValue
				if counter < len(data) - 1 {
					columnsStr = columnsStr + ", "
				}
				counter = counter + 1
			}
		}		
	}
	values = append(values, filterValues...)

	query := fmt.Sprintf("UPDATE %s SET %s %s", table.Name, columnsStr, whereStr)

	return query, values
}

func (table *Table) Update(params map[string][]string, data map[string]interface{}) laze.Exception {
	query, values := table.createUpdateQuery(data, params)

	rows, err := table.Conn.Raw(query, values...).Rows()
	defer rows.Close()

	if err != nil {
		ex := exception.FromError(err, exception.INTERNALERROR)
		return ex
	}

	return nil
}

func (table *Table) UpdateByPk(pk string, data map[string]interface{}) laze.Exception {

	err := table.recordExistsByPk(pk)
	if err != nil { return err }

	param := make(map[string][]string)
	param[table.Pk] = []string{pk}
	query, values := table.createUpdateQuery(data, param)

	rows, updateErr := table.Conn.Raw(query, values...).Rows()
	defer rows.Close()

	if updateErr != nil {
		ex := exception.FromError(updateErr, exception.INTERNALERROR)
		return ex
	}

	return nil
}