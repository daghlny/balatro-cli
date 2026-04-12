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

type SettlementScreen struct{}

func (s *SettlementScreen) Update(msg tea.Msg, game *engine.GameState) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			game.FinishSettlement()
		case "q", "ctrl+c":
			return tea.Quit
		}
	}
	return nil
}

func (s *SettlementScreen) View(game *engine.GameState, width, height int) string {
	var b strings.Builder

	blind := game.CurrentBlindType()
	base := blind.BaseReward()
	handBonus := game.RemainingHands
	interest := game.Money / 5
	if interest > 5 {
		interest = 5
	}
	total := base + handBonus + interest

	header := lipgloss.NewStyle().Bold(true).Foreground(components.ColorGold).
		Align(lipgloss.Center).Width(width).
		Render("✓ " + blindI18n(blind) + " 通过!")
	b.WriteString("\n\n" + header + "\n\n\n")

	lines := []string{
		fmt.Sprintf("  %s:        $%d", i18n.T("settlement.base"), base),
		fmt.Sprintf("  %s:  $%d  (%d × $1)", i18n.T("settlement.hand_bonus"), handBonus, game.RemainingHands),
		fmt.Sprintf("  %s:          $%d", i18n.T("settlement.interest"), interest),
		"  ─────────────────",
		fmt.Sprintf("  %s:          $%d", i18n.T("settlement.total"), total),
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(components.ColorGold).
		Padding(1, 2).
		Width(40).
		Render(strings.Join(lines, "\n"))
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(width).Render(box))
	b.WriteString("\n\n")

	moneyAfter := game.Money + total
	moneyLine := fmt.Sprintf("%s: $%d → $%d",
		i18n.T("ui.money"),
		game.Money,
		moneyAfter)
	b.WriteString(lipgloss.NewStyle().Foreground(components.ColorGreen).Bold(true).
		Align(lipgloss.Center).Width(width).Render(moneyLine))

	_ = components.FormatNumber(0) // ensure import

	prompt := lipgloss.NewStyle().Foreground(components.ColorDim).Align(lipgloss.Center).Width(width).
		Render("按 Enter 进入商店")
	b.WriteString(prompt)

	return b.String()
}

func (s *SettlementScreen) Focus() components.FocusTarget {
	return components.FocusTarget{Type: components.FocusNone}
}
