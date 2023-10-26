package main

import (
	"context"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (server *SantaclausServerImpl) findDirFromPath(ctx context.Context, dirPath string, userId primitive.ObjectID) (dir *directory, r error) {

	if dirPath == "/" { // todo get rid of that
		dir, r := server.checkRootDirExistence(ctx, userId)
		if r != nil {
			return nil, r
		}
		return dir, nil
	}
	var directories []directory
	var didFound bool = false
	dirName := filepath.Base(dirPath)
	dirBase := filepath.Dir(dirPath)
	targetDirectory := bson.D{bson.E{Key: "name", Value: dirName}, bson.E{Key: "user_id", Value: userId}}
	cur, r := server.mongoColls[DirectoriesCollName].Find(ctx, targetDirectory /*, &targetDirectoryOptions*/)

	if r != nil {
		return nil, r
	}
	r = cur.All(ctx, &directories)
	if r != nil {
		return nil, r
	}
	for _, tmpDir := range directories {
		dirName = filepath.Base(dirPath)
		dirBase = filepath.Dir(dirPath)
		for {
			if tmpDir.ParentId == primitive.NilObjectID {
				didFound = true
				break
			}
			dirName = filepath.Base(dirBase)
			dirBase = filepath.Dir(dirBase)
			targetDirectory = bson.D{bson.E{Key: "_id", Value: tmpDir.ParentId}, bson.E{Key: "name", Value: dirName}, bson.E{Key: "user_id", Value: userId}}
			res := server.mongoColls[DirectoriesCollName].FindOne(ctx, targetDirectory /*, &targetDirectoryOptions*/)
			if res == nil {
				break
			}
			r = res.Decode(&tmpDir)
			if r != nil {
				return nil, r
			}
		}
		if didFound {
			break
		}
	}
	return dir, nil
}

func (server *SantaclausServerImpl) findPathFromDir(ctx context.Context, dirId primitive.ObjectID) (dirPath string, r error) {
	var currentDir directory
	nextId := dirId

	for nextId != primitive.NilObjectID {
		r = server.mongoColls[DirectoriesCollName].FindOne(ctx, bson.D{bson.E{Key: "_id", Value: nextId}}).Decode(&currentDir)
		if r != nil {
			return "", r
		}
		dirPath = filepath.Join(currentDir.Name, dirPath)
		nextId = currentDir.ParentId
	}
	return dirPath, nil
}
