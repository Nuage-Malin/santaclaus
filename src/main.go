package main

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"fmt"
)

func main() {
	var server MaeSanta.Maestro_Santaclaus_ServiceServer

	server = GetSantaclausServerImpl()

	server.AddFile(nil, nil)
	fmt.Println("Goodbye")
}
