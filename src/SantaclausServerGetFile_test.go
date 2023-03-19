package main

// todo put this file in different directory

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"testing"
)

func TestGetFile(t *testing.T) {
	file := MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	var fileSize uint64

	createFileRequest := MaeSanta.AddFileRequest{
		File:     &file,
		FileSize: fileSize}
	createFileStatus, err := server.AddFile(ctx, &createFileRequest)
	if err != nil {
		t.Errorf(err.Error())
	}
	if createFileStatus.DiskId == "" || createFileStatus.FileId == "" {
		t.Errorf("DiskId or FileId is empty")
	}

	request := MaeSanta.GetFileRequest{FileId: createFileStatus.FileId}
	status, err := server.GetFile(ctx, &request)

	if err != nil {
		t.Error(err)
	}
	if status.File.DirPath != file.DirPath || status.File.Name != file.Name || status.File.UserId != file.UserId { // check approx metadata
		t.Errorf("Metadata about file retrieved is different from added one")
	}
	// todo check content of what I got
}

/*

func TestGetFiles(t *testing.T) {
	dir := MaeSanta.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: "/",
		UserId:  userId}
	var createDirrequest = MaeSanta.AddDirectoryRequest{Directory: &dir}

	createDirStatus, err := server.AddDirectory(ctx, &createDirrequest)
	if err != nil {
		t.Error(err)
	}
	if createDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Errorf("DirId is empty") // log and fail
	}
	file := MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	var fileSize uint64

	createFileRequest := MaeSanta.AddFileRequest{
		File:     &file,
		FileSize: fileSize}
	createFileStatus, err := server.AddFile(ctx, &createFileRequest)
	if err != nil {
		t.Errorf(err.Error())
	}
	if createFileStatus.DiskId == "" || createFileStatus.FileId == "" {
		t.Errorf("DiskId or FileId is empty")
	}

	request := MaeSanta.GetDirectoryRequest{DirId: createDirStatus.DirId, IsRecursive: true}
	status, err := server.GetDirectory(server.ctx, &request)

	for index, indexedFile := range status.SubFiles.Index {
		if index >= 1 {
			t.Errorf("Inserted only one file but several retrieved in index")
		}
		if file.DirPath != indexedFile.ApproxMetadata.DirPath || file.Name != indexedFile.ApproxMetadata.Name || file.UserId != indexedFile.ApproxMetadata.UserId {
			t.Errorf("File in index different from added one")
		}
	}
	// todo check content of what I got
}*/

//TODO test non existing file
