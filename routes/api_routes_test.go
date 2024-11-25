package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestAPIRoute(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"database":          "testdb",
			"networkCollection": "testcollectionB",
			"serviceName":       "Gaming UI",
		}

		body, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		req, err := http.NewRequest("POST", "/TrafficService", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		r := SetupRouter()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")

		var response []map[string]interface{}
		err = json.NewDecoder(rr.Body).Decode(&response)
		log.Debug().Msgf("%+v", response)
		assert.NoError(t, err, "Failed to decode response body")
		assert.Greater(t, len(response), 0, "Expected non-empty response")
	})

	t.Run("QueryParamsToBodyMiddleware", func(t *testing.T) {
		router := mux.NewRouter()

		router.HandleFunc("/TrafficService", func(w http.ResponseWriter, r *http.Request) {
			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			assert.NoError(t, err)
			ipAddress := body["ipAddress"].(string)
			assert.Equal(t, "127.0.0.1", ipAddress)
			w.WriteHeader(http.StatusOK)
		}).Methods("POST")

		router.Use(QueryParamsToBodyMiddleware)
		queryParams := map[string]string{
			"ipAddress": "127.0.0.1",
		}

		req, err := createRequest("POST", "/TrafficService", queryParams)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

}
