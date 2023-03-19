package main

// todo put this file in different directory

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* AddFile */

func TestAddFile(t *testing.T) {

	file := MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	var fileSize uint64

	request := MaeSanta.AddFileRequest{
		File:     &file,
		FileSize: fileSize}
	status, err := server.AddFile(ctx, &request)
	if err != nil {
		t.Errorf(err.Error())
	}
	if status.DiskId == primitive.NilObjectID.Hex() || status.FileId == primitive.NilObjectID.Hex() {
		t.Errorf("DiskId or FileId is empty")
	}
}

func TestAddFileSameUser(t *testing.T) {

	file := MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	var fileSize uint64

	request := MaeSanta.AddFileRequest{
		File:     &file,
		FileSize: fileSize}
	status, err := server.AddFile(ctx, &request)
	if err != nil {
		t.Errorf(err.Error())
	}
	if status.DiskId == primitive.NilObjectID.Hex() || status.FileId == primitive.NilObjectID.Hex() {
		t.Errorf("DiskId or FileId is empty") // log and fail
	}
}

// todo AddFile in directory
