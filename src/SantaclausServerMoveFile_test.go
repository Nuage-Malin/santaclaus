package main

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"path/filepath"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMoveFileLocation(t *testing.T) {
	var dir [2]MaeSanta.FileApproxMetadata
	dir[0] = MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	var createDirrequest = MaeSanta.AddDirectoryRequest{Directory: &dir[0]}

	createDirStatus, err := server.AddDirectory(ctx, &createDirrequest)
	if err != nil {
		t.Error(err)
	}
	if createDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create directory %s : DirId is nil", dir[0].Name) // log and fail
	}
	var file MaeSanta.FileApproxMetadata
	file = MaeSanta.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: filepath.Join(dir[0].DirPath, dir[0].Name),
		UserId:  userId}
	createFileRequest := MaeSanta.AddFileRequest{File: &file}
	createFileStatus, err := server.AddFile(server.ctx, &createFileRequest)
	if err != nil {
		t.Error(err)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create file '%s' : diskId or FileId is nil", file.Name)

	}

	dir[1] = MaeSanta.FileApproxMetadata{
		DirPath: filepath.Join(dir[0].DirPath, dir[0].Name),
		Name:    getUniqueName(),
		UserId:  userId}
	createDirrequest = MaeSanta.AddDirectoryRequest{Directory: &dir[1]}

	secondCreateDirStatus, err := server.AddDirectory(ctx, &createDirrequest)
	if err != nil {
		t.Error(err)
	}
	if secondCreateDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create directory %s : DirId is empty", dir[1].Name) // log and fail
	}

	request := MaeSanta.MoveFileRequest{
		FileId: createFileStatus.FileId,
		Name:   file.Name,
		DirId:  secondCreateDirStatus.DirId}
	_, err = server.MoveFile(server.ctx, &request)

	if err != nil {
		t.Error(err)
	}
	getFileRequest := MaeSanta.GetFileRequest{FileId: createFileStatus.FileId}
	fileMoved, err := server.GetFile(server.ctx, &getFileRequest)
	if err != nil {
		t.Error(err)
	}
	if fileMoved.File.DirPath != filepath.Join(dir[1].DirPath, dir[1].Name) {
		t.Errorf("File path has not been moved properly, it is %s, but should be %s", fileMoved.File.DirPath, filepath.Join(dir[1].DirPath, dir[1].Name))
	}
	if fileMoved.File.Name != createFileRequest.File.Name {
		t.Errorf("File name has not been changed properly, it is %s, but should be %s", fileMoved.File.Name, createFileRequest.File.Name)
	}
}

func TestMoveFileName(t *testing.T) {
	var dir MaeSanta.FileApproxMetadata
	dir = MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	var createDirrequest = MaeSanta.AddDirectoryRequest{Directory: &dir}

	createDirStatus, err := server.AddDirectory(ctx, &createDirrequest)
	if err != nil {
		t.Error(err)
	}
	if createDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create directory %s : DirId is nil", dir.Name) // log and fail
	}
	var file MaeSanta.FileApproxMetadata
	file = MaeSanta.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: filepath.Join(dir.DirPath, dir.Name),
		UserId:  userId}
	createFileRequest := MaeSanta.AddFileRequest{File: &file}
	createFileStatus, err := server.AddFile(server.ctx, &createFileRequest)
	if err != nil {
		t.Error(err)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create file '%s' : diskId or FileId is nil", file.Name)

	}

	request := MaeSanta.MoveFileRequest{
		FileId: createFileStatus.FileId,
		Name:   getUniqueName(),
		DirId:  createDirStatus.DirId}
	_, err = server.MoveFile(server.ctx, &request)

	if err != nil {
		t.Error(err)
	}
	getFileRequest := MaeSanta.GetFileRequest{FileId: createFileStatus.FileId}
	fileMoved, err := server.GetFile(server.ctx, &getFileRequest)
	if err != nil {
		t.Error(err)
	}
	if fileMoved.File.DirPath != filepath.Join(dir.DirPath, dir.Name) {
		t.Errorf("File path has not been moved properly, it is %s, but should be %s", fileMoved.File.DirPath, filepath.Join(dir.DirPath, dir.Name))
	}
	if fileMoved.File.Name != request.Name {
		t.Errorf("File name has not been changed properly, it is %s, but should be %s", fileMoved.File.Name, createFileRequest.File.Name)
	}
}

func TestMoveFileToFakeDirectory(t *testing.T) {
	var file MaeSanta.FileApproxMetadata
	file = MaeSanta.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: "/",
		UserId:  userId}
	createFileRequest := MaeSanta.AddFileRequest{File: &file}
	createFileStatus, err := server.AddFile(server.ctx, &createFileRequest)
	if err != nil {
		t.Error(err)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create file '%s' : diskId or FileId is nil", file.Name)

	}

	request := MaeSanta.MoveFileRequest{
		FileId: createFileStatus.FileId,
		Name:   file.Name,
		DirId:  primitive.NewObjectID().Hex()} // dir Id isn't in database as we create it now
	_, err = server.MoveFile(server.ctx, &request)

	if err == nil {
		t.Errorf("File moved to non existring directory, but error have not been returned")
	}

}

// todo change both name and location
