package main

import (
	"chat/app"
	"github.com/joho/godotenv"
)

func main() {

	//load .env file from given path
	godotenv.Load(".env")

	app.Start()
}