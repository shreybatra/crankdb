package main

import (
	"github.com/joho/godotenv"
	"github.com/shreybatra/crankdb/server"
)

func main() {
	godotenv.Load()
	server.StartServer()
}
