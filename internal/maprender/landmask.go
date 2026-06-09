package maprender

import _ "embed"

//go:embed assets/landmask.bin
var landmaskData []byte

const maskCols = 360
const maskRows = 180

// LandAt returns true if the given coordinate is over land.
func LandAt(lat, lon float64) bool {
	row := int((90.0 - lat))
	col := int((lon + 180.0))
	if row < 0 {
		row = 0
	} else if row >= maskRows {
		row = maskRows - 1
	}
	if col < 0 {
		col = 0
	} else if col >= maskCols {
		col = maskCols - 1
	}
	return landmaskData[row*maskCols+col] == 1
}
