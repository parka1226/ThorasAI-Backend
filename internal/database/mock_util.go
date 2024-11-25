package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoClientInterface is the interface that the mock client will implement
type MongoClientInterface interface {
	Connect(ctx context.Context) error
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
	Disconnect(ctx context.Context) error
}

// MockMongoClient is a mock struct that implements MongoClientInterface
type MockMongoClient struct {
	connected bool
}

// NewMockMongoClient creates and returns a new mock MongoDB client
func NewMockMongoClient() *MockMongoClient {
	return &MockMongoClient{}
}

// CreateMockClient creates and returns a mock MongoDB client.
func CreateMockClient() (MongoClientInterface, error) {
	mockClient := NewMockMongoClient()
	err := mockClient.Connect(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create mock MongoDB client: %v", err)
	}
	log.Println("Successfully created and connected mock MongoDB client")
	return mockClient, nil
}

// Connect simulates connecting to MongoDB
func (m *MockMongoClient) Connect(ctx context.Context) error {
	m.connected = true
	log.Println("Mock MongoDB client connected")
	return nil
}

// Ping simulates pinging the MongoDB server
func (m *MockMongoClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	if !m.connected {
		return fmt.Errorf("mock client is not connected")
	}
	log.Println("Mock MongoDB client pinged successfully")
	return nil
}

// Database simulates getting a MongoDB database
func (m *MockMongoClient) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	// Return a dummy database (mock behavior)
	log.Printf("Mock MongoDB client accessed database: %s\n", name)
	return &mongo.Database{}
}

// Collection simulates getting a MongoDB collection
func (m *MockMongoClient) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	// Return a dummy collection (mock behavior)
	log.Printf("Mock MongoDB client accessed collection: %s\n", name)
	return &mongo.Collection{}
}

// Disconnect simulates disconnecting the MongoDB client
func (m *MockMongoClient) Disconnect(ctx context.Context) error {
	if !m.connected {
		return fmt.Errorf("mock client is not connected: %w", errors.New("connection failure"))
	}
	m.connected = false
	log.Println("Mock MongoDB client disconnected")
	return nil
}
