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

func TestUpdateFileSuccess(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	file := MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	var fileSize uint64

	addFileRequest := MaeSanta.AddFileRequest{
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
	request := MaeSanta.UpdateFileSuccessRequest{
		FileId:      addFileStatus.FileId,
		NewFileSize: fileSize}
	_, err = server.UpdateFileSuccess(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
