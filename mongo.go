package grok

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NewMongoConnection ...
func NewMongoConnection(connectionString string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))

	if err != nil {
		logrus.WithError(err).Panic("Error connecting to MongoDB")
	}

	client.Connect(context.Background())

	err = client.Ping(context.Background(), readpref.Primary())

	if err != nil {
		logrus.WithError(err).Panic("Error pinging MongoDB")
	}

	return client
}
