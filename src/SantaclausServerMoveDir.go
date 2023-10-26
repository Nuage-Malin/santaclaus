package main

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"fmt"

	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func (server *SantaclausServerImpl) MoveDirectory(ctx context.Context, req *pb.MoveDirectoryRequest) (status *pb.MoveDirectoryStatus, r error) {
	log.Println("Request: MoveDirectory" /* todo try to get function name from variable or macro */) // todo replace with class request logger

	// Check that dir exists
	dir, r := server.GetDirFromStringId(ctx, req.GetDirId())

	if r != nil {
		return nil, r
	}

	// Check that new dir exists
	newParentDir, r := server.GetDirFromStringId(ctx, req.GetNewDirId())

	if r != nil {
		return nil, r
	}

	if dir.Id == newParentDir.Id {
		return nil, errors.New("Cannot move directory to itself, aborting move")
	}

	// Check that newParentDir does not contain dir with same name as dir to be moved
	if server.CheckDirHasChild(ctx, newParentDir.Id, dir.Name) {
		return nil, errors.New(fmt.Sprintf("Directory with name '%s' already exists in parent directory %s, aborting move", dir.Name, newParentDir.Id.Hex()))
	}

	filter := bson.D{bson.E{Key: "_id", Value: dir.Id}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "parent_id", Value: newParentDir.Id}}}}
	res, r := server.mongoColls[DirectoriesCollName].UpdateOne(ctx, filter, update)

	if r != nil {
		return nil, r
	}
	if res.MatchedCount != 1 {
		return nil, errors.New("Could not find file to be updated")
	}
	if res.ModifiedCount != 1 {
		// log.Print(res.ModifiedCount)
		return nil, errors.New("Could not modify file to be updated")
	}
	status = &pb.MoveDirectoryStatus{}
	return status, nil
}

// todo put in separate file
func (server *SantaclausServerImpl) RenameDirectory(ctx context.Context, req *pb.RenameDirectoryRequest) (status *pb.RenameDirectoryStatus, r error) {
	log.Println("Request: RenameDirectory" /* todo try to get function name from variable or macro */) // todo replace with class request logger
	status = &pb.RenameDirectoryStatus{}
	// return status, nil

	// Check that dir exists
	dir, r := server.GetDirFromStringId(ctx, req.GetDirId())

	if r != nil {
		return nil, r
	}

	newDirName := req.GetNewDirName()
	// Check that parent dir does not contain dir with same name as newDirName
	if server.CheckDirHasChild(ctx, dir.ParentId, newDirName) {
		return nil, errors.New(fmt.Sprintf("Directory with name '%s' already exists in parent directory %s, aborting rename", dir.Name, dir.Id.Hex()))
	}

	filter := bson.D{bson.E{Key: "_id", Value: dir.Id}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "name", Value: newDirName}}}}
	res, r := server.mongoColls[DirectoriesCollName].UpdateOne(ctx, filter, update)

	if r != nil {
		return nil, r
	}
	if res.MatchedCount != 1 {
		return nil, errors.New("Could not find file to be updated")
	}
	if res.ModifiedCount != 1 {
		// log.Print(res.ModifiedCount)
		return nil, errors.New("Could not modify file to be updated")
	}
	status = &pb.RenameDirectoryStatus{}
	return status, nil

}
