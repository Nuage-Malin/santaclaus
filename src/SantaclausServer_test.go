package main

// todo put this file in different directory

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	context "context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var server *SantaclausServerImpl = NewSantaclausServerImpl()
var ctx context.Context

var userId string = primitive.NewObjectID().Hex()

/* AddFile */
func TestAddFile(t *testing.T) {

	var file MaeSanta.FileApproxMetadata
	file.DirPath = "/"
	file.Name = "my_file"
	file.UserId = userId
	var fileSize uint64

	var request MaeSanta.AddFileRequest
	request.File = &file
	request.FileSize = fileSize
	status, err := server.AddFile(ctx, &request)
	if err != nil {
		t.Errorf(err.Error())
	}
	if status.DiskId == "" || status.FileId == "" {
		t.Errorf("DiskId or FileId is empty")
	}
}

// todo AddFile in directory
func TestAddFileSameUser(t *testing.T) {

	var file MaeSanta.FileApproxMetadata
	file.DirPath = "/"
	file.Name = "my_otherFile"
	file.UserId = userId
	var fileSize uint64

	var request MaeSanta.AddFileRequest
	request.File = &file
	request.FileSize = fileSize
	status, err := server.AddFile(ctx, &request)
	if err != nil {
		t.Errorf(err.Error())
	}
	if status.DiskId == "" || status.FileId == "" {
		t.Errorf("DiskId or FileId is empty") // log and fail
	}
}

/* AddDirectory */

func TestCreateDirectory(t *testing.T) {
	var dir MaeSanta.FileApproxMetadata
	dir.DirPath = "/"
	dir.Name = "my_directory"
	dir.UserId = userId
	var request = MaeSanta.AddDirectoryRequest{Directory: &dir}

	status, err := server.AddDirectory(ctx, &request)
	if err != nil {
		t.Error(err)
	}
	if status.DirId == primitive.NilObjectID.Hex() {
		t.Errorf("DirId is empty") // log and fail
	}
}

func TestCreateSubDirectoryTwice(t *testing.T) {
	var dir MaeSanta.FileApproxMetadata
	dir.DirPath = "/my_directory"
	dir.Name = "my_subDirectory"
	dir.UserId = userId
	var request = MaeSanta.AddDirectoryRequest{Directory: &dir}
	var dirId string

	status, err := server.AddDirectory(ctx, &request)
	if err != nil {
		t.Error(err)
	}
	dirId = status.DirId
	if dirId == primitive.NilObjectID.Hex() {
		t.Errorf("DirId is empty") // log and fail
	}
	status, err = server.AddDirectory(ctx, &request)
	if err != nil {
		t.Error(err) // log and fail
	}
	if dirId != status.DirId {
		t.Errorf("DirId is different from previously created same directory") // log and fail
	}
}

func TestCreateSubDirectorySameSubName(t *testing.T) {
	var dir MaeSanta.FileApproxMetadata
	dir.DirPath = "/my_directory/my_subDirectory"
	dir.Name = "my_directory"
	dir.UserId = userId
	var request = MaeSanta.AddDirectoryRequest{Directory: &dir}

	status, err := server.AddDirectory(ctx, &request)
	if err != nil {
		t.Error(err)
	}
	if status.DirId == primitive.NilObjectID.Hex() {
		t.Errorf("DirId is empty") // log and fail
	}
}
