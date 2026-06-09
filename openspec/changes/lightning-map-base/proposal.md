## Why

With the data pipeline established in PR 1, a visual canvas is needed before strike data can be overlaid. This PR introduces the terminal world map and TUI shell as a standalone, testable component — no live data yet.

## What Changes

- Bubbletea TUI application shell with lipgloss layout: full-screen map pane + status bar.
- Equirectangular world map rendered using Unicode braille cells or block characters, dynamically sized to terminal dimensions.
- Map drawn in colour: land mass outlines in green, ocean background in dark blue/grey.
- Status bar rendered with lipgloss showing placeholder connection state and strike count.
- Keyboard controls: `q`/`Ctrl+C` to quit, `+`/`-` to zoom, arrow keys to pan.
- Map redraws correctly on terminal resize.

## Capabilities

### New Capabilities

- `world-map-renderer`: Equirectangular world map projection onto braille/block character cells, with zoom, pan, and terminal-resize handling.
- `tui-app`: Bubbletea application shell with full-screen map pane and lipgloss-styled status bar.

### Modified Capabilities

## Impact

- Builds on the Go module from `blitzortung-mqtt-stdout`; no data wiring yet.
- New dependencies: `charmbracelet/bubbletea`, `charmbracelet/lipgloss`.
- A world map data source (coastline coordinates) is needed — will use a compact embedded GeoJSON or pre-rasterised ASCII map asset.
