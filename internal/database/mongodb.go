package database

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AggregateTrafficWithService(ctx context.Context, client *mongo.Client, database string, networkCollection string, serviceName string) ([]bson.M, error) {
	db := client.Database(database)
	trafficCollection := db.Collection(networkCollection)

	serviceIP, err := getServiceIPByName(ctx, db, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get service IP: %w", err)
	}

	pipeline := mongo.Pipeline{
		// Step 1: Lookup to join network traffic with service data
		bson.D{
			{Key: "$lookup", Value: bson.D{
				//testCollectionA stores service data and testCollectionB stores network traffic data
				{Key: "from", Value: "testcollectionA"},
				{Key: "let", Value: bson.D{
					{Key: "source_ip", Value: "$source_ip"},
					{Key: "destination_ip", Value: "$destination_ip"},
				}},
				{Key: "pipeline", Value: bson.A{
					bson.D{
						{Key: "$match", Value: bson.D{
							{Key: "$or", Value: bson.A{
								bson.D{{Key: "ip_address", Value: "$$source_ip"}},
								bson.D{{Key: "ip_address", Value: "$$destination_ip"}},
							}},
						}},
					},
					bson.D{
						{Key: "$project", Value: bson.D{
							{Key: "name", Value: 1},           // Service name
							{Key: "ip_address", Value: 1},     // IP address of service
							{Key: "listening_port", Value: 1}, // Port of service
						}},
					},
				}},
				{Key: "as", Value: "service_info"},
			}},
		},

		// Step 2: Unwind the result from the lookup to handle multiple matches
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "source_ip", Value: serviceIP}},
					bson.D{{Key: "destination_ip", Value: serviceIP}},
				}},
			}},
		},

		// Step 3: Project stage to select the fields we want to return
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "source_ip", Value: 1},
				{Key: "source_port", Value: 1},
				{Key: "destination_ip", Value: 1},
				{Key: "destination_port", Value: 1},
				{Key: "status", Value: 1},
				{Key: "service_name", Value: "$service_info.name"},
				{Key: "service_ip", Value: "$service_info.ip_address"},
				{Key: "service_port", Value: "$service_info.listening_port"},
			}},
		},
	}

	// Execute the aggregation pipeline
	cursor, err := trafficCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate data: %w", err)
	}
	defer cursor.Close(ctx)

	// Collect the results into a slice of documents
	var results []bson.M
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Printf("failed to decode result: %v", err)
			continue
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	return results, nil
}

func getServiceIPByName(ctx context.Context, db *mongo.Database, serviceName string) (serviceIP string, err error) {
	serviceCollection := db.Collection("testcollectionA")

	// Define the filter for the service name
	filter := bson.M{
		"name": serviceName,
	}

	cursor, err := serviceCollection.Find(ctx, filter)
	if err != nil {
		return "", fmt.Errorf("failed to query service collection: %w", err)
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and get the result
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return "", fmt.Errorf("failed to decode service result: %w", err)
		}
		serviceIP, ok := result["ip_address"].(string)
		if !ok {
			return "", fmt.Errorf("invalid IP address format for service %s", serviceName)
		}
		return serviceIP, nil
	}

	if err := cursor.Err(); err != nil {
		return "", fmt.Errorf("cursor iteration error: %w", err)
	}
	return "", fmt.Errorf("service not found: %s", serviceName)
}
