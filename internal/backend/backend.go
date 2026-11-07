package backend

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL   *url.URL
	Proxy *httputil.ReverseProxy

	mu                sync.Mutex
	activeConnections int
}

func New(rawURL string) (*Backend, error) {
	target, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if target.Scheme == "" || target.Host == "" {
		return nil, fmt.Errorf("invalid backend URL: %s", rawURL)
	}

	return &Backend{
		URL:   target,
		Proxy: httputil.NewSingleHostReverseProxy(target),
	}, nil
}

// IncrementActive increases the active request count.
func (b *Backend) IncrementActive() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.activeConnections++
}

// DecrementActive decreases the active request count.
func (b *Backend) DecrementActive() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.activeConnections > 0 {
		b.activeConnections--
	}
}

// GetActive returns the current active request count.
func (b *Backend) GetActive() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.activeConnections
}
