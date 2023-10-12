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

func TestUpdateFileSuccess(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	file := pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: userId}
	var fileSize uint64

	addFileRequest := pb.AddFileRequest{
		File:     &file,
		FileSize: fileSize}
	addFileStatus, err := server.AddFile(ctx, &addFileRequest)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if addFileStatus.DiskId == primitive.NilObjectID.Hex() || addFileStatus.FileId == primitive.NilObjectID.Hex() {
		t.Fatalf("DiskId or FileId is empty")
	}
	fileSize = 10
	request := pb.UpdateFileSuccessRequest{
		FileId:      addFileStatus.FileId,
		NewFileSize: fileSize}
	_, err = server.UpdateFileSuccess(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
