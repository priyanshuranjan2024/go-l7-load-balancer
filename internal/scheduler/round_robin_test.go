package scheduler

import (
	"testing"

	"go-l7-load-balancer/internal/backend"
)

func TestRoundRobinDistribution(t *testing.T) {
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

	rr := NewRoundRobin(backends)

	expected := []string{
		"localhost:8081",
		"localhost:8082",
		"localhost:8083",
		"localhost:8081",
		"localhost:8082",
		"localhost:8083",
	}

	for i, want := range expected {
		got := rr.NextBackend()

		if got == nil {
			t.Fatalf("Request %d: expected backend %s, got nil", i+1, want)
		}

		if got.URL.Host != want {
			t.Errorf(
				"Request %d: expected %s, got %s",
				i+1,
				want,
				got.URL.Host,
			)
		}
	}
}
