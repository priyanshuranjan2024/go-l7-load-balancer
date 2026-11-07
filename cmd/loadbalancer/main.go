package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go-l7-load-balancer/internal/backend"
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
		"Scheduling algorithm: rr or wrr",
	)

	weightsFlag := flag.String(
		"weights",
		"5,3,2",
		"Backend weights for WRR",
	)

	flag.Parse()

	// Initialize backend servers
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

	// Initialize selected scheduling algorithm
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

	default:
		log.Fatalf("Unsupported scheduling algorithm: %s", *algorithm)
	}

	// Initialize load balancer
	lb := proxy.New(selectedScheduler)

	server := &http.Server{
		Addr:    *addr,
		Handler: lb,
	}

	log.Printf("Load balancer listening on %s", *addr)
	log.Printf("Scheduling algorithm: %s", strings.ToUpper(*algorithm))

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Load balancer failed: %v", err)
	}
}
