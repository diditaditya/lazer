package table

import (
	"fmt"

	"lazer/laze"
	exception "lazer/error"
)

func (table *Table) recordExistsByPk(pkColumn string, pkValue string) laze.Exception {
	condition := fmt.Sprintf("%s = ?", pkColumn)

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


func (table *Table) UpdateByPk(pk string, data map[string]interface{}) laze.Exception {
	condition := fmt.Sprintf("%s = ?", table.Pk)

	err := table.recordExistsByPk(table.Pk, pk)
	if err != nil { return err }
	
	values := []interface{}{}
	columnsValues := ""
	counter := 0
	for key, val := range data {
		if _, ok := table.RawColumns[key]; ok {
			if key != table.Pk {
				values = append(values, val)
				columnValue := fmt.Sprintf("%s = ?", key)
				columnsValues = columnsValues + columnValue
				if counter < len(data) - 1 {
					columnsValues = columnsValues + ", "
				}
				counter = counter + 1
			}
		}		
	}
	values = append(values, pk)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table.Name, columnsValues, condition)

	rows, updateErr := table.Conn.Raw(query, values...).Rows()
	defer rows.Close()

	if updateErr != nil {
		ex := exception.FromError(updateErr, exception.INTERNALERROR)
		return ex
	}

	return nil
}