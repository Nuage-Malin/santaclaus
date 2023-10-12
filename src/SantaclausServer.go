package main

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"

	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *SantaclausServerImpl) AddFile(ctx context.Context, req *pb.AddFileRequest) (status *pb.AddFileStatus, r error) {
	/* todo check if filename already exists */
	userId, err := primitive.ObjectIDFromHex(req.File.UserId)

	if err != nil {
		return status, err
	}
	dirId, err := primitive.ObjectIDFromHex(req.File.DirId)
	if err != nil {
		return status, err
		// TODO do something
	}

	// TODO find diskId

	foundDisk, r := server.findAvailableDisk(ctx, req.FileSize, req.File.UserId)
	// todo update disk size here ? or in other function ?
	if r != nil {
		return nil, r
	}
	newFile := file{
		Id:         primitive.NewObjectID(),
		Name:       req.File.Name,
		DirId:      dirId, // TODO find dirId from dirpath
		UserId:     userId,
		Size:       req.FileSize,
		DiskId:     foundDisk,
		LastUpload: time.Now(),
		CreatedAt:  time.Now(),
		EditedAt:   time.Now(),
		Deleted:    false,
	}
	insertRes, err := server.mongoColls[FilesCollName].InsertOne(ctx, newFile)
	if err != nil {
		return status, err
	}
	newFileId, ok := insertRes.InsertedID.(primitive.ObjectID)

	if ok == false {
		log.Println("Wrong type assertion!")
		// TODO check
	}

	status = &pb.AddFileStatus{
		FileId: newFileId.Hex(),
		DiskId: newFile.DiskId.Hex()}

	return status, nil
}

func (server *SantaclausServerImpl) VirtualRemoveFile(ctx context.Context, req *pb.RemoveFileRequest) (status *pb.RemoveFileStatus, r error) {
	fileId, err := primitive.ObjectIDFromHex(req.GetFileId())

	if err != nil {
		return status, err
	}
	filter := bson.D{{"_id", fileId}}
	update := bson.D{{"$set", bson.D{{"deleted", true}}}} // only modify 'deleted' to true

	res, r := server.mongoColls[FilesCollName].UpdateOne(ctx, filter, update)
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
	status = &pb.RemoveFileStatus{}
	return status, r
}

func (server *SantaclausServerImpl) VirtualRemoveFiles(ctx context.Context, req *pb.RemoveFilesRequest) (status *pb.RemoveFilesStatus, r error) {
	var fileId primitive.ObjectID
	var tmpErr error
	var filter bson.D
	var res *mongo.UpdateResult
	update := bson.D{{"$set", bson.D{{"deleted", true}}}}

	for _, tmpFileId := range req.FileIds {
		fileId, tmpErr = primitive.ObjectIDFromHex(tmpFileId)
		if tmpErr != nil {
			log.Print(tmpErr)
			r = tmpErr
			continue
		}
		filter = bson.D{{"_id", fileId}}
		res, tmpErr = server.mongoColls[FilesCollName].UpdateOne(ctx, filter, update)
		if tmpErr != nil {
			log.Print(tmpErr)
			r = tmpErr
			continue
		}
		if res.MatchedCount != 1 {
			log.Printf("Could not find file with id %s (in order to virtually delete it)", fileId)
			continue
		}
		if res.ModifiedCount != 1 {
			log.Printf("Could not modify file with id %s (in order to virtually delete it)", fileId)
			continue
		}
	}
	return status, r
}

func (server *SantaclausServerImpl) PhysicalRemoveFile(ctx context.Context, req *pb.RemoveFileRequest) (status *pb.RemoveFileStatus, r error) {
	fileId, err := primitive.ObjectIDFromHex(req.GetFileId())

	if err != nil {
		return status, err
	}
	filter := bson.D{{"_id", fileId}}
	res, r := server.mongoColls[FilesCollName].DeleteOne(ctx, filter)

	if r != nil {
		return nil, r
	}
	if res.DeletedCount != 1 {
		return nil, fmt.Errorf("Deleted %d instead of 1", res.DeletedCount)
	}
	status = &pb.RemoveFileStatus{}
	return status, nil
}

func (server *SantaclausServerImpl) PhysicalRemoveFiles(ctx context.Context, req *pb.RemoveFilesRequest) (status *pb.RemoveFilesStatus, r error) {
	var fileId primitive.ObjectID
	var tmpErr error
	var filter bson.D
	var res *mongo.DeleteResult

	for _, tmpFileId := range req.FileIds {
		fileId, tmpErr = primitive.ObjectIDFromHex(tmpFileId)
		if tmpErr != nil {
			log.Print(tmpErr)
			r = tmpErr
			continue
		}
		filter = bson.D{{"_id", fileId}}
		res, tmpErr = server.mongoColls[FilesCollName].DeleteOne(ctx, filter)
		if tmpErr != nil {
			log.Print(tmpErr)
			r = tmpErr
			continue
		}
		if res.DeletedCount != 1 {
			log.Printf("Could not delete file with id %s", fileId)
			continue
		}
	}
	return status, r
}

func (server *SantaclausServerImpl) MoveFile(ctx context.Context, req *pb.MoveFileRequest) (status *pb.MoveFileStatus, r error) {

	// todo if nil object id for dirId, move to root dir ?
	fileId, r := primitive.ObjectIDFromHex(req.GetFileId())

	if r != nil {
		return status, r
	}
	dirId, r := primitive.ObjectIDFromHex(req.GetDirId())

	if r != nil {
		return nil, r
	}
	// If dirId is incorrect, return error
	filter := bson.D{{"_id", dirId}}
	var dir directory
	r = server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(&dir)
	if r != nil {
		return nil, r
	}
	// note: If name already exists, no problem, as id uniquely identifies the file
	// todo change the behaviour described above, cause problem for directories

	filter = bson.D{{"_id", fileId}}
	update := bson.D{{"$set", bson.D{{"name", req.NewFileName}, {"dir_id", dirId}}}}
	res, r := server.mongoColls[FilesCollName].UpdateOne(ctx, filter, update) // todo test updateById
	if r != nil {
		return nil, r
	}
	if res.MatchedCount != 1 {
		return nil, errors.New("Could not find file to be updated")
	}
	if res.ModifiedCount != 1 {
		// log.Print(res.ModifiedCount)
		return nil, errors.New("Could not modify file to be updated")
	}

	status = &pb.MoveFileStatus{}

	return status, r
}

func (server *SantaclausServerImpl) GetFile(ctx context.Context, req *pb.GetFileRequest) (*pb.GetFileStatus, error) {
	var fileFound file
	fileId, err := primitive.ObjectIDFromHex(req.FileId)

	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", fileId}}
	err = server.mongoColls[FilesCollName].FindOne(ctx, filter).Decode(&fileFound)
	if err != nil {
		return nil, err
	}

	/* status := &pb.GetFileStatus{
	File: &pb.FileApproxMetadata{
		Name:    fileFound.Name,
		DirPath: server.findPathFromDir(fileFound.DirId),
		UserId:  fileFound.UserId.Hex()},
	DiskId: fileFound.DiskId.Hex()} */

	status := &pb.GetFileStatus{
		File: &pb.FileMetadata{
			ApproxMetadata: &pb.FileApproxMetadata{
				Name:   fileFound.Name,
				DirId:  fileFound.DirId.Hex(),
				UserId: fileFound.UserId.Hex()},
			FileId:         fileFound.Id.Hex(),
			DirId:          fileFound.DirId.Hex(),
			IsDownloadable: fileFound.Downloadable,
			LastEditorId:   primitive.NilObjectID.Hex(), // todo this field is useless
			Creation:       timestamppb.New(fileFound.CreatedAt),
			LastEdit:       timestamppb.New(fileFound.LastUpload),
		},
		DiskId: fileFound.DiskId.Hex()}
	// todo think about this :
	// in the case of using fileMetadata instead of approxMetadata, but I don't think it is usefull
	// because this procedure should only be called when trying to download a file
	// but then how do we know the file is downloadable ??
	return status, nil
}

func (server *SantaclausServerImpl) UpdateFileSuccess(ctx context.Context, req *pb.UpdateFileSuccessRequest) (status *pb.UpdateFileSuccessStatus, r error) {
	fileId, err := primitive.ObjectIDFromHex(req.GetFileId())

	if err != nil {
		return status, err
	}

	filter := bson.D{{"_id", fileId}}
	update := bson.D{{"$set", bson.D{{"size", req.GetNewFileSize()}}}}

	res, r := server.mongoColls[FilesCollName].UpdateOne(ctx, filter, update)
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

func (server *SantaclausServerImpl) ChangeFileDisk(ctx context.Context, req *pb.ChangeFileDiskRequest) (status *pb.ChangeFileDiskStatus, r error) {

	filter := bson.D{primitive.E{Key: "_id", Value: req.GetFileId()}}
	// find the file in order not to put it on the same disk as it is already
	server.mongoColls[FilesCollName].FindOne(ctx, filter)

	// TODO algorithm to find new disk
	// find disk where
	//		there is some other file from this user
	// 	there is enough space for the file (and a bit more)

	// filter = bson.D{primitive.E{Key: "diskId", Value: /* value found from last request */}, primitive.E{Key: "userId", Value: /* value found from last request */}}
	// todo exclude from filter diskId that is the actual
	// server.mongoColls[FilesCollName].Find(ctx, filter, update)

	// update := bson.D{primitive.E{Key: "size", Value: /* new disk id */}}
	// server.mongoColls[FilesCollName].UpdateOne(ctx, filter, update)

	return status, r
}

// Directories
func (server *SantaclausServerImpl) AddDirectory(ctx context.Context, req *pb.AddDirectoryRequest) (status *pb.AddDirectoryStatus, r error) {
	// find parent ID from req.Directory.DirPath
	userId, err := primitive.ObjectIDFromHex(req.Directory.UserId)

	/* todo check if dirId already exists */

	if err != nil {
		return nil, err
	}
	dirId, err := primitive.ObjectIDFromHex(req.Directory.DirId)
	if err != nil {
		// TODO check error in another way than that
		return nil, err
	}
	dir, r := server.createDir(ctx, userId, dirId, req.Directory.Name)
	if r != nil {
		return nil, r
	}
	status = &pb.AddDirectoryStatus{DirId: dir.Id.Hex()}
	return status, r
}

func (server *SantaclausServerImpl) removeOneDirectory(ctx context.Context, dirId *primitive.ObjectID, status *pb.RemoveDirectoryStatus) (r error) {
	var files []file
	filter := bson.D{{"_id", dirId}}
	update := bson.D{{"$set", bson.D{{"deleted", true}}}}
	updateRes, r := server.mongoColls[DirectoriesCollName].UpdateOne(ctx, filter, update)

	if r != nil {
		return r
	}
	if updateRes == nil {
		return errors.New("Could not delete directory")
	}
	if updateRes.MatchedCount != 1 {
		return errors.New("Could not find directory to delete")
	}
	if updateRes.ModifiedCount != 1 {
		return errors.New("Could not modify directory to delete")
	}
	filter = bson.D{{"dir_id", dirId}}
	// update := bson.D{{"$set", bson.D{{"deleted", true}}}} // todo update deleted ? dirId ?
	// res, r := server.mongoColls[DirectoriesCollName].Update(ctx, filter, update)
	findRes, r := server.mongoColls[FilesCollName].Find(ctx, filter)
	if r != nil {
		return r
	}
	r = findRes.All(ctx, &files)
	if r != nil {
		return r
	}
	for _, file := range files {
		status.FileIdsToRemove = append(status.FileIdsToRemove, file.Id.Hex())
	}
	return nil
}

/**
 * recursivelly remove directories from directory
 */
func (server *SantaclausServerImpl) removeSubDirectories(ctx context.Context, parentId *primitive.ObjectID, status *pb.RemoveDirectoryStatus) (r error) {
	var dirs []directory
	filter := bson.D{{"parent_id", parentId}}
	findRes, r := server.mongoColls[DirectoriesCollName].Find(ctx, filter)

	if r != nil {
		return r
	}
	r = findRes.All(ctx, &dirs)
	if r != nil {
		return r
	}
	for _, dir := range dirs {
		server.removeSubDirectories(ctx, &dir.Id, status)
		server.removeOneDirectory(ctx, &dir.Id, status)
	}
	return nil
}

func (server *SantaclausServerImpl) RemoveDirectory(ctx context.Context, req *pb.RemoveDirectoryRequest) (status *pb.RemoveDirectoryStatus, r error) {
	//	find all subdirectories recursively
	//	in each subdirectory, add fileIds to the status (fileIdsToRemove)
	//	set directories as deleted
	dirId, err := primitive.ObjectIDFromHex(req.DirId)

	if err != nil {
		return status, err
	}
	status = &pb.RemoveDirectoryStatus{FileIdsToRemove: make([]string, 0)}
	r = server.removeSubDirectories(ctx, &dirId, status)
	if r != nil {
		return status, r
	}
	r = server.removeOneDirectory(ctx, &dirId, status)
	return status, r
}

func (server *SantaclausServerImpl) MoveDirectory(ctx context.Context, req *pb.MoveDirectoryRequest) (status *pb.MoveDirectoryStatus, r error) {
	// todo if nil object id for dirId, move to root dir ?rectoryRequest) (status *pb.MoveDirectoryStatus, r error) {
	// todo if nil object id for dirId, move to root dir ?
	dirId, r := primitive.ObjectIDFromHex(req.GetDirId())

	if r != nil {
		return nil, r
	}
	var newLocationDirId primitive.ObjectID = primitive.NilObjectID
	if req.NewLocationDirId != nil {
		newLocationDirId, r = primitive.ObjectIDFromHex(req.GetNewLocationDirId())

		if r != nil {
			return nil, r
		}
	}
	var filter bson.D
	// If dirId is incorrect, return error
	var parentDir directory

	if newLocationDirId != primitive.NilObjectID {
		filter = bson.D{{"_id", newLocationDirId}}
		r = server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(&parentDir)
		if r != nil {
			return nil, r
		}
	}
	var update bson.D

	filter = bson.D{{"_id", dirId}}
	var dir directory
	r = server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(&dir)
	if r != nil {
		return nil, r
	}

	if newLocationDirId != primitive.NilObjectID && newLocationDirId != dirId {
		// todo refactor this piece of code
		// a lot of code could be removed cause it does the same thing twice
		if req.Name != nil {
			filter = bson.D{{"name", *req.Name}, {"parent_id", newLocationDirId}}
		} else {
			filter = bson.D{{"name", dir.Name}, {"parent_id", newLocationDirId}}
		}
		r = server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(&dir)
		if r == nil {
			return nil, errors.New("Directory name already exists in this directory, aborting move")
		}
		if parentDir.ParentId == dirId {
			return nil, errors.New("Cannot store directory in its subdirectory")
		}
		if req.Name != nil {
			update = bson.D{{"$set", bson.D{{"name", *req.Name}, {"parent_id", newLocationDirId}}}}
		} else {
			update = bson.D{{"$set", bson.D{{"name", dir.Name}, {"parent_id", newLocationDirId}}}}
		}
	} else {
		if req.Name == nil {
			return nil, errors.New("No new name nor new location id provided for directory move")
		}
		filter = bson.D{{"name", *req.Name}, {"parent_id", dir.ParentId}} // todo add * for every valid req.name ?
		r = server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(&dir)
		if r == nil {
			return nil, errors.New("Directory name already exists in this directory, aborting move")
		}
		// In order not to change the location, but only the name, in the parameter, specify the same dirId as actual and new dir
		update = bson.D{{"$set", bson.D{{"name", *req.Name}}}}
	}

	filter = bson.D{{"_id", dirId}}
	res, r := server.mongoColls[DirectoriesCollName].UpdateOne(ctx, filter, update)

	if r != nil {
		return nil, r
	}
	if res.MatchedCount != 1 {
		return nil, errors.New("Could not find file to be updated")
	}
	if res.ModifiedCount != 1 {
		// log.Print(res.ModifiedCount)
		return nil, errors.New("Could not modify file to be updated")
	}
	status = &pb.MoveDirectoryStatus{}
	return status, r
}

func (server *SantaclausServerImpl) VirtualRemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (status *pb.RemoveDirectoryStatus, r error) {
	userId, r := primitive.ObjectIDFromHex(req.UserId)

	// directories
	filter := bson.D{{"user_id", userId}}
	findRes, r := server.mongoColls[DirectoriesCollName].Find(ctx, filter)
	var dirs []directory

	// todo actually remove directories from database at some point

	if r != nil {
		return status, r
	}
	r = findRes.All(ctx, &dirs)
	if r != nil {
		return status, r
	}
	var filesToRemove pb.RemoveFilesRequest
	// var removeDirStatus pb.RemoveDirectoryStatus
	for _, dir := range dirs {
		removeDirStatus, r := server.RemoveDirectory(ctx, &pb.RemoveDirectoryRequest{DirId: dir.Id.String() /* todo use other method to convert to string ? */})
		if r != nil {
			log.Print(r)
		}
		for _, fileId := range removeDirStatus.FileIdsToRemove {
			filesToRemove.FileIds = append(filesToRemove.FileIds, fileId)
		}
	}
	//
	// files
	server.VirtualRemoveFiles(ctx, &filesToRemove)
	//
	return status, r
}

func (server *SantaclausServerImpl) PhysicalRemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (status *pb.RemoveUserStatus, r error) {
	userId, r := primitive.ObjectIDFromHex(req.UserId)
	filter := bson.D{{"user_id", userId}}

	_, r = server.mongoColls[FilesCollName].DeleteMany(ctx, filter)

	if r != nil {
		return status, r
	}
	// todo physical remove directory

	return status, r
}
