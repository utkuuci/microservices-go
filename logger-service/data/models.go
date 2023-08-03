package data

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

var client *mongo.Client

type Log struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Models struct {
	Log Log
}

func New(mon *mongo.Client) Models {
	client = mon

	return Models{
		Log: Log{},
	}
}
func (l *Log) Insert(entry Log) error {
	collection := client.Database("logs").Collection("logs")

	doc := Log{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Println("Inserting log error:", err)
		return err
	}
	return nil
}
