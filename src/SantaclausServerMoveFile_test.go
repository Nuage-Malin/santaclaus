package main

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"

	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMoveFileLocation(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var dir [2]pb.FileApproxMetadata
	dir[0] = pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: userId}
	var createDirrequest = pb.AddDirectoryRequest{Directory: &dir[0]}

	createDirStatus, err := server.AddDirectory(ctx, &createDirrequest)
	if err != nil {
		t.Fatal(err)
	}
	if createDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is nil", dir[0].Name) // log and fail
	}
	var file pb.FileApproxMetadata
	file = pb.FileApproxMetadata{
		Name:   getUniqueName(),
		DirId:  createDirStatus.DirId,
		UserId: userId}
	createFileRequest := pb.AddFileRequest{File: &file}
	createFileStatus, err := server.AddFile(ctx, &createFileRequest)
	if err != nil {
		t.Fatal(err)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create file '%s' : diskId or FileId is nil", file.Name)

	}

	dir[1] = pb.FileApproxMetadata{
		DirId:  createDirStatus.DirId,
		Name:   getUniqueName(),
		UserId: userId}
	createDirrequest = pb.AddDirectoryRequest{Directory: &dir[1]}

	secondCreateDirStatus, err := server.AddDirectory(ctx, &createDirrequest)
	if err != nil {
		t.Fatal(err)
	}
	if secondCreateDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is empty", dir[1].Name) // log and fail
	}

	request := pb.MoveFileRequest{
		FileId:      createFileStatus.FileId,
		NewFileName: &file.Name,
		DirId:       &secondCreateDirStatus.DirId}
	_, err = server.MoveFile(ctx, &request)

	if err != nil {
		t.Fatal(err)
	}
	getFileRequest := pb.GetFileRequest{FileId: createFileStatus.FileId}
	fileMoved, err := server.GetFile(ctx, &getFileRequest)
	if err != nil {
		t.Fatal(err)
	}
	if fileMoved.File.ApproxMetadata.DirId != dir[1].DirId {
		t.Fatalf("File path has not been moved properly, it is in %s, but should be in %s", fileMoved.File.ApproxMetadata.DirId, dir[1].DirId)
	}
	if fileMoved.File.ApproxMetadata.Name != createFileRequest.File.Name {
		t.Fatalf("File name has not been changed properly, it is %s, but should be %s", fileMoved.File.ApproxMetadata.Name, createFileRequest.File.Name)
	}
}

func TestMoveFileName(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var dir pb.FileApproxMetadata
	dir = pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: userId}
	var createDirrequest = pb.AddDirectoryRequest{Directory: &dir}

	createDirStatus, err := server.AddDirectory(ctx, &createDirrequest)
	if err != nil {
		t.Fatal(err)
	}
	if createDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is nil", dir.Name) // log and fail
	}
	var file pb.FileApproxMetadata
	file = pb.FileApproxMetadata{
		Name:   getUniqueName(),
		DirId:  dir.DirId,
		UserId: userId}
	createFileRequest := pb.AddFileRequest{File: &file}
	createFileStatus, err := server.AddFile(ctx, &createFileRequest)
	if err != nil {
		t.Fatal(err)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create file '%s' : diskId or FileId is nil", file.Name)

	}

	newFileName := getUniqueName()
	request := pb.MoveFileRequest{
		FileId:      createFileStatus.FileId,
		NewFileName: &newFileName,
		DirId:       &createDirStatus.DirId}
	_, err = server.MoveFile(ctx, &request)

	if err != nil {
		t.Fatal(err)
	}
	getFileRequest := pb.GetFileRequest{FileId: createFileStatus.FileId}
	fileMoved, err := server.GetFile(ctx, &getFileRequest)
	if err != nil {
		t.Fatal(err)
	}
	if fileMoved.File.ApproxMetadata.DirId != *request.DirId {
		t.Fatalf("File path has not been moved properly, it is in %s, but should be in %s", fileMoved.File.ApproxMetadata.DirId, dir.DirId)
	}
	if fileMoved.File.ApproxMetadata.Name != *request.NewFileName {
		t.Fatalf("File name has not been changed properly, it is %s, but should be %s", fileMoved.File.ApproxMetadata.Name, createFileRequest.File.Name)
	}
}

func TestMoveFileToFakeDirectory(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var file pb.FileApproxMetadata
	file = pb.FileApproxMetadata{
		Name:   getUniqueName(),
		DirId:  primitive.NilObjectID.Hex(),
		UserId: userId}
	createFileRequest := pb.AddFileRequest{File: &file}
	createFileStatus, err := server.AddFile(ctx, &createFileRequest)
	if err != nil {
		t.Fatal(err)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create file '%s' : diskId or FileId is nil", file.Name)

	}

	newDirId := primitive.NewObjectID().Hex()
	request := pb.MoveFileRequest{
		FileId:      createFileStatus.FileId,
		NewFileName: &file.Name,
		DirId:       &newDirId} // dir Id isn't in database as we create it now
	_, err = server.MoveFile(ctx, &request)

	if err == nil {
		t.Fatalf("File moved to non existring directory, but error have not been returned")
	}

}

// todo change both name and location
