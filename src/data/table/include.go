package table

import (
	"fmt"

	"lazer/data/trait"
)

func (table *Table) parseFields(raw trait.Joined, fields []string) []string {
	if raw == nil { return fields }

	tableName := raw.GetTableName()
	incFields := raw.GetFields()
	pk := raw.GetTablePk()
	isPkIncluded := false
	for _, field := range incFields {
		if field == pk { isPkIncluded = true }
		fieldName := tableName + "." + field
		fields = append(fields, fieldName)
	}

	if (!isPkIncluded) {
		fieldName := tableName + "." + pk
		fields = append(fields, fieldName)
	}

	incJoined := raw.GetJoined()
	for _, joined := range incJoined {
		joinedFields := table.parseFields(joined, []string{})
		fields = append(fields, joinedFields...)
	}

	return fields
}

func (table *Table) getFields(raw trait.Joined) (fieldMarks string, fields []string) {
	fields = table.parseFields(raw, []string{})

	if len(fields) == 0 {
		for _, col := range table.ColumnNames {
			fields = append(fields, fmt.Sprintf("%s.%s", table.Name, col))
		}
	}

	for i, field := range fields {
		fieldMarks = fieldMarks + field
		if i < len(fields) - 1 {
			fieldMarks = fieldMarks + ", "
		}
	}

	return fieldMarks, fields
}

func (table *Table) parseJoined(raw trait.Joined, result []map[string]string) []map[string]string {
	if raw == nil { return result }

	incJoined := raw.GetJoined()
	for _, joined := range incJoined {
		association := map[string]string{
			"tableName": joined.GetTableName(),
			"field": joined.GetForeignKey(),
			"referencedField": joined.GetReferencedField(),
			"referencedTable": joined.GetReferencedTable(),
			"referenceType": joined.GetReferenceType(),
		}

		result = append(result, association)

		deeperJoined := joined.GetJoined()

		if len(deeperJoined) > 0 {
			deeperAssociations := table.parseJoined(joined, []map[string]string{})
			result = append(result, deeperAssociations...)
		}
	}

	return result
}

func (table *Table) getJoined(raw trait.Joined) (marks string, values []string) {
	joins := table.parseJoined(raw, []map[string]string{})

	for i, join := range joins {
		if i > 0 { marks = marks + " " }

		marks = marks + "LEFT JOIN "

		marks = marks + join["tableName"] + " "
		values = append(values, join["tableName"])

		marks = marks + "ON "
		
		marks = marks + join["tableName"]
		values = append(values, join["tableName"])

		marks = marks + "." + join["field"]
		values = append(values, join["field"])

		marks = marks + " = "

		marks = marks + join["referencedTable"]
		values = append(values, join["referencedTable"])

		marks = marks + "." + join["referencedField"]
		values = append(values, join["referencedField"])
	}

	return marks, values
}