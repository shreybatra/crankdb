package server

import "net"

type Connection struct {
	ip     net.Addr
	socket net.Conn
	data   chan []byte
}
