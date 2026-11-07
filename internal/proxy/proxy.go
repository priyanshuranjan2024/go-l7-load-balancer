package proxy

import (
	"log"
	"net/http"

	"go-l7-load-balancer/internal/scheduler"
)

type LoadBalancer struct {
	scheduler scheduler.Scheduler
}

func New(s scheduler.Scheduler) *LoadBalancer {
	return &LoadBalancer{
		scheduler: s,
	}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Select the next backend
	selected := lb.scheduler.NextBackend()

	if selected == nil {
		http.Error(
			w,
			"No available backend servers",
			http.StatusServiceUnavailable,
		)
		return
	}

	log.Printf(
		"Forwarding request %s %s to %s",
		r.Method,
		r.URL.Path,
		selected.URL.Host,
	)

	// Forward the request
	selected.Proxy.ServeHTTP(w, r)
}
