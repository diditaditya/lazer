package data

import (
	"lazer/data/trait"
)

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

func (db *DB) getAllFields(tableName string) []string {
	fields := []string{}
	if table, tableFound := db.tables[tableName]; tableFound {
		for field, _ := range table.RawColumns {
			fields = append(fields, field)
		}
	}
	return fields
}

func (db *DB) parseInclude(tableName string, include string, result trait.Joined) trait.Joined {
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
				joined := []trait.Joined{}
				prevJoined := result.GetJoined()
				if len(prevJoined) > 0 { joined = prevJoined }

				joinedRes := Included{
					tableName: field,
					pk: db.tables[field].Pk,
					fields: []string{},
					foreignKey: db.associations[field][tableName].Field,
					referencedField: db.associations[field][tableName].ReferencedField,
					referencedTable: db.associations[field][tableName].ReferencedTable,
					referenceType: db.associations[field][tableName].Type,
					joined: []trait.Joined{},
				}

				joined = append(joined, &joinedRes)

				result.SetJoined(joined)
				field = ""
			}
		case ")":
			closingBracketCounter = closingBracketCounter + 1
			if openBracketCounter == closingBracketCounter {
				isInBracket = false

				joined := []trait.Joined{}
				prevJoined := result.GetJoined()
				if len(prevJoined) > 0 { joined = prevJoined }

				joinedRes := Included{
					tableName: joinedTable,
					pk: db.tables[joinedTable].Pk,
					fields: []string{},
					foreignKey: db.associations[joinedTable][tableName].Field,
					referencedField: db.associations[joinedTable][tableName].ReferencedField,
					referencedTable: db.associations[joinedTable][tableName].ReferencedTable,
					referenceType: db.associations[joinedTable][tableName].Type,
					joined: []trait.Joined{},
				}

				res := db.parseInclude(joinedTable, joinedFieldsStr, &joinedRes)
				joined = append(joined, res)

				result.SetJoined(joined)
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
			joined := []trait.Joined{}
			prevJoined := result.GetJoined()
			if len(prevJoined) > 0 { joined = prevJoined }

			joinedRes := Included{
				tableName: field,
				pk: db.tables[field].Pk,
				fields: db.getAllFields(field),
				foreignKey: db.associations[field][tableName].Field,
				referencedField: db.associations[field][tableName].ReferencedField,
				referencedTable: db.associations[field][tableName].ReferencedTable,
				referenceType: db.associations[field][tableName].Type,
				joined: []trait.Joined{},
			}

			joined = append(joined, &joinedRes)

			result.SetJoined(joined)
		}
	}

	if len(fields) == 0 {
		fields = db.getAllFields(tableName)
	}

	result.SetTablePk(db.tables[tableName].Pk)
	result.SetFields(fields)

	return result
}