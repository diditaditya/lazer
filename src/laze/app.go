// The main app
package laze

import (
	exception "lazer/error"
)

type App struct {
	data Repo
}

func Init(data Repo) *App {
	app := App{
		data: data,
	}
	return &app
}

func (app *App) FindAll(tableName string, query map[string][]string) ([]map[string]interface{}, map[string]interface{}, *exception.Exception) {
	result, meta, err := app.data.FindAll(tableName, query)
	return result, meta, err
}

func (app *App) FindByPk(tableName string, pk string) (map[string]interface{}, *exception.Exception) {
	result, err := app.data.FindByPk(tableName, pk)
	return result, err
}

func (app *App) GetAllTables() []string {
	result := app.data.GetTableNames()
	return result
}

func (app *App) Create(tableName string, data map[string]interface{}) (map[string]interface{}, Exception) {
	result, err := app.data.Create(tableName, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (app *App) Delete(tableName string, params map[string][]string) Exception {
	err := app.data.Delete(tableName, params)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) DeleteByPk(tableName string, pk string) Exception {
	err := app.data.DeleteByPk(tableName, pk)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) Update(tableName string, params map[string][]string, data map[string]interface{}) Exception {
	err := app.data.Update(tableName, params, data)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) UpdateByPk(tableName string, pk string, data map[string]interface{}) Exception {
	err := app.data.UpdateByPk(tableName, pk, data)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) GetAssociations(tableName string) []map[string]interface{} {
	return app.data.GetTableAssociations(tableName)
}
