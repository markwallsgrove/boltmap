## Why

The embedded land mask is 360×180 (1° resolution), which produces visibly blocky coastlines — especially when zoomed in. Upgrading to 3600×1800 (0.1° resolution) eliminates the staircase artefacts and makes the map usable at all zoom levels.

## What Changes

- `scripts/genmask.py` rewritten to use `rasterio`/`fiona` for fast scan-line rasterization against Natural Earth 10m land polygons, producing a bit-packed 3600×1800 mask (~810KB).
- `internal/maprender/assets/landmask.bin` regenerated at the new resolution.
- `internal/maprender/landmask.go` updated: constants `maskCols`/`maskRows` raised to 3600/1800, `LandAt` updated to unpack bits from the new format.

## Capabilities

### New Capabilities

### Modified Capabilities

- `world-map-renderer`: Land mask resolution increases from 1°/cell to 0.1°/cell; `LandAt` reads bit-packed data instead of byte-per-cell.

## Impact

- `scripts/genmask.py`: full rewrite; adds `rasterio`, `fiona`, `numpy` as generation-time Python dependencies (not runtime).
- `internal/maprender/assets/landmask.bin`: grows from ~65KB to ~810KB.
- `internal/maprender/landmask.go`: constants and `LandAt` change; no public API changes.
- No changes to renderer, TUI, or MQTT pipeline.
