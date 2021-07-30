package main

import (
	"flag"
	"log"
	"net"

	"github.com/ahsanbarkati/crankdb/server"
	"github.com/ahsanbarkati/crankdb/utils"
	"google.golang.org/grpc"
)

func main() {

	hostport := utils.ReadServerConfig()
	flag.Parse()

	lis, err := net.Listen("tcp", hostport)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("starting server on %v", hostport)

	grpcServer := grpc.NewServer()
	server.RegisterCrankDBServer(grpcServer, &server.CrankServer{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
