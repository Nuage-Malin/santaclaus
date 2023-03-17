package main

import (
	context "context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (server *SantaclausServerImpl) setMongoClient(mongoURI string) {
	// var cancelFunc context.CancelFunc
	server.ctx /* cancelFunc */, _ = context.WithTimeout(context.Background(), 10*time.Second)
	// server.ctx = context.TODO()
	log.Printf("mongoURI: %v\n", mongoURI)
	clientOptions := options.Client().ApplyURI(mongoURI)

	var err error
	// server.mongoClient, err = mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	server.mongoClient, err = mongo.Connect(server.ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// defer func() { // TODO
	// if err := server.mongoClient.Disconnect(context.TODO()); err != nil {
	// panic(err)
	// }
	// }()

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

func (server *SantaclausServerImpl) setMongoCollections(collNames []string) {
	server.mongoColls = make(map[string]*mongo.Collection, 0)

	for _, collName := range collNames {
		server.mongoColls[collName] = server.mongoDb.Collection(collName)
		if server.mongoColls[collName] == nil {
			log.Fatalf("Could not find collection \"%s\", in database \"%s\"", collName, server.mongoDb.Name())
		} else {
			log.Printf("%s collection initialized successfully\n", collName)
		}
	}
}

func NewSantaclausServerImpl() *SantaclausServerImpl {

	var server SantaclausServerImpl
	envVarNameMongoURI := "SANTACLAUS_MONGO_URI"
	mongoURI := os.Getenv(envVarNameMongoURI)
	if mongoURI == "" {
		log.Fatalf("Missing environment variable '%s'", envVarNameMongoURI)
	}
	log.Printf("env var %s = %s\n", envVarNameMongoURI, mongoURI)
	envVarNameMongoDB := "SANTACLAUS_MONGO_DB"
	indexDbName := os.Getenv(envVarNameMongoDB)
	if indexDbName == "" {
		log.Fatalf("Missing environment variable '%s'", envVarNameMongoDB)
	}
	log.Printf("env var %s = %s\n", envVarNameMongoDB, indexDbName)

	server.setMongoClient(mongoURI)
	server.setMongoDatabase(indexDbName)
	server.setMongoCollections([]string{FilesCollName, DirectoriesCollName, DisksCollName})
	return &server
}
