package query

import (
	"fmt"
	"strings"

	def "lazer/laze/definitions"
)

type MySQL struct {}

func (q *MySQL) DescribeTable(tableName string) string {
	fields := []string{
		"COLUMN_NAME AS Field",
		"DATA_TYPE AS Type",
		"IS_NULLABLE AS 'Null'",
		"COLUMN_KEY AS 'Key'",
		"COLUMN_DEFAULT AS 'Default'",
		"EXTRA AS Extra",
	}
	base := "SELECT %s FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '%s'"
	order := "ORDER BY ORDINAL_POSITION"
	baseSubbed := fmt.Sprintf(base, strings.Join(fields[:], ", "), tableName)
	query := baseSubbed + " " + order
	return query
}

func (q *MySQL) GetTables() string {
	return "SHOW TABLES"
}

func (q *MySQL) GetAssociations(tableName string) string {
	fields := []string{
		"COLUMN_NAME AS Field",
		"REFERENCED_TABLE_NAME AS ReferencedTable",
		"REFERENCED_COLUMN_NAME AS ReferencedField",
	}
	base := "SELECT %s FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE"
	baseSubbed := fmt.Sprintf(base, strings.Join(fields[:], ", "))
	where := "WHERE TABLE_NAME = '%s' AND REFERENCED_TABLE_NAME IS NOT NULL"
	whereSubbed := fmt.Sprintf(where, tableName)
	query := baseSubbed + " " + whereSubbed
	return query
}

func (q *MySQL) CreateTable(tableDef def.TableDef) string {
	base := "CREATE"
	return base
}