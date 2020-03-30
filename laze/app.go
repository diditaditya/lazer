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

func (app *App) FindAll(tableName string, query map[string][]string) ([]map[string]interface{}, *exception.Exception) {
	result, err := app.data.FindAll(tableName, query)
	return result, err
}

func (app *App) FindByPk(tableName string, value string) (map[string]interface{}, *exception.Exception) {
	result, err := app.data.FindByPk(tableName, value)
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