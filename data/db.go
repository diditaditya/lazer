package data

import (
	"fmt"
	"os"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	dataError "lazer/data/error"
)

type DB struct {
	Config *DBConfig
	Connection *gorm.DB
	tables map[string]*Table
}

func (db *DB) GetConnString() string {
	host := db.Config.Host
	port := db.Config.Port
	user := db.Config.User
	password := db.Config.Password
	database := db.Config.Database
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, database)
}

func (db *DB) GetTableNames() []string {
	tableNames := []string{}
	for name, _ := range db.tables {
		tableNames = append(tableNames, name)
	}
	return tableNames
}

func (db *DB) FindAll(tableName string) ([]map[string]interface{}, error) {
	if table, ok := db.tables[tableName]; ok {
		return table.FindAll(), nil
	}

	return nil, dataError.New("NOTFOUND", "table not found") 
}

func (db *DB) Tables() map[string]*Table {
	return db.tables
}

func (db *DB) GetTable(tableName string) *Table {
	return db.tables[tableName]
}

func (db *DB) GetAllTables() {
	var tableNames []string
	err := db.Connection.Raw("SHOW TABLES").Pluck("Tables_in_db", &tableNames).Error

	if err != nil {
		fmt.Println("[DB] error getting tables in db")
	}

	fmt.Printf("[DB] found tables: %v\n", tableNames)

	db.tables = make(map[string]*Table)
	for i:= 0; i < len(tableNames); i++ {
		rawColumns, columnNames := db.describeTable(tableNames[i])
		table := Table{
			name: tableNames[i],
			conn: db.Connection,
			ColumnNames: columnNames,
			RawColumns: rawColumns,
		}
		db.tables[tableNames[i]] = &table
	}
	fmt.Printf("[DB] tables: %v\n", db.tables)
}

func (db *DB) describeTable(tableName string) (map[string]RawColumn, []string) {

	query := fmt.Sprintf("DESCRIBE %s", tableName)

	fmt.Printf("[DB] %s\n", query)

	rows, err := db.Connection.Raw(query).Rows()

	if err != nil {
		fmt.Println("[DB] error describing table", tableName, err)
	}

	columns := make(map[string]RawColumn)
	names := []string{}

	if err != nil {
		fmt.Println("[DB] error getting column names", tableName, err)
	}

	for rows.Next() {
		var Field string
		var Type string
		var Null string
		var Key string
		var Default string
		var Extra string

		rows.Scan(&Field, &Type, &Null, &Key, &Default, &Extra)

		column := RawColumn{
			Field: Field,
			Type: Type,
			Null: Null,
			Key: Key,
			Default: Default,
			Extra: Extra,
		}

		columns[Field] = column
		names = append(names, Field)
	}

	return columns, names
}

func (db *DB) Close() {
	fmt.Println("[DB] closing database connection..")
	db.Connection.Close()
}

func newDB(config *DBConfig) *DB {
	if config.Port == 0 {
		config.Port = 3306
	}
	db := DB{
		Config: config,
	}

	return &db
}

func Connect() *DB {
	config := DBConfig{
		Host: os.Getenv("DB_HOST"),
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	}

	db := newDB(&config)
	connString := db.GetConnString()	
	
	connection, err := gorm.Open("mysql", connString)

	if err != nil {
		panic(fmt.Sprintf("[DB] error connecting to %s\n", connString))
	}

	fmt.Printf("[DB] successfully connected to %s\n", connString)
	db.Connection = connection

	db.GetAllTables()

	return db
}