package main

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"fmt"

	"errors"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// function get all disks
// function get disk for file size : algo get disk : best fit ? first fit ?

// service with Hardware malin

func (server *SantaclausServerImpl) updateDiskBase() (r error) {
	// todo grpc client for hardware manager
	//		query getDisks
	//		update (in mongo) disks that have changed according to hardware manager
	hardwaremalinAddress := os.Getenv("HARDWAREMALIN_ADDRESS")
	options := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, r := grpc.Dial(hardwaremalinAddress, options)
	if r != nil {
		return r
	}
	defer conn.Close()
	client := pb.NewSantaclaus_HardwareMalin_ServiceClient(conn)
	request := pb.GetDisksRequest{}
	status, r := client.GetDisks(server.ctx, &request)
	if r != nil {
		return r
	}
	fmt.Printf("status.Disks: %v\n", status.Disks)
	// for _, disk := status.Disks {
	// continue
	// }
	return r
}

func (server *SantaclausServerImpl) findAvailableDisk(fileSize uint64, userId string) (found primitive.ObjectID, r error) {
	// todo query hardware malin for updates
	var disks []disk
	diskFound := disk{Id: primitive.NilObjectID, AvailableSize: 0, TotalSize: 0}
	filter := bson.D{{"available_size", bson.D{{"$gt", fileSize}}}}
	res, r := server.mongoColls[DisksCollName].Find(server.ctx, filter)

	if r != nil {
		return primitive.NilObjectID, r
	}
	if res == nil {
		return primitive.NilObjectID, errors.New("Could not find disk big enough for file")
	}
	r = res.All(server.ctx, &disks)
	if r != nil {
		return primitive.NilObjectID, r
	}

	for _, disk := range disks {
		if diskFound.AvailableSize < disk.AvailableSize {
			diskFound.AvailableSize = disk.AvailableSize
		}
	}
	if diskFound.Id == primitive.NilObjectID {
		found = primitive.NewObjectID() // tmp, todo change
		// return primitive.NilObjectID, errors.New("Could not find disk big enough for file") // todo uncomment
	}
	// todo update disk size here ? or in other function ?
	return found, nil
}
