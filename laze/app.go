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