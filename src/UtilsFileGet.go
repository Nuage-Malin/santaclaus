package main

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (server *SantaclausServerImpl) GetFileFromId(ctx context.Context, fileId primitive.ObjectID) (*file, error) {
	filter := bson.D{bson.E{Key: "_id", Value: fileId}, bson.E{Key: "deleted", Value: false}} // get all files if not deleted
	var fileFound file

	if server.mongoColls[FilesCollName].FindOne(ctx, filter).Decode(&fileFound) == nil {
		return &fileFound, nil
	}
	return nil, errors.New(fmt.Sprintf("Could not find file %s", fileId.Hex()))
}

func (server *SantaclausServerImpl) GetFileFromStringId(ctx context.Context, strFileId string) (*file, error) {

	fileId, r := primitive.ObjectIDFromHex(strFileId)

	if r != nil {
		return nil, r
	}
	filter := bson.D{bson.E{Key: "_id", Value: fileId}, bson.E{Key: "deleted", Value: false}} // get all files if not deleted
	var fileFound file

	if server.mongoColls[FilesCollName].FindOne(ctx, filter).Decode(&fileFound) == nil {
		return &fileFound, nil
	}
	return nil, errors.New(fmt.Sprintf("Could not find directory %s", fileId))
}
