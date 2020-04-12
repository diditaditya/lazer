package query

import (
	"fmt"
	"strings"
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