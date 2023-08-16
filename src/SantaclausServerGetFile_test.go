package main

// todo put this file in different directory

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"

	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetFile(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	file := pb.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	var fileSize uint64

	createFileRequest := pb.AddFileRequest{
		File:     &file,
		FileSize: fileSize}
	createFileStatus, err := server.AddFile(ctx, &createFileRequest)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if createFileStatus.DiskId == "" || createFileStatus.FileId == "" {
		t.Fatalf("DiskId or FileId is empty")
	}

	request := pb.GetFileRequest{FileId: createFileStatus.FileId}
	status, err := server.GetFile(ctx, &request)

	if err != nil {
		t.Fatal(err)
	}
	if status.File.ApproxMetadata.DirPath != file.DirPath || status.File.ApproxMetadata.Name != file.Name || status.File.ApproxMetadata.UserId != file.UserId { // check approx metadata
		t.Fatalf("Metadata about file retrieved is different from added one")
	}
	if status.File.FileId != createFileStatus.FileId || status.DiskId != createFileStatus.DiskId {
		t.Fatalf("Ids inserted and retrieved don't match :\nfileId inserted : %s\tfileId retrieved : %s\ndiskId inserted : %s\tdiskId retrieved %s\n", createFileStatus.FileId, status.File.FileId, createFileStatus.DiskId, status.DiskId)
	}
	userIdPrimitive, err := primitive.ObjectIDFromHex(userId)
	dir, err := server.findDirFromPath(ctx, file.DirPath, userIdPrimitive)
	if err != nil {
		t.Fatal(err)
	}
	if status.File.DirId != dir.Id.Hex() {
		t.Fatalf("File retrieved is in different directory than the one inserted")
	}
	// todo check content of what I got
}

/*

func TestGetFiles(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(),  4*time.Second)

	dir := pb.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: "/",
		UserId:  userId}
	var createDirrequest = pb.AddDirectoryRequest{Directory: &dir}

	createDirStatus, err := server.AddDirectory(ctx, &createDirrequest)
	if err != nil {
		t.Fatal(err)
	}
	if createDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("DirId is empty") // log and fail
	}
	file := pb.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	var fileSize uint64

	createFileRequest := pb.AddFileRequest{
		File:     &file,
		FileSize: fileSize}
	createFileStatus, err := server.AddFile(ctx, &createFileRequest)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if createFileStatus.DiskId == "" || createFileStatus.FileId == "" {
		t.Fatalf("DiskId or FileId is empty")
	}

	request := pbrectoryRequest{DirId: &createDirStatus.DirId, UserId: userId, IsRecursive: true}
	status, err := server.GetDirectory(ctx, &request)

	for index, indexedFile := range status.SubFiles.Index {
		if index >= 1 {
			t.Fatalf("Inserted only one file but several retrieved in index")
		}
		if file.DirPath != indexedFile.ApproxMetadata.DirPath || file.Name != indexedFile.ApproxMetadata.Name || file.UserId != indexedFile.ApproxMetadata.UserId {
			t.Fatalf("File in index different from added one")
		}
	}
	// todo check content of what I got
}*/

//TODO test non existing file
