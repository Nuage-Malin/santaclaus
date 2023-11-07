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
	log.Println("Request: AddFile") // todo replace with class request logger

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
	log.Println("Request: VirtualRemoveFile") // todo replace with class request logger

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
	log.Println("Request: VirtualRemoveFiles") // todo replace with class request logger

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
	log.Println("Request: PhysicalRemoveFile") // todo replace with class request logger

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
	log.Println("Request: VirtualRemoveFiles") // todo replace with class request logger

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

func (server *SantaclausServerImpl) GetFile(ctx context.Context, req *pb.GetFileRequest) (*pb.GetFileStatus, error) {
	log.Println("Request: GetFile") // todo replace with class request logger

	var fileFound file
	fileId, r := primitive.ObjectIDFromHex(req.FileId)

	if r != nil {
		return nil, r
	}
	filter := bson.D{bson.E{Key: "_id", Value: fileId}}
	r = server.mongoColls[FilesCollName].FindOne(ctx, filter).Decode(&fileFound)
	if r != nil {
		return nil, r
	}

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
	return status, nil
}

func (server *SantaclausServerImpl) UpdateFileSuccess(ctx context.Context, req *pb.UpdateFileSuccessRequest) (status *pb.UpdateFileSuccessStatus, r error) {
	log.Println("Request: UpdateFileSuccess") // todo replace with class request logger

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
	log.Println("Request: ChangeFileDisk") // todo replace with class request logger

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
	userId, r := primitive.ObjectIDFromHex(req.Directory.UserId)
	if r != nil {
		return nil, r
	}

	dirId, r := primitive.ObjectIDFromHex(req.Directory.DirId)
	if r != nil {
		return nil, r
	}

	dir, r := server.createDir(ctx, userId, dirId, req.Directory.Name)
	if r != nil {
		return nil, r
	}
	status = &pb.AddDirectoryStatus{DirId: dir.Id.Hex()}
	return status, nil
}

/*
 * If dirId is nil, return root directory
 */
func (server *SantaclausServerImpl) GetDirectory(ctx context.Context, req *pb.GetDirectoryRequest) (status *pb.GetDirectoryStatus, r error) {
	log.Println("Request: GetDirectory") // todo replace with class request logger

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
