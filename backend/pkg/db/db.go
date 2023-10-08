package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

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

func InsertSongs(collection *mongo.Collection, songName string, songArtist string) (*mongo.InsertOneResult, error) {
	previousSong, err := FetchMostRecentSong(collection)
	fmt.Printf("THE PREVIOUS SONG IS: %s ---- %s", previousSong["song"], previousSong["artist"])

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if previousSong != nil && previousSong["song"] == songName && previousSong["artist"] == songArtist {
		return nil, errors.New("song is already the most recently inserted in the playlist")
	}

	doc := bson.M{
		"song":   songName,
		"artist": songArtist,
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
