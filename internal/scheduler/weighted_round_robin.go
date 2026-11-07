package scheduler

import (
	"fmt"
	"sync"

	"go-l7-load-balancer/internal/backend"
)

type WeightedRoundRobin struct {
	backends  []*backend.Backend
	weights   []int
	current   int
	remaining int
	mu        sync.Mutex
}

func NewWeightedRoundRobin(
	backends []*backend.Backend,
	weights []int,
) (*WeightedRoundRobin, error) {

	if len(backends) == 0 {
		return nil, fmt.Errorf("no backend servers configured")
	}

	if len(backends) != len(weights) {
		return nil, fmt.Errorf("number of weights must match number of backends")
	}

	for _, weight := range weights {
		if weight <= 0 {
			return nil, fmt.Errorf("backend weights must be positive")
		}
	}

	return &WeightedRoundRobin{
		backends:  backends,
		weights:   weights,
		current:   0,
		remaining: weights[0],
	}, nil
}

func (wrr *WeightedRoundRobin) NextBackend() *backend.Backend {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	selected := wrr.backends[wrr.current]

	wrr.remaining--

	if wrr.remaining == 0 {
		wrr.current = (wrr.current + 1) % len(wrr.backends)
		wrr.remaining = wrr.weights[wrr.current]
	}

	return selected
}
