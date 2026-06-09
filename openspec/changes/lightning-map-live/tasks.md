## 1. Client channel exposure

- [x] 1.1 Verify `internal/blitzortung.Client` already exposes a `Strikes() <-chan Strike` channel (from PR 1); if not, add it
- [x] 1.2 Verify the channel is buffered with capacity 256 to absorb bursts without blocking MQTT callbacks

## 2. Strike buffer

- [x] 2.1 Define `PlottedStrike` struct in `internal/tui/buffer.go`: `Strike`, `ArrivalTime time.Time`, `Col int`, `Row int`
- [x] 2.2 Implement `StrikeBuffer` with a fixed-capacity ring of 5 000 `PlottedStrike` values
- [x] 2.3 Implement `Add(s Strike, col, row int)` that appends and evicts the oldest when at cap
- [x] 2.4 Implement `Active(now time.Time, ttl time.Duration) []PlottedStrike` returning only non-expired strikes
- [x] 2.5 Implement `Reproject(vp Viewport, cols, rows int)` that recomputes `Col`/`Row` for all buffered strikes
- [x] 2.6 Write unit tests for ring eviction and `Active` expiry logic

## 3. Strike Bubbletea message

- [x] 3.1 Define `StrikeMsg` type (wrapping `Strike`) in `internal/tui/messages.go`
- [x] 3.2 Implement `waitForStrike(ch <-chan blitzortung.Strike) tea.Cmd` that blocks on the channel and returns a `StrikeMsg`

## 4. Tick-driven redraw

- [x] 4.1 Add a 500 ms tick command (`tea.Tick`) to `Init()` in the TUI model
- [x] 4.2 Handle `tickMsg` in `Update`: increment tick counter, schedule next tick via returned `Cmd`

## 5. TUI model wiring

- [x] 5.1 Add `client *blitzortung.Client`, `buffer StrikeBuffer`, `rate *blitzortung.RateCounter`, and `connected bool` fields to the TUI `model`
- [x] 5.2 In `Init()`, start the MQTT client and return `waitForStrike(client.Strikes())` as the initial command
- [x] 5.3 In `Update`, handle `StrikeMsg`: project the strike onto the viewport, call `buffer.Add`, call `rate.Add`, schedule the next `waitForStrike` cmd
- [x] 5.4 In `Update`, handle connection-state messages (connected/disconnected) by setting `model.connected`

## 6. Strike rendering in View

- [x] 6.1 In `View()`, after rendering the base map, overlay each `PlottedStrike` from `buffer.Active(now, ttl)` onto the map string at the correct cell position
- [x] 6.2 Apply yellow (`lipgloss.Color("11")`) for `Pol == 1` strikes < 10 s old
- [x] 6.3 Apply cyan (`lipgloss.Color("14")`) for `Pol == -1` strikes < 10 s old
- [x] 6.4 Apply grey (`lipgloss.Color("240")`) for strikes 10–60 s old regardless of polarity
- [x] 6.5 On zoom/pan, call `buffer.Reproject(vp, cols, rows)` before returning the updated model

## 7. Live stats bar

- [x] 7.1 Replace placeholder status bar text with live values: MQTT status, total strikes, and strikes/min from `rate.Rate()`
- [x] 7.2 Style connection status green (`lipgloss.Color("10")`) when connected, red (`lipgloss.Color("9")`) when disconnected

## 8. TTL flag

- [x] 8.1 Add `--ttl` CLI flag (default `60s`) parsed as `time.Duration`; pass to `buffer.Active` and fade threshold

## 9. Verification

- [x] 9.1 Run `go build ./...` with no errors
- [x] 9.2 Run `go test ./...` with all tests passing
- [ ] 9.3 Start the live TUI and confirm strikes appear on the map as coloured dots
- [ ] 9.4 Confirm strikes fade to grey after 10 s and disappear after 60 s
- [ ] 9.5 Pan and zoom while strikes are visible; confirm they reproject correctly
- [ ] 9.6 Confirm the stats bar updates every 500 ms with live rate and count
