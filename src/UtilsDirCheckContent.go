package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (server *SantaclausServerImpl) CheckDirHasChild(ctx context.Context, dirId primitive.ObjectID, childName string) bool {
	// todo check if dirId is an existing directory ?
	filter := bson.D{bson.E{Key: "parent_id", Value: dirId}, bson.E{Key: "name", Value: childName}, bson.E{Key: "deleted", Value: false}} // get all files if not deleted

	res := server.mongoColls[DirectoriesCollName].FindOne(ctx, filter)

	if res.Err() == nil {
		return true
	}
	return false
}

func (server *SantaclausServerImpl) CheckDirHasFile(ctx context.Context, dirId primitive.ObjectID, fileName string) bool {
	filter := bson.D{bson.E{Key: "dir_id", Value: dirId}, bson.E{Key: "name", Value: fileName}, bson.E{Key: "deleted", Value: false}} // get all files if not deleted

	res := server.mongoColls[FilesCollName].FindOne(ctx, filter)

	if res.Err() == nil {
		return true
	}
	return false
}

// todo test
