## 1. Generation script

- [x] 1.1 Add `rasterio`, `fiona`, and `numpy` pip dependencies to `scripts/genmask.py` header comment (install instructions)
- [x] 1.2 Rewrite `scripts/genmask.py` to download `ne_10m_land.geojson` from Natural Earth, rasterize to 3600×1800 using `rasterio.features.rasterize`, and write a bit-packed binary (MSB-first within each byte) to `internal/maprender/assets/landmask.bin`
- [x] 1.3 Run the script and commit the regenerated `internal/maprender/assets/landmask.bin` (~810KB)

## 2. Go land mask decoder

- [x] 2.1 Update `maskCols` and `maskRows` constants in `internal/maprender/landmask.go` to 3600 and 1800 respectively
- [x] 2.2 Update `LandAt` to decode bit-packed data: compute `idx := row*maskCols + col`, return `landmaskData[idx/8]>>(7-(idx%8))&1 == 1`

## 3. Tests

- [x] 3.1 Add unit tests in `internal/maprender/landmask_test.go` (or extend existing) for: known land coordinate (51.5°N, 0.1°W), known ocean coordinate (0°N, 30°W), and bit-boundary cells (idx%8==0 and idx%8==7)
- [x] 3.2 Run `go test ./internal/maprender/...` and confirm all tests pass

## 4. Verification

- [x] 4.1 Run `go build ./...` with no errors
- [x] 4.2 Start the TUI (`./boltmap --tui`) and visually confirm coastlines are smooth at full-world zoom
- [x] 4.3 Zoom in to max (`+` six times) over a coastal region and confirm no 1°-block staircase artefacts
