package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/markwallsgrove/boltmap/internal/blitzortung"
	"github.com/markwallsgrove/boltmap/internal/tui"
)

const defaultBroker = "tcp://blitzortung.ha.sed.pl:1883"

func brokerURL() string {
	if v := os.Getenv("BOLTMAP_BROKER"); v != "" {
		return v
	}
	return defaultBroker
}

func main() {
	tuiMode := flag.Bool("tui", false, "start the TUI map instead of stdout printing")
	flag.Parse()

	if *tuiMode {
		p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client := blitzortung.NewClient()
	if err := client.Connect(brokerURL()); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer client.Disconnect()

	rate := &blitzortung.RateCounter{}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			os.Exit(0)
		case s := <-client.Strikes():
			PrintStrike(s)
			rate.Add()
		case <-ticker.C:
			fmt.Printf("rate: %.1f strikes/min\n", rate.Rate())
		}
	}
}
