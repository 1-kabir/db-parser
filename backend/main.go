package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Response structure for JSON responses
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Status  string      `json:"status"`
}

// enableCORS adds CORS headers to allow frontend connections
// Note: Using "*" for Allow-Origin is suitable for development only.
// In production, specify exact origins or use environment-based configuration.
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// healthHandler returns the health status of the server
func healthHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: "Server is running",
		Status:  "ok",
		Data: map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding health response: %v", err)
	}
}

// apiHandler handles API requests
func apiHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: "Welcome to the API",
		Status:  "success",
		Data: map[string]interface{}{
			"version": "1.0.0",
			"endpoints": []string{
				"/health",
				"/api",
				"/api/data",
			},
		},
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding api response: %v", err)
	}
}

// dataHandler returns sample data
func dataHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: "Sample data retrieved successfully",
		Status:  "success",
		Data: []map[string]interface{}{
			{"id": 1, "name": "Item 1", "description": "First sample item"},
			{"id": 2, "name": "Item 2", "description": "Second sample item"},
			{"id": 3, "name": "Item 3", "description": "Third sample item"},
		},
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding data response: %v", err)
	}
}

func main() {
	// Register routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/api/data", dataHandler)

	// Start server
	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("API endpoints:\n")
	fmt.Printf("  - http://localhost:%s/health\n", port)
	fmt.Printf("  - http://localhost:%s/api\n", port)
	fmt.Printf("  - http://localhost:%s/api/data\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
