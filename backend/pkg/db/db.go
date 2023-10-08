package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PeteMango/website-v2/pkg/definition"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitializeMongo() *mongo.Client {
	connectionString := "mongodb+srv://PeteMango:test123@songs.db2njvu.mongodb.net/?retryWrites=true&w=majority"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection.
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	return client
}

func InsertSongs(collection *mongo.Collection, songName string, songArtist string, songLink string) (*mongo.InsertOneResult, error) {
	previousSong, err := FetchMostRecentSong(collection)
	fmt.Printf("THE PREVIOUS SONG IS: %s ---- %s", previousSong["song"], previousSong["artist"])

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if previousSong != nil && previousSong["song"] == songName && previousSong["artist"] == songArtist {
		fmt.Printf("already in the database\n")
		return nil, nil
	}

	doc := bson.M{
		"song":   songName,
		"artist": songArtist,
		"link":   songLink,
	}

	result, err := collection.InsertOne(context.Background(), doc)
	return result, err
}

func FetchMostRecentSong(collection *mongo.Collection) (bson.M, error) {
	var result bson.M

	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})
	err := collection.FindOne(context.Background(), bson.M{}, opts).Decode(&result)

	return result, err
}

func UpsertTokens(collection *mongo.Collection, accessToken, refreshToken string, expiryTime time.Time) error {
	filter := bson.M{}
	update := bson.M{
		"$set": bson.M{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expiry_time":   expiryTime,
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	return err
}

func FetchTokens(collection *mongo.Collection) (*definition.Token, time.Time, error) {
	var token definition.Token

	err := collection.FindOne(context.Background(), bson.M{}).Decode(&token)
	if err != nil {
		return nil, time.Time{}, err
	}

	expiryTime := token.ExpiryTime

	return &token, expiryTime, nil
}

func DeleteTokens(collection *mongo.Collection) error {
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	return err
}
