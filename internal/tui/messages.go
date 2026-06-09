package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/markwallsgrove/boltmap/internal/blitzortung"
)

const tickInterval = 500 * time.Millisecond

type strikeMsg blitzortung.Strike

type connectionMsg bool

type tickMsg struct{}

func waitForStrike(ch <-chan blitzortung.Strike) tea.Cmd {
	return func() tea.Msg {
		return strikeMsg(<-ch)
	}
}

func waitForConnectionChange(ch <-chan bool) tea.Cmd {
	return func() tea.Msg {
		return connectionMsg(<-ch)
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(tickInterval, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}
