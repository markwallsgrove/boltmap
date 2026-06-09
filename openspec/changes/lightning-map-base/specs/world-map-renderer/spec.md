## ADDED Requirements

### Requirement: Equirectangular world map rendering
The renderer SHALL draw a world map using Unicode block characters, sized dynamically to the current terminal dimensions, using an equirectangular projection where longitude maps linearly to column and latitude maps linearly to row.

#### Scenario: Map fills terminal on startup
- **WHEN** the TUI starts in a terminal of any size >= 40 columns × 20 rows
- **THEN** the map fills the available pane with land cells and ocean cells distinguished by colour

#### Scenario: Terminal too small
- **WHEN** the terminal is smaller than 40 columns or 20 rows
- **THEN** a message "Terminal too small" is displayed instead of the map

### Requirement: Land and ocean colour
Land cells SHALL be rendered in bright green and ocean cells in dark blue.

#### Scenario: Land cell colour
- **WHEN** a cell corresponds to a land coordinate in the embedded map asset
- **THEN** it is rendered with a bright green foreground/background

#### Scenario: Ocean cell colour
- **WHEN** a cell corresponds to an ocean coordinate
- **THEN** it is rendered with a dark blue foreground/background

### Requirement: Terminal resize handling
The renderer SHALL redraw the map at the new dimensions within one render cycle when the terminal is resized.

#### Scenario: Window resize redraws map
- **WHEN** the terminal window is resized
- **THEN** the map redraws to fill the new dimensions without artifacts

### Requirement: Zoom
The user SHALL be able to zoom in and out with `+` and `-` keys. Zoom SHALL double or halve the visible lat/lon range, centred on the current viewport centre.

#### Scenario: Zoom in
- **WHEN** the user presses `+`
- **THEN** the visible lat/lon range halves and map detail increases

#### Scenario: Zoom out
- **WHEN** the user presses `-`
- **THEN** the visible lat/lon range doubles, up to the full world view

### Requirement: Pan
The user SHALL be able to pan the viewport with the arrow keys. Each keypress SHALL move the centre by 10% of the current visible range.

#### Scenario: Pan right
- **WHEN** the user presses the right arrow key
- **THEN** the map centre longitude increases by 10% of the current visible longitude range

#### Scenario: Pan wraps longitude
- **WHEN** panning causes the centre longitude to exceed 180°
- **THEN** it wraps to -180° so the map is continuous
