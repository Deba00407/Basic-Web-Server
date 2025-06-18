package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	dbName         = "Goserver"
	collectionName = "users"
)

var Collection *mongo.Collection

func MakeConnectionToDB() {

	// Load env file
	func() {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}()

	// check if the connection string exists in the env variables
	if _, ok := os.LookupEnv("MONGODB_URI"); !ok {
		log.Fatal("MONGODB_URI not found in environment")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MongoDB connection string not set")
		return
	}

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to mongo client: %v", err)
	}

	log.Println("MongoDB Client connection successful")

	Collection = client.Database(dbName).Collection(collectionName)
}
