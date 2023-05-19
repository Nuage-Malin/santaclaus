package main

// todo put this file in different directory

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	context "context"
	"testing"
	"time"
)

/* AddFile */

func TestVirtualRemoveFile(t *testing.T) {
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
	if addFileStatus.DiskId == "" || addFileStatus.FileId == "" {
		t.Fatalf("DiskId or FileId is empty")
	}

	request := MaeSanta.RemoveFileRequest{FileId: addFileStatus.FileId}
	_, err = server.VirtualRemoveFile(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// todo do getFile procedure to check the file
	// todo maybe use the server to query into the database and check if the file has the 'removed' field set
	getDirRequest := MaeSanta.GetDirectoryRequest{
		DirId: nil, UserId: userId, IsRecursive: true}
	getDirStatus, err := server.GetDirectory(ctx, &getDirRequest)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range getDirStatus.SubFiles.FileIndex {
		if file.FileId == addFileStatus.FileId {
			t.Fatalf("Virtually deleted file is present in index")
		}
	}
}

// todo AddFile in directory
func TestPhysicalRemoveFile(t *testing.T) {
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
	if addFileStatus.DiskId == "" || addFileStatus.FileId == "" {
		t.Fatalf("DiskId or FileId is empty")
	}

	request := MaeSanta.RemoveFileRequest{FileId: addFileStatus.FileId}
	_, err = server.VirtualRemoveFile(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// todo do getFile procedure

	_, err = server.PhysicalRemoveFile(ctx, &request)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// todo do getFile procedure
	getDirRequest := MaeSanta.GetDirectoryRequest{
		DirId: nil, UserId: userId, IsRecursive: true}
	getDirStatus, err := server.GetDirectory(ctx, &getDirRequest)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range getDirStatus.SubFiles.FileIndex {
		if file.FileId == addFileStatus.FileId {
			t.Fatalf("Virtually deleted file is present in index")
		}
	}
}
