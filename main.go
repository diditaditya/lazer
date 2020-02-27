package main

import (
	"fmt"
	data "lazer/data"
	server "lazer/server"
	laze "lazer/laze"
)

func main() {
	fmt.Println("laze to the max")
	db := data.Connect()
	defer db.Close()

	app := laze.Init(db)

	server.Start(app)
}
