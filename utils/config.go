package utils

import "os"

func ReadServerConfig() (connectString string) {
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
