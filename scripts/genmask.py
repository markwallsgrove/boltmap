#!/usr/bin/env python3
"""Generate 3600x1800 bit-packed land mask from Natural Earth 10m land polygons.

Output: internal/maprender/assets/landmask.bin
Format: 3600*1800 bits (810000 bytes), MSB-first within each byte.
Row-major: lat 89.95->-89.95 (rows), lon -179.95->179.95 (cols), 0.1 deg/cell.

Preferred: use the Go version (no extra deps, parallel, ~seconds):
    go run scripts/genmask.go

Fast path (rasterio available):
    pip install rasterio fiona numpy
    python3 scripts/genmask.py

Pure-stdlib fallback (no extra packages needed, ~30 min):
    python3 scripts/genmask.py
"""

import json
import sys
import time
import urllib.request

URL = (
    "https://raw.githubusercontent.com/nvkelso/natural-earth-vector"
    "/master/geojson/ne_10m_land.geojson"
)
COLS = 3600
ROWS = 1800
OUT = "internal/maprender/assets/landmask.bin"

# ---------------------------------------------------------------------------
# Rasterio fast path
# ---------------------------------------------------------------------------

def _run_rasterio(geojson_path: str) -> bytearray:
    import numpy as np
    import fiona
    import rasterio
    from rasterio.features import rasterize
    from rasterio.transform import from_bounds

    transform = from_bounds(-180, -90, 180, 90, COLS, ROWS)
    with fiona.open(geojson_path) as src:
        shapes = [(f["geometry"], 1) for f in src]

    raw = rasterize(
        shapes,
        out_shape=(ROWS, COLS),
        transform=transform,
        fill=0,
        dtype=np.uint8,
    )

    # Bit-pack: MSB-first
    flat = raw.flatten()
    packed = bytearray(ROWS * COLS // 8)
    for idx, v in enumerate(flat):
        if v:
            packed[idx // 8] |= 1 << (7 - idx % 8)
    return packed


# ---------------------------------------------------------------------------
# Pure-stdlib fallback
# ---------------------------------------------------------------------------

def _bbox(ring: list) -> tuple:
    xs = [p[0] for p in ring]
    ys = [p[1] for p in ring]
    return min(xs), max(xs), min(ys), max(ys)


def _in_ring(px: float, py: float, ring: list) -> bool:
    inside = False
    j = len(ring) - 1
    for i in range(len(ring)):
        xi, yi = ring[i]
        xj, yj = ring[j]
        if ((yi > py) != (yj > py)) and (px < (xj - xi) * (py - yi) / (yj - yi) + xi):
            inside = not inside
        j = i
    return inside


def _in_poly(lon: float, lat: float, ext: list, holes: list) -> bool:
    return _in_ring(lon, lat, ext) and not any(_in_ring(lon, lat, h) for h in holes)


def _build_index(geometries: list) -> dict:
    """Build a 36x18 spatial grid (10-deg cells) mapping grid_key -> list of rings."""
    index: dict = {}
    for geom in geometries:
        gtype = geom["type"]
        polys = [geom["coordinates"]] if gtype == "Polygon" else geom["coordinates"]
        for poly in polys:
            ext = poly[0]
            holes = poly[1:]
            bx0, bx1, by0, by1 = _bbox(ext)
            # Expand to all 10-deg grid cells the bbox touches
            gx0 = max(int((bx0 + 180) // 10), 0)
            gx1 = min(int((bx1 + 180) // 10), 35)
            gy0 = max(int((90 - by1) // 10), 0)
            gy1 = min(int((90 - by0) // 10), 17)
            for gy in range(gy0, gy1 + 1):
                for gx in range(gx0, gx1 + 1):
                    index.setdefault((gx, gy), []).append((bx0, bx1, by0, by1, ext, holes))
    return index


def _run_stdlib(geojson_bytes: bytes) -> bytearray:
    data = json.loads(geojson_bytes)
    geometries = [f["geometry"] for f in data["features"]]
    print(f"  {len(geometries)} features loaded, building spatial index...", file=sys.stderr)

    index = _build_index(geometries)
    print(f"  index built ({len(index)} non-empty 10-deg grid cells)", file=sys.stderr)

    packed = bytearray(ROWS * COLS // 8)
    t0 = time.time()

    for row in range(ROWS):
        lat = 89.95 - row * 0.1
        gy = min(int((90 - lat) // 10), 17)

        for col in range(COLS):
            lon = -179.95 + col * 0.1
            gx = min(int((lon + 180) // 10), 35)
            rings = index.get((gx, gy), ())

            for bx0, bx1, by0, by1, ext, holes in rings:
                if bx0 <= lon <= bx1 and by0 <= lat <= by1:
                    if _in_poly(lon, lat, ext, holes):
                        idx = row * COLS + col
                        packed[idx // 8] |= 1 << (7 - idx % 8)
                        break

        if row % 180 == 0:
            elapsed = time.time() - t0
            pct = row / ROWS * 100
            eta = (elapsed / max(row, 1)) * (ROWS - row)
            print(
                f"  {pct:5.1f}%  row {row:4d}/{ROWS}  elapsed {elapsed:5.0f}s  eta {eta:5.0f}s",
                file=sys.stderr,
            )

    return packed


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

def main() -> None:
    print(f"Downloading {URL} ...", file=sys.stderr)
    with urllib.request.urlopen(URL) as resp:
        raw_bytes = resp.read()
    print(f"  downloaded {len(raw_bytes):,} bytes", file=sys.stderr)

    packed: bytearray

    try:
        import tempfile, os
        with tempfile.NamedTemporaryFile(suffix=".geojson", delete=False) as tmp:
            tmp.write(raw_bytes)
            tmp_path = tmp.name
        try:
            print("rasterio available — using fast path", file=sys.stderr)
            packed = _run_rasterio(tmp_path)
        finally:
            os.unlink(tmp_path)
    except ImportError:
        print("rasterio not available — using pure-Python path", file=sys.stderr)
        packed = _run_stdlib(raw_bytes)

    with open(OUT, "wb") as f:
        f.write(bytes(packed))

    land = sum(bin(b).count("1") for b in packed)
    total = ROWS * COLS
    print(
        f"Written {len(packed):,} bytes to {OUT}  "
        f"({land:,} land cells, {total - land:,} ocean, {land/total*100:.1f}% land)",
        file=sys.stderr,
    )


if __name__ == "__main__":
    main()
