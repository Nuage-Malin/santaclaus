package main

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (server *SantaclausServerImpl) GetDirFromId(ctx context.Context, dirId primitive.ObjectID) (*directory, error) {

	filter := bson.D{bson.E{Key: "_id", Value: dirId}, bson.E{Key: "deleted", Value: false}} // get all files if not deleted
	var dirFound directory

	if server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(&dirFound) == nil {
		return &dirFound, nil
	}
	return nil, errors.New(fmt.Sprintf("Could not find directory %s", dirId.Hex()))
}

func (server *SantaclausServerImpl) GetDirFromStringId(ctx context.Context, strDirId string) (*directory, error) {

	dirId, r := primitive.ObjectIDFromHex(strDirId)

	if r != nil {
		return nil, r
	}
	filter := bson.D{bson.E{Key: "_id", Value: dirId}, bson.E{Key: "deleted", Value: false}} // get all files if not deleted
	var dirFound directory

	if server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(&dirFound) == nil {
		return &dirFound, nil
	}
	return nil, errors.New(fmt.Sprintf("Could not find directory %s", dirId))
}

// todo test
