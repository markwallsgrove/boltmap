//go:build ignore

// genmask generates the 3600×1800 bit-packed land mask asset.
//
// Usage:
//
//	go run scripts/genmask.go
//
// Output: internal/maprender/assets/landmask.bin
// Format: 3600*1800 bits (810 000 bytes), MSB-first within each byte.
// Row-major: lat 90→-90 (rows 0-1799), lon -180→180 (cols 0-3599), 0.1°/cell.
//
// Data source: Natural Earth 10m land polygons (ne_10m_land.geojson, ~10MB).
// No external dependencies — stdlib only.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	cols   = 3600
	rows   = 1800
	url    = "https://raw.githubusercontent.com/nvkelso/natural-earth-vector/master/geojson/ne_10m_land.geojson"
	outDir = "internal/maprender/assets"
	outFile = outDir + "/landmask.bin"
)

// ring is a closed polygon ring stored flat as alternating [lon, lat] pairs.
type ring struct {
	pts          [][2]float64
	minX, maxX   float64
	minY, maxY   float64
}

type polygon struct {
	exterior ring
	holes    []ring
}

func newRing(coords []interface{}) ring {
	r := ring{pts: make([][2]float64, len(coords))}
	r.minX, r.minY = 1e9, 1e9
	r.maxX, r.maxY = -1e9, -1e9
	for i, p := range coords {
		pt := p.([]interface{})
		x := pt[0].(float64)
		y := pt[1].(float64)
		r.pts[i] = [2]float64{x, y}
		if x < r.minX {
			r.minX = x
		}
		if x > r.maxX {
			r.maxX = x
		}
		if y < r.minY {
			r.minY = y
		}
		if y > r.maxY {
			r.maxY = y
		}
	}
	return r
}

func (r ring) contains(px, py float64) bool {
	inside := false
	pts := r.pts
	n := len(pts)
	j := n - 1
	for i := 0; i < n; i++ {
		xi, yi := pts[i][0], pts[i][1]
		xj, yj := pts[j][0], pts[j][1]
		if (yi > py) != (yj > py) {
			if px < (xj-xi)*(py-yi)/(yj-yi)+xi {
				inside = !inside
			}
		}
		j = i
	}
	return inside
}

func (p polygon) contains(lon, lat float64) bool {
	if !p.exterior.contains(lon, lat) {
		return false
	}
	for _, h := range p.holes {
		if h.contains(lon, lat) {
			return false
		}
	}
	return true
}

// spatialIndex maps 10°×10° grid cell → polygons whose exterior bbox overlaps it.
type spatialIndex [36][18][]polygon

func (si *spatialIndex) add(p polygon) {
	ext := p.exterior
	gx0 := int((ext.minX + 180) / 10)
	gx1 := int((ext.maxX + 180) / 10)
	gy0 := int((90 - ext.maxY) / 10)
	gy1 := int((90 - ext.minY) / 10)
	if gx0 < 0 {
		gx0 = 0
	}
	if gx1 > 35 {
		gx1 = 35
	}
	if gy0 < 0 {
		gy0 = 0
	}
	if gy1 > 17 {
		gy1 = 17
	}
	for gy := gy0; gy <= gy1; gy++ {
		for gx := gx0; gx <= gx1; gx++ {
			si[gx][gy] = append(si[gx][gy], p)
		}
	}
}

func (si *spatialIndex) isLand(lon, lat float64) bool {
	gx := int((lon + 180) / 10)
	gy := int((90 - lat) / 10)
	if gx < 0 {
		gx = 0
	} else if gx > 35 {
		gx = 35
	}
	if gy < 0 {
		gy = 0
	} else if gy > 17 {
		gy = 17
	}
	for _, p := range si[gx][gy] {
		ext := p.exterior
		if lon < ext.minX || lon > ext.maxX || lat < ext.minY || lat > ext.maxY {
			continue
		}
		if p.contains(lon, lat) {
			return true
		}
	}
	return false
}

func fetchGeoJSON() ([]byte, error) {
	log.Printf("Downloading %s ...", url)
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var buf []byte
	tmp := make([]byte, 1<<20)
	for {
		n, err := resp.Body.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
		}
		if err != nil {
			break
		}
	}
	log.Printf("  downloaded %d bytes", len(buf))
	return buf, nil
}

func buildIndex(data []byte) (*spatialIndex, error) {
	var geojson struct {
		Features []struct {
			Geometry struct {
				Type        string        `json:"type"`
				Coordinates []interface{} `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}
	if err := json.Unmarshal(data, &geojson); err != nil {
		return nil, err
	}

	idx := &spatialIndex{}
	totalPolys := 0
	for _, feat := range geojson.Features {
		g := feat.Geometry
		var polys [][]interface{}
		switch g.Type {
		case "Polygon":
			polys = [][]interface{}{g.Coordinates}
		case "MultiPolygon":
			for _, p := range g.Coordinates {
				polys = append(polys, p.([]interface{}))
			}
		}
		for _, rawPoly := range polys {
			poly := polygon{}
			for i, rawRing := range rawPoly {
				rr := rawRing.([]interface{})
				r := newRing(rr)
				if i == 0 {
					poly.exterior = r
				} else {
					poly.holes = append(poly.holes, r)
				}
			}
			idx.add(poly)
			totalPolys++
		}
	}
	log.Printf("  %d polygons indexed across %d features", totalPolys, len(geojson.Features))
	return idx, nil
}

func main() {
	data, err := fetchGeoJSON()
	if err != nil {
		log.Fatalf("download: %v", err)
	}

	log.Println("Parsing and building spatial index...")
	idx, err := buildIndex(data)
	if err != nil {
		log.Fatalf("build index: %v", err)
	}

	packed := make([]byte, rows*cols/8)
	nWorkers := runtime.NumCPU()
	log.Printf("Rasterizing %d×%d at %d workers...", cols, rows, nWorkers)

	rowCh := make(chan int, rows)
	for r := 0; r < rows; r++ {
		rowCh <- r
	}
	close(rowCh)

	var mu sync.Mutex
	var wg sync.WaitGroup
	t0 := time.Now()
	done := make([]int32, rows)

	for w := 0; w < nWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rowBuf := make([]byte, cols/8)
			for r := range rowCh {
				lat := 89.95 - float64(r)*0.1
				clear(rowBuf)
				for c := 0; c < cols; c++ {
					lon := -179.95 + float64(c)*0.1
					if idx.isLand(lon, lat) {
						rowBuf[c/8] |= 1 << (7 - c%8)
					}
				}
				mu.Lock()
				copy(packed[r*cols/8:], rowBuf)
				done[r] = 1
				mu.Unlock()
			}
		}()
	}

	// Progress monitor
	go func() {
		for {
			time.Sleep(10 * time.Second)
			mu.Lock()
			n := 0
			for _, v := range done {
				n += int(v)
			}
			mu.Unlock()
			if n >= rows {
				return
			}
			elapsed := time.Since(t0)
			eta := time.Duration(float64(elapsed) / float64(max(n, 1)) * float64(rows-n))
			log.Printf("  %d/%d rows (%.0f%%)  elapsed %s  eta %s",
				n, rows, float64(n)/float64(rows)*100, elapsed.Round(time.Second), eta.Round(time.Second))
		}
	}()

	wg.Wait()

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(outFile, packed, 0o644); err != nil {
		log.Fatalf("write: %v", err)
	}

	land := 0
	for _, b := range packed {
		land += int(b&1) + int((b>>1)&1) + int((b>>2)&1) + int((b>>3)&1) +
			int((b>>4)&1) + int((b>>5)&1) + int((b>>6)&1) + int((b>>7)&1)
	}
	total := rows * cols
	log.Printf("Written %d bytes to %s (%d land, %d ocean, %.1f%% land)",
		len(packed), outFile, land, total-land, float64(land)/float64(total)*100)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
