package laze

import (
	exception "lazer/error"
)

type Repo interface {
	GetTableNames() []string
	FindAll(tableName string) ([]map[string]interface{}, *exception.Exception)
	FindByPk(tableName string, pk string) (map[string]interface{}, *exception.Exception)
	Create(tableName string, data map[string]interface{}) (map[string]interface{}, *exception.Exception)
}

type Exception interface {
	Message() string
	Code() string
}