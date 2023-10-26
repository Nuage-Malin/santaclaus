package main

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"fmt"
	"log"

	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

func (server *SantaclausServerImpl) MoveFile(ctx context.Context, req *pb.MoveFileRequest) (status *pb.MoveFileStatus, r error) {
	// todo test

	log.Println("Request: MoveFile" /* todo try to get function name from variable or macro */) // todo replace with class request logger

	file, r := server.GetFileFromStringId(ctx, req.GetFileId())

	if r != nil {
		return nil, r
	}

	// todo if nil object id for dirId, move to root dir ?
	newDir, r := server.GetDirFromStringId(ctx, req.GetNewDirId())
	if r != nil {
		return nil, r
	}
	if file.DirId == newDir.Id {
		return nil, errors.New("Cannot move file to directory it already is in, aborting move")
	}
	if server.CheckDirHasChild(ctx, newDir.Id, file.Name) {
		return nil, errors.New(fmt.Sprintf("File with name '%s' already exists in parent directory %s, aborting move", file.Name, newDir.Id.Hex()))
	}

	filter := bson.D{bson.E{Key: "_id", Value: file.Id}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "dir_id", Value: newDir.Id}}}}
	res, r := server.mongoColls[FilesCollName].UpdateOne(ctx, filter, update) // todo test updateById

	if r != nil {
		return nil, r
	}
	if res.MatchedCount != 1 {
		return nil, errors.New("Could not find file to be updated")
	}
	if res.ModifiedCount != 1 {
		return nil, errors.New("Could not modify file directory")
	}
	status = &pb.MoveFileStatus{}
	return status, nil
}

func (server *SantaclausServerImpl) RenameFile(ctx context.Context, req *pb.RenameFileRequest) (status *pb.RenameFileStatus, r error) {
	// todo test
	log.Println("Request: RenameFile" /* todo try to get function name from variable or macro */) // todo replace with class request logger

	file, r := server.GetFileFromStringId(ctx, req.GetFileId())

	if r != nil {
		return nil, r
	}
	newFileName := req.GetNewFileName()

	if server.CheckDirHasChild(ctx, file.DirId, newFileName) {
		return nil, errors.New(fmt.Sprintf("File with name '%s' already exists in parent directory %s, aborting move", newFileName, file.Id.Hex()))
	}

	filter := bson.D{bson.E{Key: "_id", Value: file.Id}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "name", Value: newFileName}}}}
	res, r := server.mongoColls[FilesCollName].UpdateOne(ctx, filter, update) // todo test updateById

	if r != nil {
		return nil, r
	}
	if res.MatchedCount != 1 {
		return nil, errors.New("Could not find file to be updated")
	}
	if res.ModifiedCount != 1 {
		return nil, errors.New("Could not modify file name")
	}
	status = &pb.RenameFileStatus{}
	return status, nil
}
