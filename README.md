# boltmap

Streams real-time lightning strike data from the [Blitzortung](https://www.blitzortung.org) public MQTT broker to stdout, with ANSI-coloured output and a strikes/min rate counter.

## Requirements

- Go 1.26+
- Outbound TCP access to `blitzortung.ha.sed.pl:1883`

## Build

```sh
go build -o boltmap ./cmd/boltmap
```

## Run

```sh
./boltmap
```

Each strike is printed as a single coloured line:

```
2026-06-09T15:17:07Z 36.1601 6.1665 -
2026-06-09T15:17:08Z 51.5624 4.7033 +
```

- **Yellow** — positive polarity (`+`)
- **Cyan** — negative polarity (`-`)

Every 10 seconds a rate summary is printed:

```
rate: 142.0 strikes/min
```

Stop with `Ctrl+C` — the MQTT connection is closed cleanly before exit.

## Test

```sh
go test ./...
```

## Project layout

```
cmd/boltmap/       # main entrypoint and printer
internal/blitzortung/  # MQTT client, Strike model, rate counter
```

## Notes

- Connects to `blitzortung/1.1/#` (global feed via wildcard).
- Connection failure at startup exits with code 1 and a message on stderr.
- Blitzortung usage policy: single-client personal use only; do not proxy or relay the stream.
