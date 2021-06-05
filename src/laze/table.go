package laze

import (
	"fmt"

	def "lazer/laze/definitions"
)

func (app *App) CreateTable(table def.TableDef) {
	fmt.Println("[LAZE]] create table")
}