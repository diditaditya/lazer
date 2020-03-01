package laze

type Repo interface {
	GetTableNames() []string
	FindAll(name string) ([]map[string]interface{}, error)
}

type Exception interface {
	Message() string
	Code() string
}