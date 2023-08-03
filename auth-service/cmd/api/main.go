package main

import (
	"auth-service/data"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

const (
	MONGO_URL = "<your-mongodb-url>"
	PORT      = "80"
)

type App struct {
	Models data.Models
}

func main() {
	app := App{
		Models: data.New(ConnectToMongo()),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.Routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func ConnectToMongo() *mongo.Client {
	clientOptions := options.Client().ApplyURI(MONGO_URL)
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to the server:", err)
	}
	log.Println("Connected to the server")
	return c
}
