package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Models struct {
	AuthenticationModel AuthenticationModel
}
type AuthenticationModel struct {
	Id        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"password"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		AuthenticationModel: AuthenticationModel{},
	}
}

func (auth *AuthenticationModel) GetByEmail(email string) (*AuthenticationModel, error) {

	collection := client.Database("auth").Collection("auth")
	filter := bson.D{{"email", email}}
	var result AuthenticationModel

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (auth *AuthenticationModel) GetAll() ([]*AuthenticationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	collection := client.Database("auth").Collection("aut")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all user error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var auths []*AuthenticationModel
	for cursor.Next(ctx) {
		var item AuthenticationModel
		err := cursor.Decode(&item)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		auths = append(auths, &item)
	}

	return auths, nil
}
