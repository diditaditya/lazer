package laze

type App struct {
	data Repo
}

func Init(data Repo) *App {
	app := App{
		data: data,
	}
	return &app
}

func (app *App) FindAll(tableName string) ([]map[string]interface{}, error) {
	result, err := app.data.FindAll(tableName)
	return result, err
}

func (app *App) FindByPk(tableName string, value string) (map[string]interface{}, error) {
	result, err := app.data.FindByPk(tableName, value)
	return result, err
}

func (app *App) GetAllTables() []string {
	result := app.data.GetTableNames()
	return result
}