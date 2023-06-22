package common

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func GetDBCollection(col string) *mongo.Collection {
	return db.Collection(col)
}

// Initializes db (called by main.go) - set URI in .env file
func InitDB() error {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	db = client.Database("PaddeBillingApp")

	return nil
}

func CloseDB() error {
	return db.Client().Disconnect(context.Background())
}
