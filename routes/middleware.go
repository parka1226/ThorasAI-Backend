package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Middleware to copy query parameters into the body
func QueryParamsToBodyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		baseURL := "http://localhost:8080"
		fullURL := fmt.Sprintf("%s%s", baseURL, r.URL.Path)
		parsedURL, err := url.Parse(fullURL)
		if err != nil {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		if r.URL.Host == "" {
			r.URL.Host = parsedURL.Host
		}
		bodyBytes, _ := io.ReadAll(r.Body)
		queryParams := r.URL.Query()
		for key, values := range queryParams {
			if len(values) > 0 {
				bodyBytes = append(bodyBytes, []byte(fmt.Sprintf("\"%s\":\"%s\"", key, values[0]))...)
			}
		}
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		next.ServeHTTP(w, r)
	})
}

// Helper function to create a request with JSON body
func createRequest(method, url string, body interface{}) (*http.Request, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
