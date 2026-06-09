## Context

PR 2 in a three-part series. The MQTT data pipeline exists (PR 1); this PR adds the visual TUI shell and a static world map renderer. No live data is wired yet — the goal is a solid, resizable canvas that PR 3 can overlay strikes onto.

## Goals / Non-Goals

**Goals:**
- Bubbletea TUI that fills the terminal with a coloured world map.
- Correct equirectangular projection that maps lat/lon to character-cell coordinates at any terminal size.
- Zoom and pan keyboard controls that work without live data.
- Status bar with placeholder text, styled with lipgloss.
- Map redraws correctly on terminal resize.

**Non-Goals:**
- Live strike data (PR 3).
- Sub-cell precision rendering (one strike per cell is sufficient).
- Scrollable history or recording.

## Decisions

**TUI framework: `charmbracelet/bubbletea`**
Bubbletea's Elm-architecture model (Msg/Cmd/Update/View) fits this use case cleanly: the map is pure state → view, and resize/key events are standard messages. Alternative (`tcell` directly) would require more boilerplate for layout and styling. `tview` was rejected because it is widget-tree based and harder to do pixel-precise canvas rendering with.

**Styling: `charmbracelet/lipgloss`**
Lipgloss provides terminal-width-aware layout (status bar clamped to terminal width) and consistent ANSI colour handling across platforms. It pairs naturally with Bubbletea.

**Map data: embedded low-resolution coastline raster**
A pre-rasterised ASCII/binary world map at ~360×180 resolution (1°/cell) embedded via `go:embed` is the simplest approach. Alternatives:
- GeoJSON coastlines + runtime rasterisation: accurate but complex; adds a geometry dependency.
- `ericm/go-worldmap` or similar: introduces a package dependency for a one-time asset.
  
A compact binary blob (~6 KB) of land/ocean booleans at 1° resolution is generated once offline and committed to the repo. At render time it is scaled to the terminal size bilinearly.

**Projection: equirectangular**
Maps `lon` → x and `lat` → y linearly. Simple, fast, no special-case poles. Sufficient for a lightning map where geographic accuracy at the coastline level is not critical.

**Colour scheme**
- Land: bright green (`#00FF00` / lipgloss `Color("10")`) on first pass; can be toned down in PR 3.
- Ocean: dark blue (`Color("17")`).
- Status bar background: `Color("235")` (dark grey), foreground white.

**Zoom/pan state**
`model.viewport` holds `centerLat`, `centerLon`, and `zoom` (float64 scale multiplier). Pan moves center by a fixed fraction of the visible range per keypress. Zoom doubles/halves the visible range.

## Risks / Trade-offs

[Braille vs block characters] Braille cells give 2×4 subpixels per cell (higher resolution) but require font support. Block characters (`█`, `▄`, `▀`) are universally supported but lower resolution. → Default to block characters; braille can be a `--braille` flag.

[Terminal size minimum] Very small terminals (<40 cols or <20 rows) will produce a distorted map. → Guard with a minimum size check and print a "terminal too small" message.

[Map asset accuracy] A 1°-resolution land mask has visible pixelation at high zoom. → Acceptable for a lightning strike overlay tool; high-zoom coastline fidelity is not a goal.

## Open Questions

- Which specific embedded map asset to use? Will generate from Natural Earth 110m land polygons rasterised to 360×180 offline.
