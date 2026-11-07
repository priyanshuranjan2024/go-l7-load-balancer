package scheduler

import (
	"sync"

	"go-l7-load-balancer/internal/backend"
)

type RoundRobin struct {
	backends []*backend.Backend
	current  int
	mu       sync.Mutex
}

func NewRoundRobin(backends []*backend.Backend) *RoundRobin {
	return &RoundRobin{
		backends: backends,
		current:  0,
	}
}

func (rr *RoundRobin) NextBackend() *backend.Backend {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if len(rr.backends) == 0 {
		return nil
	}

	selected := rr.backends[rr.current]

	rr.current = (rr.current + 1) % len(rr.backends)

	return selected
}
