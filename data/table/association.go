package table

import (
	"fmt"
)

func (table *Table) SetAssociations(data map[string]map[string]string) {
	table.Associations = data
}

func (table *Table) GetAssociations() []map[string]any {
	associations := []map[string]any{}
	for _, raw := range table.Associations {
		rel := make(map[string]any)
		rel["field"] = raw["field"]
		rel["referencedTable"] = raw["referencedTable"]
		rel["referencedField"] = raw["referencedField"]
		rel["type"] = raw["type"]
		fmt.Printf("%v\n", rel)
		associations = append(associations, rel)
	}

	return associations
}
