## ADDED Requirements

### Requirement: Land mask precision
The embedded land mask SHALL have a resolution of 3600×1800 cells (0.1° per cell), sourced from Natural Earth 10m land polygons. The mask SHALL be stored in bit-packed format (1 bit per cell, MSB-first within each byte).

#### Scenario: Sub-degree coastline accuracy
- **WHEN** the viewport is zoomed in to 1/16 of the world
- **THEN** coastline cells form smooth outlines without visible 1°-block staircase artefacts

#### Scenario: Bit boundary cells
- **WHEN** `LandAt` is called for coordinates whose flat index falls on byte boundaries (idx % 8 == 0 or idx % 8 == 7)
- **THEN** it returns the correct land/ocean value without off-by-one bit errors

#### Scenario: Known land coordinate
- **WHEN** `LandAt` is called with a coordinate known to be land (e.g. 51.5°N, 0.1°W — central London)
- **THEN** it returns true

#### Scenario: Known ocean coordinate
- **WHEN** `LandAt` is called with a coordinate known to be ocean (e.g. 0°N, -30°W — mid-Atlantic)
- **THEN** it returns false
