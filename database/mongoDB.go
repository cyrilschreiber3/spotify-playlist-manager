package database

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/cyrilschreiber3/spotify-playlist-manager/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var dbClient *mongo.Client = nil
var collections = make(map[string]*mongo.Collection)

func mongoInit() {
	dbUri := os.Getenv("MONGO_URI")
	if dbUri == "" {
		log.Fatal("MONGO_URI is not set")
	}

	clientOptions := options.Client().ApplyURI(dbUri)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	dbClient = client
}

func mongoClose() {
	if dbClient != nil {
		if err := dbClient.Disconnect(context.Background()); err != nil {
			log.Fatal("Error disconnecting from MongoDB:", err)
		}
	}
}

func openCollection(collectionName string) (*mongo.Collection, error) {
	collection, exists := collections[collectionName]
	if exists {
		return collection, nil
	}

	databaseName := os.Getenv("MONGO_DB")
	if databaseName == "" {
		log.Fatal("MONGO_DB is not set")
	}

	if dbClient == nil {
		log.Fatal("Database connection is not established")
	}

	collection = dbClient.Database(databaseName).Collection(collectionName)
	if collection == nil {
		return nil, errors.New("failed to open collection: " + collectionName)
	}
	collections[collectionName] = collection
	return collection, nil
}

func GetUserByID(userID string) (models.User, error) {
	userCollection, err := openCollection("users")
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	err = userCollection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
