## Why

A foundation is needed for ingesting real-time Blitzortung lightning strike data before any visualisation can be built. This PR establishes the Go module, MQTT client, and data model so subsequent PRs can build on a working, tested data pipeline.

## What Changes

- New Go module scaffold (`go.mod`, `cmd/boltmap/main.go`).
- MQTT client connecting to `blitzortung.ha.sed.pl:1883`, subscribing to `blitzortung/1.1/#`.
- `Strike` struct parsed from JSON payloads: `time` (nanoseconds since Unix epoch), `lat`, `lon`, `alt`, `pol` (polarity: `+1`/`-1`), `mds`, `scs`.
- Coloured stdout output per strike: positive polarity printed in yellow, negative in cyan, using ANSI escape codes.
- Strike rate counter printed every 10 seconds (strikes/min).
- Graceful shutdown on `Ctrl+C` with clean MQTT disconnect.

## Capabilities

### New Capabilities

- `blitzortung-client`: MQTT connection to the public Blitzortung broker, wildcard topic subscription, JSON parsing into typed `Strike` structs, and reconnect handling.

### Modified Capabilities

## Impact

- New Go module — no existing code affected.
- Dependencies: `eclipse/paho.mqtt.golang` for MQTT.
- Requires outbound TCP to `blitzortung.ha.sed.pl:1883`.
- Blitzortung usage policy: single-client personal use is permitted; the app must not proxy or relay the stream.
