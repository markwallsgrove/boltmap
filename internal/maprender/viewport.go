package maprender

// Viewport describes the visible region of the map.
// Zoom 1.0 = full world visible; 0.5 = half the world, etc.
type Viewport struct {
	CenterLat float64
	CenterLon float64
	Zoom      float64
}

// DefaultViewport returns a full-world centered viewport.
func DefaultViewport() Viewport {
	return Viewport{CenterLat: 0, CenterLon: 0, Zoom: 1.0}
}

const (
	minZoom = 1.0 / 16.0
	maxZoom = 1.0
)

// ZoomIn halves the visible range (doubles detail). Clamped to minZoom.
func (v *Viewport) ZoomIn() {
	v.Zoom /= 2.0
	if v.Zoom < minZoom {
		v.Zoom = minZoom
	}
}

// ZoomOut doubles the visible range. Clamped to full world.
func (v *Viewport) ZoomOut() {
	v.Zoom *= 2.0
	if v.Zoom > maxZoom {
		v.Zoom = maxZoom
	}
}

// Pan moves the viewport centre. Longitude wraps at ±180°.
func (v *Viewport) Pan(dLat, dLon float64) {
	v.CenterLat += dLat
	if v.CenterLat > 90 {
		v.CenterLat = 90
	} else if v.CenterLat < -90 {
		v.CenterLat = -90
	}

	v.CenterLon += dLon
	for v.CenterLon > 180 {
		v.CenterLon -= 360
	}
	for v.CenterLon < -180 {
		v.CenterLon += 360
	}
}

// LatLonToCell maps a lat/lon to a terminal cell using equirectangular
// projection within the viewport. Returns ok=false if outside the viewport.
func LatLonToCell(vp Viewport, lat, lon float64, cols, rows int) (col, row int, ok bool) {
	halfLat := 90.0 * vp.Zoom
	halfLon := 180.0 * vp.Zoom

	minLat := vp.CenterLat - halfLat
	maxLat := vp.CenterLat + halfLat
	minLon := vp.CenterLon - halfLon
	maxLon := vp.CenterLon + halfLon

	if lat < minLat || lat > maxLat || lon < minLon || lon > maxLon {
		return 0, 0, false
	}

	col = int((lon - minLon) / (maxLon - minLon) * float64(cols))
	row = int((maxLat - lat) / (maxLat - minLat) * float64(rows))

	if col >= cols {
		col = cols - 1
	}
	if row >= rows {
		row = rows - 1
	}
	return col, row, true
}
