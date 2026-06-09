package maprender

import "testing"

func TestLatLonToCell_FullWorld_Center(t *testing.T) {
	vp := DefaultViewport()
	col, row, ok := LatLonToCell(vp, 0, 0, 360, 180)
	if !ok {
		t.Fatal("expected ok for center of world")
	}
	if col != 180 {
		t.Errorf("col: got %d, want 180", col)
	}
	if row != 90 {
		t.Errorf("row: got %d, want 90", row)
	}
}

func TestLatLonToCell_FullWorld_TopLeft(t *testing.T) {
	vp := DefaultViewport()
	col, row, ok := LatLonToCell(vp, 90, -180, 360, 180)
	if !ok {
		t.Fatal("expected ok for top-left corner")
	}
	if col != 0 {
		t.Errorf("col: got %d, want 0", col)
	}
	if row != 0 {
		t.Errorf("row: got %d, want 0", row)
	}
}

func TestLatLonToCell_FullWorld_BottomRight(t *testing.T) {
	vp := DefaultViewport()
	col, row, ok := LatLonToCell(vp, -90, 180, 360, 180)
	if !ok {
		t.Fatal("expected ok for bottom-right corner")
	}
	if col != 359 {
		t.Errorf("col: got %d, want 359", col)
	}
	if row != 179 {
		t.Errorf("row: got %d, want 179", row)
	}
}

func TestLatLonToCell_OutsideViewport(t *testing.T) {
	vp := Viewport{CenterLat: 0, CenterLon: 0, Zoom: 0.5}
	_, _, ok := LatLonToCell(vp, 80, 170, 100, 50)
	if ok {
		t.Error("expected ok=false for point outside zoomed viewport")
	}
}

func TestZoomIn_ClampsToMin(t *testing.T) {
	vp := Viewport{Zoom: minZoom}
	vp.ZoomIn()
	if vp.Zoom != minZoom {
		t.Errorf("ZoomIn past min: got %f, want %f", vp.Zoom, minZoom)
	}
}

func TestZoomOut_ClampsToMax(t *testing.T) {
	vp := Viewport{Zoom: maxZoom}
	vp.ZoomOut()
	if vp.Zoom != maxZoom {
		t.Errorf("ZoomOut past max: got %f, want %f", vp.Zoom, maxZoom)
	}
}

func TestPan_LongitudeWraps(t *testing.T) {
	vp := Viewport{CenterLat: 0, CenterLon: 170, Zoom: 1.0}
	vp.Pan(0, 20)
	if vp.CenterLon != -170 {
		t.Errorf("longitude wrap: got %f, want -170", vp.CenterLon)
	}
}

func TestPan_LatitudeClamped(t *testing.T) {
	vp := Viewport{CenterLat: 85, CenterLon: 0, Zoom: 1.0}
	vp.Pan(10, 0)
	if vp.CenterLat != 90 {
		t.Errorf("latitude clamp: got %f, want 90", vp.CenterLat)
	}
}
