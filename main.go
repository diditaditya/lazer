package main

import (
	"fmt"
	"github.com/joho/godotenv"
	data "lazer/data"
	laze "lazer/laze"
	server "lazer/server"
	"log"
)

func getEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	getEnv()

	fmt.Println("[LAZER] laze to the max")
	db := data.Connect()
	defer db.Close()

	app := laze.Init(db)

	server.Start(app)
}
