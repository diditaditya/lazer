package definitions

type TableDef struct {
	name    string
	columns []ColumnDef
	options TableOpts
}

type ColumnDef struct {
	name         string
	dtype        string
	allowNull    bool
	defaultValue any
}

type TableOpts struct {
	primaryKey string
}
