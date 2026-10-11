package data

import (
	"lazer/data/trait"
)

type Included struct {
	tableName				string
	pk							string
	fields 					[]string
	foreignKey			string
	referencedField	string
	referencedTable	string
	referenceType		string
	joined					[]trait.Joined
}

func (inc *Included) GetTableName() string {
	return inc.tableName
}

func (inc *Included) GetTablePk() string {
	return inc.pk
}

func (inc *Included) SetTablePk(pk string) {
	inc.pk = pk
}

func (inc *Included) GetFields() []string {
	return inc.fields
}

func (inc *Included) SetFields(fields []string) {
	inc.fields = fields
}

func (inc *Included) GetForeignKey() string {
	return inc.foreignKey
}

func (inc *Included) GetReferencedField() string {
	return inc.referencedField
}

func (inc *Included) GetReferencedTable() string {
	return inc.referencedTable
}

func (inc *Included) GetReferenceType() string {
	return inc.referenceType
}

func (inc *Included) GetJoined() []trait.Joined {
	return inc.joined
}

func (inc *Included) SetJoined(joined []trait.Joined) {
	inc.joined = joined
}