package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (server *SantaclausServerImpl) createDir(ctx context.Context, userId primitive.ObjectID, parentId primitive.ObjectID, name string) (dir directory, err error) {
	// Check if directory already exists

	var dirFound directory
	err = server.mongoColls[DirectoriesCollName].FindOne(ctx,
		bson.D{
			bson.E{Key: "parent_id", Value: parentId},
			bson.E{Key: "name", Value: name},
			bson.E{Key: "user_id", Value: userId}}).Decode(&dirFound)

	if err == nil {
		return dirFound, errors.New("Directory name already exists in this directory, aborting directory creation")
	}
	/* Check if parent dir exists */
	if name == "/" {
		// If creating root dir, do not check for parent
	} else if parentId == primitive.NilObjectID { // nil object id means parent will be root
		dirFound, err = server.checkRootDirExistence(ctx, userId)
		if err != nil {
			return dir, err
		}
	} else { // Parent is not root

		err = server.mongoColls[DirectoriesCollName].FindOne(ctx,
			bson.D{
				bson.E{Key: "_id", Value: parentId},
				bson.E{Key: "user_id", Value: userId}}).Decode(&dirFound)

		if err != nil {
			return dir, errors.New("Parent directory does not exist, aborting directory creation")
		}
	}

	dir = directory{
		Id:        primitive.NewObjectID(),
		Name:      name,
		UserId:    userId,
		ParentId:  dirFound.Id,
		CreatedAt: time.Now(),
		EditedAt:  time.Now(),
		Deleted:   false}
	insertRes, err := server.mongoColls[DirectoriesCollName].InsertOne(ctx, dir)

	if insertRes == nil || err != nil {
		return dir, err
	}
	dir.Id = insertRes.InsertedID.(primitive.ObjectID)
	return dir, nil
}

func (server *SantaclausServerImpl) createRootDir(ctx context.Context, userId primitive.ObjectID) (directory, error) {

	return server.createDir(ctx, userId, primitive.NilObjectID, "/")
}
