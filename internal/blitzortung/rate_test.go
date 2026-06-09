package blitzortung

import (
	"testing"
	"time"
)

func TestRateCounter_Empty(t *testing.T) {
	r := &RateCounter{}
	if got := r.Rate(); got != 0 {
		t.Errorf("empty window: got %f, want 0", got)
	}
}

func TestRateCounter_SteadyRate(t *testing.T) {
	r := &RateCounter{}
	for i := 0; i < 10; i++ {
		r.Add()
	}
	if got := r.Rate(); got != 10 {
		t.Errorf("steady rate: got %f, want 10", got)
	}
}

func TestRateCounter_WindowExpiry(t *testing.T) {
	r := &RateCounter{}

	old := time.Now().Add(-61 * time.Second)
	r.mu.Lock()
	r.timestamps = append(r.timestamps, old, old)
	r.mu.Unlock()

	r.Add()

	if got := r.Rate(); got != 1 {
		t.Errorf("after expiry: got %f, want 1", got)
	}
}
