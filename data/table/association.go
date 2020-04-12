package table

import (
	"fmt"
)

func (table *Table) SetAssociations(data map[string]map[string]string) {
	fmt.Printf("[table][%s] setting associations: %v\n", table.Name, data)
	table.Associations = data
}

func (table *Table) GetAssociations() []map[string]interface{} {
	fmt.Printf("[table][%s] getting associations\n", table.Name)
	fmt.Printf("[table][%s] associations: %v\n", table.Name, table.Associations)
	associations := []map[string]interface{}{}
	for _, raw := range table.Associations {
		rel := make(map[string]interface{})
		rel["field"] = raw["field"]
		rel["referencedTable"] = raw["referencedTable"]
		rel["referencedField"] = raw["referencedField"]
		rel["type"] = raw["type"]
		fmt.Printf("%v\n", rel)
		associations = append(associations, rel)
	}

	return associations
}