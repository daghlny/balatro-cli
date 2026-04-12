package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yinuoli/balatro-cli/engine"
	"github.com/yinuoli/balatro-cli/i18n"
	"github.com/yinuoli/balatro-cli/ui/components"
)

type GameOverScreen struct{}

func (s *GameOverScreen) Update(msg tea.Msg, game *engine.GameState) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			// Restart
			*game = *engine.NewGame(engine.NewRNG())
			game.StartGame()
		case "q", "ctrl+c":
			return tea.Quit
		}
	}
	return nil
}

func (s *GameOverScreen) View(game *engine.GameState, width, height int) string {
	var b strings.Builder

	topPad := height / 3
	b.WriteString(strings.Repeat("\n", topPad))

	var msg string
	var color lipgloss.Color
	if game.Won {
		msg = i18n.T("gameover.win")
		color = components.ColorGold
	} else {
		msg = i18n.T("gameover.lose")
		color = components.ColorRed
	}

	title := lipgloss.NewStyle().Bold(true).Foreground(color).
		Align(lipgloss.Center).Width(width).
		Render(msg)
	b.WriteString(title + "\n\n")

	prompt := lipgloss.NewStyle().Foreground(components.ColorDim).
		Align(lipgloss.Center).Width(width).
		Render(i18n.T("gameover.play_again"))
	b.WriteString(prompt)

	return b.String()
}

func (s *GameOverScreen) Focus() components.FocusTarget {
	return components.FocusTarget{Type: components.FocusNone}
}
