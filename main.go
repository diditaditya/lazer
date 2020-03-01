package main

import (
	"fmt"
	data "lazer/data"
	laze "lazer/laze"
	server "lazer/server"
)

func main() {
	fmt.Println("laze to the max")
	db := data.Connect()
	defer db.Close()

	app := laze.Init(db)

	server.Start(app)
}
