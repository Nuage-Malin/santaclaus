package main

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	"path/filepath"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMoveDirectoryLocation(t *testing.T) {
	var dir [2]MaeSanta.FileApproxMetadata
	dir[0] = MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}

	createDirReq := MaeSanta.AddDirectoryRequest{Directory: &dir[0]}

	var createDirStatus [2]*MaeSanta.AddDirectoryStatus
	var err error
	createDirStatus[0], err = server.AddDirectory(ctx, &createDirReq)
	if err != nil {
		t.Error(err)
	}
	if createDirStatus[0].DirId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create directory %s : DirId is nil", dir[0].Name) // log and fail
	}

	dir[1] = MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	createDirReq = MaeSanta.AddDirectoryRequest{Directory: &dir[1]}

	createDirStatus[1], err = server.AddDirectory(ctx, &createDirReq)
	if err != nil {
		t.Error(err)
	}
	if createDirStatus[1].DirId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create directory %s : DirId is empty", dir[1].Name) // log and fail
	}

	request := MaeSanta.MoveDirectoryRequest{
		DirId:            createDirStatus[0].DirId,
		Name:             dir[0].Name,
		NewLocationDirId: createDirStatus[1].DirId}
	_, err = server.MoveDirectory(server.ctx, &request)

	if err != nil {
		t.Error(err)
	}
	getDirReq := MaeSanta.GetDirectoryRequest{DirId: request.NewLocationDirId, IsRecursive: false}
	index, err := server.GetDirectory(server.ctx, &getDirReq)
	if err != nil {
		t.Error(err)
	}
	for _, dir := range index.SubFiles.DirIndex {
		if dir.DirId == request.DirId {
			if dir.ApproxMetadata.DirPath != filepath.Join(createDirReq.Directory.DirPath, createDirReq.Directory.Name) {
				t.Errorf("Directory path has not been moved properly, it is %s, but should be %s", dir.ApproxMetadata.DirPath, filepath.Join(createDirReq.Directory.DirPath, createDirReq.Directory.Name))
			}
		}
	}
}

func TestMoveDirectoryName(t *testing.T) {
	var dir MaeSanta.FileApproxMetadata
	dir = MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	createDirReq := MaeSanta.AddDirectoryRequest{Directory: &dir}

	createDirStatus, err := server.AddDirectory(ctx, &createDirReq)
	if err != nil {
		t.Error(err)
	}
	if createDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create directory %s : DirId is nil", dir.Name) // log and fail
	}

	request := MaeSanta.MoveDirectoryRequest{
		DirId:            createDirStatus.DirId,
		Name:             getUniqueName(),
		NewLocationDirId: createDirStatus.DirId}
	_, err = server.MoveDirectory(server.ctx, &request)

	if err != nil {
		t.Error(err)
	}
	getDirReq := MaeSanta.GetDirectoryRequest{DirId: request.DirId, IsRecursive: false}
	index, err := server.GetDirectory(server.ctx, &getDirReq)
	if err != nil {
		t.Error(err)
	}
	for _, dir := range index.SubFiles.DirIndex {
		if dir.DirId == request.DirId {
			if dir.ApproxMetadata.Name != request.Name {
				t.Errorf("File name has not been changed properly, it is %s, but should be %s", dir.ApproxMetadata.Name, request.Name)
			}
		}
	}
}

func TestMoveDirectoryToFakeDirectory(t *testing.T) {
	var dir MaeSanta.FileApproxMetadata
	dir = MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	createDirReq := MaeSanta.AddDirectoryRequest{Directory: &dir}

	var createDirStatus *MaeSanta.AddDirectoryStatus
	var err error
	createDirStatus, err = server.AddDirectory(ctx, &createDirReq)
	if err != nil {
		t.Error(err)
	}
	if createDirStatus.DirId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create directory %s : DirId is nil", dir.Name) // log and fail
	}

	request := MaeSanta.MoveDirectoryRequest{
		DirId:            createDirStatus.DirId,
		Name:             dir.Name,
		NewLocationDirId: primitive.NewObjectID().Hex()}
	_, err = server.MoveDirectory(server.ctx, &request)

	if err == nil {
		t.Errorf("Moved directory to fake directory without returning an error")
	}
}

// todo change name to name that already exists
// todo change both name and location

func TestMoveDirSameName(t *testing.T) {
	var dir [2]MaeSanta.FileApproxMetadata
	var createDirReq [2]MaeSanta.AddDirectoryRequest
	var createDirStatus [2]*MaeSanta.AddDirectoryStatus
	var err error
	dir[0] = MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	createDirReq[0] = MaeSanta.AddDirectoryRequest{Directory: &dir[0]}
	createDirStatus[0], err = server.AddDirectory(ctx, &createDirReq[0])

	if err != nil {
		t.Error(err)
	}
	if createDirStatus[0].DirId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create directory %s : DirId is nil", dir[0].Name) // log and fail
	}
	dir[1] = MaeSanta.FileApproxMetadata{
		DirPath: "/",
		Name:    getUniqueName(),
		UserId:  userId}
	createDirReq[1] = MaeSanta.AddDirectoryRequest{Directory: &dir[1]}

	createDirStatus[1], err = server.AddDirectory(ctx, &createDirReq[1])
	if err != nil {
		t.Error(err)
	}
	if createDirStatus[1].DirId == primitive.NilObjectID.Hex() {
		t.Errorf("Could not create directory %s : DirId is empty", dir[1].Name) // log and fail
	}

	request := MaeSanta.MoveDirectoryRequest{
		DirId:            createDirStatus[0].DirId,
		Name:             createDirReq[1].Directory.Name,
		NewLocationDirId: createDirStatus[0].DirId}
	_, err = server.MoveDirectory(server.ctx, &request)

	if err == nil {
		t.Errorf("Should not be able to move directory to the same name as ")
	}

}
