package main

// todo put this file in different directory

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"

	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* AddFile */

func TestVirtualRemoveFiles(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var addFileStatus *pb.AddFileStatus
	var err error
	var request pb.RemoveFilesRequest
	var file pb.FileApproxMetadata
	var fileSize uint64 // zero value

	for i := 0; i <= 10; i++ {
		file = pb.FileApproxMetadata{
			DirId:  primitive.NilObjectID.Hex(),
			Name:   getUniqueName(),
			UserId: userId}

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

	// request := pb.RemoveFilesRequest{FileIds: addFileStatuses}
	_, err = server.VirtualRemoveFiles(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// todo do getFile procedure
	// todo maybe use the server to query into the database and check if the file has the 'removed' field set
	getDirRequest := pb.GetDirectoryRequest{
		DirId: nil, UserId: userId, IsRecursive: true}
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

func TestPhysicalRemoveFiles(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var addFileStatus *pb.AddFileStatus
	var err error
	var request pb.RemoveFilesRequest
	var file pb.FileApproxMetadata
	var fileSize uint64 // zero value

	for i := 0; i <= 10; i++ {
		file = pb.FileApproxMetadata{
			DirId:  primitive.NilObjectID.Hex(),
			Name:   getUniqueName(),
			UserId: userId}

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

	// request := pb.RemoveFilesRequest{FileIds: addFileStatuses}
	_, err = server.VirtualRemoveFiles(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = server.PhysicalRemoveFiles(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// todo do getFile procedure
	// todo use the server to query into the database and check if the file has the 'removed' field set
	// 		cause this test is useless as it is same as previous with both VirtualRemoveFiles
	getDirRequest := pb.GetDirectoryRequest{
		DirId: nil, UserId: userId, IsRecursive: true}
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
