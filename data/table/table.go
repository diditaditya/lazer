package table

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	exception "lazer/error"
	"lazer/laze"
)

type Table struct {
	Name        string
	Conn        *gorm.DB
	ColumnNames []string
	RawColumns  map[string]RawColumn
	Pk          string
}

type RawColumn struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

func (table *Table) GetPkColumn() {
	for k, v := range table.RawColumns {
		if v.Key == "PRI" {
			table.Pk = k
			break
		}
	}
}

func (table *Table) recordExistsByPk(pkValue string) laze.Exception {
	condition := fmt.Sprintf("%s = ?", table.Pk)

	found, err := table.Conn.Table(table.Name).Where(condition, pkValue).Rows()

	if err != nil {
		ex := exception.FromError(err, exception.INTERNALERROR)
		return ex
	}

	if !found.Next() {
		ex := exception.New(exception.UNPROCESSABLE, "record not found")
		return ex
	}
	return nil
}
