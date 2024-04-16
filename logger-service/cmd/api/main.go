package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	webPort  = "8080"
	prcPort  = "5001"
	mongoUrl = "mongodb://mongo:27017"
	gRPCPort = "50051"
)

var client *mongo.Client

type Config struct {
	WebPort string
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()
}

func connectToMongo() (*mongo.Client, error) {
	// create a connection options
	clientOptions := options.Client().ApplyURI(mongoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	// connect to mongo
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println("Error connecting to mongo: ", err)
		return nil, err
	}
	return client, nil
}
