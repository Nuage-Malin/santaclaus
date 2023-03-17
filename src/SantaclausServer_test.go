package main

// todo put this file in different directory

import (
	context "context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var server *SantaclausServerImpl = NewSantaclausServerImpl()
var ctx context.Context

var userId string = primitive.NewObjectID().Hex()

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // eneable line printing with logs // TODO put somewhere else
}

var unique_count int = 0

func getUniqueName() string {
	unique_count += 1
	return fmt.Sprintf("name_%d", unique_count)
}
