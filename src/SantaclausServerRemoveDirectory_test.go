package main

// todo put this file in different directory

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"

	"context"
	"path/filepath"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* AddFile */

func TestRemoveDirectory(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var err error
	var file pb.FileApproxMetadata
	var fileSize uint64 // zero value
	const nbFilesInDir = 10
	var addFileStatuses [nbFilesInDir]*pb.AddFileStatus

	addDirReq := pb.AddDirectoryRequest{
		Directory: &pb.FileApproxMetadata{
			Name:    "directoryToBeRemoved",
			DirPath: "/",
			UserId:  userId,
		}}
	addDirStatus, err := server.AddDirectory(ctx, &addDirReq)
	if err != nil {
		t.Fatal(err)
	}
	if addDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not add dir, status contains nil dirId")
	}
	for i := 0; i < nbFilesInDir; i++ { // todo an other test with recursive directory creation
		file = pb.FileApproxMetadata{
			DirPath: filepath.Join(addDirReq.Directory.DirPath, addDirReq.Directory.Name),
			Name:    getUniqueName(),
			UserId:  userId}

		addFileRequest := pb.AddFileRequest{
			File:     &file,
			FileSize: fileSize}
		addFileStatuses[i], err = server.AddFile(ctx, &addFileRequest)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if addFileStatuses[i].DiskId == "" || addFileStatuses[i].FileId == "" {
			t.Fatalf("DiskId or FileId is empty, file name : %s", file.Name)
		}
	}
	request := pb.RemoveDirectoryRequest{DirId: addDirStatus.DirId}
	_, err = server.RemoveDirectory(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// todo do getFile procedure
	// todo maybe use the server to query into the database and check if the directory has been removed
	getDirReq := pb.GetDirectoryRequest{DirId: &addDirStatus.DirId, UserId: userId}
	getDirStatus, err := server.GetDirectory(ctx, &getDirReq)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			t.Fatal(err)
		}
	}
	for _, dir := range getDirStatus.SubFiles.DirIndex {
		if dir.DirId == addDirStatus.DirId {
			t.Fatalf("Got directory supposently removed")
		}
	}
}

/*
func TestPhysicalRemoveDirectory(t *testing.T) {
	var addFileStatus *pb.AddFileStatus
	var err error
	var request pb.RemoveDirectoryRequest
	var file pb.FileApproxMetadata
	var fileSize uint64 // zero value

	for i := 0; i <= 10; i++ {
		file = pb.FileApproxMetadata{
			DirPath: "/",
			Name:    getUniqueName(),
			UserId:  userId}

		addFileRequest := pb.AddFileRequest{
			File:     &file,
			FileSize: fileSize}
		addFileStatus, err = server.AddFile(ctx, &addFileRequest)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if addFileStatus.DiskId == "" || addFileStatus.FileId == "" {
			t.Fatalf("DiskId or FileId is empty, file name : %s", file.Name)
		}
		request.FileIds = append(request.FileIds, addFileStatus.FileId)
	}

	// request := pb.RemoveDirectoryRequest{FileIds: addFileStatuses}
	_, err = server.VirtualRemoveDirectory(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = server.PhysicalRemoveDirectory(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// todo do getFile procedure
	// todo use the server to query into the database and check if the file has the 'removed' field set
	// 		cause this test is useless as it is same as previous with both VirtualRemoveDirectory
	getDirRequest := pbrectoryRequest{
		DirId: &primitive.NilObjectID.Hex(), UserId: userId, IsRecursive: true}
	getDirStatus, err := server.GetDirectory(ctx, &getDirRequest)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range getDirStatus.SubFiles.FileIndex {
		for _, fileId := range request.FileIds {
			if file.FileId == fileId {
				t.Fatalf("Virtually deleted file is present in index, id : %s", fileId)
			}
		}
	}
}
*/
