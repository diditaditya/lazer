package laze

type Repo interface {
	GetTableNames() []string
	FindAll(tableName string) ([]map[string]interface{}, error)
	FindByPk(tableName string, pk string) (map[string]interface{}, error)
	Create(tableName string, data map[string]interface{}) (map[string]interface{}, error)
}

type Exception interface {
	Message() string
	Code() string
}