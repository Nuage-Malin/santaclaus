package main

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"

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
