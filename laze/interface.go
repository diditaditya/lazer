package laze

import (
	exception "lazer/error"
)

type Repo interface {
	GetTableNames() []string
	FindAll(tableName string, params map[string][]string) ([]map[string]interface{}, *exception.Exception)
	FindByPk(tableName string, pk string) (map[string]interface{}, *exception.Exception)
	Create(tableName string, data map[string]interface{}) (map[string]interface{}, Exception)
	Delete(tableName string, params map[string][]string) Exception
	DeleteByPk(tableName string, pk string) Exception
	Update(tableName string, params map[string][]string, data map[string]interface{}) Exception
	UpdateByPk(tableName string, pk string, data map[string]interface{}) Exception
}

type Exception interface {
	error
	Message() string
	Name() string
	Trace() string
}