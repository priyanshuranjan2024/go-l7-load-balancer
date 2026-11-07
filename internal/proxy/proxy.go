package proxy

import (
	"log"
	"net/http"

	"go-l7-load-balancer/internal/backend"
	"go-l7-load-balancer/internal/scheduler"
)

type LoadBalancer struct {
	scheduler scheduler.Scheduler
}

// Optional interface for schedulers that track active requests.
type backendReleaser interface {
	ReleaseBackend(*backend.Backend)
}

func New(s scheduler.Scheduler) *LoadBalancer {
	return &LoadBalancer{
		scheduler: s,
	}
}

func (lb *LoadBalancer) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Select a backend.
	selected := lb.scheduler.NextBackend()

	if selected == nil {
		http.Error(
			w,
			"No available backend servers",
			http.StatusServiceUnavailable,
		)
		return
	}

	// Release the backend after the request finishes,
	// but only for schedulers that track active requests.
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
