package main

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"fmt"

	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (server *SantaclausServerImpl) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (status *pb.RemoveDirectoryStatus, r error) {
	log.Println("Request: RemoveUser") // todo replace with class request logger

	userId, r := primitive.ObjectIDFromHex(req.UserId)

	// directories
	filter := bson.D{bson.E{Key: "user_id", Value: userId}}
	findRes, r := server.mongoColls[FilesCollName].Find(ctx, filter) // Check if user exists by getting a file of his
	if r != nil {
		return nil, r
	}
	if findRes == nil || findRes.RemainingBatchLength() == 0 {
		return nil, fmt.Errorf("User %s doesn't exist\n", userId)
	}
	findRes, r = server.mongoColls[DirectoriesCollName].Find(ctx, filter) /* todo do the same for files ? */
	var dirs []directory

	if r != nil {
		return nil, r
	}
	r = findRes.All(ctx, &dirs)
	if r != nil {
		return nil, r
	}
	var filesToRemove pb.RemoveFilesRequest
	// var removeDirStatus pb.RemoveDirectoryStatus
	for _, dir := range dirs {
		removeDirStatus, r := server.RemoveDirectory(ctx, &pb.RemoveDirectoryRequest{DirId: dir.Id.String() /* todo use other method to convert to string ? */})
		if r != nil {
			log.Print(r)
			// Will print when removing sub directory of an already deleted directory, not a big problem but have to have it in mind when reading logs
			continue
		}
		for _, fileId := range removeDirStatus.FileIdsToRemove {
			filesToRemove.FileIds = append(filesToRemove.FileIds, fileId)
			status.FileIdsToRemove = append(status.FileIdsToRemove, fileId)
		}
	}
	//
	// files
	server.VirtualRemoveFiles(ctx, &filesToRemove)
	//
	status = &pb.RemoveDirectoryStatus{FileIdsToRemove: filesToRemove.FileIds}
	return status, nil

	// todo Test !
}
