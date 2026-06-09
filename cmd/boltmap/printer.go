package main

import (
	"fmt"
	"time"

	"github.com/markwallsgrove/boltmap/internal/blitzortung"
)

const (
	ansiYellow = "\033[33m"
	ansiCyan   = "\033[36m"
	ansiReset  = "\033[0m"
)

func PrintStrike(s blitzortung.Strike) {
	color := ansiCyan
	polSym := "-"
	if s.Pol == 1 {
		color = ansiYellow
		polSym = "+"
	}
	t := time.Unix(0, s.Time).UTC()
	fmt.Printf("%s%s %.4f %.4f %s%s\n", color, t.Format("2006-01-02T15:04:05Z"), s.Lat, s.Lon, polSym, ansiReset)
}
