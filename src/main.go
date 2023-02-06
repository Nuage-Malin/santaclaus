package main

import (
    "context"
//     "time"
//
    "log"
//
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
//     "go.mongodb.org/mongo-driver/mongo/readpref"

    "os"

    "fmt"
)
var collection *mongo.Collection
var ctx = context.TODO()

func main() {
    mongoURL := os.Getenv("MAESTRO_MONGO_URL")
    fmt.Printf("MAESTRO_MONGO_URL = %s\n", mongoURL)

    clientOptions := options.Client().ApplyURI(mongoURL)
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    collection = client.Database("santaclaus").Collection("clients")
}