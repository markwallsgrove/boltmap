package tui

import (
	"time"

	"github.com/markwallsgrove/boltmap/internal/blitzortung"
	"github.com/markwallsgrove/boltmap/internal/maprender"
)

const bufferCap = 5000

type PlottedStrike struct {
	Strike      blitzortung.Strike
	ArrivalTime time.Time
	Col         int
	Row         int
}

type StrikeBuffer struct {
	buf  [bufferCap]PlottedStrike
	head int
	size int
}

func (sb *StrikeBuffer) Add(s blitzortung.Strike, col, row int) {
	sb.buf[sb.head] = PlottedStrike{
		Strike:      s,
		ArrivalTime: time.Now(),
		Col:         col,
		Row:         row,
	}
	sb.head = (sb.head + 1) % bufferCap
	if sb.size < bufferCap {
		sb.size++
	}
}

func (sb *StrikeBuffer) Active(now time.Time, ttl time.Duration) []PlottedStrike {
	result := make([]PlottedStrike, 0, sb.size)
	start := (sb.head - sb.size + bufferCap) % bufferCap
	for i := 0; i < sb.size; i++ {
		ps := sb.buf[(start+i)%bufferCap]
		if now.Sub(ps.ArrivalTime) < ttl {
			result = append(result, ps)
		}
	}
	return result
}

func (sb *StrikeBuffer) Reproject(vp maprender.Viewport, cols, rows int) {
	start := (sb.head - sb.size + bufferCap) % bufferCap
	for i := 0; i < sb.size; i++ {
		idx := (start + i) % bufferCap
		ps := &sb.buf[idx]
		col, row, ok := maprender.LatLonToCell(vp, ps.Strike.Lat, ps.Strike.Lon, cols, rows)
		if ok {
			ps.Col = col
			ps.Row = row
		} else {
			ps.Col = -1
			ps.Row = -1
		}
	}
}
