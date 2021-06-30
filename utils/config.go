package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func ReadServerConfig() (connectString string) {
	godotenv.Load()
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "9876"
	}
	hosts, ok := os.LookupEnv("HOSTS")
	if !ok {
		hosts = "localhost"
	}

	return hosts + ":" + port
}
