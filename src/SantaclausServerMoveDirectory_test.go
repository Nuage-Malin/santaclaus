package main

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"

	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMoveDirectoryLocation(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var dir [2]pb.FileApproxMetadata
	dir[0] = pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: TestUserId}

	createDirReq := pb.AddDirectoryRequest{Directory: &dir[0]}

	var createDirStatus [2]*pb.AddDirectoryStatus
	var err error
	createDirStatus[0], err = TestServer.AddDirectory(ctx, &createDirReq)
	if err != nil {
		t.Fatal(err)
	}
	if createDirStatus[0].DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is nil", dir[0].Name) // log and fail
	}

	dir[1] = pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: TestUserId}
	createDirReq = pb.AddDirectoryRequest{Directory: &dir[1]}

	createDirStatus[1], err = TestServer.AddDirectory(ctx, &createDirReq)
	if err != nil {
		t.Fatal(err)
	}
	if createDirStatus[1].DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is empty", dir[1].Name) // log and fail
	}

	request := pb.MoveDirectoryRequest{
		DirId:    createDirStatus[0].DirId,
		NewDirId: createDirStatus[1].DirId}
	_, err = TestServer.MoveDirectory(ctx, &request)

	if err != nil {
		t.Fatal(err)
	}
	getDirReq := pb.GetDirectoryRequest{DirId: &request.NewDirId, UserId: TestUserId, IsRecursive: false}
	index, err := TestServer.GetDirectory(ctx, &getDirReq)
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range index.SubFiles.DirIndex {
		if dir.DirId == request.DirId {
			if dir.ApproxMetadata.DirId != createDirStatus[0].DirId { // todo is it good ids ?
				t.Fatalf("Directory path has not been moved properly, it is in %s, but should be in %s", dir.ApproxMetadata.DirId, createDirStatus[0].DirId)
			}
		}
	}
}

func TestMoveDirectoryToItself(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var dir pb.FileApproxMetadata
	dir = pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: TestUserId}
	createDirReq := pb.AddDirectoryRequest{Directory: &dir}

	createDirStatus, err := TestServer.AddDirectory(ctx, &createDirReq)
	if err != nil {
		t.Fatal(err)
	}
	if createDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is nil", dir.Name) // log and fail
	}

	request := pb.MoveDirectoryRequest{
		DirId:    createDirStatus.DirId,
		NewDirId: createDirStatus.DirId}
	_, err = TestServer.MoveDirectory(ctx, &request)

	if err == nil {
		t.Fatal("Moved directory to itself when it should not be possible")
	}
}

func TestMoveDirectoryToFakeDirectory(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var dir pb.FileApproxMetadata
	dir = pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: TestUserId}
	createDirReq := pb.AddDirectoryRequest{Directory: &dir}

	var createDirStatus *pb.AddDirectoryStatus
	var err error
	createDirStatus, err = TestServer.AddDirectory(ctx, &createDirReq)
	if err != nil {
		t.Fatal(err)
	}
	if createDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is nil", dir.Name) // log and fail
	}

	newLocationDirId := primitive.NewObjectID().Hex()
	request := pb.MoveDirectoryRequest{
		DirId:    createDirStatus.DirId,
		NewDirId: newLocationDirId}
	_, err = TestServer.MoveDirectory(ctx, &request)

	if err == nil {
		t.Fatalf("Moved directory to fake directory without returning an error")
	}
}

// todo rename

func TestRenameDir(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var dir pb.FileApproxMetadata
	dir = pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: TestUserId}
	createDirReq := pb.AddDirectoryRequest{Directory: &dir}

	var createDirStatus *pb.AddDirectoryStatus
	var err error
	createDirStatus, err = TestServer.AddDirectory(ctx, &createDirReq)
	if err != nil {
		t.Fatal(err)
	}
	if createDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is nil", dir.Name) // log and fail
	}
	request := pb.RenameDirectoryRequest{
		DirId:      createDirStatus.DirId,
		NewDirName: getUniqueName(),
	}
	_, err = TestServer.RenameDirectory(ctx, &request)

	if err != nil {
		t.Fatalf(err.Error())
	}
	dirFound, err := TestServer.GetDirFromStringId(ctx, createDirStatus.DirId)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if dirFound.Name != request.NewDirName {
		t.Fatalf("Directory has been wrongly renamed")
	}
}

func TestRenameDirWhenNameAlreadyExists(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var dir pb.FileApproxMetadata
	dir = pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: TestUserId}
	createDirReq := pb.AddDirectoryRequest{Directory: &dir}

	var createDirStatus0 *pb.AddDirectoryStatus
	var err error
	createDirStatus0, err = TestServer.AddDirectory(ctx, &createDirReq)
	if err != nil {
		t.Fatal(err)
	}
	if createDirStatus0.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is nil", dir.Name) // log and fail
	}
	dir.Name = getUniqueName()
	var createDirStatus1 *pb.AddDirectoryStatus
	createDirStatus1, err = TestServer.AddDirectory(ctx, &createDirReq)
	if err != nil {
		t.Fatal(err)
	}
	if createDirStatus1.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is nil", dir.Name) // log and fail
	}
	request := pb.RenameDirectoryRequest{
		DirId:      createDirStatus0.DirId,
		NewDirName: dir.Name,
	}
	_, err = TestServer.RenameDirectory(ctx, &request)

	if err == nil {
		t.Fatalf("Renamed directory without error when should've been one : rename with name already existing in this directory")
	}
}

func TestRenameDirWithSameName(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var dir pb.FileApproxMetadata
	dir = pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: TestUserId}
	createDirReq := pb.AddDirectoryRequest{Directory: &dir}

	var createDirStatus *pb.AddDirectoryStatus
	var err error
	createDirStatus, err = TestServer.AddDirectory(ctx, &createDirReq)
	if err != nil {
		t.Fatal(err)
	}
	if createDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is nil", dir.Name) // log and fail
	}
	request := pb.RenameDirectoryRequest{
		DirId:      createDirStatus.DirId,
		NewDirName: dir.Name,
	}
	_, err = TestServer.RenameDirectory(ctx, &request)

	if err == nil {
		t.Fatalf("Renamed directory without error when should've been one : rename with same name")
	}
}
