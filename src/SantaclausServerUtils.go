package main

import (
	context "context"
	"log"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (server *SantaclausServerImpl) createDir(ctx context.Context, userId primitive.ObjectID, parentId primitive.ObjectID, name string) (dir directory, err error) {
	// check if directory already exists
	findRes := server.mongoColls[DirectoriesCollName].FindOne(ctx, primitive.D{{"parent_id", parentId}, {"name", name}, {"user_id", userId}})

	if findRes.Err() == nil {
		// if found existing directory, return it
		err = findRes.Decode(&dir)
		if err != nil {
			return dir, err
		}
		return dir, nil
	}

	dir = directory{
		Id:        primitive.NewObjectID(),
		Name:      name,
		UserId:    userId,
		ParentId:  parentId,
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

/**
 * creates root dir if doesn't exists, otherwise return existing root dir
 */
func (server *SantaclausServerImpl) checkRootDirExistence(ctx context.Context, userId primitive.ObjectID) (rootDir directory, err error) {

	targetDir := bson.D{{"name", "/"}, {"user_id", userId}}
	res := server.mongoColls[DirectoriesCollName].FindOne(ctx, targetDir)

	if res.Err() != nil {
		if err == mongo.ErrNoDocuments {

			log.Println("Couldn't find root dir, creating it") // todo do logs so it prints only when debuging
			return server.createRootDir(ctx, userId)           // If the root directory doesn't exist, we create it
		}
		return rootDir, res.Err()
	}
	err = res.Decode(&rootDir)
	return rootDir, err
}

func (server *SantaclausServerImpl) findDirFromPath(ctx context.Context, dirPath string, userId primitive.ObjectID) (dir directory, err error) {

	if dirPath == "/" { // todo get rid of that
		dir, err = server.checkRootDirExistence(ctx, userId)
		if dir.Id == primitive.NilObjectID || err != nil {
			return dir, err
		}
	}
	var directories []directory
	var tmpDir directory
	var didFind bool = false
	dirName := filepath.Base(dirPath)
	dirBase := filepath.Dir(dirPath)
	targetDirectory := bson.D{{"name", dirName}, {"user_id", userId}}
	cur, err := server.mongoColls[DirectoriesCollName].Find(ctx, targetDirectory /*, &targetDirectoryOptions*/)

	if err != nil {
		return dir, err
	}
	err = cur.All(ctx, &directories)
	if err != nil {
		return dir, err
	}
	for _, dir = range directories {
		dirName = filepath.Base(dirPath)
		dirBase = filepath.Dir(dirPath)
		tmpDir = dir
		for {
			if tmpDir.ParentId == primitive.NilObjectID {
				didFind = true
				break
			}
			dirName = filepath.Base(dirBase)
			dirBase = filepath.Dir(dirBase)
			targetDirectory = bson.D{{"_id", tmpDir.ParentId}, {"name", dirName}, {"user_id", userId}}
			res := server.mongoColls[DirectoriesCollName].FindOne(ctx, targetDirectory /*, &targetDirectoryOptions*/)
			if res == nil {
				break
			}
			err = res.Decode(&tmpDir)
			if err != nil {
				return dir, err
			}
		}
		if didFind {
			break
		}
	}
	return dir, nil
}

func (server *SantaclausServerImpl) findPathFromDir(ctx context.Context, dirId primitive.ObjectID) (dirPath string, err error) {
	var currentDir directory
	nextId := dirId

	for nextId != primitive.NilObjectID {
		err = server.mongoColls[DirectoriesCollName].FindOne(ctx, bson.D{{"_id", nextId}}).Decode(&currentDir)
		if err != nil {
			return "", err
		}
		dirPath = filepath.Join(currentDir.Name, dirPath)
		nextId = currentDir.ParentId
	}
	return dirPath, nil
}
