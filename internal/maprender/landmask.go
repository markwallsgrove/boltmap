package maprender

import _ "embed"

//go:embed assets/landmask.bin
var landmaskData []byte

const maskCols = 3600
const maskRows = 1800

// LandAt returns true if the given coordinate is over land.
func LandAt(lat, lon float64) bool {
	row := int((90.0 - lat) * 10.0)
	col := int((lon + 180.0) * 10.0)
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
	idx := row*maskCols + col
	return landmaskData[idx/8]>>(7-(idx%8))&1 == 1
}
