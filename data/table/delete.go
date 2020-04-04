package table

import (
	"fmt"

	"lazer/laze"
	exception "lazer/error"
)

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

func (table *Table) DeleteByPk(value string) laze.Exception {
	condition := fmt.Sprintf("%s = ?", table.Pk)

	found, err := table.Conn.Table(table.Name).Where(condition, value).Rows()

	if err != nil {
		fmt.Println("[DB] error finding entry from ", table.Name)
		fmt.Println(err)
		ex := exception.FromError(err, exception.INTERNALERROR)
		return ex
	}

	if !found.Next() {
		ex := exception.New(exception.UNPROCESSABLE, "record not found")
		return ex
	}

	fmt.Printf("%v\n", *found)

	rawQuery := "DELETE FROM "
	rawQuery = rawQuery + table.Name
	rawQuery = rawQuery + " WHERE " + condition

	rows, err := table.Conn.Raw(rawQuery, value).Rows()

	defer rows.Close()

	if err != nil {
		fmt.Println("[DB] error deleting from ", table.Name)
		fmt.Println(err)
		ex := exception.FromError(err, exception.INTERNALERROR)
		return ex
	}

	return nil
}