package main

// todo put this file in different directory

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* AddFile */

func TestVirtualRemoveFiles(t *testing.T) {
	var addFileStatus *MaeSanta.AddFileStatus
	var err error
	var request MaeSanta.RemoveFilesRequest
	var file MaeSanta.FileApproxMetadata
	var fileSize uint64 // zero value

	for i := 0; i <= 10; i++ {
		file = MaeSanta.FileApproxMetadata{
			DirPath: "/",
			Name:    getUniqueName(),
			UserId:  userId}

		addFileRequest := MaeSanta.AddFileRequest{
			File:     &file,
			FileSize: fileSize}
		addFileStatus, err = server.AddFile(ctx, &addFileRequest)
		if err != nil {
			t.Errorf(err.Error())
		}
		if addFileStatus.DiskId == "" || addFileStatus.FileId == "" {
			t.Errorf("DiskId or FileId is empty, file name : %s", file.Name)
		}
		request.FileIds = append(request.FileIds, addFileStatus.FileId)
	}

	// request := MaeSanta.RemoveFilesRequest{FileIds: addFileStatuses}
	_, err = server.VirtualRemoveFiles(server.ctx, &request)
	if err != nil {
		t.Errorf(err.Error())
	}
	// todo do getFile procedure
	// todo maybe use the server to query into the database and check if the file has the 'removed' field set
	getDirRequest := MaeSanta.GetDirectoryRequest{
		DirId: primitive.NilObjectID.Hex(), IsRecursive: true}
	getDirStatus, err := server.GetDirectory(server.ctx, &getDirRequest)
	if err != nil {
		t.Error(err)
	}
	for _, file := range getDirStatus.SubFiles.FileIndex {
		for _, fileId := range request.FileIds {
			if file.FileId == fileId {
				t.Errorf("Virtually deleted file is present in index, id : %s", fileId)
			}
		}
	}
}

func TestPhysicalRemoveFiles(t *testing.T) {
	var addFileStatus *MaeSanta.AddFileStatus
	var err error
	var request MaeSanta.RemoveFilesRequest
	var file MaeSanta.FileApproxMetadata
	var fileSize uint64 // zero value

	for i := 0; i <= 10; i++ {
		file = MaeSanta.FileApproxMetadata{
			DirPath: "/",
			Name:    getUniqueName(),
			UserId:  userId}

		addFileRequest := MaeSanta.AddFileRequest{
			File:     &file,
			FileSize: fileSize}
		addFileStatus, err = server.AddFile(ctx, &addFileRequest)
		if err != nil {
			t.Errorf(err.Error())
		}
		if addFileStatus.DiskId == "" || addFileStatus.FileId == "" {
			t.Errorf("DiskId or FileId is empty, file name : %s", file.Name)
		}
		request.FileIds = append(request.FileIds, addFileStatus.FileId)
	}

	// request := MaeSanta.RemoveFilesRequest{FileIds: addFileStatuses}
	_, err = server.VirtualRemoveFiles(server.ctx, &request)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = server.PhysicalRemoveFiles(server.ctx, &request)
	if err != nil {
		t.Errorf(err.Error())
	}
	// todo do getFile procedure
	// todo use the server to query into the database and check if the file has the 'removed' field set
	// 		cause this test is useless as it is same as previous with both VirtualRemoveFiles
	getDirRequest := MaeSanta.GetDirectoryRequest{
		DirId: primitive.NilObjectID.Hex(), IsRecursive: true}
	getDirStatus, err := server.GetDirectory(server.ctx, &getDirRequest)
	if err != nil {
		t.Error(err)
	}
	for _, file := range getDirStatus.SubFiles.FileIndex {
		for _, fileId := range request.FileIds {
			if file.FileId == fileId {
				t.Errorf("Virtually deleted file is present in index, id : %s", fileId)
			}
		}
	}
}
