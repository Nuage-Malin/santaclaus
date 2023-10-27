package main

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *SantaclausServerImpl) getChildrenDirectories(ctx context.Context, dirId primitive.ObjectID, recursive bool, status *pb.GetDirectoryStatus) (*pb.GetDirectoryStatus, error) {
	/* find all children directories thanks to a request with their parent ID (which is the current dirId) */

	var dirs []directory
	filter := bson.D{bson.E{Key: "parent_id", Value: dirId}, bson.E{Key: "deleted", Value: false}}
	childDirIds, r := server.mongoColls[DirectoriesCollName].Find(ctx, filter)

	if r != nil {
		return status, r
	}

	r = childDirIds.All(ctx, &dirs)
	if r != nil {
		return status, r
	}
	for _, dir := range dirs {
		if recursive {
			status, r = server.getOneDirectory(ctx, dir.Id, recursive, status)
		} else {
			status, r = server.addOneDirectoryToIndex(ctx, dir.Id, status)
		}
		if r != nil {
			return status, r
		}
	}
	return status, nil
}

func (server *SantaclausServerImpl) addFilesToIndex(ctx context.Context, dirId primitive.ObjectID, status *pb.GetDirectoryStatus) (*pb.GetDirectoryStatus, error) {
	var files []file
	filter := bson.D{bson.E{Key: "dir_id", Value: dirId}, bson.E{Key: "deleted", Value: false}} // get all files if not deleted
	cur, err := server.mongoColls[FilesCollName].Find(ctx, filter)

	if err != nil {
		return status, err
	}
	err = cur.All(ctx, &files)

	if err != nil {
		return status, err
	}

	for _, file := range files {
		metadata := &pb.FileMetadata{
			ApproxMetadata: &pb.FileApproxMetadata{Name: file.Name, DirId: dirId.Hex(), UserId: file.UserId.Hex()},
			FileId:         file.Id.Hex(),
			State:          pb.FileState_UNKNOWN,
			LastEditorId:   file.UserId.Hex(), /* TODO ? */
			Creation:       &timestamppb.Timestamp{Seconds: file.CreatedAt.Unix()},
			LastEdit:       &timestamppb.Timestamp{Seconds: file.EditedAt.Unix()}}
		status.SubFiles.FileIndex = append(status.SubFiles.FileIndex, metadata)
	}
	return status, nil
}

func (server *SantaclausServerImpl) addOneDirectoryToIndex(ctx context.Context, dirId primitive.ObjectID, status *pb.GetDirectoryStatus) (*pb.GetDirectoryStatus, error) {
	dir, r := server.GetDirFromId(ctx, dirId)

	if r != nil {
		return status, r
	}
	status.SubFiles.DirIndex = append(status.SubFiles.DirIndex,
		&pb.DirMetadata{
			ApproxMetadata: &pb.FileApproxMetadata{
				Name:   dir.Name,
				DirId:  dir.Id.Hex(),
				UserId: dir.UserId.Hex(),
			},
			DirId: dir.Id.Hex(),
			Creation: &timestamppb.Timestamp{Seconds: dir.CreatedAt.Unix()},
			LastEdit: &timestamppb.Timestamp{Seconds: dir.EditedAt.Unix()}})
	return status, nil
}

func (server *SantaclausServerImpl) getOneDirectory(ctx context.Context, dirId primitive.ObjectID, recursive bool, status *pb.GetDirectoryStatus) (*pb.GetDirectoryStatus, error) {

	// todo special case for root dir :
	// if queried with nil dirId, queries root dir
	// var dirs []directory
	status, err := server.addOneDirectoryToIndex(ctx, dirId, status)

	if err != nil {
		return status, err
	}
	status, err = server.addFilesToIndex(ctx, dirId, status)
	if err != nil {
		return status, err
	}
	status, err = server.getChildrenDirectories(ctx, dirId, recursive, status)
	// print(err)
	return status, err
}

func (server *SantaclausServerImpl) GetRootDirectory(ctx context.Context, recursive bool, userId primitive.ObjectID, status *pb.GetDirectoryStatus) (*pb.GetDirectoryStatus, error) {
	rootDir, err := server.checkRootDirExistence(ctx, userId) // creates root dir if inexistant, else return existant

	if err != nil {
		return nil, err
	}
	return server.getOneDirectory(ctx, rootDir.Id, recursive, status)
}

// todo func to check if a dir exists and get it (as a directory var)
