#!/usr/bin/env python3
"""Generate 360x180 binary land mask from Natural Earth 110m land polygons.

Output: assets/landmask.bin
Format: 360*180 bytes, one byte per degree cell (1=land, 0=ocean)
Row-major: lat 90→-90 (rows), lon -180→180 (cols)
"""

import json
import struct
import sys
import urllib.request


URL = "https://raw.githubusercontent.com/nvkelso/natural-earth-vector/master/geojson/ne_110m_land.geojson"
COLS = 360
ROWS = 180


def point_in_ring(px: float, py: float, ring: list) -> bool:
    inside = False
    n = len(ring)
    j = n - 1
    for i in range(n):
        xi, yi = ring[i]
        xj, yj = ring[j]
        if ((yi > py) != (yj > py)) and (px < (xj - xi) * (py - yi) / (yj - yi) + xi):
            inside = not inside
        j = i
    return inside


def point_in_polygon(lon: float, lat: float, geometry: dict) -> bool:
    gtype = geometry["type"]
    coords = geometry["coordinates"]

    if gtype == "Polygon":
        polys = [coords]
    elif gtype == "MultiPolygon":
        polys = coords
    else:
        return False

    for poly in polys:
        exterior = poly[0]
        if point_in_ring(lon, lat, exterior):
            # check holes
            inside_hole = any(point_in_ring(lon, lat, hole) for hole in poly[1:])
            if not inside_hole:
                return True
    return False


def main():
    print(f"Downloading {URL} ...", file=sys.stderr)
    with urllib.request.urlopen(URL) as resp:
        data = json.load(resp)

    features = data["features"]
    geometries = [f["geometry"] for f in features]
    print(f"Loaded {len(geometries)} land polygons", file=sys.stderr)

    mask = bytearray(ROWS * COLS)
    for row in range(ROWS):
        lat = 89.5 - row  # 89.5 down to -89.5
        for col in range(COLS):
            lon = -179.5 + col  # -179.5 up to 179.5
            for geom in geometries:
                if point_in_polygon(lon, lat, geom):
                    mask[row * COLS + col] = 1
                    break
        if row % 18 == 0:
            pct = int(row / ROWS * 100)
            print(f"  {pct}% (row {row}/{ROWS})", file=sys.stderr)

    out = "assets/landmask.bin"
    with open(out, "wb") as f:
        f.write(bytes(mask))
    land = sum(mask)
    print(f"Written {len(mask)} bytes to {out} ({land} land cells, {len(mask)-land} ocean)", file=sys.stderr)


if __name__ == "__main__":
    main()
