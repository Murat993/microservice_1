package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"log-service/cmd/api/routes"
	"log-service/data"
	"net/http"
	"net/rpc"
	"time"
)

const (
	webPort  = "80"
	mongoUrl = "mongodb://mongo:27017"
)

var client *mongo.Client

func main() {
	log.Println("Starting the application")
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

	app := routes.Config{
		Models: data.New(client),
	}

	//Register RPC server
	err = rpc.Register(new(RPCServer))
	go app.RPCListen()

	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.Routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
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
