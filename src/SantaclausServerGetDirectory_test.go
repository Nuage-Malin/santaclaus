package main

// todo put this file in different directory

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"path/filepath"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetDirectory(t *testing.T) {
	dir := MaeSanta.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: "/",
		UserId:  userId}
	var createDirrequest = MaeSanta.AddDirectoryRequest{Directory: &dir}

	addDirStatus, err := server.AddDirectory(ctx, &createDirrequest)
	if err != nil {
		t.Error(err)
	}
	if addDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create directory %s : DirId returned is nil", dir.Name) // log and fail
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
	if createFileStatus.DiskId == primitive.NilObjectID.Hex() || createFileStatus.FileId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create file '%s' : diskId or FileId is nil", file.Name)
	}

	request := MaeSanta.GetDirectoryRequest{DirId: addDirStatus.DirId, IsRecursive: true}
	status, err := server.GetDirectory(server.ctx, &request)
	if err != nil {
		t.Error(err)
	}
	if status == nil || status.SubFiles == nil || status.SubFiles.FileIndex == nil || status.SubFiles.DirIndex == nil {
		t.Errorf("Could not retrive index : file or directory index is empty")
	}
	for index, indexedFile := range status.SubFiles.FileIndex {
		if index >= 1 {
			t.Errorf("Inserted only one file but several retrieved in index")
		}
		if file.DirPath != indexedFile.ApproxMetadata.DirPath || file.Name != indexedFile.ApproxMetadata.Name {
			t.Errorf("File in index different from added one")
		}
	}
	for index, indexedDir := range status.SubFiles.DirIndex {
		if index >= 1 {
			t.Errorf("Inserted only one directory but several retrieved in index")
		}
		if addDirStatus.DirId != indexedDir.DirId {
			t.Errorf("DirId of directory in index different from added one")
		}
		if dir.DirPath != indexedDir.ApproxMetadata.DirPath || dir.Name != indexedDir.ApproxMetadata.Name {
			t.Errorf("Directory in index different from added one : \n got \"%s\", expected \"%s\"",
				filepath.Join(indexedDir.ApproxMetadata.DirPath, indexedDir.ApproxMetadata.Name),
				filepath.Join(dir.DirPath, dir.Name))
		}
	}
}

// todo recursive with multiple directories and recursive = false

func TestGetSubDirectories(t *testing.T) {
	var dir [2]MaeSanta.FileApproxMetadata
	dir[0] = MaeSanta.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: "/",
		UserId:  userId}
	var createDirRequest [2]MaeSanta.AddDirectoryRequest
	createDirRequest[0] = MaeSanta.AddDirectoryRequest{Directory: &dir[0]}
	var addDirStatus [2]*MaeSanta.AddDirectoryStatus
	var err error
	addDirStatus[0], err = server.AddDirectory(ctx, &createDirRequest[0])
	if err != nil {
		t.Error(err)
	}
	if addDirStatus[0].DirId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create directory %s : DirId is nil", dir[0].Name) // log and fail
	}
	var file [2]MaeSanta.FileApproxMetadata
	file[0] = MaeSanta.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: filepath.Join(dir[0].DirPath, dir[0].Name),
		UserId:  userId}
	createFileRequest := MaeSanta.AddFileRequest{File: &file[0]}
	createFileStatus, err := server.AddFile(server.ctx, &createFileRequest)
	if err != nil {
		t.Error(err)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create file '%s' : diskId or FileId is nil", file[0].Name)

	}

	dir[1] = MaeSanta.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: filepath.Join(dir[0].DirPath, dir[0].Name),
		UserId:  userId}
	createDirRequest[1] = MaeSanta.AddDirectoryRequest{Directory: &dir[1]}

	addDirStatus[1], err = server.AddDirectory(ctx, &createDirRequest[1])
	if err != nil {
		t.Error(err)
	}
	if addDirStatus[1].DirId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create directory %s : DirId is empty", dir[1].Name) // log and fail
	}
	file[1] = MaeSanta.FileApproxMetadata{
		Name:    getUniqueName(),
		DirPath: filepath.Join(dir[1].DirPath, dir[1].Name),
		UserId:  userId}
	createFileRequest = MaeSanta.AddFileRequest{File: &file[1]}
	createFileStatus, err = server.AddFile(server.ctx, &createFileRequest)
	if err != nil {
		t.Error(err)
	}
	if createFileStatus.FileId == primitive.NilObjectID.Hex() || createFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create file '%s' : diskId or FileId is nil", file[1].Name)

	}

	request := MaeSanta.GetDirectoryRequest{DirId: addDirStatus[0].DirId, IsRecursive: true}
	status, err := server.GetDirectory(server.ctx, &request)

	if err != nil {
		t.Error(err)
	}
	if status == nil || status.SubFiles == nil || status.SubFiles.FileIndex == nil || status.SubFiles.DirIndex == nil {
		t.Errorf("Empty index")
	}
	for index, indexedFile := range status.SubFiles.FileIndex {
		if file[index].DirPath != indexedFile.ApproxMetadata.DirPath || file[index].Name != indexedFile.ApproxMetadata.Name {
			t.Errorf("File in index different from added one")
		}
	}
	for index, indexedDir := range status.SubFiles.DirIndex {
		if index >= 2 {
			t.Errorf("Inserted only one directory but several retrieved in index")
		}
		if addDirStatus[index].DirId != indexedDir.DirId {
			t.Errorf("DirId of directory in index different from added one")
		}
		if dir[index].DirPath != indexedDir.ApproxMetadata.DirPath || dir[index].Name != indexedDir.ApproxMetadata.Name {
			t.Errorf("Directory in index different from added one")
		}
	}
}

func TestGetRootDirectory(t *testing.T) {
	request := MaeSanta.GetDirectoryRequest{DirId: primitive.NilObjectID.Hex(), IsRecursive: true}
	status, err := server.GetDirectory(server.ctx, &request)
	if err != nil {
		t.Error(err)
	}
	if status == nil || status.SubFiles == nil || status.SubFiles.FileIndex == nil || status.SubFiles.DirIndex == nil {
		t.Errorf("Could not retrive index : file or directory index is empty")
	}
	if len(status.SubFiles.FileIndex) <= 1 {
		t.Errorf("File index from root dir is empty")
	}
	if len(status.SubFiles.DirIndex) <= 1 {
		t.Errorf("Directory index from root dir is empty")
	}
}
