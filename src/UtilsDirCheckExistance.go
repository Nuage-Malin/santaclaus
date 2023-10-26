package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/**
 * creates root dir if doesn't exists, otherwise return existing root dir
 */
// todo change name to getRootDir
func (server *SantaclausServerImpl) checkRootDirExistence(ctx context.Context, userId primitive.ObjectID) (rootDir *directory, err error) {

	targetDir := bson.D{bson.E{Key: "name", Value: "/"}, bson.E{Key: "user_id", Value: userId}}
	res := server.mongoColls[DirectoriesCollName].FindOne(ctx, targetDir)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {

			log.Println("Couldn't find root dir, creating it") // todo do logs so it prints only when debuging
			return server.createRootDir(ctx, userId)           // If the root directory doesn't exist, we create it
		}
		return nil, res.Err()
	}
	err = res.Decode(&rootDir)
	if err != nil {
		return nil, err
	}
	return rootDir, nil
}

// todo test
