package main

// todo put this file in different directory

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	context "context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* AddFile */

func TestAddFile(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	file := MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	var fileSize uint64 = 1

	request := MaeSanta.AddFileRequest{
		File:     &file,
		FileSize: fileSize}
	status, err := server.AddFile(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if status.DiskId == primitive.NilObjectID.Hex() || status.FileId == primitive.NilObjectID.Hex() {
		t.Fatalf("DiskId or FileId is empty")
	}
}

func TestAddFileSameUser(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	file := MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	var fileSize uint64 = 1

	request := MaeSanta.AddFileRequest{
		File:     &file,
		FileSize: fileSize}
	status, err := server.AddFile(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if status.DiskId == primitive.NilObjectID.Hex() || status.FileId == primitive.NilObjectID.Hex() {
		t.Fatalf("DiskId or FileId is empty") // log and fail
	}
}

// todo AddFile in directory
