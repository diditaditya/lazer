package laze

import (
	exception "lazer/error"
)

type Repo interface {
	GetTableNames() []string
	FindAll(tableName string, params map[string][]string) ([]map[string]any, map[string]any, *exception.Exception)
	FindByPk(tableName string, pk string) (map[string]any, *exception.Exception)
	Create(tableName string, data map[string]any) (map[string]any, Exception)
	Delete(tableName string, params map[string][]string) Exception
	DeleteByPk(tableName string, pk string) Exception
	Update(tableName string, params map[string][]string, data map[string]any) Exception
	UpdateByPk(tableName string, pk string, data map[string]any) Exception
	GetTableAssociations(tableName string) []map[string]any
}

type Exception interface {
	error
	Message() string
	Name() string
	Trace() string
}
