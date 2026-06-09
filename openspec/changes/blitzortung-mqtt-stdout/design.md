## Context

This is the first PR in a three-part series. It establishes the Go module, MQTT data pipeline, and `Strike` data model that subsequent PRs depend on. There is no existing code — the repo is empty.

## Goals / Non-Goals

**Goals:**
- Working Go binary that connects to Blitzortung's public MQTT broker and streams strike data to stdout.
- Clean `Strike` struct and parser reusable by later PRs.
- Coloured output (ANSI) to confirm polarity visually.
- Graceful shutdown and MQTT disconnect on `Ctrl+C`.

**Non-Goals:**
- Any TUI or map rendering (PR 2/3).
- Persistence or historical data.
- Reconnect retry with backoff (can be added later; initial connect failure exits).

## Decisions

**Module layout: flat `cmd/boltmap` + `internal/` packages**
A `cmd/boltmap/main.go` entry point with `internal/blitzortung` (client + model) keeps the first PR small and makes the packages importable by later PRs without a public API surface. Alternative (single `main.go` file) was rejected because it would require restructuring when the TUI is added.

**MQTT library: `eclipse/paho.mqtt.golang`**
The Eclipse Paho client is the de-facto standard Go MQTT library, actively maintained, and supports MQTT 3.1.1 which the Blitzortung broker uses. Alternative (`hivemq/hivemq-mqtt-client`) is Java-only.

**Colour: raw ANSI escape codes, not a colour library**
For stdout-only output a thin wrapper over `\033[33m` (yellow) / `\033[36m` (cyan) / `\033[0m` (reset) avoids a dependency. A proper colour library (`fatih/color`, `mgutz/ansi`) will be unnecessary once lipgloss is introduced in PR 2.

**Strike rate: computed over a 10-second sliding window**
A simple ring buffer of arrival timestamps gives a low-overhead strikes/min estimate without goroutine complexity.

**MQTT topic: `blitzortung/1.1/#`**
Subscribes to all geohash sub-topics in one wildcard subscription. The geohash structure (`blitzortung/1.1/{hash1}/{hash2}/...`) is how the broker partitions the global feed; subscribing to `#` receives everything without needing to enumerate regions.

## Risks / Trade-offs

[Broker availability] The public broker `blitzortung.ha.sed.pl` is community-run with no SLA. → Binary exits with a clear error message on connection failure; reconnect logic deferred to a later PR.

[High message rate] Global strike rate can reach hundreds per second during storm seasons. → Stdout is fast enough; rate-limiting is not needed for this PR.

[ANSI on Windows] Raw ANSI codes don't render on older Windows terminals. → Acceptable for now; lipgloss (PR 2) handles this transparently.

## Open Questions

- Should the `Strike` struct expose `sig` (detector array) fields, or strip them? Likely omit for now — they add size but aren't used until/unless a detector map view is added.
