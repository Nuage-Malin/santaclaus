package main

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	port := os.Getenv("SANTACLAUS_LISTENING_PORT")

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	var server MaeSanta.Maestro_Santaclaus_ServiceServer
	server = GetSantaclausServerImpl()

	grpcServer := grpc.NewServer(opts...)
	MaeSanta.RegisterMaestro_Santaclaus_ServiceServer(grpcServer, server)
	fmt.Println("Hello")
	grpcServer.Serve(listener)

	// server.AddFile(nil, nil)
	fmt.Println("Goodbye")
}
