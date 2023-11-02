package main

// todo put this file in different directory

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"

	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRemoveUser(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var userId string = primitive.NewObjectID().Hex()

	// Create user through Directory creation
	addDirRequest := pb.AddDirectoryRequest{Directory: &pb.FileApproxMetadata{Name: getUniqueName(), DirId: primitive.NilObjectID.Hex(), UserId: userId}}
	addDirStatus, r := TestServer.AddDirectory(ctx, &addDirRequest)

	if r != nil {
		t.Fatal(r)
	}
	if addDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Dir creation returned nil DirId")
	}

	addFileRequest := pb.AddFileRequest{File: &pb.FileApproxMetadata{Name: getUniqueName(), DirId: addDirStatus.DirId, UserId: userId}, FileSize: 100}
	addFileStatus, r := TestServer.AddFile(ctx, &addFileRequest)

	if r != nil {
		t.Fatal(r)
	}
	if addFileStatus.FileId == primitive.NilObjectID.Hex() || addFileStatus.DiskId == primitive.NilObjectID.Hex() {
		t.Fatalf("File creation returned nil FileId or DiskId")
	}

	// Remove Dir and File through removeUserRequest
	rmUserRequest := pb.RemoveUserRequest{UserId: userId}
	rmUserStatus, r := TestServer.RemoveUser(ctx, &rmUserRequest)

	if r != nil {
		t.Fatal(r)
	}
	if len(rmUserStatus.FileIdsToRemove) == 0 {
		t.Fatalf("Remove User Status should contain fileId of file to be physically removed")
	}
	if rmUserStatus.FileIdsToRemove[0] != addFileStatus.FileId {
		t.Fatalf("Remove User Status should contain fileId that user created, to then be able to remove it")
	}
	// Check that directory is removed // at least virtually
	dirFound, r := TestServer.GetDirFromStringId(ctx, addDirStatus.DirId)

	if r == nil {
		t.Fatalf("Got dir that should've been removed by removeUserRequest")
	}
	if dirFound != nil {
		t.Fatalf("Got dir that should've been removed by removeUserRequest")
	}
	// Check that file is removed // at least virtually
	fileFound, r := TestServer.GetFileFromStringId(ctx, addFileStatus.FileId)

	if r == nil {
		t.Fatalf("Got file that should've been removed by removeUserRequest")
	}
	if fileFound != nil {
		t.Fatalf("Got file that should've been removed by removeUserRequest")
	}
}
