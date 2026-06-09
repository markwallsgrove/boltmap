package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/markwallsgrove/boltmap/internal/maprender"
)

var statusBarStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("235")).
	Foreground(lipgloss.Color("255"))

type model struct {
	viewport maprender.Viewport
	width    int
	height   int
}

func NewModel() model {
	return model{viewport: maprender.DefaultViewport()}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "+":
			m.viewport.ZoomIn()
		case "-":
			m.viewport.ZoomOut()
		case "up":
			m.viewport.Pan(panStep(m.viewport)*0.1, 0)
		case "down":
			m.viewport.Pan(-panStep(m.viewport)*0.1, 0)
		case "left":
			m.viewport.Pan(0, -panLonStep(m.viewport)*0.1)
		case "right":
			m.viewport.Pan(0, panLonStep(m.viewport)*0.1)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	mapRows := m.height - 1
	mapStr := maprender.Render(m.viewport, m.width, mapRows)
	statusBar := statusBarStyle.Width(m.width).Render("Connecting… | 0 strikes")
	return mapStr + "\n" + statusBar
}

func panStep(vp maprender.Viewport) float64 {
	return 180.0 * vp.Zoom
}

func panLonStep(vp maprender.Viewport) float64 {
	return 360.0 * vp.Zoom
}
