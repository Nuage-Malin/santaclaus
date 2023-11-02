package main

// todo put this file in different directory

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"

	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* AddFile */

func TestAddFile(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	file := pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: TestUserId}
	var fileSize uint64 = 1

	request := pb.AddFileRequest{
		File:     &file,
		FileSize: fileSize}
	status, err := TestServer.AddFile(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if status.DiskId == primitive.NilObjectID.Hex() || status.FileId == primitive.NilObjectID.Hex() {
		t.Fatalf("DiskId or FileId is empty")
	}
}

func TestAddFileSameUser(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	file := pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: TestUserId}
	var fileSize uint64 = 1

	request := pb.AddFileRequest{
		File:     &file,
		FileSize: fileSize}
	status, err := TestServer.AddFile(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if status.DiskId == primitive.NilObjectID.Hex() || status.FileId == primitive.NilObjectID.Hex() {
		t.Fatalf("DiskId or FileId is empty") // log and fail
	}
}

// todo AddFile in directory
