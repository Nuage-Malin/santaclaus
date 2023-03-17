package main

import (
	"log"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (server *SantaclausServerImpl) createDir(userId primitive.ObjectID, parentId primitive.ObjectID, name string) (dir directory) {
	// check if directory already exists
	findRes := server.mongoColls[DirectoriesCollName].FindOne(server.ctx, primitive.D{{"parent_id", parentId}, {"name", name}, {"user_id", userId}})

	if findRes.Err() == nil {
		// if found existing directory, return it
		err := findRes.Decode(&dir)
		if err != nil {
			log.Fatal(err)
		}
		return dir
	}

	dir = directory{
		Id:        primitive.NewObjectID(),
		Name:      name,
		UserId:    userId,
		ParentId:  parentId,
		CreatedAt: time.Now(),
		EditedAt:  time.Now(),
		Deleted:   false}
	insertRes, err := server.mongoColls[DirectoriesCollName].InsertOne(server.ctx, dir)

	if insertRes == nil || err != nil {
		log.Fatal(err)
	}
	dir.Id = insertRes.InsertedID.(primitive.ObjectID)
	return dir
}

func (server *SantaclausServerImpl) createRootDir(userId primitive.ObjectID) directory {

	return server.createDir(userId, primitive.NilObjectID, "/")
}

func (server *SantaclausServerImpl) checkRootDirExistance(userId primitive.ObjectID) directory {

	targetDir := bson.D{{"name", "/"}, {"user_id", userId}}
	res := server.mongoColls[DirectoriesCollName].FindOne(server.ctx, targetDir)

	if res.Err() != nil {
		log.Println(res.Err().Error())
		log.Println("Couldn't find root dir, creating it") // todo do logs so it prints only when debuging
		return server.createRootDir(userId)                // If the root directory doesn't exist, we create it
	}
	var rootDir directory
	err := res.Decode(&rootDir)
	if err != nil {
		log.Fatal(err)
	}
	return rootDir
}

func (server *SantaclausServerImpl) findDirFromPath(dirPath string, userId primitive.ObjectID) (directory, error) {

	var dir directory
	if dirPath == "/" { // todo get rid of that
		dir = server.checkRootDirExistance(userId)
		if dir.Id == primitive.NilObjectID {
			log.Fatalf("Error while getting root directory")
		}
	}
	var directories []directory
	var tmpDir directory
	var didFind bool = false
	dirName := filepath.Base(dirPath)
	dirBase := filepath.Dir(dirPath)
	targetDirectory := bson.D{{"name", dirName}, {"user_id", userId}}
	cur, err := server.mongoColls[DirectoriesCollName].Find(server.ctx, targetDirectory /*, &targetDirectoryOptions*/)

	if err != nil {
		log.Fatal(err)
	}
	err = cur.All(server.ctx, &directories)
	if err != nil {
		log.Fatal(err)
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
			res := server.mongoColls[DirectoriesCollName].FindOne(server.ctx, targetDirectory /*, &targetDirectoryOptions*/)
			if res == nil {
				break
			}
			err = res.Decode(&tmpDir)
			if err != nil {
				log.Fatal(err)
			}
		}
		if didFind {
			break
		}
	}
	return dir, err // TODO learn about creating errors
}

func (server *SantaclausServerImpl) findPathFromDir(dirId primitive.ObjectID) (dirPath string) {
	var currentDir directory
	nextId := dirId

	for nextId != primitive.NilObjectID {
		err := server.mongoColls[DirectoriesCollName].FindOne(server.ctx, bson.D{{"_id", nextId}}).Decode(&currentDir)
		if err != nil {
			log.Fatal(err)
		}
		dirPath = filepath.Join(currentDir.Name, dirPath)
		nextId = currentDir.ParentId
	}
	return dirPath
}

func (server *SantaclausServerImpl) findAvailableDisk(fileSize uint64, userId string) (found disk /* TODO change by disk id*/) {
	filter := bson.D{}
	targetDiskOptions := options.FindOneOptions{Min: bson.D{primitive.E{Key: "availableSize", Value: fileSize}}} // TODO refine search

	/*res := */
	server.mongoColls[DisksCollName].FindOne(server.ctx, filter, &targetDiskOptions)

	return found
}
