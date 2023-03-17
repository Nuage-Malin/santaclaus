package main

// todo put this file in different directory

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"path/filepath"

	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* AddDirectory */

var directoryNames [4]string = [4]string{"/", getUniqueName(), getUniqueName(), getUniqueName()}

func TestCreateDirectory(t *testing.T) {
	dir := MaeSanta.FileApproxMetadata{
		DirPath: directoryNames[0],
		Name:    directoryNames[1],
		UserId:  userId}
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
	dir := MaeSanta.FileApproxMetadata{
		DirPath: filepath.Join(directoryNames[0], directoryNames[1]),
		Name:    directoryNames[2],
		UserId:  userId}
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
	dir := MaeSanta.FileApproxMetadata{
		DirPath: filepath.Join(directoryNames[0], directoryNames[1], directoryNames[2]),
		Name:    directoryNames[3],
		UserId:  userId}
	var request = MaeSanta.AddDirectoryRequest{Directory: &dir}

	status, err := server.AddDirectory(ctx, &request)
	if err != nil {
		t.Error(err)
	}
	if status.DirId == primitive.NilObjectID.Hex() {
		t.Errorf("DirId is empty") // log and fail
	}
}
