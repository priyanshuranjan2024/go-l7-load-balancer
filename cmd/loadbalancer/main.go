package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-l7-load-balancer/internal/backend"
	"go-l7-load-balancer/internal/metrics"
	"go-l7-load-balancer/internal/proxy"
	"go-l7-load-balancer/internal/scheduler"
)

func parseWeights(raw string) ([]int, error) {
	parts := strings.Split(raw, ",")
	weights := make([]int, 0, len(parts))

	for _, part := range parts {
		weight, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("invalid weight %q: %w", part, err)
		}

		if weight <= 0 {
			return nil, fmt.Errorf("weights must be positive")
		}

		weights = append(weights, weight)
	}

	return weights, nil
}

func main() {
	addr := flag.String("addr", ":8080", "Load balancer address")

	backendList := flag.String(
		"backends",
		"http://localhost:8081,http://localhost:8082,http://localhost:8083",
		"Comma-separated backend URLs",
	)

	algorithm := flag.String(
		"algorithm",
		"rr",
		"Scheduling algorithm: rr, wrr or lc",
	)

	weightsFlag := flag.String(
		"weights",
		"5,3,2",
		"Backend weights for WRR",
	)

	flag.Parse()

	// Initialize backend servers.
	var backends []*backend.Backend

	for _, rawURL := range strings.Split(*backendList, ",") {
		rawURL = strings.TrimSpace(rawURL)

		if rawURL == "" {
			continue
		}

		b, err := backend.New(rawURL)
		if err != nil {
			log.Fatalf("Failed to initialize backend: %v", err)
		}

		backends = append(backends, b)
		log.Printf("Registered backend: %s", b.URL)
	}

	if len(backends) == 0 {
		log.Fatal("No backend servers configured")
	}

	// Initialize selected scheduling algorithm.
	var selectedScheduler scheduler.Scheduler

	switch strings.ToLower(*algorithm) {
	case "rr":
		selectedScheduler = scheduler.NewRoundRobin(backends)

	case "wrr":
		weights, err := parseWeights(*weightsFlag)
		if err != nil {
			log.Fatalf("Invalid weights: %v", err)
		}

		if len(weights) != len(backends) {
			log.Fatal("Number of weights must match number of backends")
		}

		wrr, err := scheduler.NewWeightedRoundRobin(backends, weights)
		if err != nil {
			log.Fatalf("Failed to initialize WRR: %v", err)
		}

		selectedScheduler = wrr

	case "lc":
		selectedScheduler = scheduler.NewLeastConnections(backends)

	default:
		log.Fatalf("Unsupported scheduling algorithm: %s", *algorithm)
	}

	// Initialize metrics collector.
	metricsCollector := metrics.New()

	// Initialize load balancer with scheduler and metrics.
	lb := proxy.New(selectedScheduler, metricsCollector)

	// Configure HTTP routes.
	mux := http.NewServeMux()

	// Forward normal requests to the load balancer.
	mux.Handle("/", lb)

	// Expose performance metrics.
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"total_requests": metricsCollector.TotalRequests(),
			"average_latency_ms": float64(
				metricsCollector.AverageLatency(),
			) / float64(time.Millisecond),
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode metrics: %v", err)
		}
	})

	// Start HTTP server.
	server := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	log.Printf("Load balancer listening on %s", *addr)
	log.Printf("Scheduling algorithm: %s", strings.ToUpper(*algorithm))

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Load balancer failed: %v", err)
	}
}
