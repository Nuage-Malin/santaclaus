package main

import (
	MaeSanta "NuageMalin/Santaclaus/third_parties/protobuf-interfaces/generated"
	context "context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// const indexDbName = "Santaclaus"
const filesCollName = "Children"
const directoriesCollName = "Rooms"
const disksCollName = "Houses"

type SantaclausServerImpl struct { // implements Maestro_Santaclaus_ServiceClient interface
	mongoClient *mongo.Client
	mongoDb     *mongo.Database
	mongoColls  map[string]*mongo.Collection
	ctx         context.Context
	MaeSanta.UnimplementedMaestro_Santaclaus_ServiceServer
	// proto.UnimplementedGreeterServer
}

func (server *SantaclausServerImpl) setMongoClient(mongoURI string) {
	// var cancelFunc context.CancelFunc
	// server.ctx /* cancelFunc */, _ = context.WithTimeout(context.Background(), 10*time.Second)
	// server.ctx = context.TODO()
	fmt.Printf("mongoURI: %v\n", mongoURI)
	clientOptions := options.Client().ApplyURI(mongoURI)

	var err error
	// server.mongoClient, err = mongo.NewClient(clientOptions)
	fmt.Println("hello")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("continue")

	fmt.Printf("clientOptions.Auth: %v\n", clientOptions.Auth)
	fmt.Printf("clientOptions.AppName: %v\n", clientOptions.AppName)
	server.mongoClient, err = mongo.Connect(server.ctx, clientOptions)
	if err != nil {
		fmt.Println("error 1")
		log.Fatal(err)
	}
	// defer func() {
	// if err := server.mongoClient.Disconnect(context.TODO()); err != nil {
	// panic(err)
	// }
	// }()

	// todo error with the username and password (when connecting without them, it works)
	// err = server.mongoClient.Ping(server.ctx, nil)
	// if err != nil {
	// fmt.Println("error 2")
	// log.Fatal(err)
	// }
	fmt.Println("finish")

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
			fmt.Printf("%s collection initialized successfully\n", collName)
		}
	}
}

func NewSantaclausServerImpl() SantaclausServerImpl {
	var server SantaclausServerImpl
	envVarNameMongoURI := "SANTACLAUS_MONGO_URI"
	mongoURI := os.Getenv(envVarNameMongoURI)
	if mongoURI == "" {
		log.Fatalf("Missing environment variable '%s'", envVarNameMongoURI)
	}
	fmt.Printf("env var %s = %s\n", envVarNameMongoURI, mongoURI)
	envVarNameMongoDB := "SANTACLAUS_MONGO_DB"
	indexDbName := os.Getenv(envVarNameMongoDB)
	if indexDbName == "" {
		log.Fatalf("Missing environment variable '%s'", envVarNameMongoDB)
	}
	fmt.Printf("env var %s = %s\n", envVarNameMongoDB, indexDbName)

	server.setMongoClient(mongoURI)
	server.setMongoDatabase(indexDbName)
	server.setMongoCollections([]string{filesCollName, directoriesCollName, disksCollName})
	return server
}

type file struct {
	name       string
	dirId      primitive.ObjectID // `bson:"_id"` // can be undefined
	userId     string
	size       uint64
	diskId     string
	lastUpload time.Time // `bson:"created_at"` // Lorsqu'il est virtuel: undefined, lorsqu'il est sur le disque dur: date
	createdAt  time.Time // `bson:"created_at"`
	editedAt   time.Time // `bson:"updated_at"`
	deleted    bool
}

type directory struct {
	name      string
	userId    string
	parentId  primitive.ObjectID // can be undefined
	dirId     primitive.ObjectID // can be undefined // should remove that cause automatically created by mongo
	createdAt time.Time
	editedAt  time.Time
}

type disk struct {
	_id           string // disk Serial Number
	totalSize     uint64
	availableSize uint64
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (server SantaclausServerImpl) findDirFromPath(dirPath string, userId string) (directory, error) {

	dirName := filepath.Base(dirPath)
	dirBase := filepath.Dir(dirPath)
	var tmp directory
	var didFind bool = false

	targetDirectory := bson.D{primitive.E{Key: "name", Value: dirName}, primitive.E{Key: "userId", Value: userId}}
	cur, err := server.mongoColls[directoriesCollName].Find(server.ctx, targetDirectory /*, &targetDirectoryOptions*/)
	checkErr(err)
	for i := cur; i != nil; i.Next(server.ctx) {
		dirName = filepath.Base(dirPath)
		dirBase = filepath.Dir(dirPath)
		err = i.Decode(tmp)
		checkErr(err) // TODO this is not the solution, I guess ?
		for {
			if tmp.parentId == primitive.NilObjectID {
				fmt.Println("Found directory !")
				didFind = true
			}
			dirName = filepath.Base(dirBase)
			dirBase = filepath.Dir(dirBase)
			targetDirectory = bson.D{primitive.E{Key: "dirId", Value: tmp.parentId}, primitive.E{Key: "name", Value: dirName}, primitive.E{Key: "userId", Value: userId}}
			res := server.mongoColls[directoriesCollName].FindOne(server.ctx, targetDirectory /*, &targetDirectoryOptions*/)

			if res == nil {
				break
			}

			err = res.Decode(tmp)
			checkErr(err)

			// 	dirID should be formed from hashing name and parentId

			// targetDirectory := bson.D{{Key: "name", Value: filepath.Dir(dirPath)}, primitive.E{Key: "userId", Value: userId}}
			// targetDirectoryOptions := options.FindOneOptions{} // TODO refine search

		}
		if didFind {
			i.Decode(tmp)
			return tmp, nil
		}
	}
	return tmp, err // TODO learn about creating errors
}

func (server SantaclausServerImpl) findAvailableDisk(fileSize uint64, userId string) (found disk /* TODO change by disk id*/) {
	filter := bson.D{}
	targetDiskOptions := options.FindOneOptions{Min: bson.D{primitive.E{Key: "availableSize", Value: fileSize}}} // TODO refine search

	/*res := */
	server.mongoColls[disksCollName].FindOne(server.ctx, filter, &targetDiskOptions)

	return found
}

func (server SantaclausServerImpl) AddFile(ctx context.Context, req *MaeSanta.AddFileRequest) (status *MaeSanta.AddFileStatus, r error) {
	foundDirectory, err := server.findDirFromPath(req.File.DirPath, req.File.UserId)
	if err != nil {
		// TODO do something
	}

	// TODO find diskId

	foundDisk := server.findAvailableDisk(req.FileSize, req.File.UserId)
	newFile := file{
		name:       req.File.Name,
		dirId:      foundDirectory.dirId, // TODO find dirId from dirpath
		userId:     req.File.UserId,
		size:       req.FileSize,
		diskId:     foundDisk._id,
		lastUpload: time.Now(),
		createdAt:  time.Now(),
		editedAt:   time.Now(),
		deleted:    false,
	}
	insertRes, _ := server.mongoColls[filesCollName].InsertOne(server.ctx, newFile)
	newFileId, ok := insertRes.InsertedID.(string)

	if ok == false {
		fmt.Println("Wrong type assertion!")
		// TODO check
	}

	status = &MaeSanta.AddFileStatus{
		FileId: newFileId,
		DiskId: newFile.diskId}

	return status, nil
}

func (server SantaclausServerImpl) VirtualRemoveFile(ctx context.Context, req *MaeSanta.RemoveFileRequest) (status *MaeSanta.RemoveFileStatus, r error) {
	filter := bson.D{primitive.E{Key: "_id", Value: req.GetFileId()}}
	update := bson.D{primitive.E{Key: "deleted", Value: true}} // only modify 'deleted' to true

	server.mongoColls[filesCollName].UpdateOne(server.ctx, filter, update)
	return status, r
}

func (server SantaclausServerImpl) PhysicalRemoveFile(ctx context.Context, req *MaeSanta.RemoveFileRequest) (status *MaeSanta.RemoveFileStatus, r error) {
	filter := bson.D{primitive.E{Key: "_id", Value: req.GetFileId()}}

	// TODO find out more about contexts !
	server.mongoColls[filesCollName].DeleteOne(server.ctx, filter)
	return status, r
}

func (server SantaclausServerImpl) MoveFile(_ context.Context, req *MaeSanta.MoveFileRequest) (status *MaeSanta.MoveFileStatus, r error) {

	// filter := bson.D{primitive.E{Key: "_id", Value: req.GetFileId()}}
	// file := server.mongoColls[filesCollName].FindOne(server.ctx, filter)
	// server.findDirFromPath(file.dirBase(req.GetFilepath()), /* file. find user id from file */)

	// update := bson.D{primitive.E{Key: "dirId", Value: /* new directory id */}}

	// modify dir Id
	// server.mongoColls[filesCollName].UpdateOne(server.ctx, filter, update)

	return status, r
}

func (server SantaclausServerImpl) GetFile(_ context.Context, req *MaeSanta.GetFileRequest) (status *MaeSanta.GetFileStatus, r error) {
	filter := bson.D{primitive.E{Key: "_id", Value: req.GetFileId()}}

	/* file := */
	server.mongoColls[filesCollName].FindOne(server.ctx, filter)
	return status, r
}

func (server SantaclausServerImpl) UpdateFileSuccess(_ context.Context, req *MaeSanta.UpdateFileSuccessRequest) (status *MaeSanta.UpdateFileSuccessStatus, r error) {
	filter := bson.D{primitive.E{Key: "_id", Value: req.GetFileId()}}
	update := bson.D{primitive.E{Key: "size", Value: req.GetNewFileSize()}}

	server.mongoColls[filesCollName].UpdateOne(server.ctx, filter, update)
	return status, r
}

func (server SantaclausServerImpl) ChangeFileDisk(_ context.Context, req *MaeSanta.ChangeFileDiskRequest) (status *MaeSanta.ChangeFileDiskStatus, r error) {

	filter := bson.D{primitive.E{Key: "_id", Value: req.GetFileId()}}
	// find the file in order not to put it on the same disk as it is already
	server.mongoColls[filesCollName].FindOne(server.ctx, filter)

	// TODO algorithm to find new disk
	// find disk where
	//		there is some other file from this user
	// 	there is enough space for the file (and a bit more)

	// filter = bson.D{primitive.E{Key: "diskId", Value: /* value found from last request */}, primitive.E{Key: "userId", Value: /* value found from last request */}}
	// todo exclude from filter diskId that is the actual
	// server.mongoColls[filesCollName].Find(server.ctx, filter, update)

	// update := bson.D{primitive.E{Key: "size", Value: /* new disk id */}}
	// server.mongoColls[filesCollName].UpdateOne(server.ctx, filter, update)

	return status, r
}

// Directories
func (server SantaclausServerImpl) AddDirectory(_ context.Context, req *MaeSanta.AddDirectoryRequest) (status *MaeSanta.AddDirectoryStatus, r error) {
	// TODO algo to create new directory :
	// 		hash of parentId and name

	// find parent ID from req.Directory.DirPath
	parentDir, err := server.findDirFromPath(filepath.Dir(req.Directory.DirPath), req.Directory.UserId)
	// TODO check error other than that
	checkErr(err)
	newDirectory := directory{
		name:      filepath.Base(req.Directory.Name),
		userId:    req.Directory.UserId,
		parentId:  parentDir.dirId,
		dirId:     primitive.NewObjectID(), // todo remove cause useless cause automatically created by mongo (the id)
		createdAt: time.Now(),
		editedAt:  time.Now()}
	res, err := server.mongoColls[directoriesCollName].InsertOne(server.ctx, newDirectory)
	if err != nil {
		// log.Logger()
	}

	// res.InsertedID.Decode(status.DirId)

	var ok bool
	status.DirId, ok = res.InsertedID.(string)
	if !ok {
		// err
	}
	return status, r
}
func (server SantaclausServerImpl) RemoveDirectory(context.Context, *MaeSanta.RemoveDirectoryRequest) (status *MaeSanta.RemoveDirectoryStatus, r error) {
	// remove all files
	// server.server.mongoColls[filesCollsName].FindAndDelete(/* filter with dirID */)
	// server.server.mongoColls[directoriesCollsName].FindAndDelete(/* fileter with dirID */)
	return status, r
}
func (server SantaclausServerImpl) MoveDirectory(context.Context, *MaeSanta.MoveDirectoryRequest) (status *MaeSanta.MoveDirectoryStatus, r error) {
	// - add directory
	// - change files' directory Id
	// - remove directory
	return status, r
}

func (server SantaclausServerImpl) getOneDirectory(dirId primitive.ObjectID, recursive bool, dirPath string, status *MaeSanta.GetDirectoryStatus) *MaeSanta.GetDirectoryStatus {

	filter := bson.D{primitive.E{Key: "dirId", Value: dirId}}
	dirFound := server.mongoColls[filesCollName].FindOne(server.ctx, filter)

	var currentDir directory
	dirFound.Decode(currentDir)
	currentMetadata := MaeSanta.FileApproxMetadata{Name: currentDir.name, DirPath: dirPath /* TODO if not exists, find directory path from dir ID */, UserId: currentDir.userId}
	status.Directories = append(status.Directories, &currentMetadata)
	if recursive {
		/* find all children directories thanks to a request with their parent ID (which is the current dirId) */
		filter := bson.D{primitive.E{Key: "parentId", Value: dirId}}
		childDirIds, err := server.mongoColls[directoriesCollName].Find(server.ctx, filter)
		checkErr(err)
		for i := childDirIds; i != nil; i.Next(server.ctx) {
			i.Decode(currentDir)
			status = server.getOneDirectory(currentDir.dirId, recursive, filepath.Join(dirPath, currentDir.name), status)
		}
	}
	return status
}

func (server SantaclausServerImpl) GetDirectory(_ context.Context, req *MaeSanta.GetDirectoryRequest) (status *MaeSanta.GetDirectoryStatus, r error) {
	// TODO fetch all the files
	if !primitive.IsValidObjectID(req.DirId) {
		// err
	}
	objID, err := primitive.ObjectIDFromHex(req.DirId)
	checkErr(err)
	status = server.getOneDirectory(objID, req.IsRecursive, "", status)
	return status, r
}

// func (server SantaclausServerImpl) mustEmbedUnimplementedMaestro_Santaclaus_ServiceServer() {
// fmt.Println("hello")
// }
