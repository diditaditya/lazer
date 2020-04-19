package data

func (db *DB) isField(tableName string, field string) bool {
	if table, tableFound := db.tables[tableName]; tableFound {
		if _, fieldFound := table.RawColumns[field]; fieldFound {
			return true
		}
	}
	return false
}

func (db *DB) isAssociation(sourceTable string, targetTable string) bool {
	if source, sourceFound := db.associations[sourceTable]; sourceFound {
		if _, targetFound := source[targetTable]; targetFound {
			return true
		}
	}
	return false
}

func (db *DB) parseInclude(tableName string, include string, result map[string]interface{}) map[string]interface{} {
	isInBracket := false
	openBracketCounter := 0
	closingBracketCounter := 0
	fields := []string{}
	field := ""
	joinedTable := ""
	joinedFieldsStr := ""
	for _, char := range include {
		str := string(char)
		switch str {
		case "(":
			openBracketCounter = openBracketCounter + 1
			isInBracket = true
			if (db.isAssociation(tableName, field)) {
				joinedTable = field
				field = ""
			}
			if openBracketCounter > closingBracketCounter + 1 {
				joinedFieldsStr = joinedFieldsStr + str
			}
		case ",":
			if isInBracket {
				joinedFieldsStr = joinedFieldsStr + str
			} else if db.isField(tableName, field) {
				fields = append(fields, field)
				field = ""
			} else if db.isAssociation(tableName, field) {
				joined := []map[string]interface{}{}
				if prev, ok := result["joined"]; ok {
					if prevJoined, isValid := prev.([]map[string]interface{}); isValid {
						joined = prevJoined
					}
				}
				joinedRes := map[string]interface{}{
					"tableName": field,
					"fields": []string{},
					"joined": []map[string]interface{}{},
				}
				joined = append(joined, joinedRes)
				result["joined"] = joined
				field = ""
			}
		case ")":
			closingBracketCounter = closingBracketCounter + 1
			if openBracketCounter == closingBracketCounter {
				isInBracket = false
				joined := []map[string]interface{}{}
				if prev, ok := result["joined"]; ok {
					if prevJoined, isValid := prev.([]map[string]interface{}); isValid {
						joined = prevJoined
					}
				}
				joinedRes := map[string]interface{}{
					"tableName": joinedTable,
					"fields": []string{},
					"joined": []map[string]interface{}{},
				}
				res := db.parseInclude(joinedTable, joinedFieldsStr, joinedRes)
				joined = append(joined, res)

				result["joined"] = joined
				joinedTable = ""
				joinedFieldsStr = ""
			} else {
				joinedFieldsStr = joinedFieldsStr + str
			}
		default:
			if !isInBracket {
				field = field + str
			} else {
				joinedFieldsStr = joinedFieldsStr + str
			}
		}
	}

	if len(field) > 0 {
		if db.isField(tableName, field) {
			fields = append(fields, field)
			field = ""
		} else if db.isAssociation(tableName, field) {
			joined := []map[string]interface{}{}
			if prev, ok := result["joined"]; ok {
				if prevJoined, isValid := prev.([]map[string]interface{}); isValid {
					joined = prevJoined
				}
			}
			joinedRes := map[string]interface{}{
				"tableName": field,
				"fields": []string{},
				"joined": []map[string]interface{}{},
			}
			joined = append(joined, joinedRes)
			result["joined"] = joined
		}
	}

	result["fields"] = fields

	return result
}