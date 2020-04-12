package data

import (
	"fmt"

	"lazer/data/constants"
)

type Association struct {
	Field string
	ReferencedTable string
	ReferencedField string
	Type string
}

func (db *DB) getAssociations(tableName string) {
	query := db.query.GetAssociations(tableName)
	rows, err := db.Connection.Raw(query).Rows()

	defer rows.Close()

	if err != nil {
		fmt.Println("[DB] error getting table association", err)
		return
	}

	rels := make(map[string]*Association)

	for rows.Next() {
		var Field string
		var ReferencedTable string
		var ReferencedField string

		rows.Scan(&Field, &ReferencedTable, &ReferencedField)

		rel := Association{
			Field: Field,
			ReferencedTable: ReferencedTable,
			ReferencedField: ReferencedField,
			Type: constants.AssociationBelongsTo,
		}
		rels[ReferencedTable] = &rel
	}
	db.associations[tableName] = rels
}

func (db *DB) crossCheckAssociation(tableName string) {
	rels := db.associations[tableName]

	for relatedTable, rel := range rels {
		if _, ok := db.associations[relatedTable][tableName]; !ok {
			newRel := Association{
				Field: rel.ReferencedField,
				ReferencedTable: tableName,
				ReferencedField: rel.Field,
				Type: constants.AssociationHasMany,
			}
			db.associations[relatedTable][tableName] = &newRel
		}
	}
}

func (db *DB) completeAssociations() {
	for tableName := range db.associations {
		db.crossCheckAssociation(tableName)
	}

	for tableName, rels := range db.associations {
		associations := map[string]map[string]string{}
		for _, rel := range rels {
			data := map[string]string{
				"field": rel.Field,
				"referencedTable": rel.ReferencedTable,
				"referencedField": rel.ReferencedField,
				"type": rel.Type,		
			}
			associations[rel.ReferencedTable] = data
		}
		
		db.tables[tableName].SetAssociations(associations)
	}
}

func (db *DB) GetTableAssociations(tableName string) (rels []map[string]interface{}) {
	if table, ok := db.tables[tableName]; ok {
		rels = table.GetAssociations()
		return rels
	}
	return rels
}