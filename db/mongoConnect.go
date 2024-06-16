package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
	// Log the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current working directory:", err)
	}
	log.Println("Current working directory:", cwd)

	// Load the .env file from the root directory
	err = godotenv.Load("../../.env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Log the value of MONGODB_URI
	MONGODB_URI := os.Getenv("MONGODB_URI")
	log.Println("MONGODB_URI:", MONGODB_URI)

	if MONGODB_URI == "" {
		log.Fatal("MONGODB_URI is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	fmt.Println("Client", client)
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("Video-Transcoder").Collection(collectionName)
}
