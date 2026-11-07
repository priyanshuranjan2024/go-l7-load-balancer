package scheduler

import (
	"testing"

	"go-l7-load-balancer/internal/backend"
)

func TestLeastConnectionsSelection(t *testing.T) {
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

	lc := NewLeastConnections(backends)

	// First three requests should be distributed
	// across the three backends.
	first := lc.NextBackend()
	second := lc.NextBackend()
	third := lc.NextBackend()

	if first.URL.Host != "localhost:8081" {
		t.Errorf("Expected backend 1, got %s", first.URL.Host)
	}

	if second.URL.Host != "localhost:8082" {
		t.Errorf("Expected backend 2, got %s", second.URL.Host)
	}

	if third.URL.Host != "localhost:8083" {
		t.Errorf("Expected backend 3, got %s", third.URL.Host)
	}

	// All backends now have one active request.
	// Release backend 2 so it has the fewest connections.
	lc.ReleaseBackend(second)

	next := lc.NextBackend()

	if next.URL.Host != "localhost:8082" {
		t.Errorf(
			"Expected backend 2, got %s",
			next.URL.Host,
		)
	}

	// Release all acquired requests.
	lc.ReleaseBackend(first)
	lc.ReleaseBackend(third)
	lc.ReleaseBackend(next)
}
