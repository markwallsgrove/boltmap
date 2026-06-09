package maprender

import (
	"testing"
)

func TestLandAt_KnownLand(t *testing.T) {
	// Central London: should be land.
	if !LandAt(51.5, -0.1) {
		t.Error("expected LandAt(51.5, -0.1) == true (London)")
	}
}

func TestLandAt_KnownOcean(t *testing.T) {
	// Mid-Atlantic: should be ocean.
	if LandAt(0, -30) {
		t.Error("expected LandAt(0, -30) == false (mid-Atlantic)")
	}
}

// TestLandAt_BitBoundaries verifies the bit-unpack logic at byte-boundary indices.
//
// Cell (row=0, col=0) → idx=0 → MSB of byte 0 (idx%8 == 0).
// Cell (row=0, col=7) → idx=7 → LSB of byte 0 (idx%8 == 7).
// Both are in the Arctic Ocean and must return false.
func TestLandAt_BitBoundaries(t *testing.T) {
	// idx=0: lat≈89.95°, lon≈-179.95° (Arctic, MSB of first byte)
	if LandAt(89.95, -179.95) {
		t.Error("expected LandAt(89.95, -179.95) == false (Arctic Ocean, idx%8==0)")
	}
	// idx=7: lat≈89.95°, lon≈-179.25° (Arctic, LSB of first byte)
	if LandAt(89.95, -179.25) {
		t.Error("expected LandAt(89.95, -179.25) == false (Arctic Ocean, idx%8==7)")
	}
	// idx=8: lat≈89.95°, lon≈-179.15° (Arctic, MSB of second byte)
	if LandAt(89.95, -179.15) {
		t.Error("expected LandAt(89.95, -179.15) == false (Arctic Ocean, idx%8==0 second byte)")
	}
}

func TestLandAt_Clamping(t *testing.T) {
	// Coordinates at/beyond extremes should not panic.
	LandAt(90, 180)
	LandAt(-90, -180)
	LandAt(91, 181)
	LandAt(-91, -181)
}
