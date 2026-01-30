package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Response structure for JSON responses
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Status  string      `json:"status"`
}

// ListFilesRequest for listing files
type ListFilesRequest struct {
	Path string `json:"path"`
}

// ProcessRequest for processing files
type ProcessRequest struct {
	Path      string   `json:"path"`
	Files     []string `json:"files"`
	Separator string   `json:"separator"`
}

// FileInfo represents a file's information
type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"isDir"`
	ModTime string `json:"modTime"`
}

// SSE client connection
type SSEClient struct {
	ID       string
	Messages chan string
}

// SSE broker
type SSEBroker struct {
	clients    map[string]*SSEClient
	mu         sync.RWMutex
	register   chan *SSEClient
	unregister chan *SSEClient
	broadcast  chan string
}

var broker *SSEBroker

func init() {
	broker = &SSEBroker{
		clients:    make(map[string]*SSEClient),
		register:   make(chan *SSEClient),
		unregister: make(chan *SSEClient),
		broadcast:  make(chan string, 100),
	}
	go broker.run()
}

func (b *SSEBroker) run() {
	for {
		select {
		case client := <-b.register:
			b.mu.Lock()
			b.clients[client.ID] = client
			b.mu.Unlock()
			log.Printf("SSE client %s connected", client.ID)

		case client := <-b.unregister:
			b.mu.Lock()
			if _, ok := b.clients[client.ID]; ok {
				close(client.Messages)
				delete(b.clients, client.ID)
				log.Printf("SSE client %s disconnected", client.ID)
			}
			b.mu.Unlock()

		case message := <-b.broadcast:
			b.mu.RLock()
			for _, client := range b.clients {
				select {
				case client.Messages <- message:
				default:
					// Client buffer full, skip
				}
			}
			b.mu.RUnlock()
		}
	}
}

func (b *SSEBroker) SendMessage(message string) {
	select {
	case b.broadcast <- message:
	default:
		// Broadcast channel full, skip
	}
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

// listFilesHandler lists files in a directory
func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ListFilesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate path
	if req.Path == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	// Check if path exists
	info, err := os.Stat(req.Path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Message: fmt.Sprintf("Path does not exist: %v", err),
			Status:  "error",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if !info.IsDir() {
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Message: "Path is not a directory",
			Status:  "error",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Read directory
	entries, err := os.ReadDir(req.Path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Message: fmt.Sprintf("Error reading directory: %v", err),
			Status:  "error",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Filter for .txt files only
	var fileInfos []FileInfo
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".txt") {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			fileInfos = append(fileInfos, FileInfo{
				Name:    entry.Name(),
				Size:    info.Size(),
				IsDir:   false,
				ModTime: info.ModTime().Format(time.RFC3339),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: fmt.Sprintf("Found %d text files", len(fileInfos)),
		Status:  "success",
		Data:    fileInfos,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// processFilesHandler processes files
func processFilesHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Path == "" || len(req.Files) == 0 {
		http.Error(w, "Path and files are required", http.StatusBadRequest)
		return
	}

	if req.Separator == "" {
		req.Separator = ":"
	}

	// Send immediate response
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: fmt.Sprintf("Processing started for %d files", len(req.Files)),
		Status:  "success",
		Data: map[string]interface{}{
			"filesCount": len(req.Files),
		},
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}

	// Process files in background
	go func() {
		broker.SendMessage(fmt.Sprintf("Starting processing of %d files...", len(req.Files)))
		broker.SendMessage(fmt.Sprintf("Using separator: '%s'", req.Separator))
		
		// Determine worker count (use CPU count or 8, whichever is less)
		workerCount := runtime.NumCPU()
		if workerCount > 8 {
			workerCount = 8
		}
		broker.SendMessage(fmt.Sprintf("Using %d worker threads", workerCount))

		// Create processor
		processor := NewFileProcessor(req.Separator, workerCount, func(msg string) {
			broker.SendMessage(msg)
		})

		// Process files
		startTime := time.Now()
		results := processor.ProcessFiles(req.Path, req.Files)
		duration := time.Since(startTime)

		// Send summary
		totalValid := 0
		totalInvalid := 0
		errorCount := 0
		for _, result := range results {
			if result.Error != nil {
				errorCount++
				broker.SendMessage(fmt.Sprintf("Error processing %s: %v", filepath.Base(result.OriginalPath), result.Error))
			} else {
				totalValid += result.ValidCount
				totalInvalid += result.InvalidCount
				broker.SendMessage(fmt.Sprintf("Completed: %s - Valid: %d, Invalid: %d", 
					filepath.Base(result.OriginalPath), result.ValidCount, result.InvalidCount))
			}
		}

		broker.SendMessage(fmt.Sprintf("Processing complete in %v", duration))
		broker.SendMessage(fmt.Sprintf("Total valid emails: %d", totalValid))
		broker.SendMessage(fmt.Sprintf("Total invalid/duplicate emails: %d", totalInvalid))
		broker.SendMessage(fmt.Sprintf("Files processed: %d", len(results)-errorCount))
		if errorCount > 0 {
			broker.SendMessage(fmt.Sprintf("Files with errors: %d", errorCount))
		}
		broker.SendMessage("DONE")
	}()
}

// sseHandler handles SSE connections
func sseHandler(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create client
	client := &SSEClient{
		ID:       fmt.Sprintf("client-%d", time.Now().UnixNano()),
		Messages: make(chan string, 100),
	}

	broker.register <- client
	defer func() {
		broker.unregister <- client
	}()

	// Send initial connection message
	fmt.Fprintf(w, "data: Connected to server\n\n")
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// Listen for messages
	for {
		select {
		case message, ok := <-client.Messages:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", message)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			return
		}
	}
}

func main() {
	// Serve frontend static files
	// Try relative path first (when run from extracted folder)
	frontendPath := "frontend"
	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		// Fallback: try relative to binary location
		frontendPath = filepath.Join(filepath.Dir(os.Args[0]), "..", "frontend")
	}
	fs := http.FileServer(http.Dir(frontendPath))
	http.Handle("/", fs)
	
	// Register routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/api/data", dataHandler)
	http.HandleFunc("/api/list-files", listFilesHandler)
	http.HandleFunc("/api/process", processFilesHandler)
	http.HandleFunc("/api/events", sseHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("API endpoints:\n")
	fmt.Printf("  - http://localhost:%s/health\n", port)
	fmt.Printf("  - http://localhost:%s/api\n", port)
	fmt.Printf("  - http://localhost:%s/api/data\n", port)
	fmt.Printf("  - http://localhost:%s/api/list-files\n", port)
	fmt.Printf("  - http://localhost:%s/api/process\n", port)
	fmt.Printf("  - http://localhost:%s/api/events\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
