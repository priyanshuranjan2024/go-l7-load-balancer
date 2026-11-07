package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ServerResponse struct {
	ServerID string `json:"server_id"`
	Message  string `json:"message"`
	Port     int    `json:"port"`
	Time     string `json:"time"`
}

type HealthResponse struct {
	ServerID string `json:"server_id"`
	Status   string `json:"status"`
}

func main() {
	serverID := flag.String("id", "server-1", "Backend server ID")
	port := flag.Int("port", 8081, "Backend server port")

	flag.Parse()

	mux := http.NewServeMux()

	// Main endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		response := ServerResponse{
			ServerID: *serverID,
			Message:  "Request successfully handled by backend",
			Port:     *port,
			Time:     time.Now().Format(time.RFC3339),
		}

		json.NewEncoder(w).Encode(response)
	})

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		response := HealthResponse{
			ServerID: *serverID,
			Status:   "healthy",
		}

		json.NewEncoder(w).Encode(response)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}

	log.Printf("Starting %s on port %d", *serverID, *port)

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
