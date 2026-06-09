## ADDED Requirements

### Requirement: Full-screen TUI layout
The application SHALL render a full-screen Bubbletea TUI with two regions: a map pane occupying all rows except the last, and a single-row status bar at the bottom.

#### Scenario: Layout on startup
- **WHEN** the application starts
- **THEN** the map pane occupies all but the bottom row, and the status bar is visible on the last row

### Requirement: Status bar
The status bar SHALL display placeholder text "Connecting…" for MQTT status and "0 strikes" for the strike count, styled with a dark grey background and white text using lipgloss.

#### Scenario: Status bar visible
- **WHEN** the TUI renders
- **THEN** the status bar shows styled placeholder text spanning the full terminal width

### Requirement: Quit keyboard shortcut
The user SHALL be able to quit the application by pressing `q` or `Ctrl+C`.

#### Scenario: Quit with q
- **WHEN** the user presses `q`
- **THEN** the TUI exits cleanly and the terminal is restored

#### Scenario: Quit with Ctrl+C
- **WHEN** the user presses Ctrl+C
- **THEN** the TUI exits cleanly and the terminal is restored
