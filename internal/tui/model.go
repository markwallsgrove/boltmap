package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/markwallsgrove/boltmap/internal/blitzortung"
	"github.com/markwallsgrove/boltmap/internal/maprender"
)

const fadeTTL = 10 * time.Second

var (
	statusBarStyle    = lipgloss.NewStyle().Background(lipgloss.Color("235")).Foreground(lipgloss.Color("255"))
	connectedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	disconnectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

type Config struct {
	Client *blitzortung.Client
	TTL    time.Duration
}

type model struct {
	viewport  maprender.Viewport
	width     int
	height    int
	client    *blitzortung.Client
	buffer    *StrikeBuffer
	rate      *blitzortung.RateCounter
	total     int
	connected bool
	ttl       time.Duration
}

func NewModel(cfg Config) model {
	return model{
		viewport: maprender.DefaultViewport(),
		client:   cfg.Client,
		buffer:   &StrikeBuffer{},
		rate:     &blitzortung.RateCounter{},
		ttl:      cfg.TTL,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForStrike(m.client.Strikes()),
		waitForConnectionChange(m.client.ConnectionState()),
		tickCmd(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "+":
			m.viewport.ZoomIn()
			m.buffer.Reproject(m.viewport, m.width, m.height-1)
		case "-":
			m.viewport.ZoomOut()
			m.buffer.Reproject(m.viewport, m.width, m.height-1)
		case "up":
			m.viewport.Pan(panStep(m.viewport)*0.1, 0)
			m.buffer.Reproject(m.viewport, m.width, m.height-1)
		case "down":
			m.viewport.Pan(-panStep(m.viewport)*0.1, 0)
			m.buffer.Reproject(m.viewport, m.width, m.height-1)
		case "left":
			m.viewport.Pan(0, -panLonStep(m.viewport)*0.1)
			m.buffer.Reproject(m.viewport, m.width, m.height-1)
		case "right":
			m.viewport.Pan(0, panLonStep(m.viewport)*0.1)
			m.buffer.Reproject(m.viewport, m.width, m.height-1)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.buffer.Reproject(m.viewport, m.width, m.height-1)
	case strikeMsg:
		s := blitzortung.Strike(msg)
		col, row := -1, -1
		if m.width > 0 && m.height > 1 {
			if c, r, ok := maprender.LatLonToCell(m.viewport, s.Lat, s.Lon, m.width, m.height-1); ok {
				col, row = c, r
			}
		}
		m.buffer.Add(s, col, row)
		m.rate.Add()
		m.total++
		return m, waitForStrike(m.client.Strikes())
	case connectionMsg:
		m.connected = bool(msg)
		return m, waitForConnectionChange(m.client.ConnectionState())
	case tickMsg:
		return m, tickCmd()
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	mapRows := m.height - 1
	now := time.Now()
	active := m.buffer.Active(now, m.ttl)

	overlays := make(map[[2]int]lipgloss.Color, len(active))
	for _, ps := range active {
		if ps.Col < 0 {
			continue
		}
		age := now.Sub(ps.ArrivalTime)
		var color lipgloss.Color
		switch {
		case age >= fadeTTL:
			color = lipgloss.Color("240")
		case ps.Strike.Pol == 1:
			color = lipgloss.Color("11")
		default:
			color = lipgloss.Color("14")
		}
		overlays[[2]int{ps.Col, ps.Row}] = color
	}

	mapStr := maprender.RenderWithOverlay(m.viewport, m.width, mapRows, overlays)
	return mapStr + "\n" + m.renderStatusBar()
}

func (m model) renderStatusBar() string {
	var connStr string
	if m.connected {
		connStr = connectedStyle.Render("CONNECTED")
	} else {
		connStr = disconnectedStyle.Render("DISCONNECTED")
	}
	text := fmt.Sprintf("%s | %d strikes | %.1f/min", connStr, m.total, m.rate.Rate())
	return statusBarStyle.Width(m.width).Render(text)
}

func panStep(vp maprender.Viewport) float64 {
	return 180.0 * vp.Zoom
}

func panLonStep(vp maprender.Viewport) float64 {
	return 360.0 * vp.Zoom
}
