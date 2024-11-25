package database

import (
	"context"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestMongoMock(t *testing.T) {
	client, err := CreateMockClient()
	if err != nil {
		log.Fatalf("Error creating mock MongoDB client: %v", err)
	}

	client.Database("mockDatabase")
	client.Collection("mockCollection")

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatalf("Error disconnecting mock MongoDB client: %v", err)
	}
}
