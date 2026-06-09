## 1. Module scaffold

- [ ] 1.1 Initialise Go module (`go mod init github.com/markwallsgrove/boltmap`)
- [ ] 1.2 Create directory structure: `cmd/boltmap/`, `internal/blitzortung/`
- [ ] 1.3 Add `eclipse/paho.mqtt.golang` dependency (`go get`)

## 2. Strike model

- [ ] 2.1 Define `Strike` struct in `internal/blitzortung/strike.go` with fields: `Time int64`, `Lat float64`, `Lon float64`, `Alt float64`, `Pol int`, `Mds int64`, `Scs int`
- [ ] 2.2 Implement `ParseStrike(payload []byte) (Strike, error)` using `encoding/json`
- [ ] 2.3 Write unit tests for `ParseStrike` covering valid payload, missing fields, and malformed JSON

## 3. MQTT client

- [ ] 3.1 Implement `Client` struct in `internal/blitzortung/client.go` with `Connect(broker string) error` and `Strikes() <-chan Strike`
- [ ] 3.2 Connect to `blitzortung.ha.sed.pl:1883` with a generated unique client ID
- [ ] 3.3 Subscribe to `blitzortung/1.1/#` with QoS 0; parse each message via `ParseStrike` and send to the strikes channel (buffered, cap 256)
- [ ] 3.4 Implement `Disconnect()` for clean shutdown
- [ ] 3.5 Return a clear error from `Connect` if the broker is unreachable at startup

## 4. Rate counter

- [ ] 4.1 Implement `RateCounter` in `internal/blitzortung/rate.go` using a 60-second sliding window of arrival timestamps
- [ ] 4.2 Expose `Add()` and `Rate() float64` (strikes/min) methods
- [ ] 4.3 Write unit tests for `RateCounter` covering empty window, steady rate, and window expiry

## 5. Coloured stdout output

- [ ] 5.1 Implement `PrintStrike(s Strike)` in `cmd/boltmap/printer.go` using raw ANSI codes: yellow (`\033[33m`) for `Pol == 1`, cyan (`\033[36m`) for `Pol == -1`
- [ ] 5.2 Format each line as: `<UTC time> <lat> <lon> <+/->` followed by reset code `\033[0m`

## 6. Main entrypoint and shutdown

- [ ] 6.1 Write `cmd/boltmap/main.go`: connect client, start reading strikes channel, call `PrintStrike` and `RateCounter.Add` for each strike
- [ ] 6.2 Print rate summary every 10 seconds using a `time.Ticker`
- [ ] 6.3 Handle `SIGINT`/`SIGTERM` with `signal.NotifyContext`; call `client.Disconnect()` and exit 0
- [ ] 6.4 Exit with code 1 and stderr message if `Connect` returns an error

## 7. Verification

- [ ] 7.1 Run `go build ./...` with no errors
- [ ] 7.2 Run `go test ./...` with all tests passing
- [ ] 7.3 Run the binary for 30 seconds and confirm coloured strike lines appear on stdout
