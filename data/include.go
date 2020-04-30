package data

import (
	"lazer/data/trait"
)

type Included struct {
	tableName				string
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

	result.SetFields(fields)

	return result
}