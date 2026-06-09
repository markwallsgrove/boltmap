## ADDED Requirements

### Requirement: MQTT connection to Blitzortung broker
The client SHALL connect to `blitzortung.ha.sed.pl` on port `1883` using MQTT 3.1.1 with a unique client ID and subscribe to topic `blitzortung/1.1/#`.

#### Scenario: Successful connection and subscription
- **WHEN** the binary is started with network access to `blitzortung.ha.sed.pl:1883`
- **THEN** the client connects, subscribes to `blitzortung/1.1/#`, and begins receiving messages within 5 seconds

#### Scenario: Connection failure on startup
- **WHEN** the broker is unreachable at startup
- **THEN** the binary exits with a non-zero status code and prints a human-readable error to stderr

### Requirement: Strike JSON parsing
The client SHALL parse each MQTT message payload as a JSON object into a typed `Strike` struct containing: `Time` (int64, nanoseconds since Unix epoch), `Lat` (float64), `Lon` (float64), `Alt` (float64), `Pol` (int, -1 or +1), `Mds` (int64), `Scs` (int).

#### Scenario: Valid strike payload
- **WHEN** a well-formed JSON message arrives on the subscribed topic
- **THEN** it is parsed into a `Strike` struct with all fields populated correctly

#### Scenario: Malformed payload
- **WHEN** a message arrives that is not valid JSON or is missing required fields
- **THEN** the message is silently dropped and processing continues

### Requirement: Coloured stdout output
The client SHALL print each received strike to stdout as a single line. Positive polarity strikes (`Pol == 1`) SHALL be printed in yellow (ANSI code `\033[33m`). Negative polarity strikes (`Pol == -1`) SHALL be printed in cyan (ANSI code `\033[36m`). Each line SHALL include timestamp (UTC), latitude, longitude, and polarity symbol (`+`/`-`).

#### Scenario: Positive polarity strike printed in yellow
- **WHEN** a strike with `Pol == 1` is received
- **THEN** a yellow-coloured line is printed to stdout containing the UTC time, lat, lon, and `+`

#### Scenario: Negative polarity strike printed in cyan
- **WHEN** a strike with `Pol == -1` is received
- **THEN** a cyan-coloured line is printed to stdout containing the UTC time, lat, lon, and `-`

### Requirement: Strike rate reporting
The client SHALL print a strike rate summary to stdout every 10 seconds, showing strikes per minute computed over the preceding 60-second sliding window.

#### Scenario: Rate printed periodically
- **WHEN** 10 seconds have elapsed since the last rate report
- **THEN** a line is printed to stdout showing the current strikes/min value

### Requirement: Graceful shutdown
The client SHALL handle `SIGINT` (Ctrl+C) and `SIGTERM` by disconnecting cleanly from the MQTT broker and exiting with status code 0.

#### Scenario: Ctrl+C triggers clean exit
- **WHEN** the user sends SIGINT while the client is running
- **THEN** the MQTT connection is closed cleanly and the process exits with code 0
