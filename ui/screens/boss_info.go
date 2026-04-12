package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yinuoli/balatro-cli/engine"
	"github.com/yinuoli/balatro-cli/i18n"
	"github.com/yinuoli/balatro-cli/ui/components"
)

type BossInfoScreen struct{}

func (s *BossInfoScreen) Update(msg tea.Msg, game *engine.GameState) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			game.ConfirmBossInfo()
		case "q", "ctrl+c":
			return tea.Quit
		}
	}
	return nil
}

func (s *BossInfoScreen) View(game *engine.GameState, width, height int) string {
	var b strings.Builder

	boss := game.CurrentBoss
	if boss == nil {
		return "No boss"
	}

	name := i18n.T(boss.NameKey())
	desc := i18n.T(boss.DescKey())
	target := game.TargetScore
	if target == 0 {
		target = engine.CalcTargetScore(game.Ante, engine.BossBlind, boss.TargetMultiplier())
	}

	header := lipgloss.NewStyle().Bold(true).Foreground(components.ColorRed).
		Align(lipgloss.Center).Width(width).
		Render("⚠ " + i18n.T("blind.boss") + " ⚠")
	b.WriteString("\n\n" + header + "\n\n")

	info := fmt.Sprintf("  %s\n\n  %s\n\n  %s: %s",
		lipgloss.NewStyle().Bold(true).Foreground(components.ColorGold).Render(name),
		lipgloss.NewStyle().Foreground(components.ColorWhite).Render(desc),
		i18n.T("ui.target"),
		lipgloss.NewStyle().Bold(true).Foreground(components.ColorRed).Render(components.FormatNumber(target)),
	)

	infoBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(components.ColorRed).
		Padding(1, 3).
		Width(50).
		Render(info)
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(width).Render(infoBox))
	b.WriteString("\n\n")

	prompt := lipgloss.NewStyle().Foreground(components.ColorDim).Align(lipgloss.Center).Width(width).
		Render("按 Enter " + i18n.T("action.continue"))
	b.WriteString(prompt)

	return b.String()
}

func (s *BossInfoScreen) Focus() components.FocusTarget {
	return components.FocusTarget{Type: components.FocusNone}
}
