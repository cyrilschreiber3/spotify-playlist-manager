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

	log.Println("Connecting to MongoDB...")
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	dbClient = client

	log.Println("Connected to MongoDB")
}

func mongoClose() {
	if dbClient != nil {
		log.Println("Disconnecting from MongoDB...")
		if err := dbClient.Disconnect(context.Background()); err != nil {
			log.Fatal("Error disconnecting from MongoDB:", err)
		}
		log.Println("Disconnected from MongoDB")
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

	log.Println("Opening collection:", collectionName)
	collection = dbClient.Database(databaseName).Collection(collectionName)
	if collection == nil {
		return nil, errors.New("failed to open collection: " + collectionName)
	}
	collections[collectionName] = collection
	return collection, nil
}

func GetUserByID(ctx context.Context, userID string) (models.User, error) {
	userCollection, err := openCollection("users")
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func CreateSession(ctx context.Context, session *models.UserSession) (*mongo.InsertOneResult, error) {
	sessionCollection, err := openCollection("sessions")
	if err != nil {
		return nil, err
	}

	result, err := sessionCollection.InsertOne(ctx, session)
	return result, err
}

func CreateUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	userCollection, err := openCollection("users")
	if err != nil {
		return nil, err
	}

	result, err := userCollection.InsertOne(ctx, user)
	return result, err
}

func UpdateUserByID(ctx context.Context, userID string, update bson.M) (*mongo.UpdateResult, error) {
	userCollection, err := openCollection("users")
	if err != nil {
		return nil, err
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": update})
	return result, err
}

func GetSessionByID(ctx context.Context, sessionID string) (models.UserSession, error) {
	sessionCollection, err := openCollection("sessions")
	if err != nil {
		return models.UserSession{}, err
	}

	var session models.UserSession
	err = sessionCollection.FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&session)
	if err != nil {
		return models.UserSession{}, err
	}
	return session, nil
}

func UpdateSessionByID(ctx context.Context, sessionID string, update bson.M) (*mongo.UpdateResult, error) {
	sessionCollection, err := openCollection("sessions")
	if err != nil {
		return nil, err
	}

	result, err := sessionCollection.UpdateOne(ctx, bson.M{"session_id": sessionID}, bson.M{"$set": update})
	return result, err
}
