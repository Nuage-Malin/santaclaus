package main

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	context "context"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *SantaclausServerImpl) getChildrenDirectories(dirId primitive.ObjectID, recursive bool, dirPath string, status *MaeSanta.GetDirectoryStatus) (*MaeSanta.GetDirectoryStatus, error) {
	/* find all children directories thanks to a request with their parent ID (which is the current dirId) */

	var dirs []directory
	filter := bson.D{primitive.E{Key: "parent_id", Value: dirId}, {"deleted", false}}
	childDirIds, err := server.mongoColls[DirectoriesCollName].Find(server.ctx, filter)

	if err != nil {
		// log.Fatal(err)
		return status, err
	}

	err = childDirIds.All(server.ctx, &dirs)
	if err != nil {
		// log.Fatal(err)
		return status, err
	}
	for _, dir := range dirs {
		status, err = server.getOneDirectory(dir.Id, recursive, filepath.Join(dirPath, dir.Name), status)
		if err != nil {
			return status, err
		}
	}
	return status, nil
}

func (server *SantaclausServerImpl) addFilesToIndex(dirId primitive.ObjectID, dirPath string, status *MaeSanta.GetDirectoryStatus) (*MaeSanta.GetDirectoryStatus, error) {
	var files []file
	filter := bson.D{{"dir_id", dirId}, {"deleted", false}} // get all files if not deleted
	cur, err := server.mongoColls[FilesCollName].Find(server.ctx, filter)

	if err != nil {
		return status, err
	}
	err = cur.All(server.ctx, &files)

	if err != nil {
		return status, err
	}

	for _, file := range files {
		metadata := &MaeSanta.FileMetadata{
			ApproxMetadata: &MaeSanta.FileApproxMetadata{Name: file.Name, DirPath: dirPath, UserId: file.UserId.Hex()},
			FileId:         file.Id.Hex(),
			IsDownloadable: false,             /* TODO change by real stored field */
			LastEditorId:   file.UserId.Hex(), /* TODO ? */
			Creation:       &timestamppb.Timestamp{Seconds: 0 /* TODO file.CreatedAt */},
			LastEdit:       &timestamppb.Timestamp{Seconds: 0 /* TODO file.EditedAt */}}
		status.SubFiles.FileIndex = append(status.SubFiles.FileIndex, metadata)
	}
	return status, nil
}

func (server *SantaclausServerImpl) addOneDirectoryToIndex(dirId primitive.ObjectID, status *MaeSanta.GetDirectoryStatus) (*MaeSanta.GetDirectoryStatus, error) {
	var dir directory
	filter := bson.D{{"_id", dirId}, {"deleted", false}} // get the directory if exists and not deleted
	err := server.mongoColls[DirectoriesCollName].FindOne(server.ctx, filter).Decode(&dir)

	if err != nil {
		return status, err
	}
	status.SubFiles.DirIndex = append(status.SubFiles.DirIndex,
		&MaeSanta.DirMetadata{
			ApproxMetadata: &MaeSanta.FileApproxMetadata{
				Name:    dir.Name,
				DirPath: filepath.Dir(server.findPathFromDir(dir.Id)),
				UserId:  dir.UserId.Hex(),
			},
			DirId: dir.Id.Hex()})
	return status, nil
}

func (server *SantaclausServerImpl) getOneDirectory(dirId primitive.ObjectID, recursive bool, dirPath string, status *MaeSanta.GetDirectoryStatus) (*MaeSanta.GetDirectoryStatus, error) {

	// todo special case for root dir :
	// if queried with nil dirId, queries root dir
	// var dirs []directory
	status, err := server.addOneDirectoryToIndex(dirId, status)

	if err != nil {
		return status, err
	}
	status, err = server.addFilesToIndex(dirId, dirPath, status)
	if err != nil {
		return status, err
	}
	if recursive {
		status, err = server.getChildrenDirectories(dirId, recursive, dirPath, status)
	}
	return status, err
}

func (server *SantaclausServerImpl) GetRootDirectory(recursive bool, userId primitive.ObjectID, status *MaeSanta.GetDirectoryStatus) (*MaeSanta.GetDirectoryStatus, error) {
	var rootDir directory = server.checkRootDirExistance(userId) // creates root dir if inexistant, else return existant

	return server.getOneDirectory(rootDir.Id, recursive, "/", status)
}

func (server *SantaclausServerImpl) GetDirectory(_ context.Context, req *MaeSanta.GetDirectoryRequest) (status *MaeSanta.GetDirectoryStatus, r error) {
	userId, err := primitive.ObjectIDFromHex(req.UserId)

	status = &MaeSanta.GetDirectoryStatus{
		SubFiles: &MaeSanta.FilesIndex{
			FileIndex: []*MaeSanta.FileMetadata{},
			DirIndex:  []*MaeSanta.DirMetadata{}}}

	if req.DirId == nil {
		return server.GetRootDirectory(req.IsRecursive, userId, status)
	}
	dirId, err := primitive.ObjectIDFromHex(*req.DirId)

	if err != nil {
		return nil, r
	}
	if dirId == primitive.NilObjectID { // todo does it actually ever branches through that ?
		return server.GetRootDirectory(req.IsRecursive, userId, status)
	}

	if dirId != primitive.NilObjectID {
		return server.getOneDirectory(dirId, req.IsRecursive, server.findPathFromDir(dirId), status)
	} else { // todo does it actually ever branches through that ?
		return server.GetRootDirectory(req.IsRecursive, userId, status)
	}
}
