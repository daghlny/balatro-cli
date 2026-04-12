package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yinuoli/balatro-cli/engine"
	"github.com/yinuoli/balatro-cli/i18n"
	"github.com/yinuoli/balatro-cli/ui/components"
)

type TitleScreen struct{}

func (s *TitleScreen) Update(msg tea.Msg, game *engine.GameState) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			game.StartGame()
		case "q", "ctrl+c":
			return tea.Quit
		}
	}
	return nil
}

func (s *TitleScreen) View(game *engine.GameState, width, height int) string {
	var b strings.Builder

	title := components.TitleStyle.Width(width).Render(i18n.T("title.name"))
	subtitle := lipgloss.NewStyle().
		Foreground(components.ColorDim).
		Align(lipgloss.Center).
		Width(width).
		Render(i18n.T("title.subtitle"))
	prompt := lipgloss.NewStyle().
		Foreground(components.ColorWhite).
		Align(lipgloss.Center).
		Width(width).
		Render(i18n.T("title.press_start"))

	// Center vertically
	topPad := height / 3
	b.WriteString(strings.Repeat("\n", topPad))
	b.WriteString(title + "\n\n")
	b.WriteString(subtitle + "\n\n\n")
	b.WriteString(prompt + "\n")

	return b.String()
}

func (s *TitleScreen) Focus() components.FocusTarget {
	return components.FocusTarget{Type: components.FocusNone}
}
