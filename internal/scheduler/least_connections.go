package scheduler

import (
	"sync"

	"go-l7-load-balancer/internal/backend"
)

type LeastConnections struct {
	backends  []*backend.Backend
	nextIndex int
	mu        sync.Mutex
}

func NewLeastConnections(
	backends []*backend.Backend,
) *LeastConnections {
	return &LeastConnections{
		backends:  backends,
		nextIndex: 0,
	}
}

func (lc *LeastConnections) NextBackend() *backend.Backend {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if len(lc.backends) == 0 {
		return nil
	}

	minConnections := -1
	var selected *backend.Backend
	selectedIndex := -1

	// Find the backend with the fewest active requests.
	// Start from nextIndex to rotate ties.
	for i := 0; i < len(lc.backends); i++ {
		index := (lc.nextIndex + i) % len(lc.backends)
		current := lc.backends[index]

		active := current.GetActive()

		if minConnections == -1 || active < minConnections {
			minConnections = active
			selected = current
			selectedIndex = index
		}
	}

	// Reserve this backend before releasing the scheduler lock.
	selected.IncrementActive()

	// Start tie-breaking from the next backend.
	lc.nextIndex = (selectedIndex + 1) % len(lc.backends)

	return selected
}

// ReleaseBackend decreases the active request count.
func (lc *LeastConnections) ReleaseBackend(
	b *backend.Backend,
) {
	b.DecrementActive()
}
