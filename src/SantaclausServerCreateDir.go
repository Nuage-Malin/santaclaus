package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (server *SantaclausServerImpl) createDir(ctx context.Context, userId primitive.ObjectID, parentId primitive.ObjectID, name string) (dir *directory, r error) {

	var tmpParentDir *directory
	parendId := primitive.NilObjectID

	/// Check if parent dir exists
	if name == "/" {
		// function called from createRootDir
		// If creating root dir, do not check for parent
	} else {
		if parentId == primitive.NilObjectID { // nil object id means parent will be root
			// function called from gRPC request
			tmpParentDir, r = server.checkRootDirExistence(ctx, userId)
			if r != nil {
				return nil, r
			}
			parendId = tmpParentDir.Id
		} else { // Parent is not root
			tmpParentDir, r = server.GetDirFromId(ctx, parentId)
			if r != nil {
				return nil, r
			}
			parendId = tmpParentDir.Id
		}
		/// Check if directory name already exists in this directory
		if server.CheckDirHasChild(ctx, parendId, name) {
			return nil, errors.New("Directory with same name already exists in parent directory, aborting dir creation")
		}
	}

	dir = &directory{
		Id:        primitive.NewObjectID(),
		Name:      name,
		UserId:    userId,
		ParentId:  parendId,
		CreatedAt: time.Now(),
		EditedAt:  time.Now(),
		Deleted:   false}
	insertRes, r := server.mongoColls[DirectoriesCollName].InsertOne(ctx, dir)

	if insertRes == nil || r != nil {
		if r == nil {
			return nil, errors.New("Could not insert new directory")
		}
		return nil, r
	}
	dir.Id = insertRes.InsertedID.(primitive.ObjectID)
	return dir, nil
}

func (server *SantaclausServerImpl) createRootDir(ctx context.Context, userId primitive.ObjectID) (*directory, error) {

	return server.createDir(ctx, userId, primitive.NilObjectID, "/")
}
