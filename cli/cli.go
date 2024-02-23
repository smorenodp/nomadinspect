package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smorenodp/nomadinspect/screens"
)

type Model struct {
	screen     screens.MainScreen
	listScreen screens.ListScreen
	width      int
	height     int
	quitting   bool
	err        error
}

func New(namespaces, matches []string, and bool) Model {
	s := screens.NewSpinnerScreen(namespaces, matches, and)
	return Model{screen: s}
}

func (m Model) Run() error {
	var err error

	if _, err = tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		return err
	}

	return nil
}

func (m Model) Init() tea.Cmd {
	return m.screen.Start()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			screen, cmd := m.screen.Update(msg)
			m.screen = screen
			return m, cmd
		}
	case error:
		m.err = msg
		return m, nil
	case screens.OutputMessage:
		listScreen := screens.NewListScreen(msg.Jobs, m.width, m.height)
		m.screen = listScreen
		m.listScreen = listScreen
		return m, nil
	case screens.JobView:
		viewScreen := screens.NewResourceScreen(screens.Job(msg), m.width, m.height)
		m.screen = &viewScreen
		return m, viewScreen.Start()
	case screens.EndView:
		m.screen = m.listScreen
		return m, nil
	default:
		screen, cmd := m.screen.Update(msg)
		m.screen = screen
		return m, cmd
	}
}

func (m Model) View() string {
	if !m.quitting {
		if m.err != nil {
			return fmt.Sprintf("Error: %s", m.err)
		} else {
			return m.screen.View()
		}
	}
	return ""
}
