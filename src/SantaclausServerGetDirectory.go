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
	filter := bson.D{primitive.E{Key: "parent_id", Value: dirId}, bson.E{Key: "deleted", Value: false}}
	childDirIds, err := server.mongoColls[DirectoriesCollName].Find(ctx, filter)

	if err != nil {
		return status, err
	}

	err = childDirIds.All(ctx, &dirs)
	if err != nil {
		return status, err
	}
	for _, dir := range dirs {
		if recursive {
			status, err = server.getOneDirectory(ctx, dir.Id, recursive, status)
		} else {
			status, err = server.addOneDirectoryToIndex(ctx, dir.Id, status)
		}
		if err != nil {
			return status, err
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
			State: 			pb.FileState_UNKNOWN,
			LastEditorId:   file.UserId.Hex(), /* TODO ? */
			Creation:       &timestamppb.Timestamp{Seconds: file.CreatedAt.Unix()},
			LastEdit:       &timestamppb.Timestamp{Seconds: file.EditedAt.Unix()}}
		status.SubFiles.FileIndex = append(status.SubFiles.FileIndex, metadata)
	}
	return status, nil
}

func (server *SantaclausServerImpl) addOneDirectoryToIndex(ctx context.Context, dirId primitive.ObjectID, status *pb.GetDirectoryStatus) (*pb.GetDirectoryStatus, error) {
	var dir directory
	filter := bson.D{bson.E{Key: "_id", Value: dirId}, bson.E{Key: "deleted", Value: false}} // get the directory if exists and not deleted
	err := server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(&dir)

	if err != nil {
		return status, err
	}
	status.SubFiles.DirIndex = append(status.SubFiles.DirIndex,
<<<<<<< HEAD
		&pb.DirMetadata{
			ApproxMetadata: &pb.FileApproxMetadata{
				Name:   dir.Name,
				DirId:  dir.Id.Hex(),
				UserId: dir.UserId.Hex(),
=======
		&MaeSanta.DirMetadata{
			ApproxMetadata: &MaeSanta.FileApproxMetadata{
				Name:    dir.Name,
				DirId: dir.Id.Hex(),
				UserId:  dir.UserId.Hex(),
>>>>>>> b69dc041435fac2675eca1cfebac41bf1ba3d99a
			},
			DirId: dir.Id.Hex()})
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
