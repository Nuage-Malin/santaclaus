package main

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	context "context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *SantaclausServerImpl) AddFile(ctx context.Context, req *MaeSanta.AddFileRequest) (status *MaeSanta.AddFileStatus, r error) {
	userId, err := primitive.ObjectIDFromHex(req.File.UserId)
	if err != nil {
		log.Fatal(err)
	}
	foundDirectory, err := server.findDirFromPath(req.File.DirPath, userId)
	if err != nil {
		log.Fatal(err)
		// TODO do something
	}

	// TODO find diskId

	foundDisk := server.findAvailableDisk(req.FileSize, req.File.UserId)
	newFile := file{
		Id:         primitive.NewObjectID(),
		Name:       req.File.Name,
		DirId:      foundDirectory.Id, // TODO find dirId from dirpath
		UserId:     userId,
		Size:       req.FileSize,
		DiskId:     foundDisk.Id,
		LastUpload: time.Now(),
		CreatedAt:  time.Now(),
		EditedAt:   time.Now(),
		Deleted:    false,
	}
	insertRes, err := server.mongoColls[FilesCollName].InsertOne(server.ctx, newFile)
	if err != nil {
		log.Fatal(err)
	}
	newFileId, ok := insertRes.InsertedID.(primitive.ObjectID)

	if ok == false {
		log.Println("Wrong type assertion!")
		// TODO check
	}

	status = &MaeSanta.AddFileStatus{
		FileId: newFileId.Hex(),
		DiskId: newFile.DiskId.String()}

	return status, nil
}

func (server *SantaclausServerImpl) VirtualRemoveFile(ctx context.Context, req *MaeSanta.RemoveFileRequest) (status *MaeSanta.RemoveFileStatus, r error) {
	fileId, err := primitive.ObjectIDFromHex(req.GetFileId())

	if err != nil {
		log.Fatal(err)
	}
	filter := bson.D{{"_id", fileId}}
	update := bson.D{{"$set", bson.D{{"deleted", true}}}} // only modify 'deleted' to true

	res, r := server.mongoColls[FilesCollName].UpdateOne(server.ctx, filter, update)
	if r != nil {
		return nil, r
	}
	if res.MatchedCount != 1 {
		return nil, fmt.Errorf("Could not find file %s\n", fileId)
	}
	if res.ModifiedCount != 1 {
		return nil, fmt.Errorf("Could not modify file %s\n", fileId)
	}
	// if res.UpsertedCount != 1 {
	// return nil, fmt.Errorf("Could not upsert file %s\n", fileId)
	// }
	status = &MaeSanta.RemoveFileStatus{}
	return status, r
}

func (server *SantaclausServerImpl) PhysicalRemoveFile(ctx context.Context, req *MaeSanta.RemoveFileRequest) (status *MaeSanta.RemoveFileStatus, r error) {
	fileId, err := primitive.ObjectIDFromHex(req.GetFileId())
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.D{{"_id", fileId}}
	// TODO find out more about contexts !!
	res, r := server.mongoColls[FilesCollName].DeleteOne(server.ctx, filter)
	if r != nil {
		return nil, r
	}
	if res.DeletedCount != 1 {
		return nil, fmt.Errorf("Deleted %d instead of 1", res.DeletedCount)
	}
	status = &MaeSanta.RemoveFileStatus{}
	return status, nil
}

func (server *SantaclausServerImpl) MoveFile(_ context.Context, req *MaeSanta.MoveFileRequest) (status *MaeSanta.MoveFileStatus, r error) {

	// filter := bson.D{primitive.E{Key: "_id", Value: req.GetFileId()}}
	// file := server.mongoColls[FilesCollName].FindOne(server.ctx, filter)
	// server.findDirFromPath(file.dirBase(req.GetFilepath()), /* file. find user id from file */)

	// update := bson.D{primitive.E{Key: "dirId", Value: /* new directory id */}}

	// modify dir Id
	// server.mongoColls[FilesCollName].UpdateOne(server.ctx, filter, update)

	return status, r
}

func (server *SantaclausServerImpl) GetFile(_ context.Context, req *MaeSanta.GetFileRequest) (*MaeSanta.GetFileStatus, error) {
	var fileFound file
	fileId, err := primitive.ObjectIDFromHex(req.FileId)

	if err != nil {
		log.Fatal(err)
	}
	filter := bson.D{{"_id", fileId}}
	err = server.mongoColls[FilesCollName].FindOne(server.ctx, filter).Decode(&fileFound)
	if err != nil {
		return nil, err
	}
	status := &MaeSanta.GetFileStatus{
		File: &MaeSanta.FileApproxMetadata{
			Name:    fileFound.Name,
			DirPath: server.findPathFromDir(fileFound.DirId),
			UserId:  fileFound.UserId.Hex()},
		DiskId: fileFound.DiskId.Hex()}
	return status, err
	/* todo is this the way of returning errors ? */
}

func (server *SantaclausServerImpl) UpdateFileSuccess(_ context.Context, req *MaeSanta.UpdateFileSuccessRequest) (status *MaeSanta.UpdateFileSuccessStatus, r error) {
	fileId, err := primitive.ObjectIDFromHex(req.GetFileId())
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.D{{"_id", fileId}}
	update := bson.D{{"$set", bson.D{{"size", req.GetNewFileSize()}}}}

	res, r := server.mongoColls[FilesCollName].UpdateOne(server.ctx, filter, update)
	if r != nil {
		return nil, r
	}
	if res.MatchedCount != 1 {
		return nil, fmt.Errorf("Could not find file %s\n", fileId)
	}
	if res.ModifiedCount != 1 {
		return nil, fmt.Errorf("Could not modify file %s\n", fileId)
	}
	// if res.UpsertedCount != 1 {
	// return nil, fmt.Errorf("Could not upsert file %s\n", fileId)
	// }
	return status, r
}

func (server *SantaclausServerImpl) ChangeFileDisk(_ context.Context, req *MaeSanta.ChangeFileDiskRequest) (status *MaeSanta.ChangeFileDiskStatus, r error) {

	filter := bson.D{primitive.E{Key: "_id", Value: req.GetFileId()}}
	// find the file in order not to put it on the same disk as it is already
	server.mongoColls[FilesCollName].FindOne(server.ctx, filter)

	// TODO algorithm to find new disk
	// find disk where
	//		there is some other file from this user
	// 	there is enough space for the file (and a bit more)

	// filter = bson.D{primitive.E{Key: "diskId", Value: /* value found from last request */}, primitive.E{Key: "userId", Value: /* value found from last request */}}
	// todo exclude from filter diskId that is the actual
	// server.mongoColls[FilesCollName].Find(server.ctx, filter, update)

	// update := bson.D{primitive.E{Key: "size", Value: /* new disk id */}}
	// server.mongoColls[FilesCollName].UpdateOne(server.ctx, filter, update)

	return status, r
}

// Directories
func (server *SantaclausServerImpl) AddDirectory(_ context.Context, req *MaeSanta.AddDirectoryRequest) (status *MaeSanta.AddDirectoryStatus, r error) {
	// find parent ID from req.Directory.DirPath
	userId, err := primitive.ObjectIDFromHex(req.Directory.UserId)
	if err != nil {
		log.Fatal(err)
	}
	parentDir, err := server.findDirFromPath(req.Directory.DirPath, userId)
	if err != nil {
		// TODO check error in another way than that
		log.Fatal(err)
	}
	dir := server.createDir(userId, parentDir.Id, req.Directory.Name)
	status = &MaeSanta.AddDirectoryStatus{DirId: dir.Id.Hex()}
	return status, r
}

func (server *SantaclausServerImpl) RemoveDirectory(context.Context, *MaeSanta.RemoveDirectoryRequest) (status *MaeSanta.RemoveDirectoryStatus, r error) {
	// remove all files
	// server.server.mongoColls[filesCollsName].FindAndDelete(/* filter with dirID */)
	// server.server.mongoColls[directoriesCollsName].FindAndDelete(/* fileter with dirID */)
	return status, r
}
func (server *SantaclausServerImpl) MoveDirectory(context.Context, *MaeSanta.MoveDirectoryRequest) (status *MaeSanta.MoveDirectoryStatus, r error) {
	// - add directory
	// - change files' directory Id
	// - remove directory
	return status, r
}

func (server *SantaclausServerImpl) getOneDirectory(dirId primitive.ObjectID, recursive bool, dirPath string, status *MaeSanta.GetDirectoryStatus) (*MaeSanta.GetDirectoryStatus, error) {

	// todo special case for root dir :
	// if queried with nil dirId, queries root dir
	filter := bson.D{{"dir_id", dirId}}
	var files []file
	var dir directory
	var dirs []directory
	cur, err := server.mongoColls[FilesCollName].Find(server.ctx, filter)

	if err != nil {
		return status, err
	}
	err = cur.All(server.ctx, &files)

	if err != nil {
		return status, err
	}

	for _, file := range files {

		if file.Deleted == false {
			metadata := &MaeSanta.FileMetadata{
				ApproxMetadata: &MaeSanta.FileApproxMetadata{Name: file.Name, DirPath: dirPath, UserId: file.UserId.Hex()},
				FileId:         file.Id.Hex(),
				IsDownloadable: false,             /* TODO change by real stored field */
				LastEditorId:   file.UserId.Hex(), /* TODO ? */
				Creation:       &timestamppb.Timestamp{Seconds: 0 /* TODO file.CreatedAt */},
				LastEdit:       &timestamppb.Timestamp{Seconds: 0 /* TODO file.EditedAt */}}
			status.SubFiles.FileIndex = append(status.SubFiles.FileIndex, metadata)
		}
	}

	// todo find all files in current directory

	filter = bson.D{{Key: "_id", Value: dirId}}
	err = server.mongoColls[DirectoriesCollName].FindOne(server.ctx, filter).Decode(&dir)

	status.SubFiles.DirIndex = append(status.SubFiles.DirIndex,
		&MaeSanta.DirMetadata{
			ApproxMetadata: &MaeSanta.FileApproxMetadata{
				Name:    dir.Name,
				DirPath: filepath.Dir(server.findPathFromDir(dir.Id)),
				UserId:  dir.UserId.Hex(),
			},
			DirId: dir.Id.Hex()})
	if err == nil {
		return status, err
	}
	if recursive {
		/* find all children directories thanks to a request with their parent ID (which is the current dirId) */
		filter := bson.D{primitive.E{Key: "parent_id", Value: dirId}}
		childDirIds, err := server.mongoColls[DirectoriesCollName].Find(server.ctx, filter)
		if err != nil {
			log.Fatal(err)
		}

		err = childDirIds.All(server.ctx, &dirs)
		if err != nil {
			log.Fatal(err)
		}
		for _, dir := range dirs {
			status, err = server.getOneDirectory(dir.Id, recursive, filepath.Join(dirPath, dir.Name), status)
			if err != nil {
				return status, err
			}
		}
	}
	return status, nil
}

func (server *SantaclausServerImpl) GetDirectory(_ context.Context, req *MaeSanta.GetDirectoryRequest) (status *MaeSanta.GetDirectoryStatus, r error) {
	// TODO fetch all the files
	dirId, err := primitive.ObjectIDFromHex(req.DirId)
	if err != nil {
		log.Fatal(err)
	}
	status = &MaeSanta.GetDirectoryStatus{
		SubFiles: &MaeSanta.FilesIndex{
			FileIndex: []*MaeSanta.FileMetadata{},
			DirIndex:  []*MaeSanta.DirMetadata{}}}

	// if !primitive.IsValidObjectID(req.DirId) { // TODO use that instead of ObjectIDFromHex ?
	// err
	// }
	dirPath := server.findPathFromDir(dirId)
	status, r = server.getOneDirectory(dirId, req.IsRecursive, dirPath, status)
	return status, r
}
