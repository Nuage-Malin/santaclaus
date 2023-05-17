package main

// todo put this file in different directory

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"testing"
)

/* AddFile */

func TestVirtualRemoveFile(t *testing.T) {
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
	if addFileStatus.DiskId == "" || addFileStatus.FileId == "" {
		t.Errorf("DiskId or FileId is empty")
	}

	request := MaeSanta.RemoveFileRequest{FileId: addFileStatus.FileId}
	_, err = server.VirtualRemoveFile(server.ctx, &request)
	if err != nil {
		t.Errorf(err.Error())
	}
	// todo do getFile procedure to check the file
	// todo maybe use the server to query into the database and check if the file has the 'removed' field set
	getDirRequest := MaeSanta.GetDirectoryRequest{
		DirId: nil, UserId: userId, IsRecursive: true}
	getDirStatus, err := server.GetDirectory(server.ctx, &getDirRequest)
	if err != nil {
		t.Error(err)
	}
	for _, file := range getDirStatus.SubFiles.FileIndex {
		if file.FileId == addFileStatus.FileId {
			t.Errorf("Virtually deleted file is present in index")
		}
	}
}

// todo AddFile in directory
func TestPhysicalRemoveFile(t *testing.T) {

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
	if addFileStatus.DiskId == "" || addFileStatus.FileId == "" {
		t.Errorf("DiskId or FileId is empty")
	}

	request := MaeSanta.RemoveFileRequest{FileId: addFileStatus.FileId}
	_, err = server.VirtualRemoveFile(server.ctx, &request)
	if err != nil {
		t.Errorf(err.Error())
	}
	// todo do getFile procedure

	_, err = server.PhysicalRemoveFile(server.ctx, &request)
	if err != nil {
		t.Errorf(err.Error())
	}
	// todo do getFile procedure
	getDirRequest := MaeSanta.GetDirectoryRequest{
		DirId: nil, UserId: userId, IsRecursive: true}
	getDirStatus, err := server.GetDirectory(server.ctx, &getDirRequest)
	if err != nil {
		t.Error(err)
	}
	for _, file := range getDirStatus.SubFiles.FileIndex {
		if file.FileId == addFileStatus.FileId {
			t.Errorf("Virtually deleted file is present in index")
		}
	}
}
