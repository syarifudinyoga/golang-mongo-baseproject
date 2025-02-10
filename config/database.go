package config

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

var DB *mongo.Client

func ConnectDB() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("URL_DB")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	DB = client
	log.Println("Connected to MongoDB!")
}

func GetCollection(collectionName string) *mongo.Collection {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
	return DB.Database(os.Getenv("COLLECTION_DB")).Collection(collectionName)
}
