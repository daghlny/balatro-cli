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

type BlindSelectScreen struct {
	cursor int // 0=challenge, 1=skip
}

func (s *BlindSelectScreen) Update(msg tea.Msg, game *engine.GameState) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		blind := game.CurrentBlindType()
		switch key.String() {
		case "left", "h":
			s.cursor = 0
		case "right", "l":
			if blind != engine.BossBlind {
				s.cursor = 1
			}
		case "enter":
			if s.cursor == 0 {
				game.EnterBlind()
			} else if blind != engine.BossBlind {
				game.SkipBlind()
			}
		case "q", "ctrl+c":
			return tea.Quit
		}
	}
	return nil
}

func (s *BlindSelectScreen) View(game *engine.GameState, width, height int) string {
	var b strings.Builder

	blind := game.CurrentBlindType()
	blindName := blindI18n(blind)
	target := engine.CalcTargetScore(game.Ante, blind, 0)

	header := lipgloss.NewStyle().Bold(true).Foreground(components.ColorGold).
		Align(lipgloss.Center).Width(width).
		Render(fmt.Sprintf("%s %d/8", i18n.T("ui.ante"), game.Ante))
	b.WriteString("\n" + header + "\n\n")

	// Blind info box
	info := fmt.Sprintf("  %s\n  %s: %s\n  %s: $%d",
		blindName,
		i18n.T("ui.target"), components.FormatNumber(target),
		i18n.T("ui.money"), game.Money)
	if blind == engine.BossBlind {
		info += "\n  (Boss效果将在下一屏展示)"
	}

	infoBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(components.ColorBorder).
		Padding(1, 2).
		Width(40).
		Align(lipgloss.Center).
		Render(info)
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(width).Render(infoBox))
	b.WriteString("\n\n")

	// Buttons
	challengeStyle := buttonStyle(s.cursor == 0)
	challengeBtn := challengeStyle.Render(" " + i18n.T("action.challenge") + " ")

	var buttons string
	if blind == engine.BossBlind {
		buttons = challengeBtn
	} else {
		skipStyle := buttonStyle(s.cursor == 1)
		skipBtn := skipStyle.Render(" " + i18n.T("action.skip") + " ")
		buttons = challengeBtn + "    " + skipBtn
	}

	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(width).Render(buttons))
	return b.String()
}

func (s *BlindSelectScreen) Focus() components.FocusTarget {
	return components.FocusTarget{Type: components.FocusNone}
}

func buttonStyle(active bool) lipgloss.Style {
	s := lipgloss.NewStyle().Padding(0, 2).Border(lipgloss.RoundedBorder())
	if active {
		return s.BorderForeground(components.ColorSelected).Foreground(components.ColorGold).Bold(true)
	}
	return s.BorderForeground(components.ColorDim).Foreground(components.ColorDim)
}

func blindI18n(b engine.BlindType) string {
	switch b {
	case engine.SmallBlind:
		return i18n.T("blind.small")
	case engine.BigBlind:
		return i18n.T("blind.big")
	case engine.BossBlind:
		return i18n.T("blind.boss")
	}
	return ""
}
