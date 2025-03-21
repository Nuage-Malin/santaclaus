package main

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"

	"context"
	"log"

	"errors"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// function get all disks
// function get disk for file size : algo get disk : best fit ? first fit ?

// service with Hardware malin

func getDisksFromBugle(ctx context.Context) ([]*pb.Disk, error) {
	//  grpc client for hardware manager
	//  query getDisks
	//	update (in mongo) disks that have changed according to hardware manager
	bugleAddress := os.Getenv("SANTACLAUS_BUGLE_URI")
	grpcOpts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, r := grpc.Dial(bugleAddress, grpcOpts)

	if r != nil {
		// todo query DB containing disks
		log.Println("Fail to reach bugle to update disks")
		return nil, r
	}
	defer conn.Close()
	client := pb.NewSantaclaus_HardwareMalin_ServiceClient(conn)
	request := pb.GetDisksRequest{}
	status, r := client.GetDisks(ctx, &request)

	if r != nil {
		return nil, r
	}
	return status.GetDisks(), nil
}

func (server *SantaclausServerImpl) updateDiskBase(ctx context.Context) (r error) {
	disks, r := getDisksFromBugle(ctx)

	if r != nil {
		return r
	}
	filter := bson.D{}
	update := bson.D{}
	// opts := options.Update().SetUpsert(true)

	for _, newDisk := range disks {
		filter = bson.D{bson.E{Key: "physical_id", Value: newDisk.Id}}
		findRes := server.mongoColls[DisksCollName].FindOne(ctx, filter)
		var foundDisk disk

		if findRes != nil {
			actFindRes := findRes.Decode(&foundDisk)
			if actFindRes != mongo.ErrNoDocuments {
				continue
			}
		}

		update = bson.D{
			bson.E{Key: "_id", Value: primitive.NewObjectID()},
			bson.E{Key: "physical_id", Value: newDisk.Id},
			bson.E{Key: "total_size", Value: 1000000000},     // todo put real value
			bson.E{Key: "available_size", Value: 1000000000}} // todo put real values

		// update DB containing disks
		insertRes, err := server.mongoColls[DisksCollName].InsertOne(ctx, update)
		if err != nil {
			log.Println(err)
			continue
		}
		if insertRes == nil {
			continue
		}
	}
	return r
}

func (server *SantaclausServerImpl) findAvailableDisk(ctx context.Context, fileSize uint64, userId string) (found primitive.ObjectID, r error) {
	// todo query hardware malin for updates
	server.updateDiskBase(ctx) // ignore error : if bugle can't update, we use the existing disk DB

	var disks []disk
	diskFound := disk{Id: primitive.NilObjectID, AvailableSize: 0, TotalSize: 0}
	filter := bson.D{bson.E{Key: "available_size", Value: bson.D{bson.E{Key: "$gt", Value: fileSize}}}}
	res, r := server.mongoColls[DisksCollName].Find(ctx, filter)

	if r != nil {
		return primitive.NilObjectID, r
	}
	if res == nil {
		return primitive.NilObjectID, errors.New("Could not find disk big enough for file")
	}
	r = res.All(ctx, &disks)
	if r != nil {
		return primitive.NilObjectID, r
	}
	// todo change this algorithm to something more accurate, with a real choice that's optimised
	for _, disk := range disks {
		if diskFound.AvailableSize < disk.AvailableSize {
			diskFound.AvailableSize = disk.AvailableSize
			found = disk.Id
		}
	}
	if found == primitive.NilObjectID {
		return primitive.NilObjectID, errors.New("Could not find disk big enough for file") // todo uncomment
	}
	// todo update disk size here ? or in other function ?
	return found, nil
}
