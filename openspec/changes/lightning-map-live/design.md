## Context

PR 3 — the final piece. The MQTT client (`internal/blitzortung`) and the Bubbletea TUI with world map renderer both exist. This PR wires them together so strikes appear on the map in real time with colour and fade.

## Goals / Non-Goals

**Goals:**
- MQTT strike stream fed into the Bubbletea model via a Go channel.
- Each strike projected onto the current viewport and rendered as a coloured dot.
- Age-based fade: yellow/cyan → grey → removed over a configurable TTL (default 60 s).
- Live stats bar: strikes/min, total count, connection state — lipgloss-styled.
- Memory-bounded strike buffer (cap 5 000).

**Non-Goals:**
- Strike clustering or storm-cell grouping.
- Historical playback or data persistence.
- User-configurable colour themes (can be added later).

## Decisions

**Strike delivery: unbuffered channel from MQTT goroutine → Bubbletea `Program.Send`**
Bubbletea's `Program.Send(msg)` is goroutine-safe and the idiomatic way to inject external events. The MQTT message handler calls `program.Send(StrikeMsg{strike})` directly; no intermediate queue needed for this data rate. A buffered channel (capacity 256) is used as a safety valve to avoid blocking the MQTT callback during a burst.

**Strike buffer: slice with circular eviction**
`[]PlottedStrike` capped at 5 000 entries. When full, the oldest entry is overwritten (ring buffer). This bounds memory without a heap allocation per eviction. Alternative (linked list) adds GC pressure for negligible benefit.

**Fade rendering: age buckets**
Strikes are bucketed by age into three visual states:
- `< 10 s`: full colour (yellow positive / cyan negative)
- `10–60 s`: grey (`Color("240")`)
- `> 60 s`: evicted

Age is computed from `strike.Time` (nanoseconds epoch) relative to `time.Now()` on each render pass. No separate expiry goroutine needed.

**Tick-driven redraw: `tea.Tick(500ms)`**
A 500 ms tick message triggers a re-render to animate fades without reacting to every individual strike. Strike arrival also triggers an immediate redraw. This avoids rendering at the full MQTT message rate (potentially hundreds/sec) while keeping the display responsive.

**Stats rate calculation: reuse PR 1's sliding window**
The `blitzortung.RateCounter` from PR 1 is exposed on the client struct and read by the TUI model on each tick to populate the stats bar.

**Projection: same equirectangular transform as PR 2**
`PlottedStrike` stores pre-projected cell coordinates (x, y) computed once on arrival. On viewport change (zoom/pan) all strikes are re-projected. The re-projection cost for 5 000 strikes is negligible (<1 ms).

## Risks / Trade-offs

[High burst rate] During large storms, hundreds of strikes/sec may arrive. At 500 ms ticks, up to 250 strikes may queue between renders. → The buffered channel (cap 256) absorbs bursts; excess messages are dropped with a counter incremented for diagnostics.

[Strike–viewport sync] Re-projecting all strikes on every pan/zoom could cause a brief visual stutter. → With a 5 000-strike cap and simple arithmetic, re-projection takes <1 ms on any modern CPU. No caching needed.

[Time drift] `strike.Time` is a server-assigned nanosecond timestamp; client clock skew could cause strikes to appear already-faded or not-yet-visible. → Use arrival wall-clock time (`time.Now()`) for fade age, not the embedded `strike.Time`. The embedded time is still displayed in the stats bar for reference.

## Open Questions

- Should pan/zoom reset the strike buffer (avoids showing strikes in wrong positions from before the viewport change)? Leaning yes — re-projection handles it automatically since coordinates are recomputed.
