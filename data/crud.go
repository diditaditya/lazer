package data

import (
	"strings"

	exception "lazer/error"
	"lazer/laze"

	"lazer/data/table"
	"lazer/data/trait"
)

func (db *DB) FindAll(tableName string, params map[string][]string) ([]map[string]interface{}, map[string]interface{}, *exception.Exception) {
	var include trait.Joined
	if inc, ok := params["include"]; ok {
		result := Included{
			tableName: tableName,
			fields: []string{},
			foreignKey: "",
			referencedField: "",
			referencedTable: "",
			referenceType: "",
			joined: []trait.Joined{},
		}
		incStr := strings.Join(inc, ",")
		include = db.parseInclude(tableName, incStr, &result)
	}

	if table, ok := db.tables[tableName]; ok {
		result, meta, err := table.FindAll(params, include)
		if err != nil {
			ex := exception.FromError(err, exception.INTERNALERROR)
			return nil, nil, ex
		}
		return result, meta, nil
	}

	ex := exception.New(exception.NOTFOUND, "table not found")

	return nil, nil, ex
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

func (db *DB) DeleteByPk(tableName string, value string) laze.Exception {
	var table *table.Table
	if checked, ok := db.tables[tableName]; ok {
		table = checked
	} else {
		ex := exception.New(exception.NOTFOUND, "table not found")
		return ex
	}

	err := table.DeleteByPk(value)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateByPk(tableName string, pkValue string, data map[string]interface{}) laze.Exception {
	var table *table.Table
	if checked, ok := db.tables[tableName]; ok {
		table = checked
	} else {
		ex := exception.New(exception.NOTFOUND, "table not found")
		return ex
	}

	err := table.UpdateByPk(pkValue, data)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Update(tableName string, params map[string][]string, data map[string]interface{}) laze.Exception {
	err := db.tableExists(tableName)
	if err != nil {
		return err
	}

	table := db.tables[tableName]
	err = table.Update(params, data)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) tableExists(tableName string) laze.Exception {
	if _, ok := db.tables[tableName]; !ok {
		ex := exception.New(exception.NOTFOUND, "table not found")
		return ex
	}
	return nil
}
