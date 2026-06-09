package maprender

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	landStyle  = lipgloss.NewStyle().Background(lipgloss.Color("10")).Foreground(lipgloss.Color("10"))
	oceanStyle = lipgloss.NewStyle().Background(lipgloss.Color("17")).Foreground(lipgloss.Color("17"))
)

const cellChar = " "

// Render returns a lipgloss-styled string representing the world map.
// Returns a "Terminal too small" message when cols < 40 or rows < 20.
func Render(vp Viewport, cols, rows int) string {
	if cols < 40 || rows < 20 {
		return "Terminal too small"
	}

	var sb strings.Builder
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			lat, lon := cellToLatLon(vp, c, r, cols, rows)
			if LandAt(lat, lon) {
				sb.WriteString(landStyle.Render(cellChar))
			} else {
				sb.WriteString(oceanStyle.Render(cellChar))
			}
		}
		if r < rows-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

// cellToLatLon is the inverse of LatLonToCell — maps a cell centre to lat/lon.
func cellToLatLon(vp Viewport, col, row, cols, rows int) (lat, lon float64) {
	halfLat := 90.0 * vp.Zoom
	halfLon := 180.0 * vp.Zoom

	minLat := vp.CenterLat - halfLat
	maxLat := vp.CenterLat + halfLat
	minLon := vp.CenterLon - halfLon
	maxLon := vp.CenterLon + halfLon

	lon = minLon + (float64(col)+0.5)/float64(cols)*(maxLon-minLon)
	lat = maxLat - (float64(row)+0.5)/float64(rows)*(maxLat-minLat)
	return lat, lon
}
