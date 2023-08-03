package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	data2 "logger-serivce/data"
	"net/http"
)

const (
	WEB_PORT  = "80"
	MONGO_URL = "<your-mongo-db-url>"
)

type App struct {
	Models data2.Models
	Port   string
}

func main() {

	app := App{
		Models: data2.New(ConnectMongo()),
		Port:   WEB_PORT,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.Port),
		Handler: app.Routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func ConnectMongo() *mongo.Client {
	clientOptions := options.Client().ApplyURI(MONGO_URL)
	//clientOptions.SetAuth(options.Credential{
	//	Username: "admin",
	//	Password: "password",
	//})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
	}
	log.Println("Connected to the db")
	return c
}
