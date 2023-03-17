package main

// todo put this file in different directory

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* AddFile */

func TestUpdateFileSuccess(t *testing.T) {

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
		t.Errorf(err.Error())
	}
	if addFileStatus.DiskId == primitive.NilObjectID.Hex() || addFileStatus.FileId == primitive.NilObjectID.Hex() {
		t.Errorf("DiskId or FileId is empty")
	}
	fileSize = 10
	request := MaeSanta.UpdateFileSuccessRequest{
		FileId:      addFileStatus.FileId,
		NewFileSize: fileSize}
	_, err = server.UpdateFileSuccess(ctx, &request)
	if err != nil {
		t.Errorf(err.Error())
	}
}
