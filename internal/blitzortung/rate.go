package blitzortung

import (
	"sync"
	"time"
)

const windowDuration = 60 * time.Second

type RateCounter struct {
	mu         sync.Mutex
	timestamps []time.Time
}

func (r *RateCounter) Add() {
	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.timestamps = append(r.timestamps, now)
	r.evict(now)
}

func (r *RateCounter) Rate() float64 {
	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.evict(now)
	return float64(len(r.timestamps))
}

func (r *RateCounter) evict(now time.Time) {
	cutoff := now.Add(-windowDuration)
	i := 0
	for i < len(r.timestamps) && r.timestamps[i].Before(cutoff) {
		i++
	}
	r.timestamps = r.timestamps[i:]
}
