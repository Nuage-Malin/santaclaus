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

	createDirStatus, r := server.AddDirectory(ctx, &createDirrequest)
	if r != nil {
		t.Fatal(r)
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
	createFileStatus, r := server.AddFile(ctx, &createFileRequest)
	if r != nil {
		t.Fatal(r)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create file '%s' : diskId or FileId is nil", file.Name)

	}

	dir[1] = pb.FileApproxMetadata{
		DirId:  createDirStatus.DirId,
		Name:   getUniqueName(),
		UserId: userId}
	createDirrequest = pb.AddDirectoryRequest{Directory: &dir[1]}

	secondCreateDirStatus, r := server.AddDirectory(ctx, &createDirrequest)
	if r != nil {
		t.Fatal(r)
	}
	if secondCreateDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is empty", dir[1].Name) // log and fail
	}

	request := pb.MoveFileRequest{
		FileId:   createFileStatus.FileId,
		NewDirId: secondCreateDirStatus.DirId}
	_, r = server.MoveFile(ctx, &request)

	if r != nil {
		t.Fatal(r)
	}
	getFileRequest := pb.GetFileRequest{FileId: createFileStatus.FileId}
	fileMoved, r := server.GetFile(ctx, &getFileRequest)
	if r != nil {
		t.Fatal(r)
	}
	if fileMoved.File.ApproxMetadata.DirId != dir[1].DirId {
		t.Fatalf("File path has not been moved properly, it is in %s, but should be in %s", fileMoved.File.ApproxMetadata.DirId, dir[1].DirId)
	}
	if fileMoved.File.ApproxMetadata.Name != createFileRequest.File.Name {
		t.Fatalf("File name has not been changed properly, it is %s, but should be %s", fileMoved.File.ApproxMetadata.Name, createFileRequest.File.Name)
	}
}

func TestMoveFileToSameDirectory(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var dir [2]pb.FileApproxMetadata
	dir[0] = pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: userId}
	var createDirrequest = pb.AddDirectoryRequest{Directory: &dir[0]}

	createDirStatus, r := server.AddDirectory(ctx, &createDirrequest)
	if r != nil {
		t.Fatal(r)
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
	createFileStatus, r := server.AddFile(ctx, &createFileRequest)
	if r != nil {
		t.Fatal(r)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create file '%s' : diskId or FileId is nil", file.Name)

	}

	request := pb.MoveFileRequest{
		FileId:   createFileStatus.FileId,
		NewDirId: createDirStatus.DirId}
	_, r = server.MoveFile(ctx, &request)

	if r == nil {
		t.Fatal("Moved file to same directory when should've been an ror")
	}
}

func TestMoveFileToFakeDirectory(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	file := pb.FileApproxMetadata{
		Name:   getUniqueName(),
		DirId:  primitive.NilObjectID.Hex(),
		UserId: userId}
	createFileRequest := pb.AddFileRequest{File: &file}
	createFileStatus, r := server.AddFile(ctx, &createFileRequest)
	if r != nil {
		t.Fatal(r)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create file '%s' : diskId or FileId is nil", file.Name)
	}

	newDirId := primitive.NewObjectID().Hex()
	request := pb.MoveFileRequest{
		FileId:   createFileStatus.FileId,
		NewDirId: newDirId} // dir Id isn't in database as we create it now
	_, r = server.MoveFile(ctx, &request)

	if r == nil {
		t.Fatalf("File moved to non existring directory, but ror have not been returned")
	}
}

func TestRenameFile(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	file := pb.FileApproxMetadata{
		Name:   getUniqueName(),
		DirId:  primitive.NilObjectID.Hex(),
		UserId: userId}
	createFileRequest := pb.AddFileRequest{File: &file}
	createFileStatus, r := server.AddFile(ctx, &createFileRequest)

	if r != nil {
		t.Fatal(r)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create file '%s' : diskId or FileId is nil", file.Name)
	}
	renameFileRequest := pb.RenameFileRequest{FileId: createFileStatus.FileId, NewFileName: getUniqueName()}
	_, r = server.RenameFile(ctx, &renameFileRequest)

	if r != nil {
		t.Fatal(r)
	}
	fileFound, r := server.GetFileFromStringId(ctx, createFileStatus.FileId)

	if r != nil {
		t.Fatal(r)
	}
	if fileFound.Name != renameFileRequest.NewFileName {
		t.Fatalf("File renamed does not have the new file name")
	}
}

func TestRenameFileToSameSame(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	file := pb.FileApproxMetadata{
		Name:   getUniqueName(),
		DirId:  primitive.NilObjectID.Hex(),
		UserId: userId}
	createFileRequest := pb.AddFileRequest{File: &file}
	createFileStatus, r := server.AddFile(ctx, &createFileRequest)

	if r != nil {
		t.Fatal(r)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create file '%s' : diskId or FileId is nil", file.Name)
	}
	renameFileRequest := pb.RenameFileRequest{FileId: createFileStatus.FileId, NewFileName: createFileRequest.File.Name}
	_, r = server.RenameFile(ctx, &renameFileRequest)

	if r == nil {
		t.Fatal("File renamed with same name as current name")
	}
}

func TestRenameFileToAlreadyExistingFilename(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	file := pb.FileApproxMetadata{
		Name:   getUniqueName(),
		DirId:  primitive.NilObjectID.Hex(),
		UserId: userId}
	createFileRequest := pb.AddFileRequest{File: &file}
	createFileStatus0, r := server.AddFile(ctx, &createFileRequest)

	if r != nil {
		t.Fatal(r)
	}
	if createFileStatus0.FileId == primitive.NilObjectID.Hex() || createFileStatus0.DiskId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create file '%s' : diskId or FileId is nil", file.Name)
	}
	file.Name = getUniqueName()
	createFileStatus1, r := server.AddFile(ctx, &createFileRequest)

	if r != nil {
		t.Fatal(r)
	}
	if createFileStatus1.FileId == primitive.NilObjectID.Hex() || createFileStatus1.DiskId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create file '%s' : diskId or FileId is nil", file.Name)
	}
	renameFileRequest := pb.RenameFileRequest{FileId: createFileStatus0.FileId, NewFileName: createFileRequest.File.Name}
	_, r = server.RenameFile(ctx, &renameFileRequest)

	if r == nil {
		t.Fatal("File renamed with name of other file")
	}
}
