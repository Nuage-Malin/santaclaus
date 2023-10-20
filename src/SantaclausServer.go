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

// Files

func (server *SantaclausServerImpl) AddFile(ctx context.Context, req *pb.AddFileRequest) (status *pb.AddFileStatus, r error) {
	log.Println("Request: AddFile" /* todo try to get function name from variable or macro */) // todo replace with class request logger

	userId, r := primitive.ObjectIDFromHex(req.File.UserId)

	if r != nil {
		return status, r
	}
	dirId, r := primitive.ObjectIDFromHex(req.File.DirId)
	if r != nil {
		return status, r
	}
	/* Check if filename already exists */
	filter := bson.D{bson.E{Key: "name", Value: req.File.Name}, bson.E{Key: "dir_id", Value: dirId}}
	// var fileFound file

	var fileFound file
	if server.mongoColls[FilesCollName].FindOne(ctx, filter).Decode(&fileFound) == nil {

		status = &pb.AddFileStatus{
			FileId: fileFound.Id.Hex(),
			DiskId: fileFound.DiskId.Hex()}
		return nil, errors.New("File name already exists in this directory, aborting file creation")
	}

	foundDisk, r := server.findAvailableDisk(ctx, req.FileSize, req.File.UserId)
	// todo update disk size here ? or in other function ?
	if r != nil {
		return nil, r
	}
	newFile := file{
		Id:         primitive.NewObjectID(),
		Name:       req.File.Name,
		DirId:      dirId,
		UserId:     userId,
		Size:       req.FileSize,
		DiskId:     foundDisk,
		LastUpload: time.Now(),
		CreatedAt:  time.Now(),
		EditedAt:   time.Now(),
		Deleted:    false,
	}
	insertRes, r := server.mongoColls[FilesCollName].InsertOne(ctx, newFile)
	if r != nil {
		return status, r
	}
	newFileId, ok := insertRes.InsertedID.(primitive.ObjectID)

	if ok == false {
		log.Println("Wrong type assertion!")
		return
	}

	status = &pb.AddFileStatus{
		FileId: newFileId.Hex(),
		DiskId: newFile.DiskId.Hex()}

	return status, nil
}

func (server *SantaclausServerImpl) VirtualRemoveFile(ctx context.Context, req *pb.RemoveFileRequest) (status *pb.RemoveFileStatus, r error) {
	log.Println("Request: VirtualRemoveFile" /* todo try to get function name from variable or macro */) // todo replace with class request logger

	fileId, r := primitive.ObjectIDFromHex(req.GetFileId())

	if r != nil {
		return status, r
	}
	filter := bson.D{bson.E{Key: "_id", Value: fileId}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "deleted", Value: true}}}} // only modify 'deleted' to true

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
	return status, nil
}

func (server *SantaclausServerImpl) VirtualRemoveFiles(ctx context.Context, req *pb.RemoveFilesRequest) (status *pb.RemoveFilesStatus, r error) {
	log.Println("Request: VirtualRemoveFiles" /* todo try to get function name from variable or macro */) // todo replace with class request logger

	var fileId primitive.ObjectID
	var tmpErr error
	var filter bson.D
	var res *mongo.UpdateResult
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "deleted", Value: true}}}}

	for _, tmpFileId := range req.FileIds {
		fileId, tmpErr = primitive.ObjectIDFromHex(tmpFileId)
		if tmpErr != nil {
			log.Print(tmpErr)
			r = tmpErr
			continue
		}
		filter = bson.D{bson.E{Key: "_id", Value: fileId}}
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
	return status, nil
}

func (server *SantaclausServerImpl) PhysicalRemoveFile(ctx context.Context, req *pb.RemoveFileRequest) (status *pb.RemoveFileStatus, r error) {
	log.Println("Request: PhysicalRemoveFile" /* todo try to get function name from variable or macro */) // todo replace with class request logger

	fileId, r := primitive.ObjectIDFromHex(req.GetFileId())

	if r != nil {
		return status, r
	}
	filter := bson.D{bson.E{Key: "_id", Value: fileId}}
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
	log.Println("Request: VirtualRemoveFiles" /* todo try to get function name from variable or macro */) // todo replace with class request logger

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
		filter = bson.D{bson.E{Key: "_id", Value: fileId}}
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
	return status, nil
}

// todo differentiate move file and rename file
func (server *SantaclausServerImpl) MoveFile(ctx context.Context, req *pb.MoveFileRequest) (status *pb.MoveFileStatus, r error) {
	log.Println("Request: MoveFile" /* todo try to get function name from variable or macro */) // todo replace with class request logger

	if req.DirId == nil && req.NewFileName == nil {
		return nil, errors.New("No new directory and file name, abortin file move")
	}
	// todo if nil object id for dirId, move to root dir ?
	fileId, r := primitive.ObjectIDFromHex(req.GetFileId())

	if r != nil {
		return status, r
	}
	/// Check if file of fileId exists
	filter := bson.D{bson.E{Key: "_id", Value: fileId}}
	var currentFile file
	r = server.mongoColls[FilesCollName].FindOne(ctx, filter).Decode(&currentFile)
	if r != nil {
		return status, r
	}

	var tmpFileFound file
	var update bson.D

	// todo another function
	if req.DirId != nil {
		dirId, r := primitive.ObjectIDFromHex(req.GetDirId())

		if r != nil {
			return nil, r
		}
		// If dirId is incorrect, return error
		filter = bson.D{bson.E{Key: "_id", Value: dirId}}
		var dir directory
		r = server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(&dir)
		if r != nil {
			return nil, r
		}
		// Check if file with this name exists in the new directory
		filter = bson.D{bson.E{Key: "name", Value: currentFile.Name}, bson.E{Key: "dir_id", Value: dirId}}
		r = server.mongoColls[FilesCollName].FindOne(ctx, filter).Decode(&tmpFileFound)
		if r == nil {
			return nil, errors.New("File with this name already exists in the new directory")
		}
		filter = bson.D{bson.E{Key: "_id", Value: fileId}}
		update = bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "dir_id", Value: dirId}}}}
		res, r := server.mongoColls[FilesCollName].UpdateOne(ctx, filter, update) // todo test updateById
		if r != nil {
			return nil, r
		}
		if res.MatchedCount != 1 {
			return nil, errors.New("Could not find file to be updated")
		}
		if res.ModifiedCount != 1 {
			return nil, errors.New("Could not modify file directory")
		}
	}
	// todo another function
	if req.NewFileName != nil {
		newFileName := req.GetNewFileName()

		if r != nil {
			return nil, r
		}
		// Check if file with new name already exists in this directory
		filter = bson.D{bson.E{Key: "dir_id", Value: currentFile.DirId}, bson.E{Key: "name", Value: newFileName}}
		r = server.mongoColls[FilesCollName].FindOne(ctx, filter).Decode(&tmpFileFound)
		if r == nil {
			return nil, errors.New("File with this new name already exists, aborting move")
		}

		filter = bson.D{bson.E{Key: "_id", Value: fileId}}
		update = bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "name", Value: newFileName}}}}
		res, r := server.mongoColls[FilesCollName].UpdateOne(ctx, filter, update) // todo test updateById
		if r != nil {
			return nil, r
		}
		if res.MatchedCount != 1 {
			return nil, errors.New("Could not find file to be updated")
		}
		if res.ModifiedCount != 1 {
			return nil, errors.New("Could not modify file name")
		}
	}
	status = &pb.MoveFileStatus{}

	return status, nil
}

func (server *SantaclausServerImpl) GetFile(ctx context.Context, req *pb.GetFileRequest) (*pb.GetFileStatus, error) {
	log.Println("Request: GetFile" /* todo try to get function name from variable or macro */) // todo replace with class request logger

	var fileFound file
	fileId, r := primitive.ObjectIDFromHex(req.FileId)

	if r != nil {
		return nil, r
	}
	// filter := bson.D{bson.E{Key: "_id", Value: fileId}}
	filter := bson.D{bson.E{Key: "_id", Value: fileId}}
	r = server.mongoColls[FilesCollName].FindOne(ctx, filter).Decode(&fileFound)
	if r != nil {
		return nil, r
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
			FileId:       fileFound.Id.Hex(),
			DirId:        fileFound.DirId.Hex(),
			State:        pb.FileState_UNKNOWN,
			LastEditorId: primitive.NilObjectID.Hex(), // todo this field is useless
			Creation:     timestamppb.New(fileFound.CreatedAt),
			LastEdit:     timestamppb.New(fileFound.LastUpload),
		},
		DiskId: fileFound.DiskId.Hex()}
	// todo think about this :
	// in the case of using fileMetadata instead of approxMetadata, but I don't think it is usefull
	// because this procedure should only be called when trying to download a file
	// but then how do we know the file is downloadable ??
	return status, nil
}

func (server *SantaclausServerImpl) UpdateFileSuccess(ctx context.Context, req *pb.UpdateFileSuccessRequest) (status *pb.UpdateFileSuccessStatus, r error) {
	log.Println("Request: UpdateFileSuccess" /* todo try to get function name from variable or macro */) // todo replace with class request logger

	fileId, r := primitive.ObjectIDFromHex(req.GetFileId())

	if r != nil {
		return status, r
	}

	filter := bson.D{bson.E{Key: "_id", Value: fileId}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "size", Value: req.GetNewFileSize()}}}}

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
	return status, nil
}

// Disks

func (server *SantaclausServerImpl) ChangeFileDisk(ctx context.Context, req *pb.ChangeFileDiskRequest) (status *pb.ChangeFileDiskStatus, r error) {
	log.Println("Request: ChangeFileDisk" /* todo try to get function name from variable or macro */) // todo replace with class request logger

	filter := bson.D{bson.E{Key: "_id", Value: req.GetFileId()}}
	// find the file in order not to put it on the same disk as it is already
	server.mongoColls[FilesCollName].FindOne(ctx, filter)

	// TODO algorithm to find new disk
	// find disk where
	//		there is some other file from this user
	// 	there is enough space for the file (and a bit more)

	// filter = bson.D{bson.E{Key: "diskId", Value: /* value found from last request */}, bson.E{Key: "userId", Value: /* value found from last request */}}
	// todo exclude from filter diskId that is the actual
	// server.mongoColls[FilesCollName].Find(ctx, filter, update)

	// update := bson.D{bson.E{Key: "size", Value: /* new disk id */}}
	// server.mongoColls[FilesCollName].UpdateOne(ctx, filter, update)

	return status, nil
}

// Directories
func (server *SantaclausServerImpl) AddDirectory(ctx context.Context, req *pb.AddDirectoryRequest) (status *pb.AddDirectoryStatus, r error) {
	userId, err := primitive.ObjectIDFromHex(req.Directory.UserId)
	if err != nil {
		return nil, err
	}

	dirId, err := primitive.ObjectIDFromHex(req.Directory.DirId)
	if err != nil {
		return nil, err
	}

	dir, r := server.createDir(ctx, userId, dirId, req.Directory.Name)
	if r != nil {
		return nil, r
	}
	status = &pb.AddDirectoryStatus{DirId: dir.Id.Hex()}
	return status, nil
}

func (server *SantaclausServerImpl) RemoveDirectory(ctx context.Context, req *pb.RemoveDirectoryRequest) (status *pb.RemoveDirectoryStatus, r error) {
	log.Println("Request: RemoveDirectory" /* todo try to get function name from variable or macro */) // todo replace with class request logger

	//	find all subdirectories recursively
	//	in each subdirectory, add fileIds to the status (fileIdsToRemove)
	//	set directories as deleted
	dirId, r := primitive.ObjectIDFromHex(req.DirId)

	if r != nil {
		return nil, r
	}
	// Check if dir exists
	filter := bson.D{bson.E{Key: "_id", Value: dirId}}
	r = server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(nil)
	if r != nil {
		return nil, r
	}

	status = &pb.RemoveDirectoryStatus{FileIdsToRemove: make([]string, 0)}
	// Mark all subdirectories as deleted (Virtual)
	r = server.removeSubDirectories(ctx, &dirId, status)
	if r != nil {
		return status, r
	}
	// Mark this directory as deleted (Virtual)
	r = server.removeOneDirectory(ctx, &dirId, status)
	// Remove all directories that have been marked as deleted in recursive sub functions (Physical)
	filter = bson.D{bson.E{Key: "deleted", Value: true}}
	_, r = server.mongoColls[DirectoriesCollName].DeleteMany(ctx, filter)
	if r != nil {
		return nil, r
	}
	return status, nil
}

func (server *SantaclausServerImpl) MoveDirectory(ctx context.Context, req *pb.MoveDirectoryRequest) (status *pb.MoveDirectoryStatus, r error) {
	log.Println("Request: MoveDirectory" /* todo try to get function name from variable or macro */) // todo replace with class request logger

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
	var parentDir directory

	// todo check if newLocationDirId is a directory that exists and is a directory of this user ?
	if newLocationDirId != primitive.NilObjectID {
		filter = bson.D{bson.E{Key: "_id", Value: newLocationDirId}}
		r = server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(&parentDir)
		if r != nil {
			return nil, r
		}
	}
	var update bson.D

	filter = bson.D{bson.E{Key: "_id", Value: dirId}}
	var dir directory
	r = server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(&dir)
	if r != nil {
		return nil, r
	}

	if newLocationDirId != primitive.NilObjectID && newLocationDirId != dirId {
		// todo refactor this piece of code
		// a lot of code could be removed cause it does the same thing twice
		if req.Name != nil {
			filter = bson.D{bson.E{Key: "name", Value: *req.Name}, bson.E{Key: "parent_id", Value: newLocationDirId}}
		} else {
			filter = bson.D{bson.E{Key: "name", Value: dir.Name}, bson.E{Key: "parent_id", Value: newLocationDirId}}
		}
		r = server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(&dir)
		if r == nil {
			return nil, errors.New("Directory name already exists in this directory, aborting move")
		}
		if parentDir.ParentId == dirId {
			return nil, errors.New("Cannot store directory in its subdirectory")
		}
		if req.Name != nil {
			update = bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "name", Value: *req.Name}, bson.E{Key: "parent_id", Value: newLocationDirId}}}}
		} else {
			update = bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "name", Value: dir.Name}, bson.E{Key: "parent_id", Value: newLocationDirId}}}}
		}
	} else {
		if req.Name == nil {
			return nil, errors.New("No new name nor new location id provided for directory move")
		}
		filter = bson.D{bson.E{Key: "name", Value: *req.Name}, bson.E{Key: "parent_id", Value: dir.ParentId}} // todo add * for every valid req.name ?
		r = server.mongoColls[DirectoriesCollName].FindOne(ctx, filter).Decode(&dir)
		if r == nil {
			return nil, errors.New("Directory name already exists in this directory, aborting move")
		}
		// In order not to change the location, but only the name, in the parameter, specify the same dirId as actual and new dir
		update = bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "name", Value: *req.Name}}}}
	}

	filter = bson.D{bson.E{Key: "_id", Value: dirId}}
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
	return status, nil
}

/*
 * If dirId is nil, return root directory
 */
func (server *SantaclausServerImpl) GetDirectory(ctx context.Context, req *pb.GetDirectoryRequest) (status *pb.GetDirectoryStatus, r error) {
	log.Println("Request: GetDirectory" /* todo try to get function name from variable or macro */) // todo replace with class request logger

	userId, r := primitive.ObjectIDFromHex(req.UserId)

	if r != nil {
		return nil, r
	}

	status = &pb.GetDirectoryStatus{
		SubFiles: &pb.FilesIndex{
			FileIndex: []*pb.FileMetadata{},
			DirIndex:  []*pb.DirMetadata{},
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
	} else {
		return server.getOneDirectory(ctx, dirId, req.IsRecursive, status)
	}
}

// Users

func (server *SantaclausServerImpl) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (status *pb.RemoveDirectoryStatus, r error) {
	log.Println("Request: RemoveUser" /* todo try to get function name from variable or macro */) // todo replace with class request logger

	userId, r := primitive.ObjectIDFromHex(req.UserId)

	// directories
	filter := bson.D{bson.E{Key: "user_id", Value: userId}}
	findRes, r := server.mongoColls[DirectoriesCollName].Find(ctx, filter) /* todo do the same for files ? */
	var dirs []directory

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
			// Will print when removing sub directory of an already deleted directory, not a big problem but have to have it in mind when reading logs
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
	return status, nil
}
