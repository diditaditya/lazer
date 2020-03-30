package data

import (
	"lazer/laze"
	exception "lazer/error"

	"lazer/data/table"
)

func (db *DB) FindAll(tableName string, params map[string][]string) ([]map[string]interface{}, *exception.Exception) {
	if table, ok := db.tables[tableName]; ok {
		result, err := table.FindAll(params)
		if err != nil {
			ex := exception.FromError(err, exception.INTERNALERROR)
			return nil, ex
		}
		return result, nil
	}

	ex := exception.New(exception.NOTFOUND, "table not found")

	return nil, ex
}

func (db *DB) FindByPk(tableName string, value string) (map[string]interface{}, *exception.Exception) {
	if table, ok := db.tables[tableName]; ok {
		return table.FindByPk(value), nil
	}

	ex := exception.New(exception.NOTFOUND, "table not found")

	return nil, ex
}

func (db *DB) Create(tableName string, data map[string]interface{}) (map[string]interface{}, laze.Exception) {
	var table *table.Table
	if checked, ok := db.tables[tableName]; ok {
		table = checked
	} else {
		ex := exception.New(exception.NOTFOUND, "table not found")
		return nil, ex
	}

	result := table.Create(data)

	return result, nil
}

func (db *DB) Delete(tableName string, params map[string][]string) laze.Exception {
	var table *table.Table
	if checked, ok := db.tables[tableName]; ok {
		table = checked
	} else {
		ex := exception.New(exception.NOTFOUND, "table not found")
		return ex
	}

	err := table.Delete(params)
	if err != nil {
		return err
	}
	
	return nil
}