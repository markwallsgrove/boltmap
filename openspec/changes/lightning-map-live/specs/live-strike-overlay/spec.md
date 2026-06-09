## ADDED Requirements

### Requirement: Strike channel bridge
The `blitzortung-client` SHALL expose a Go channel of `Strike` values. The TUI model SHALL read from this channel via `Program.Send` so strikes arrive as Bubbletea messages without blocking the MQTT goroutine.

#### Scenario: Strike delivered to TUI
- **WHEN** the MQTT client receives a strike
- **THEN** it is delivered to the Bubbletea model as a message within 100 ms

#### Scenario: Burst handling
- **WHEN** strikes arrive faster than the TUI can render
- **THEN** the channel buffer absorbs the burst and no MQTT callback is blocked

### Requirement: Strike projection onto viewport
Each strike SHALL be projected onto the current map viewport using the same equirectangular transform as the world map renderer, producing a cell coordinate (col, row). Strikes outside the current viewport SHALL not be rendered.

#### Scenario: Strike within viewport rendered
- **WHEN** a strike's lat/lon falls within the current viewport bounds
- **THEN** a coloured dot appears at the corresponding cell on the next render

#### Scenario: Strike outside viewport not rendered
- **WHEN** a strike's lat/lon falls outside the current viewport bounds
- **THEN** no dot is rendered for that strike

### Requirement: Strike colour by polarity
Strikes with positive polarity (`Pol == 1`) SHALL be rendered in yellow. Strikes with negative polarity (`Pol == -1`) SHALL be rendered in cyan.

#### Scenario: Positive polarity strike colour
- **WHEN** a positive strike is projected onto the map
- **THEN** its cell is coloured yellow

#### Scenario: Negative polarity strike colour
- **WHEN** a negative strike is projected onto the map
- **THEN** its cell is coloured cyan

### Requirement: Age-based strike fade
Strikes SHALL fade from their polarity colour to grey after 10 seconds and be removed from the render buffer after 60 seconds (configurable via `--ttl` flag, default 60).

#### Scenario: Fresh strike in polarity colour
- **WHEN** a strike is less than 10 seconds old
- **THEN** it is rendered in its polarity colour (yellow or cyan)

#### Scenario: Ageing strike turns grey
- **WHEN** a strike is between 10 and 60 seconds old
- **THEN** it is rendered in grey

#### Scenario: Expired strike removed
- **WHEN** a strike exceeds the configured TTL
- **THEN** it is removed from the buffer and no longer rendered

### Requirement: Strike buffer cap
The in-memory strike buffer SHALL hold a maximum of 5 000 strikes. When the cap is reached, the oldest strike SHALL be evicted to make room for the new one.

#### Scenario: Buffer eviction at cap
- **WHEN** the buffer contains 5 000 strikes and a new strike arrives
- **THEN** the oldest strike is evicted and the new strike is added

### Requirement: Live stats bar
The status bar SHALL display the current MQTT connection state, total strikes received, and strikes per minute — all updated on each render tick. The stats bar SHALL use lipgloss colours: connection state green when connected, red when disconnected.

#### Scenario: Stats update on tick
- **WHEN** the 500 ms render tick fires
- **THEN** the stats bar reflects the latest strike count and rate

#### Scenario: Disconnected state shown in red
- **WHEN** the MQTT connection is lost
- **THEN** the connection status in the stats bar is rendered in red

#### Scenario: Connected state shown in green
- **WHEN** the MQTT connection is active
- **THEN** the connection status in the stats bar is rendered in green

### Requirement: Strike re-projection on viewport change
When the user zooms or pans, all strikes in the buffer SHALL be re-projected to the new viewport coordinates before the next render.

#### Scenario: Pan re-projects strikes
- **WHEN** the user presses an arrow key to pan
- **THEN** all buffered strikes are re-projected and appear at their correct positions in the new viewport
