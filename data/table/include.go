package table

func (table *Table) parseFields(raw map[string]interface{}, fields []string) []string {
	tableName := ""
	incTable, isTableOk := raw["tableName"].(string)
	if isTableOk { tableName = incTable }
	incFields, isFieldsOk := raw["fields"].([]string)
	if isFieldsOk {
		for _, field := range incFields {
			fieldName := tableName + "." + field
			fields = append(fields, fieldName)
		}
	}

	incJoined, isJoinedOk := raw["joined"].([]map[string]interface{})
	if isJoinedOk {
		for _, joined := range incJoined {
			joinedFields := table.parseFields(joined, []string{})
			fields = append(fields, joinedFields...)
		}
	}

	return fields
}

func (table *Table) getFields(raw map[string]interface{}) (fieldMarks string, fields []string) {
	fields = table.parseFields(raw, []string{})

	for i, _ := range fields {
		fieldMarks = fieldMarks + "?"
		if i < len(fields) - 1 {
			fieldMarks = fieldMarks + ", "
		}
	}

	return fieldMarks, fields
}

func (table *Table) parseJoined(raw map[string]interface{}, result []map[string]string) []map[string]string {
	incJoined, isJoinedOk := raw["joined"].([]map[string]interface{})
	if isJoinedOk {
		for _, joined := range incJoined {
			association := map[string]string{}
			tableName, isTableNameOk := joined["tableName"].(string)
			if isTableNameOk {
				association["tableName"] = tableName
			}
			field, isAssociatedOk := joined["foreignKey"].(string)
			if isAssociatedOk {
				association["field"] = field
			}
			referencedField, isRefFieldOk := joined["referencedField"].(string)
			if isRefFieldOk {
				association["referencedField"] = referencedField
			}
			referencedTable, isRefTableOk := raw["tableName"].(string)
			if isRefTableOk {
				association["referencedTable"] = referencedTable
			}

			result = append(result, association)

			deeperJoin, isDeepJoinOk := joined["joined"].([]map[string]interface{})
			if isDeepJoinOk {
				if len(deeperJoin) > 0 {
					deeperAssociations := table.parseJoined(joined, []map[string]string{})
					result = append(result, deeperAssociations...)
				}
			}
		}
	}

	return result
}

func (table *Table) getJoined(raw map[string]interface{}) (marks string, values []string) {
	joins := table.parseJoined(raw, []map[string]string{})

	for i, join := range joins {
		if i > 0 { marks = marks + " " }

		marks = marks + "JOIN "

		marks = marks + "? "
		values = append(values, join["tableName"])

		marks = marks + "ON "
		
		marks = marks + "?"
		values = append(values, join["tableName"])

		marks = marks + ".?"
		values = append(values, join["field"])

		marks = marks + " = "

		marks = marks + "?"
		values = append(values, join["referencedTable"])

		marks = marks + ".?"
		values = append(values, join["referencedField"])
	}

	return marks, values
}