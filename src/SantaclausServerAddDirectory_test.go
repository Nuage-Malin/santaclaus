package main

// todo put this file in different directory

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"

	"context"
	"time"

	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* AddDirectory */

var directoryIds [4]string = [4]string{primitive.NilObjectID.Hex(), primitive.NilObjectID.Hex(), primitive.NilObjectID.Hex(), primitive.NilObjectID.Hex()}

func TestCreateDirectory(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	dir := pb.FileApproxMetadata{
		DirId:  primitive.NilObjectID.Hex(),
		Name:   getUniqueName(),
		UserId: TestUserId}
	var request = pb.AddDirectoryRequest{Directory: &dir}

	status, err := TestServer.AddDirectory(ctx, &request)
	if err != nil {
		t.Fatal(err)
	}
	if status == nil {
		t.Fatalf("Status is nil")
	}
	if status.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("DirId is empty") // log and fail
	}
	directoryIds[1] = status.DirId
}

func TestCreateSubDirectoryTwice(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	dir := pb.FileApproxMetadata{
		DirId:  directoryIds[1],
		Name:   getUniqueName(),
		UserId: TestUserId}
	var request = pb.AddDirectoryRequest{Directory: &dir}
	var dirId string

	status, err := TestServer.AddDirectory(ctx, &request)
	if err != nil {
		t.Fatal(err)
	}
	if status == nil {
		t.Fatalf("Status is nil")
	}
	dirId = status.DirId
	if dirId == primitive.NilObjectID.Hex() {
		t.Fatalf("DirId is empty") // log and fail
	}
	directoryIds[2] = dirId

	status, err = TestServer.AddDirectory(ctx, &request)
	if err == nil {
		t.Fatalf("Error: directory was created twice without error") // log and fail
	}
}

func TestCreateSubDirectorySameSubName(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	dir := pb.FileApproxMetadata{
		DirId:  directoryIds[3],
		Name:   getUniqueName(),
		UserId: TestUserId}
	var request = pb.AddDirectoryRequest{Directory: &dir}

	status, err := TestServer.AddDirectory(ctx, &request)
	if err != nil {
		t.Fatal(err)
	}
	if status.DirId == primitive.NilObjectID.Hex() {
		t.Fatalf("DirId is empty") // log and fail
	}
}
