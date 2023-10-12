package main

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"

	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// func TestMoveDirectoryLocation(t *testing.T) {
// 	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

// 	var dir [2]pb.FileApproxMetadata
// 	dir[0] = pb.FileApproxMetadata{
// 		DirId:  primitive.NilObjectID.Hex(),
// 		Name:   getUniqueName(),
// 		UserId: userId}

// 	createDirReq := pb.AddDirectoryRequest{Directory: &dir[0]}

// 	var createDirStatus [2]*pb.AddDirectoryStatus
// 	var err error
// 	createDirStatus[0], err = server.AddDirectory(ctx, &createDirReq)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if createDirStatus[0].DirId == primitive.NilObjectID.Hex() {
// 		t.Fatalf("Could not create directory %s : DirId is nil", dir[0].Name) // log and fail
// 	}

// 	dir[1] = pb.FileApproxMetadata{
// 		DirId:  primitive.NilObjectID.Hex(),
// 		Name:   getUniqueName(),
// 		UserId: userId}
// 	createDirReq = pb.AddDirectoryRequest{Directory: &dir[1]}

// 	createDirStatus[1], err = server.AddDirectory(ctx, &createDirReq)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if createDirStatus[1].DirId == primitive.NilObjectID.Hex() {
// 		t.Fatalf("Could not create directory %s : DirId is empty", dir[1].Name) // log and fail
// 	}

// 	request := pb.MoveDirectoryRequest{
// 		DirId:            createDirStatus[0].DirId,
// 		Name:             &dir[0].Name,
// 		NewLocationDirId: &createDirStatus[1].DirId}
// 	_, err = server.MoveDirectory(ctx, &request)

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	getDirReq := pb.GetDirectoryRequest{DirId: request.NewLocationDirId, UserId: userId, IsRecursive: false}
// 	index, err := server.GetDirectory(ctx, &getDirReq)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	for _, dir := range index.SubFiles.DirIndex {
// 		if dir.DirId == request.DirId {
// 			if dir.ApproxMetadata.DirId != createDirStatus[0].DirId {
// 				t.Fatalf("Directory path has not been moved properly, it is in %s, but should be in %s", dir.ApproxMetadata.DirId, createDirStatus[0].DirId)
// 			}
// 			// todo add parent ID check
// 		}
// 	}
// }

func TestMoveDirectoryName(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var dir pb.FileApproxMetadata
	dir = pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: userId}
	createDirReq := pb.AddDirectoryRequest{Directory: &dir}

	createDirStatus, err := server.AddDirectory(ctx, &createDirReq)
	if err != nil {
		t.Fatal(err)
	}
	if createDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("Could not create directory %s : DirId is nil", dir.Name) // log and fail
	}

	moveDirName := getUniqueName()
	request := pb.MoveDirectoryRequest{
		DirId:            createDirStatus.DirId,
		Name:             &moveDirName,
		NewLocationDirId: &createDirStatus.DirId}
	_, err = server.MoveDirectory(ctx, &request)

	if err != nil {
		t.Fatal(err)
	}
	getDirReq := pb.GetDirectoryRequest{DirId: &request.DirId, UserId: userId, IsRecursive: false}
	index, err := server.GetDirectory(ctx, &getDirReq)
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range index.SubFiles.DirIndex {
		if dir.DirId == request.DirId {
			if dir.ApproxMetadata.Name != *request.Name {
				t.Fatalf("File name has not been changed properly, it is %s, but should be %s", dir.ApproxMetadata.Name, *request.Name)
			}
		}
	}
}

// func TestMoveDirectoryToFakeDirectory(t *testing.T) {
// 	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

// 	var dir pb.FileApproxMetadata
// 	dir = pb.FileApproxMetadata{
// 		DirId:  primitive.NilObjectID.Hex(),
// 		Name:   getUniqueName(),
// 		UserId: userId}
// 	createDirReq := pb.AddDirectoryRequest{Directory: &dir}

// 	var createDirStatus *pb.AddDirectoryStatus
// 	var err error
// 	createDirStatus, err = server.AddDirectory(ctx, &createDirReq)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if createDirStatus.DirId == primitive.NilObjectID.Hex() {
// 		t.Fatalf("Could not create directory %s : DirId is nil", dir.Name) // log and fail
// 	}

// 	newLocationDirId := primitive.NewObjectID().Hex()
// 	request := pb.MoveDirectoryRequest{
// 		DirId:            createDirStatus.DirId,
// 		Name:             &dir.Name,
// 		NewLocationDirId: &newLocationDirId}
// 	_, err = server.MoveDirectory(ctx, &request)

// 	if err == nil {
// 		t.Fatalf("Moved directory to fake directory without returning an error")
// 	}
// }

// // todo change name to name that already exists
// // todo change both name and location

// func TestMoveDirSameName(t *testing.T) {
// 	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

// 	var dir [2]pb.FileApproxMetadata
// 	var createDirReq [2]pb.AddDirectoryRequest
// 	var createDirStatus [2]*pb.AddDirectoryStatus
// 	var err error
// 	dir[0] = pb.FileApproxMetadata{
// 		DirId:  primitive.NilObjectID.Hex(),
// 		Name:   getUniqueName(),
// 		UserId: userId}
// 	createDirReq[0] = pb.AddDirectoryRequest{Directory: &dir[0]}
// 	createDirStatus[0], err = server.AddDirectory(ctx, &createDirReq[0])

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if createDirStatus[0].DirId == primitive.NilObjectID.Hex() {
// 		t.Fatalf("Could not create directory %s : DirId is nil", dir[0].Name) // log and fail
// 	}
// 	dir[1] = pb.FileApproxMetadata{
// 		DirId:  primitive.NilObjectID.Hex(),
// 		Name:   getUniqueName(),
// 		UserId: userId}
// 	createDirReq[1] = pb.AddDirectoryRequest{Directory: &dir[1]}

// 	createDirStatus[1], err = server.AddDirectory(ctx, &createDirReq[1])
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if createDirStatus[1].DirId == primitive.NilObjectID.Hex() {
// 		t.Fatalf("Could not create directory %s : DirId is empty", dir[1].Name) // log and fail
// 	}

// 	request := pb.MoveDirectoryRequest{
// 		DirId:            createDirStatus[0].DirId,
// 		Name:             &createDirReq[1].Directory.Name,
// 		NewLocationDirId: &createDirStatus[0].DirId}
// 	_, err = server.MoveDirectory(ctx, &request)

// 	if err == nil {
// 		t.Fatalf("Should not be able to move directory to the same name as ")
// 	}

// }
