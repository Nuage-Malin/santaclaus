package main

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"flag"
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
	defer listener.Close()
	var santaClausServer MaeSanta.Maestro_Santaclaus_ServiceServer = NewSantaclausServerImpl()
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()
	MaeSanta.RegisterMaestro_Santaclaus_ServiceServer(grpcServer, santaClausServer)
	grpcServer.Serve(listener)
}
