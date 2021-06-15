package server

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/shreybatra/crankdb/utils"
)

func parseCommand(message []byte, length int) (string, string) {

	var i int
	for i = range message {
		if message[i] == ' ' {
			break
		}
	}

	if i >= length {
		panic("Too long a command.")
	}

	return string(message[0:i]), string(message[i+1:])

}

func executeCommand(command string, arguments string) (response interface{}) {

	switch command {
	case "set":
		return set(arguments)
	case "get":
		return get(arguments)
	case "find":
		return find(arguments)
	case "del":
		return del(arguments)
	default:
		return "invalid command"
	}
}

func startConnection(connection *Connection) {
	fmt.Println("[Connection Opened] - " + connection.ip.String())
	defer connection.socket.Close()

	for {
		message := make([]byte, 4096)
		length, err := connection.socket.Read(message)

		if err != nil {
			fmt.Println("[Connection Closed] - " + connection.ip.String())
			connection.socket.Close()
			break
		}
		if length > 0 {
			message = message[:length]
			command, argsAndData := parseCommand(message, length)

			response := executeCommand(command, argsAndData)
			resp, _ := json.Marshal(response)
			_, err := connection.socket.Write(resp)

			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func StartServer() {

	// To be updated later for complex storage based on Protobuf.
	database = make(map[string]interface{})

	connectString := utils.ReadServerConfig()

	fmt.Println("Starting server... Accepting connections on -", connectString)

	listener, err := net.Listen("tcp", connectString)
	if err != nil {
		fmt.Println(err)
	}
	defer listener.Close()

	for {
		socketConn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		connection := &Connection{ip: socketConn.RemoteAddr(), socket: socketConn, data: make(chan []byte)}
		go startConnection(connection)
	}
}
