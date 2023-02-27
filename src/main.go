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
	listeningAddress := os.Getenv("SANTACLAUS_ADDRESS")
	listener, err := net.Listen("tcp", listeningAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	var server MaeSanta.Maestro_Santaclaus_ServiceServer = NewSantaclausServerImpl()

	grpcServer := grpc.NewServer(opts...)
	MaeSanta.RegisterMaestro_Santaclaus_ServiceServer(grpcServer, server)
	fmt.Println("Hello")
	grpcServer.Serve(listener)

	fmt.Println("Goodbye")
}
