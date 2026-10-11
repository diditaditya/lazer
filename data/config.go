package data

import (
	"errors"
)

type DbType int

const (
	MySql DbType = iota
)

var dbTypeToString = map[DbType]string{
	MySql: "mysql",
}

var stringToDbType = map[string]DbType{
	"mysql": MySql,
}

var defaultPort = map[DbType]int{
	MySql: 3306,
}

type DBConfig struct {
	Type     DbType
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func (conf *DBConfig) GetType() string {
	return dbTypeToString[conf.Type]
}

func NewDbConfig(dbType string, dbHost string, dbPort int, dbUser string, dbPwd string, dbName string) (*DBConfig, error) {
	if len(dbType) == 0 {
		return nil, errors.New("db type is required")
	}

	if len(dbHost) == 0 {
		return nil, errors.New("db host is required")
	}

	if len(dbUser) == 0 {
		return nil, errors.New("db user is required")
	}

	if len(dbPwd) == 0 {
		return nil, errors.New("db password is required")
	}

	if len(dbName) == 0 {
		return nil, errors.New("db name is required")
	}

	dbTyped, ok := stringToDbType[dbType]
	if !ok {
		return nil, errors.New("Invalid db type")
	}

	if dbPort == 0 {
		dbPort = defaultPort[dbTyped]
	}

	dbConfig := DBConfig{
		Type:     dbTyped,
		Port:     dbPort,
		Host:     dbHost,
		User:     dbUser,
		Password: dbPwd,
		Database: dbName,
	}

	return &dbConfig, nil
}
