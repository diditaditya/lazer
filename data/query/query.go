package query

import (
	"fmt"

	"lazer/data/trait"
)

var Query = map[string]trait.SQLQuery{}

func Get(db string) trait.SQLQuery {
	mysql := MySQL{}

	Query["mysql"] = &mysql
	Query["MySQL"] = &mysql

	if querySet, ok := Query[db]; ok {
		return querySet
	} else {
		msg := fmt.Sprintf("invalid DB env variable %s, only mysql is supported for now", db)
		panic(msg)
	}
}