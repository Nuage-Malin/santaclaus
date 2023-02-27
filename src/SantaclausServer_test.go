package main

// todo put this file in different directory

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	context "context"
	"testing"
)

var server SantaclausServerImpl = NewSantaclausServerImpl()
var ctx context.Context

func TestAddFile(t *testing.T) {

	var file MaeSanta.FileApproxMetadata
	file.DirPath = "./"
	file.Name = "my_file"
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
