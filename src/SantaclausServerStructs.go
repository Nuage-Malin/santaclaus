package main

import (
	pb "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SantaclausServerImpl struct { // implements Maestro_Santaclaus_ServiceClient interface
	mongoClient *mongo.Client
	mongoDb     *mongo.Database
	mongoColls  map[string]*mongo.Collection
	pb.UnimplementedMaestro_Santaclaus_ServiceServer
}

// Names of collections, stored in SantaclausServerImpl.mongoColls
const FilesCollName = "Children"
const DirectoriesCollName = "Rooms"
const DisksCollName = "Houses"

type file struct {
	Id         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name" json:"name"`
	DirId      primitive.ObjectID `bson:"dir_id" json:"dir_id"`
	UserId     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Size       uint64             `bson:"size"`
	DiskId     string             `bson:"disk_id"`
	LastUpload time.Time          // Lorsqu'il est virtuel: undefined, lorsqu'il est sur le disque dur: date
	CreatedAt  time.Time          `bson:"created_at"`
	EditedAt   time.Time          `bson:"updated_at"`
	Deleted    bool               `bons:"deleted"`
	// todo last editor id
}

type directory struct {
	Id        primitive.ObjectID `bson:"_id"` // `bson:"_id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	UserId    primitive.ObjectID `bson:"user_id" json:"user_id"`
	ParentId  primitive.ObjectID `bson:"parent_id"` // If nil, root directory
	CreatedAt time.Time          `bson:"created_at"`
	EditedAt  time.Time          `bson:"updated_at"`
	Deleted   bool               `bons:"deleted"`
}

type disk struct {
	Id            string `bson:"id"`
	TotalSize     uint64 `bson:"total_size"`
	AvailableSize uint64 `bson:"available_size"`
}
