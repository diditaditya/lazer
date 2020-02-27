package laze

type Repo interface {
	GetTableNames() []string
	FindAll(name string) ([]map[string]interface{}, error)
}

type RepoException interface {
	Message() string
	Code() string
}