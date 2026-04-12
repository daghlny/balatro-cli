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

type ScoringScreen struct{}

func (s *ScoringScreen) Update(msg tea.Msg, game *engine.GameState) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter", " ":
			game.FinishScoring()
		case "q", "ctrl+c":
			return tea.Quit
		}
	}
	return nil
}

func (s *ScoringScreen) View(game *engine.GameState, width, height int) string {
	bd := game.LastBreakdown
	if bd == nil {
		return "No score data"
	}

	var b strings.Builder

	htKey := handTypeKeyFromType(bd.HandType)
	htName := i18n.T(htKey)

	header := lipgloss.NewStyle().Bold(true).Foreground(components.ColorCyan).
		Align(lipgloss.Center).Width(width).
		Render(fmt.Sprintf("%s Lv.%d", htName, bd.HandLevel))
	b.WriteString("\n" + header + "\n\n")

	// Base values
	baseLine := fmt.Sprintf("  %s: %d   %s: %d",
		i18n.T("ui.chips"), bd.BaseChips,
		i18n.T("ui.mult"), bd.BaseMult)
	b.WriteString(lipgloss.NewStyle().Foreground(components.ColorDim).Align(lipgloss.Center).Width(width).Render(baseLine))
	b.WriteString("\n\n")

	// Card contributions
	for _, cd := range bd.CardDetails {
		if cd.Debuffed {
			line := fmt.Sprintf("  %s  (削弱, +0)", cd.Card.String())
			b.WriteString(lipgloss.NewStyle().Foreground(components.ColorDim).Render(line) + "\n")
		} else if cd.ChipsAdd > 0 {
			line := fmt.Sprintf("  %s  +%d %s", cd.Card.String(), cd.ChipsAdd, i18n.T("ui.chips"))
			b.WriteString(lipgloss.NewStyle().Foreground(components.ColorWhite).Render(line) + "\n")
		}
	}
	b.WriteString("\n")

	// Joker contributions
	for _, jd := range bd.JokerDetails {
		name := i18n.T("joker." + jd.JokerID + ".name")
		var effects []string
		if jd.Mod.PlusChips != 0 {
			effects = append(effects, fmt.Sprintf("+%d %s", jd.Mod.PlusChips, i18n.T("ui.chips")))
		}
		if jd.Mod.PlusMult != 0 {
			effects = append(effects, fmt.Sprintf("+%d %s", jd.Mod.PlusMult, i18n.T("ui.mult")))
		}
		if jd.Mod.MultMult > 1.0 {
			effects = append(effects, fmt.Sprintf("×%.1f %s", jd.Mod.MultMult, i18n.T("ui.mult")))
		}
		if len(effects) > 0 {
			line := fmt.Sprintf("  %s: %s", name, strings.Join(effects, ", "))
			b.WriteString(lipgloss.NewStyle().Foreground(components.ColorGreen).Render(line) + "\n")
		}
	}

	// Final
	finalLine := fmt.Sprintf("%d × %d = %s",
		bd.FinalChips, int(bd.FinalMult),
		components.FormatNumber(bd.TotalScore))
	b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(components.ColorGold).
		Align(lipgloss.Center).Width(width).Render(finalLine))

	totalLine := fmt.Sprintf("%s: %s / %s",
		i18n.T("ui.score"),
		components.FormatNumber(game.RoundScore),
		components.FormatNumber(game.TargetScore))
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(width).
		Foreground(components.ColorWhite).Render(totalLine))

	prompt := lipgloss.NewStyle().Foreground(components.ColorDim).Align(lipgloss.Center).Width(width).
		Render("按 Enter " + i18n.T("action.continue"))
	b.WriteString(prompt)

	return b.String()
}

func (s *ScoringScreen) Focus() components.FocusTarget {
	return components.FocusTarget{Type: components.FocusNone}
}
