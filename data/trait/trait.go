package trait

type SQLQuery interface {
	GetTables() string
	DescribeTable(tableName string) string
	GetAssociations(tableName string) string
}

type Joined interface {
	GetTableName() string
	GetTablePk() string
	SetTablePk(string)
	GetFields() []string
	GetForeignKey() string
	GetReferencedField() string
	GetReferencedTable() string
	GetReferenceType() string
	GetJoined() []Joined
	SetFields([]string)
	SetJoined([]Joined)
}