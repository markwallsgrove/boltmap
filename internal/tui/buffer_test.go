package tui

import (
	"testing"
	"time"

	"github.com/markwallsgrove/boltmap/internal/blitzortung"
)

func TestStrikeBuffer_Eviction(t *testing.T) {
	var buf StrikeBuffer
	for i := 0; i < bufferCap+1; i++ {
		buf.Add(blitzortung.Strike{Lat: float64(i)}, 0, 0)
	}
	if buf.size != bufferCap {
		t.Errorf("size after overflow: got %d, want %d", buf.size, bufferCap)
	}
	active := buf.Active(time.Now(), time.Hour)
	if len(active) != bufferCap {
		t.Errorf("active count: got %d, want %d", len(active), bufferCap)
	}
	if active[0].Strike.Lat != 1.0 {
		t.Errorf("oldest after eviction: got %f, want 1.0", active[0].Strike.Lat)
	}
	if active[bufferCap-1].Strike.Lat != float64(bufferCap) {
		t.Errorf("newest: got %f, want %f", active[bufferCap-1].Strike.Lat, float64(bufferCap))
	}
}

func TestStrikeBuffer_Active_Expiry(t *testing.T) {
	var buf StrikeBuffer
	buf.Add(blitzortung.Strike{Lat: 1}, 0, 0)
	buf.Add(blitzortung.Strike{Lat: 2}, 0, 0)
	buf.Add(blitzortung.Strike{Lat: 3}, 0, 0)

	// Backdate the oldest entry past the TTL.
	start := (buf.head - buf.size + bufferCap) % bufferCap
	buf.buf[start].ArrivalTime = time.Now().Add(-2 * time.Minute)

	active := buf.Active(time.Now(), time.Minute)
	if len(active) != 2 {
		t.Errorf("active after expiry: got %d, want 2", len(active))
	}
}

func TestStrikeBuffer_Active_Empty(t *testing.T) {
	var buf StrikeBuffer
	active := buf.Active(time.Now(), time.Hour)
	if len(active) != 0 {
		t.Errorf("empty buffer: got %d active, want 0", len(active))
	}
}

func TestStrikeBuffer_Active_AllExpired(t *testing.T) {
	var buf StrikeBuffer
	buf.Add(blitzortung.Strike{}, 0, 0)
	start := (buf.head - buf.size + bufferCap) % bufferCap
	buf.buf[start].ArrivalTime = time.Now().Add(-2 * time.Minute)

	active := buf.Active(time.Now(), time.Minute)
	if len(active) != 0 {
		t.Errorf("all expired: got %d active, want 0", len(active))
	}
}
