package scheduler

import (
	"testing"

	"go-l7-load-balancer/internal/backend"
)

func TestWeightedRoundRobinDistribution(t *testing.T) {
	var backends []*backend.Backend

	urls := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	for _, rawURL := range urls {
		b, err := backend.New(rawURL)
		if err != nil {
			t.Fatalf("Failed to create backend: %v", err)
		}

		backends = append(backends, b)
	}

	weights := []int{5, 3, 2}

	wrr, err := NewWeightedRoundRobin(backends, weights)
	if err != nil {
		t.Fatalf("Failed to create WRR scheduler: %v", err)
	}

	counts := make(map[string]int)

	for i := 0; i < 100; i++ {
		selected := wrr.NextBackend()

		if selected == nil {
			t.Fatal("Expected a backend, got nil")
		}

		counts[selected.URL.Host]++
	}

	expected := map[string]int{
		"localhost:8081": 50,
		"localhost:8082": 30,
		"localhost:8083": 20,
	}

	for host, want := range expected {
		if counts[host] != want {
			t.Errorf(
				"Backend %s: expected %d requests, got %d",
				host,
				want,
				counts[host],
			)
		}
	}
}

func TestWeightedRoundRobinInvalidWeights(t *testing.T) {
	b, err := backend.New("http://localhost:8081")
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewWeightedRoundRobin([]*backend.Backend{b}, []int{0})

	if err == nil {
		t.Error("Expected error for zero weight")
	}
}
