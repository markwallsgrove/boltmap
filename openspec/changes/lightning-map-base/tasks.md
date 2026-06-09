## 1. Dependencies

- [ ] 1.1 Add `charmbracelet/bubbletea` and `charmbracelet/lipgloss` to `go.mod` (`go get`)

## 2. Map asset

- [ ] 2.1 Generate a 360×180 binary land mask from Natural Earth 110m land polygons (one bool per degree cell, row-major, lat 90→-90, lon -180→180); commit as `assets/landmask.bin`
- [ ] 2.2 Embed `assets/landmask.bin` in `internal/maprender/landmask.go` using `//go:embed`
- [ ] 2.3 Write `LandAt(lat, lon float64) bool` that looks up the embedded mask

## 3. Viewport and projection

- [ ] 3.1 Define `Viewport` struct in `internal/maprender/viewport.go`: `CenterLat`, `CenterLon float64`, `Zoom float64` (1.0 = full world)
- [ ] 3.2 Implement `LatLonToCell(lat, lon float64, cols, rows int) (col, row int, ok bool)` using equirectangular projection within the viewport
- [ ] 3.3 Implement `ZoomIn()`, `ZoomOut()` (halve/double visible range, min zoom = 1/16 world, max = full world)
- [ ] 3.4 Implement `Pan(dLat, dLon float64)` moving centre; wrap longitude at ±180°
- [ ] 3.5 Write unit tests for projection: known lat/lon maps to expected cell at full-world zoom

## 4. Map renderer

- [ ] 4.1 Implement `Render(vp Viewport, cols, rows int) string` in `internal/maprender/renderer.go` — iterates cells, calls `LandAt`, returns a lipgloss-styled string
- [ ] 4.2 Apply bright green (`lipgloss.Color("10")`) to land cells and dark blue (`lipgloss.Color("17")`) to ocean cells
- [ ] 4.3 Return a "Terminal too small" message string when `cols < 40` or `rows < 20`

## 5. Bubbletea TUI shell

- [ ] 5.1 Define `model` struct in `internal/tui/model.go`: holds `Viewport`, terminal `Width`/`Height`
- [ ] 5.2 Implement `Init()` returning `nil`
- [ ] 5.3 Implement `Update(msg)`: handle `tea.KeyMsg` for `q`/`ctrl+c` (quit), `+`/`-` (zoom), arrow keys (pan); handle `tea.WindowSizeMsg` to update `Width`/`Height`
- [ ] 5.4 Implement `View()`: call `maprender.Render` for the map pane (all rows minus 1); render the status bar on the last row with lipgloss dark grey background and placeholder text "Connecting… | 0 strikes"
- [ ] 5.5 Style status bar with `lipgloss.NewStyle().Background(lipgloss.Color("235")).Foreground(lipgloss.Color("255")).Width(width)`

## 6. Main entrypoint update

- [ ] 6.1 Add a `--tui` flag to `cmd/boltmap/main.go`; when set, start the Bubbletea program instead of stdout printing
- [ ] 6.2 Confirm `q` and `Ctrl+C` exit cleanly and restore the terminal

## 7. Verification

- [ ] 7.1 Run `go build ./...` with no errors
- [ ] 7.2 Run `go test ./...` with all tests passing
- [ ] 7.3 Start the TUI (`./boltmap --tui`), confirm the world map renders in colour and fills the terminal
- [ ] 7.4 Test zoom (`+`/`-`) and pan (arrow keys) work without artefacts
- [ ] 7.5 Resize the terminal window and confirm the map redraws correctly
