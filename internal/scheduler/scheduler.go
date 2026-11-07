package scheduler

import "go-l7-load-balancer/internal/backend"

type Scheduler interface {
	NextBackend() *backend.Backend
}
