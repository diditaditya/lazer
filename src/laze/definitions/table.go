package definitions

type TableDef struct {
	name		string
	columns	[]ColumnDef
	options TableOpts
}

type ColumnDef struct {
	name					string
	dtype					string
	allowNull			bool
	defaultValue	interface{}
}

type TableOpts struct {
	primaryKey	string
}