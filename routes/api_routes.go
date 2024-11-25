package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"example.com/m/internal/database"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// API endpoint handler
// HTTP handler for the endpoint
func GetTrafficWithService(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusBadRequest)
		return
	}

	db, ok := body["database"].(string)
	if !ok || db == "" {
		http.Error(w, "Missing or invalid 'database' parameter", http.StatusBadRequest)
		return
	}
	networkCollection, ok := body["networkCollection"].(string)
	if !ok || networkCollection == "" {
		http.Error(w, "Missing or invalid 'networkCollection' parameter", http.StatusBadRequest)
		return
	}
	serviceName, ok := body["serviceName"].(string)
	if !ok || networkCollection == "" {
		http.Error(w, "Missing or invalid 'serviceName' parameter", http.StatusBadRequest)
		return
	}

	mongoURI := os.Getenv("MONGO_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		http.Error(w, "Invalid Client Connection", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		http.Error(w, "Failed to connect into MongoDB", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)
	results, err := database.AggregateTrafficWithService(ctx, client, db, networkCollection, serviceName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error running aggregation query: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, fmt.Sprintf("Failed to send response: %v", err), http.StatusInternalServerError)
	}
}

// SetupRouter with CORS enabled
func SetupRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/TrafficService", GetTrafficWithService).Methods("POST")
	r.Use(QueryParamsToBodyMiddleware)

	// Set up CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(r)

	return corsHandler // Return the CORS middleware-wrapped handler
}
