package main

// todo put this file in different directory

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"path/filepath"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* AddFile */

func TestRemoveDirectory(t *testing.T) {
	var err error
	var file MaeSanta.FileApproxMetadata
	var fileSize uint64 // zero value
	const nbFilesInDir = 10
	var addFileStatuses [nbFilesInDir]*MaeSanta.AddFileStatus

	addDirReq := MaeSanta.AddDirectoryRequest{
		Directory: &MaeSanta.FileApproxMetadata{
			Name:    "directoryToBeRemoved",
			DirPath: "/",
			UserId:  userId,
		}}
	addDirStatus, err := server.AddDirectory(server.ctx, &addDirReq)
	if err != nil {
		t.Error(err)
	}
	if addDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not add dir, status contains nil dirId")
	}
	for i := 0; i < nbFilesInDir; i++ { // todo an other test with recursive directory creation
		file = MaeSanta.FileApproxMetadata{
			DirPath: filepath.Join(addDirReq.Directory.DirPath, addDirReq.Directory.Name),
			Name:    getUniqueName(),
			UserId:  userId}

		addFileRequest := MaeSanta.AddFileRequest{
			File:     &file,
			FileSize: fileSize}
		addFileStatuses[i], err = server.AddFile(ctx, &addFileRequest)
		if err != nil {
			t.Errorf(err.Error())
		}
		if addFileStatuses[i].DiskId == "" || addFileStatuses[i].FileId == "" {
			t.Errorf("DiskId or FileId is empty, file name : %s", file.Name)
		}
	}
	request := MaeSanta.RemoveDirectoryRequest{DirId: addDirStatus.DirId}
	_, err = server.RemoveDirectory(server.ctx, &request)
	if err != nil {
		t.Errorf(err.Error())
	}
	// todo do getFile procedure
	// todo maybe use the server to query into the database and check if the directory has been removed
	getDirReq := MaeSanta.GetDirectoryRequest{DirId: addDirStatus.DirId}
	getDirStatus, err := server.GetDirectory(server.ctx, &getDirReq)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			t.Error(err)
		}
	}
	for _, dir := range getDirStatus.SubFiles.DirIndex {
		if dir.DirId == addDirStatus.DirId {
			t.Errorf("Got directory supposently removed")
		}
	}
}

/*
func TestPhysicalRemoveDirectory(t *testing.T) {
	var addFileStatus *MaeSanta.AddFileStatus
	var err error
	var request MaeSanta.RemoveDirectoryRequest
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

	// request := MaeSanta.RemoveDirectoryRequest{FileIds: addFileStatuses}
	_, err = server.VirtualRemoveDirectory(server.ctx, &request)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = server.PhysicalRemoveDirectory(server.ctx, &request)
	if err != nil {
		t.Errorf(err.Error())
	}
	// todo do getFile procedure
	// todo use the server to query into the database and check if the file has the 'removed' field set
	// 		cause this test is useless as it is same as previous with both VirtualRemoveDirectory
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
*/
