## Why

With the MQTT data pipeline (PR 1) and the world map TUI (PR 2) both working independently, this PR wires them together to deliver the complete live lightning map experience.

## What Changes

- MQTT client from PR 1 integrated into the Bubbletea event loop via a channel-based message bridge.
- Each incoming `Strike` projected onto the map using the current viewport transform (lat/lon → cell coordinates).
- Strikes rendered as coloured dots: yellow for positive polarity (`+1`), cyan for negative (`-1`).
- Strikes fade to grey over a configurable TTL (default 60 s), then are removed from the render buffer.
- Stats bar updated live: strikes/min rate, total strike count, active MQTT broker — all styled with lipgloss colours.
- Maximum in-memory strike buffer capped (default 5 000) to bound memory use.

## Capabilities

### New Capabilities

- `live-strike-overlay`: Bridges the MQTT strike stream into the TUI render loop; projects strikes onto the map viewport; manages the age-based fade and eviction buffer.

### Modified Capabilities

- `blitzortung-client`: Expose strikes via a Go channel in addition to stdout (required for TUI integration).
- `tui-app`: Accept a strike channel and tick-driven redraw to animate fading strikes.
- `world-map-renderer`: Accept a list of positioned, coloured strike dots to overlay on each render pass.

## Impact

- Completes the Go module started in PR 1 and extended in PR 2.
- No new external dependencies beyond those introduced in PRs 1 and 2.
- Final binary is the deliverable: `boltmap` — a self-contained live lightning TUI.
