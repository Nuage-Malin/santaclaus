package main

import (
	"fmt"
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

func (server *SantaclausServerImpl) createRootDir(userId primitive.ObjectID) directory {

	var rootDir directory = directory{
		Id:        primitive.NewObjectID(), // todo make sure this is filled up by mongo
		Name:      "/",
		UserId:    userId,
		ParentId:  primitive.NilObjectID,
		CreatedAt: time.Now(),
		EditedAt:  time.Now(),
		Deleted:   false}
	// fmt.Printf("rootDir.Id.Hex(): %v\n", rootDir.Id.Hex())
	// fmt.Printf("rootDir.Id.String(): %v\n", rootDir.Id.String())
	res, err := server.mongoColls[DirectoriesCollName].InsertOne(server.ctx, rootDir)
	if res == nil || err != nil {
		log.Fatal(err)
	}
	rootDir.Id = res.InsertedID.(primitive.ObjectID)
	return rootDir
}

func (server *SantaclausServerImpl) checkRootDirExistance(userId primitive.ObjectID) directory {

	targetDir := bson.D{{"name", "/"}, {"user_id", userId}}
	res := server.mongoColls[DirectoriesCollName].FindOne(server.ctx, targetDir)
	if res.Err() != nil {

		fmt.Println(res.Err().Error())
		fmt.Println("Couldn't find root dir, creating it") // todo do logs so it prints only when debuging
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
	if dirPath == "/" {
		dir = server.checkRootDirExistance(userId)
		if dir.Id == primitive.NilObjectID {
			log.Fatalf("Error while getting root directory")
		}
	}
	dirName := filepath.Base(dirPath)
	dirBase := filepath.Dir(dirPath)
	tmp := make(map[string]interface{})
	var didFind bool = false

	targetDirectory := bson.D{primitive.E{Key: "Name", Value: dirName}, primitive.E{Key: "UserId", Value: userId}}
	cur, err := server.mongoColls[DirectoriesCollName].Find(server.ctx, targetDirectory /*, &targetDirectoryOptions*/)
	// checkErr(err)
	if err != nil {
		log.Fatal(err)
	}
	for i := cur; i != nil && i.Current != nil; i.Next(server.ctx) {
		dirName = filepath.Base(dirPath)
		dirBase = filepath.Dir(dirPath)
		err = i.Decode(tmp)
		// checkErr(err) // TODO this is not the solution, I guess ?
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("yo what's up")
				continue
			} else {
				log.Fatal(err)
			}
		}
		for {
			// if tmp.parentId == primitive.NilObjectID {
			if tmp["parentId"] == primitive.NilObjectID {
				fmt.Println("Found directory !")
				didFind = true
			}
			dirName = filepath.Base(dirBase)
			dirBase = filepath.Dir(dirBase)
			// targetDirectory = bson.D{primitive.E{Key: "dirId", Value: tmp.parentId}, primitive.E{Key: "name", Value: dirName}, primitive.E{Key: "userId", Value: userId}}
			targetDirectory = bson.D{primitive.E{Key: "_id", Value: tmp["ParentId"]}, primitive.E{Key: "Name", Value: dirName}, primitive.E{Key: "UserId", Value: userId}}
			res := server.mongoColls[DirectoriesCollName].FindOne(server.ctx, targetDirectory /*, &targetDirectoryOptions*/)

			if res == nil {
				break
			}

			err = res.Decode(tmp)
			// checkErr(err)
			if err != nil {
				log.Fatal(err)
			}

			// 	dirID should be formed from hashing name and parentId

			// targetDirectory := bson.D{{Key: "name", Value: filepath.Dir(dirPath)}, primitive.E{Key: "userId", Value: userId}}
			// targetDirectoryOptions := options.FindOneOptions{} // TODO refine search

		}
		if didFind {
			i.Decode(tmp)
			dir.Id = tmp["_id"].(primitive.ObjectID)
			dir.Name = tmp["name"].(string)
			dir.UserId = tmp["userId"].(primitive.ObjectID)
			dir.ParentId = tmp["parentId"].(primitive.ObjectID)
			// dir.DirId = tmp["dirId"].(primitive.ObjectID)
			dir.CreatedAt = tmp["createdAt"].(time.Time)
			dir.EditedAt = tmp["editedAt"].(time.Time)
			return dir, nil
		}
	}
	return dir, err // TODO learn about creating errors
}

func (server *SantaclausServerImpl) findAvailableDisk(fileSize uint64, userId string) (found disk /* TODO change by disk id*/) {
	filter := bson.D{}
	targetDiskOptions := options.FindOneOptions{Min: bson.D{primitive.E{Key: "availableSize", Value: fileSize}}} // TODO refine search

	/*res := */
	server.mongoColls[DisksCollName].FindOne(server.ctx, filter, &targetDiskOptions)

	return found
}
