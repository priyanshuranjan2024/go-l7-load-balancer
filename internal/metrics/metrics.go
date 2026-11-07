package metrics

import (
	"sync/atomic"
	"time"
)

type Metrics struct {
	totalRequests atomic.Int64
	totalDuration atomic.Int64
}

func New() *Metrics {
	return &Metrics{}
}

func (m *Metrics) RecordRequest(duration time.Duration) {
	m.totalRequests.Add(1)
	m.totalDuration.Add(int64(duration))
}

func (m *Metrics) TotalRequests() int64 {
	return m.totalRequests.Load()
}

func (m *Metrics) AverageLatency() time.Duration {
	total := m.totalRequests.Load()

	if total == 0 {
		return 0
	}

	duration := m.totalDuration.Load()

	return time.Duration(duration / total)
}
