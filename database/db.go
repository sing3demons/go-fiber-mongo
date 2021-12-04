package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitDB() *mongo.Database {
	port := os.Getenv("MDB_PORT")
	username := os.Getenv("MDB_USERNAME")
	password := os.Getenv("MDB_PASSWORD")
	host := os.Getenv("MDB_HOST")
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")

	db := client.Database("go-product")

	return db
}
