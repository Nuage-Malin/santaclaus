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
		fmt.Println("Wrong type assertion!")
		// TODO check
	}

	status = &MaeSanta.AddFileStatus{
		FileId: newFileId.Hex(),
		DiskId: newFile.DiskId.String()}

	return status, nil
}

func (server *SantaclausServerImpl) VirtualRemoveFile(ctx context.Context, req *MaeSanta.RemoveFileRequest) (status *MaeSanta.RemoveFileStatus, r error) {
	filter := bson.D{primitive.E{Key: "_id", Value: req.GetFileId()}}
	update := bson.D{primitive.E{Key: "deleted", Value: true}} // only modify 'deleted' to true

	server.mongoColls[FilesCollName].UpdateOne(server.ctx, filter, update)
	return status, r
}

func (server *SantaclausServerImpl) PhysicalRemoveFile(ctx context.Context, req *MaeSanta.RemoveFileRequest) (status *MaeSanta.RemoveFileStatus, r error) {
	filter := bson.D{primitive.E{Key: "_id", Value: req.GetFileId()}}

	// TODO find out more about contexts !
	server.mongoColls[FilesCollName].DeleteOne(server.ctx, filter)
	return status, r
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

func (server *SantaclausServerImpl) GetFile(_ context.Context, req *MaeSanta.GetFileRequest) (status *MaeSanta.GetFileStatus, r error) {
	filter := bson.D{primitive.E{Key: "_id", Value: req.GetFileId()}}

	/* file := */
	server.mongoColls[FilesCollName].FindOne(server.ctx, filter)
	return status, r
}

func (server *SantaclausServerImpl) UpdateFileSuccess(_ context.Context, req *MaeSanta.UpdateFileSuccessRequest) (status *MaeSanta.UpdateFileSuccessStatus, r error) {
	filter := bson.D{primitive.E{Key: "_id", Value: req.GetFileId()}}
	update := bson.D{primitive.E{Key: "size", Value: req.GetNewFileSize()}}

	server.mongoColls[FilesCollName].UpdateOne(server.ctx, filter, update)
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
	// TODO algo to create new directory :
	// 		hash of parentId and name

	// find parent ID from req.Directory.DirPath
	userId, err := primitive.ObjectIDFromHex(req.Directory.UserId)
	if err != nil {
		log.Fatal(err)
	}
	parentDir, err := server.findDirFromPath(filepath.Dir(req.Directory.DirPath), userId)
	// TODO check error other than that
	// checkErr(err)
	if err != nil {
		log.Fatal(err)
	}
	newDirectory := directory{
		Name:      filepath.Base(req.Directory.Name),
		UserId:    userId,
		ParentId:  parentDir.Id, // todo make sure ID is correct
		CreatedAt: time.Now(),
		EditedAt:  time.Now()}
	res, err := server.mongoColls[DirectoriesCollName].InsertOne(server.ctx, newDirectory)
	if err != nil {
		// log.Logger()
	}

	// res.InsertedID.Decode(status.DirId)

	var ok bool
	status.DirId, ok = res.InsertedID.(string)
	if !ok {
		// err
	}
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

func (server *SantaclausServerImpl) getOneDirectory(dirId primitive.ObjectID, recursive bool, dirPath string, status *MaeSanta.GetDirectoryStatus) *MaeSanta.GetDirectoryStatus {

	filter := bson.D{primitive.E{Key: "dirId", Value: dirId}}
	dirFound := server.mongoColls[FilesCollName].FindOne(server.ctx, filter)

	var currentDir directory
	dirFound.Decode(currentDir) // todo make sure that its ID is correct within field Id
	// todo not use Decode but marshall
	currentMetadata := MaeSanta.FileApproxMetadata{Name: currentDir.Name, DirPath: dirPath /* TODO if not exists, find directory path from dir ID */, UserId: currentDir.UserId.String()}
	status.Directories = append(status.Directories, &currentMetadata)
	if recursive {
		/* find all children directories thanks to a request with their parent ID (which is the current dirId) */
		filter := bson.D{primitive.E{Key: "parentId", Value: dirId}}
		childDirIds, err := server.mongoColls[DirectoriesCollName].Find(server.ctx, filter)
		// checkErr(err)
		if err != nil {
			log.Fatal(err)
		}
		for i := childDirIds; i != nil; i.Next(server.ctx) {
			i.Decode(currentDir)
			status = server.getOneDirectory(currentDir.Id, recursive, filepath.Join(dirPath, currentDir.Name), status)
		}
	}
	return status
}

func (server *SantaclausServerImpl) GetDirectory(_ context.Context, req *MaeSanta.GetDirectoryRequest) (status *MaeSanta.GetDirectoryStatus, r error) {
	// TODO fetch all the files
	if !primitive.IsValidObjectID(req.DirId) {
		// err
	}
	objID, err := primitive.ObjectIDFromHex(req.DirId)
	// checkErr(err)
	if err != nil {
		log.Fatal(err)
	}
	status = server.getOneDirectory(objID, req.IsRecursive, "", status)
	return status, r
}

// func (server *SantaclausServerImpl) mustEmbedUnimplementedMaestro_Santaclaus_ServiceServer() {
// fmt.Println("hello")
// }
