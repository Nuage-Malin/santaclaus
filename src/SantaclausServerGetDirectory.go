package main

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	context "context"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *SantaclausServerImpl) getChildrenDirectories(ctx context.Context, dirId primitive.ObjectID, recursive bool, dirPath string, status *MaeSanta.GetDirectoryStatus) (*MaeSanta.GetDirectoryStatus, error) {
	/* find all children directories thanks to a request with their parent ID (which is the current dirId) */

	var dirs []directory
	filter := bson.D{primitive.E{Key: "parent_id", Value: dirId}, {"deleted", false}}
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
			status, err = server.getOneDirectory(ctx, dir.Id, recursive, filepath.Join(dirPath, dir.Name), status)
		} else {
			status, err = server.addOneDirectoryToIndex(ctx, dir.Id, status)
		}
		if err != nil {
			return status, err
		}
	}
	return status, nil
}

func (server *SantaclausServerImpl) addFilesToIndex(ctx context.Context, dirId primitive.ObjectID, dirPath string, status *MaeSanta.GetDirectoryStatus) (*MaeSanta.GetDirectoryStatus, error) {
	var files []file
	filter := bson.D{{"dir_id", dirId}, {"deleted", false}} // get all files if not deleted
	cur, err := server.mongoColls[FilesCollName].Find(ctx, filter)

	if err != nil {
		return status, err
	}
	err = cur.All(ctx, &files)

	if err != nil {
		return status, err
	}

	for _, file := range files {
		metadata := &MaeSanta.FileMetadata{
			ApproxMetadata: &MaeSanta.FileApproxMetadata{Name: file.Name, DirPath: dirPath, UserId: file.UserId.Hex()},
			FileId:         file.Id.Hex(),
			IsDownloadable: false,             /* TODO change by real stored field */
			LastEditorId:   file.UserId.Hex(), /* TODO ? */
			Creation:       &timestamppb.Timestamp{Seconds: file.CreatedAt.Unix()},
			LastEdit:       &timestamppb.Timestamp{Seconds: file.EditedAt.Unix()}}
		status.SubFiles.FileIndex = append(status.SubFiles.FileIndex, metadata)
	}
	return status, nil
}

func (server *SantaclausServerImpl) addOneDirectoryToIndex(ctx context.Context, dirId primitive.ObjectID, status *MaeSanta.GetDirectoryStatus) (*MaeSanta.GetDirectoryStatus, error) {
	var dir directory
	filter := bson.D{{"_id", dirId}, {"deleted", false}} // get the directory if exists and not deleted
	err := server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(&dir)

	if err != nil {
		return status, err
	}
	dirPath, err := server.findPathFromDir(ctx, dir.Id)
	if err != nil {
		return nil, err
	}
	dirPath = filepath.Dir(dirPath)
	status.SubFiles.DirIndex = append(status.SubFiles.DirIndex,
		&MaeSanta.DirMetadata{
			ApproxMetadata: &MaeSanta.FileApproxMetadata{
				Name:    dir.Name,
				DirPath: dirPath,
				UserId:  dir.UserId.Hex(),
			},
			DirId: dir.Id.Hex()})
	return status, nil
}

func (server *SantaclausServerImpl) getOneDirectory(ctx context.Context, dirId primitive.ObjectID, recursive bool, dirPath string, status *MaeSanta.GetDirectoryStatus) (*MaeSanta.GetDirectoryStatus, error) {

	// todo special case for root dir :
	// if queried with nil dirId, queries root dir
	// var dirs []directory
	status, err := server.addOneDirectoryToIndex(ctx, dirId, status)

	if err != nil {
		return status, err
	}
	status, err = server.addFilesToIndex(ctx, dirId, dirPath, status)
	if err != nil {
		return status, err
	}
	status, err = server.getChildrenDirectories(ctx, dirId, recursive, dirPath, status)
	// print(err)
	return status, err
}

func (server *SantaclausServerImpl) GetRootDirectory(ctx context.Context, recursive bool, userId primitive.ObjectID, status *MaeSanta.GetDirectoryStatus) (*MaeSanta.GetDirectoryStatus, error) {
	rootDir, err := server.checkRootDirExistence(ctx, userId) // creates root dir if inexistant, else return existant

	if err != nil {
		return nil, err
	}
	return server.getOneDirectory(ctx, rootDir.Id, recursive, "/", status)
}

/*
 * If dirId is nil, return root directory
 */
func (server *SantaclausServerImpl) GetDirectory(ctx context.Context, req *MaeSanta.GetDirectoryRequest) (status *MaeSanta.GetDirectoryStatus, r error) {
	println("Requets: GetDirectory")
	userId, r := primitive.ObjectIDFromHex(req.UserId)

	if r != nil {
		return nil, r
	}

	status = &MaeSanta.GetDirectoryStatus{
		SubFiles: &MaeSanta.FilesIndex{
			FileIndex: []*MaeSanta.FileMetadata{},
			DirIndex:  []*MaeSanta.DirMetadata{},
		},
	}

	if req.DirId == nil {
		return server.GetRootDirectory(ctx, req.IsRecursive, userId, status)
	}
	dirId, r := primitive.ObjectIDFromHex(*req.DirId)

	if r != nil {
		return nil, r
	}
	if dirId == primitive.NilObjectID { // todo does it actually ever branches through that ?
		return server.GetRootDirectory(ctx, req.IsRecursive, userId, status)
	}

	if dirId != primitive.NilObjectID {
		dirPath, err := server.findPathFromDir(ctx, dirId)
		if err != nil {
			return nil, err
		}

		return server.getOneDirectory(ctx, dirId, req.IsRecursive, dirPath, status)
	} else { // todo does it actually ever branches through that ?
		return server.GetRootDirectory(ctx, req.IsRecursive, userId, status)
	}
}
