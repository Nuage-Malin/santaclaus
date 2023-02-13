package main

import (
	// File "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated/NuageMalin/File"
	// MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated/NuageMalin/Maestro_Santaclaus"
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	context "context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "go.mongodb.org/mongo-driver/bson/primitive"

	// MaeSanta "NuageMalin/Maestro_Santaclaus"

	// MaeSanta "third_parties/protobuf-interfaces/NuageMalin/Maestro_Santaclaus/Maestro_Santaclaus"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type SantaclausServerImpl struct { // implements Maestro_Santaclaus_ServiceClient interface
	mongoClient *mongo.Client
	mongoDb     *mongo.Database
	mongoColl   *mongo.Collection
	ctx         context.Context
	MaeSanta.UnimplementedMaestro_Santaclaus_ServiceServer
	// proto.UnimplementedGreeterServer
}

func (server *SantaclausServerImpl) setMongoClient(mongoURI string) {
	clientOptions := options.Client().ApplyURI(mongoURI)
	var err error
	server.mongoClient, err = mongo.Connect(server.ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = server.mongoClient.Ping(server.ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *SantaclausServerImpl) setMongoDatabase(dbName string) {
	server.mongoDb = server.mongoClient.Database(dbName)
	if server.mongoDb == nil {
		log.Fatalf("Could not find database \"%s\"", dbName)
	}

}

func (server *SantaclausServerImpl) setMongoCollection(collName string) {
	server.mongoColl = server.mongoDb.Collection(collName)
	if server.mongoColl == nil {
		log.Fatalf("Could not find collection \"%s\", in database \"%s\"", collName, server.mongoDb.Name())

	}
}

func NewSantaclausServerImpl() SantaclausServerImpl {
	var server SantaclausServerImpl
	envVarnameMongoURI := "SANTACLAUS_MONGO_URI"
	mongoURI := os.Getenv(envVarnameMongoURI)
	if mongoURI == "" {
		log.Fatalf("Missing environment variable '%s'", envVarnameMongoURI)
	}
	fmt.Printf("env var %s = %s\n", envVarnameMongoURI, mongoURI)

	server.setMongoClient(mongoURI)
	server.setMongoDatabase("Santaclaus")
	server.setMongoCollection("Children")
	return server
}

// Files

func (server SantaclausServerImpl) AddFile(ctx context.Context, in *MaeSanta.AddFileRequest) (*MaeSanta.AddFileStatus, error) {
	// server.mongoColl.InsertOne(server.ctx, )
	return nil, status.Errorf(codes.Unimplemented, "method AddFile not implemented")
}

func (server SantaclausServerImpl) VirtualRemoveFile(context.Context, *MaeSanta.RemoveFileRequest) (status *MaeSanta.RemoveFileStatus, r error) {
	return status, r
}
func (server SantaclausServerImpl) PhysicalRemoveFile(context.Context, *MaeSanta.RemoveFileRequest) (status *MaeSanta.RemoveFileStatus, r error) {
	return status, r
}
func (server SantaclausServerImpl) MoveFile(context.Context, *MaeSanta.MoveFileRequest) (status *MaeSanta.MoveFileStatus, r error) {
	return status, r
}
func (server SantaclausServerImpl) GetFile(context.Context, *MaeSanta.GetFileRequest) (status *MaeSanta.GetFileStatus, r error) {
	return status, r
}
func (server SantaclausServerImpl) UpdateFileSuccess(context.Context, *MaeSanta.UpdateFileSuccessRequest) (status *MaeSanta.UpdateFileSuccessStatus, r error) {
	return status, r
}
func (server SantaclausServerImpl) ChangeFileDisk(context.Context, *MaeSanta.ChangeFileDiskRequest) (status *MaeSanta.ChangeFileDiskStatus, r error) {
	return status, r
}

// Directories
func (server SantaclausServerImpl) AddDirectory(context.Context, *MaeSanta.AddDirectoryRequest) (status *MaeSanta.AddDirectoryStatus, r error) {
	return status, r
}
func (server SantaclausServerImpl) RemoveDirectory(context.Context, *MaeSanta.RemoveDirectoryRequest) (status *MaeSanta.RemoveDirectoryStatus, r error) {
	return status, r
}
func (server SantaclausServerImpl) MoveDirectory(context.Context, *MaeSanta.MoveDirectoryRequest) (status *MaeSanta.MoveDirectoryStatus, r error) {
	return status, r
}
func (server SantaclausServerImpl) GetDirectory(context.Context, *MaeSanta.GetDirectoryRequest) (status *MaeSanta.GetDirectoryStatus, r error) {
	return status, r
}

// func (server SantaclausServerImpl) mustEmbedUnimplementedMaestro_Santaclaus_ServiceServer() {
// fmt.Println("hello")
// }
