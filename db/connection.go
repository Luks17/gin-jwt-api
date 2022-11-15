package db

import (
	"context"
	"golang-jwt/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = Connect()

func Connect() *mongo.Client {
	uri := config.GetUri()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	return client
}

// defer or place at the end of main
func CloseConnection(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}
}

func OpenCollection(client *mongo.Client, collection_name string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cluster0").Collection(collection_name)

	return collection
}
