## Context

The embedded land mask is a 360×180 byte array (1° per cell, ~65KB). At max zoom (1/16 world, ~22°×45° visible), each 1° mask cell maps to 4–5 terminal columns, producing obvious staircase artefacts on coastlines. Upgrading to 3600×1800 (0.1°/cell) eliminates this at all supported zoom levels.

The current generation script (`scripts/genmask.py`) uses a per-cell point-in-polygon test: O(cells × polygons). At 64,800 cells this is tolerable (~minutes); at 6,480,000 cells it would take hours.

## Goals / Non-Goals

**Goals:**
- Land mask at 3600×1800 (0.1°/cell) using Natural Earth 10m polygons.
- Bit-packed binary asset (~810KB) to keep the embedded binary size acceptable.
- Fast generation script (seconds, not hours) using rasterio scan-line rasterization.
- No changes to public API (`LandAt`, `Render`, `RenderWithOverlay`).

**Non-Goals:**
- Braille/half-block sub-cell rendering.
- Runtime vector rasterization or dynamic resolution.
- Reducing binary size beyond bit-packing (e.g. zlib compression would require decompression at startup).

## Decisions

### Bit-packing the mask

**Decision**: Store 1 bit per cell, not 1 byte.

Raw 3600×1800 = 6,480,000 bytes (6.2MB). Bit-packed = 810,000 bytes (~810KB). Embedding 6.2MB in a terminal binary is wasteful; 810KB is reasonable.

`LandAt` cost is unchanged: one array index + one bit shift. No measurable hot-path impact.

**Alternative considered**: zlib-compress the raw mask at embed time, decompress on `init()`. Rejected — adds startup latency and complexity for ~2× extra savings over bit-packing.

### rasterio + fiona for generation

**Decision**: Replace the per-cell point-in-polygon loop with `rasterio.features.rasterize`, which uses scan-line fill internally.

At 3600×1800 this runs in seconds. The script remains a one-off offline tool; the output is committed as a binary asset.

**Alternative considered**: `gdal_rasterize` CLI. Works, but `rasterio` is easier to drive from Python and produces output we can bit-pack directly in the same script.

### Natural Earth 10m data

**Decision**: Use `ne_10m_land` polygons (vs current `ne_110m_land`).

At 0.1°/cell, 110m generalisation omits islands and blurs narrow peninsulas that are clearly resolvable. 10m polygons are large (~20MB GeoJSON) but only used during generation; the script downloads and discards them.

**Alternative considered**: 50m polygons. Better than 110m but still visibly misses small islands at max zoom.

### Bit layout

Row-major, same as current: row 0 = lat 89.95°→89.85°…, col 0 = lon -179.95°. Within each byte, MSB = lower index (bit 7 = cell 0 within the byte). This matches the natural left-to-right reading order and keeps the decode mask simple (`>> (7 - idx%8) & 1`).

## Risks / Trade-offs

- **Generation environment**: `rasterio` and `fiona` require compiled native extensions (GDAL). If a developer doesn't have them, the script won't run. Mitigation: document the pip install command in the script header; the pre-generated `landmask.bin` is committed so re-generation is only needed if the data source changes.
- **Asset size growth**: ~65KB → ~810KB (+745KB). Acceptable for a terminal binary, but worth noting if the project is later distributed with size constraints.
- **Test coverage of `LandAt` with new format**: The unit tests already validate known land/ocean coordinates; they cover the decode logic implicitly. Add an explicit test for bit-boundary coordinates (idx % 8 == 0 and == 7).
