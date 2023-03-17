package main

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	context "context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SantaclausServerImpl struct { // implements Maestro_Santaclaus_ServiceClient interface
	mongoClient *mongo.Client
	mongoDb     *mongo.Database
	mongoColls  map[string]*mongo.Collection
	ctx         context.Context
	MaeSanta.UnimplementedMaestro_Santaclaus_ServiceServer
}

// Names of collections, stored in SantaclausServerImpl.mongoColls
const FilesCollName = "Children"
const DirectoriesCollName = "Rooms"
const DisksCollName = "Houses"

type file struct {
	Id         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name" json:"name"`
	DirId      primitive.ObjectID `bson:"dir_id" json:"dir_id"`   // todo change to camelCase !
	UserId     primitive.ObjectID `bson:"user_id" json:"user_id"` // todo change to camelCase !
	Size       uint64             `bson:"size"`
	DiskId     primitive.ObjectID `bson:"disk_id"` // todo change to camelCase !
	LastUpload time.Time          // Lorsqu'il est virtuel: undefined, lorsqu'il est sur le disque dur: date
	CreatedAt  time.Time          `bson:"created_at"` // todo change to camelCase !
	EditedAt   time.Time          `bson:"updated_at"` // todo change to camelCase !
	Deleted    bool
	// Available  bool todo
}

type directory struct {
	Id        primitive.ObjectID `bson:"_id"` // `bson:"_id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	UserId    primitive.ObjectID `bson:"user_id" json:"user_id"`
	ParentId  primitive.ObjectID `bson:"parent_id"` // If nil, root directory
	CreatedAt time.Time          `bson:"created_at"`
	EditedAt  time.Time          `bson:"updated_at"`
	Deleted   bool
}

type disk struct {
	Id            primitive.ObjectID `bson:"_id"`
	TotalSize     uint64             `bson:"total_size"`
	AvailableSize uint64             `bson:"available_size"`
}
