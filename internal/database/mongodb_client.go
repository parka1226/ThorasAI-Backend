package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoClient is the real MongoDB client
type MongoClient struct {
	client *mongo.Client
}

// Disconnect cleans up the MongoDB connection.
func (mc *MongoClient) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mc.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	log.Println("Disconnected from MongoDB")
	return nil
}

// GetDatabase returns a reference to the specified database.
func (mc *MongoClient) GetDatabase(name string) *mongo.Database {
	return mc.client.Database(name)
}

// CreateMongoClient creates and returns a MongoDB client connected to the database.
func CreateMongoClient() (*MongoClient, error) {
	mongoURI := os.Getenv("MONGO_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &MongoClient{
		client: client,
	}, nil
}

// InsertAPIData inserts service or network api data into a MongoDB collection.
func (m *MongoClient) InsertAPIData(database, collection string, serviceData interface{}) error {
	db := m.client.Database(database)
	coll := db.Collection(collection)

	_, err := coll.InsertOne(context.Background(), serviceData)
	if err != nil {
		return fmt.Errorf("failed to insert service data: %v", err)
	}
	return nil
}

// GetAllDataFromMongo fetches all data from the specified MongoDB collection.
func (m *MongoClient) GetAllDataFromMongo(database, collection string) ([]interface{}, error) {
	coll := m.client.Database(database).Collection(collection)
	ctx := context.Background()
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer cursor.Close(ctx)

	var results []interface{}

	// Iterate through the cursor and decode each document
	for cursor.Next(ctx) {
		var result interface{}
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return results, nil
}

// InsertJSONData parses JSON from a file and inserts it into the MongoDB collection
func (m *MongoClient) InsertJSONData(dbName, collectionName, filePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %v", err)
	}

	// Parse the JSON data into a slice of bson.M
	var documents []bson.M
	err = json.Unmarshal(jsonData, &documents)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	// Convert documents (which is of type []bson.M) into []interface{}
	var interfaces []interface{}
	for _, doc := range documents {
		interfaces = append(interfaces, doc)
	}

	collection := m.client.Database(dbName).Collection(collectionName)

	_, err = collection.InsertMany(ctx, interfaces)
	if err != nil {
		return fmt.Errorf("failed to insert data into MongoDB: %v", err)
	}
	return nil
}
