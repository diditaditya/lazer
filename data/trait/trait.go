package trait

type SQLQuery interface {
	GetTables() string
	DescribeTable(tableName string) string
	GetAssociations(tableName string) string
}