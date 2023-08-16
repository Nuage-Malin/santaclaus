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

func TestGetDirectory(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	dir := pb.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: "/",
		UserId:  userId}
	var createDirrequest = pb.AddDirectoryRequest{Directory: &dir}

	addDirStatus, err := server.AddDirectory(ctx, &createDirrequest)
	if err != nil {
		t.Fatal(err)
	}
	if addDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId returned is nil", dir.Name) // log and fail
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
	if createFileStatus.DiskId == primitive.NilObjectID.Hex() || createFileStatus.FileId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create file '%s' : diskId or FileId is nil", file.Name)
	}

	request := pb.GetDirectoryRequest{DirId: &addDirStatus.DirId, UserId: userId, IsRecursive: true}
	status, err := server.GetDirectory(ctx, &request)
	if err != nil {
		t.Fatal(err)
	}
	if status == nil || status.SubFiles == nil || status.SubFiles.FileIndex == nil || status.SubFiles.DirIndex == nil {
		t.Fatalf("Could not retrive index : file or directory index is empty")
	}
	for index, indexedFile := range status.SubFiles.FileIndex {
		if index >= 1 {
			t.Fatalf("Inserted only one file but several retrieved in index")
		}
		if file.DirPath != indexedFile.ApproxMetadata.DirPath || file.Name != indexedFile.ApproxMetadata.Name {
			t.Fatalf("File in index different from added one")
		}
	}
	for index, indexedDir := range status.SubFiles.DirIndex {
		if index >= 1 {
			t.Fatalf("Inserted only one directory but several retrieved in index")
		}
		if addDirStatus.DirId != indexedDir.DirId {
			t.Fatalf("DirId of directory in index different from added one")
		}
		if dir.DirPath != indexedDir.ApproxMetadata.DirPath || dir.Name != indexedDir.ApproxMetadata.Name {
			t.Fatalf("Directory in index different from added one : \n got \"%s\", expected \"%s\"",
				filepath.Join(indexedDir.ApproxMetadata.DirPath, indexedDir.ApproxMetadata.Name),
				filepath.Join(dir.DirPath, dir.Name))
		}
	}
}

// todo recursive with multiple directories and recursive = false

func TestGetSubDirectories(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var dir [2]pb.FileApproxMetadata
	dir[0] = pb.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: "/",
		UserId:  userId}
	var createDirRequest [2]pb.AddDirectoryRequest
	createDirRequest[0] = pb.AddDirectoryRequest{Directory: &dir[0]}
	var addDirStatus [2]*pb.AddDirectoryStatus
	var err error
	addDirStatus[0], err = server.AddDirectory(ctx, &createDirRequest[0])
	if err != nil {
		t.Fatal(err)
	}
	if addDirStatus[0].DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is nil", dir[0].Name) // log and fail
	}
	var file [2]pb.FileApproxMetadata
	file[0] = pb.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: filepath.Join(dir[0].DirPath, dir[0].Name),
		UserId:  userId}
	createFileRequest := pb.AddFileRequest{File: &file[0]}
	createFileStatus, err := server.AddFile(ctx, &createFileRequest)
	if err != nil {
		t.Fatal(err)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create file '%s' : diskId or FileId is nil", file[0].Name)

	}

	dir[1] = pb.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: filepath.Join(dir[0].DirPath, dir[0].Name),
		UserId:  userId}
	createDirRequest[1] = pb.AddDirectoryRequest{Directory: &dir[1]}

	addDirStatus[1], err = server.AddDirectory(ctx, &createDirRequest[1])
	if err != nil {
		t.Fatal(err)
	}
	if addDirStatus[1].DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is empty", dir[1].Name) // log and fail
	}
	file[1] = pb.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: filepath.Join(dir[1].DirPath, dir[1].Name),
		UserId:  userId}
	createFileRequest = pb.AddFileRequest{File: &file[1]}
	createFileStatus, err = server.AddFile(ctx, &createFileRequest)
	if err != nil {
		t.Fatal(err)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create file '%s' : diskId or FileId is nil", file[1].Name)

	}

	request := pb.GetDirectoryRequest{DirId: &addDirStatus[0].DirId, UserId: userId, IsRecursive: true}
	status, err := server.GetDirectory(ctx, &request)

	if err != nil {
		t.Fatal(err)
	}
	if status == nil || status.SubFiles == nil || status.SubFiles.FileIndex == nil || status.SubFiles.DirIndex == nil {
		t.Fatalf("Empty index")
	}
	for index, indexedFile := range status.SubFiles.FileIndex {
		if file[index].DirPath != indexedFile.ApproxMetadata.DirPath || file[index].Name != indexedFile.ApproxMetadata.Name {
			t.Fatalf("File in index different from added one")
		}
	}
	for index, indexedDir := range status.SubFiles.DirIndex {
		if index >= 2 {
			t.Fatalf("Inserted only one directory but several retrieved in index")
		}
		if addDirStatus[index].DirId != indexedDir.DirId {
			t.Fatalf("DirId of directory in index different from added one")
		}
		if dir[index].DirPath != indexedDir.ApproxMetadata.DirPath || dir[index].Name != indexedDir.ApproxMetadata.Name {
			t.Fatalf("Directory in index different from added one")
		}
	}
}

func TestGetRootDirectory(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	request := pb.GetDirectoryRequest{DirId: nil, UserId: userId, IsRecursive: true}
	status, err := server.GetDirectory(ctx, &request)

	if err != nil {
		t.Fatal(err)
	}
	if status == nil || status.SubFiles == nil || status.SubFiles.FileIndex == nil || status.SubFiles.DirIndex == nil {
		t.Fatalf("Could not retrive index : file or directory index is empty")
	}
	if len(status.SubFiles.FileIndex) <= 1 {
		t.Fatalf("File index from root dir is empty")
	}
	if len(status.SubFiles.DirIndex) <= 1 {
		t.Fatalf("Directory index from root dir is empty")
	}
}
