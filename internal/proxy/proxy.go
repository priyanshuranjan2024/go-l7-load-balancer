package proxy

import (
	"log"
	"net/http"
	"time"

	"go-l7-load-balancer/internal/backend"
	"go-l7-load-balancer/internal/metrics"
	"go-l7-load-balancer/internal/scheduler"
)

type LoadBalancer struct {
	scheduler scheduler.Scheduler
	metrics   *metrics.Metrics
}

type backendReleaser interface {
	ReleaseBackend(*backend.Backend)
}

func New(s scheduler.Scheduler, m *metrics.Metrics) *LoadBalancer {
	return &LoadBalancer{
		scheduler: s,
		metrics:   m,
	}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	defer func() {
		lb.metrics.RecordRequest(time.Since(start))
	}()

	selected := lb.scheduler.NextBackend()

	if selected == nil {
		http.Error(w, "No available backend servers", http.StatusServiceUnavailable)
		return
	}

	if releaser, ok := lb.scheduler.(backendReleaser); ok {
		defer releaser.ReleaseBackend(selected)
	}

	log.Printf(
		"Forwarding request %s %s to %s",
		r.Method,
		r.URL.Path,
		selected.URL.Host,
	)

	selected.Proxy.ServeHTTP(w, r)
}
