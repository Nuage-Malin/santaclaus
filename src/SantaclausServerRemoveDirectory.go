package main

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"fmt"
	"log"

	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (server *SantaclausServerImpl) removeOneDirectory(ctx context.Context, dirId *primitive.ObjectID, status *pb.RemoveDirectoryStatus) (r error) {
	var files []file
	filter := bson.D{bson.E{Key: "_id", Value: dirId}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "deleted", Value: true}}}}
	updateRes, r := server.mongoColls[DirectoriesCollName].UpdateOne(ctx, filter, update)

	if r != nil {
		return r
	}
	if updateRes == nil {
		return errors.New("Could not delete directory")
	}
	if updateRes.MatchedCount != 1 {
		return errors.New("Could not find directory to delete")
	}
	if updateRes.ModifiedCount != 1 {
		return errors.New("Could not modify directory to delete")
	}
	filter = bson.D{bson.E{Key: "dir_id", Value: dirId}}
	findRes, r := server.mongoColls[FilesCollName].Find(ctx, filter)
	if r != nil {
		return r
	}
	r = findRes.All(ctx, &files)
	if r != nil {
		return r
	}
	for _, file := range files {
		status.FileIdsToRemove = append(status.FileIdsToRemove, file.Id.Hex())
	}
	return nil
}

/**
 * recursivelly remove children directories of a directory
 */
func (server *SantaclausServerImpl) removeSubDirectories(ctx context.Context, parentId *primitive.ObjectID, status *pb.RemoveDirectoryStatus) (r error) {
	var dirs []directory
	filter := bson.D{bson.E{Key: "parent_id", Value: parentId}}
	findRes, r := server.mongoColls[DirectoriesCollName].Find(ctx, filter)

	if r != nil {
		return r
	}
	r = findRes.All(ctx, &dirs)
	if r != nil {
		return r
	}
	for _, dir := range dirs {
		server.removeSubDirectories(ctx, &dir.Id, status)
		server.removeOneDirectory(ctx, &dir.Id, status)
	}
	return nil
}

func (server *SantaclausServerImpl) RemoveDirectory(ctx context.Context, req *pb.RemoveDirectoryRequest) (status *pb.RemoveDirectoryStatus, r error) {
	log.Println("Request: RemoveDirectory") // todo replace with class request logger

	//	find all subdirectories recursively
	//	in each subdirectory, add fileIds to the status (fileIdsToRemove)
	//	set directories as deleted
	dirId, r := primitive.ObjectIDFromHex(req.DirId)

	if r != nil {
		return nil, r
	}
	// Check if dir exists
	filter := bson.D{bson.E{Key: "_id", Value: dirId}}
	findRes := server.mongoColls[DirectoriesCollName].FindOne(ctx, filter)
	if findRes == nil || findRes.Err() != nil {
		return nil, fmt.Errorf("Could not remove directory : %s, because it doesn't exist", dirId)
	}
	var tmpDir directory
	r = findRes.Decode(&tmpDir)
	if r != nil {
		return nil, r
	}

	status = &pb.RemoveDirectoryStatus{FileIdsToRemove: make([]string, 0)}
	// Mark all subdirectories as deleted (Virtual)
	r = server.removeSubDirectories(ctx, &dirId, status)
	if r != nil {
		return nil, r
	}
	r = server.removeOneDirectory(ctx, &dirId, status)
	if r != nil {
	}
	// Remove all directories that have been marked as deleted in recursive sub functions (Physical)
	filter = bson.D{bson.E{Key: "deleted", Value: true}}
	_, r = server.mongoColls[DirectoriesCollName].DeleteMany(ctx, filter)
	if r != nil {
		return nil, r
	}
	return status, nil
}
