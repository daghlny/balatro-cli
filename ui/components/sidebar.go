package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yinuoli/balatro-cli/engine"
	"github.com/yinuoli/balatro-cli/i18n"
)

// ASCII art logo, pure text, no emoji dependency
var logoLines = []string{
	"╔╗ ╔═╗╦  ╔═╗╔╦╗╦═╗╔═╗",
	"╠╩╗╠═╣║  ╠═╣ ║ ╠╦╝║ ║",
	"╚═╝╩ ╩╩═╝╩ ╩ ╩ ╩╚═╚═╝",
}

// RenderSidebar renders the left sidebar with game state info.
func RenderSidebar(game *engine.GameState, height int) string {
	var b strings.Builder

	// Logo
	logoStyle := lipgloss.NewStyle().Foreground(ColorGold).Bold(true)
	for _, line := range logoLines {
		b.WriteString(logoStyle.Render(line) + "\n")
	}
	b.WriteString("──────────────────────\n")

	blindName := ""
	switch game.CurrentBlindType() {
	case engine.SmallBlind:
		blindName = i18n.T("blind.small")
	case engine.BigBlind:
		blindName = i18n.T("blind.big")
	case engine.BossBlind:
		blindName = i18n.T("blind.boss")
	}

	// Ante & Blind
	writeLabel(&b, i18n.T("ui.ante"), fmt.Sprintf("%d / 8", game.Ante))
	b.WriteString("  " + blindName + "\n")

	if game.CurrentBoss != nil {
		bossName := i18n.T(game.CurrentBoss.NameKey())
		b.WriteString("  Boss: ")
		b.WriteString(lipgloss.NewStyle().Foreground(ColorRed).Bold(true).Render(bossName))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Target
	writeLabel(&b, i18n.T("ui.target"), FormatNumber(game.TargetScore))
	b.WriteString("\n")

	// --- Big Chips × Mult display ---
	if game.LastBreakdown != nil && game.Phase == engine.PhasePlay {
		bd := game.LastBreakdown
		b.WriteString(renderSidebarScore(bd.FinalChips, bd.FinalMult, bd.TotalScore))
		b.WriteString("\n")
	}

	// Cumulative round score
	writeLabel(&b, i18n.T("ui.score"), FormatNumber(game.RoundScore))
	b.WriteString("\n")

	// Hands & Discards
	handsColor := ColorWhite
	if game.RemainingHands <= 1 {
		handsColor = ColorRed
	}
	b.WriteString(fmt.Sprintf("  %s  %s / 4\n",
		lipgloss.NewStyle().Foreground(ColorDim).Render(i18n.T("ui.hands")),
		lipgloss.NewStyle().Foreground(handsColor).Bold(true).Render(fmt.Sprintf("%d", game.RemainingHands)),
	))
	b.WriteString(fmt.Sprintf("  %s  %s / 3\n",
		lipgloss.NewStyle().Foreground(ColorDim).Render(i18n.T("ui.discards")),
		lipgloss.NewStyle().Foreground(ColorWhite).Bold(true).Render(fmt.Sprintf("%d", game.RemainingDisc)),
	))
	b.WriteString("\n")

	// Money & Interest
	writeLabel(&b, i18n.T("ui.money"), fmt.Sprintf("$%d", game.Money))
	interest := game.Money / 5
	if interest > 5 {
		interest = 5
	}
	writeLabel(&b, i18n.T("ui.interest"), fmt.Sprintf("$%d", interest))
	b.WriteString("\n")

	// Hand type reference: always show all 10 hand types with their current
	// (level-adjusted) chips × mult so the player has a scoring cheat sheet.
	// Upgraded rows (Lv > 1) are highlighted in gold.
	b.WriteString("──────────────────────\n")
	title := lipgloss.NewStyle().Foreground(ColorCyan).Bold(true).
		Render("  " + i18n.T("ui.hand_reference"))
	b.WriteString(title + "\n")

	nameWidth := 10 // visual columns, fits "皇家同花顺" (10 cols)
	valWidth := 10  // visual columns for "  100 × 8"
	for _, ht := range engine.AllHandTypes() {
		level := game.HandLevels[ht]
		data := engine.HandTypeDataTable[ht]
		chips := data.BaseChips + (level-1)*data.LevelChips
		mult := data.BaseMult + (level-1)*data.LevelMult

		name := i18n.T(handTypeKey(ht))
		nameCell := lipgloss.NewStyle().Width(nameWidth).Render(name)
		valCell := lipgloss.NewStyle().Width(valWidth).Align(lipgloss.Right).
			Render(fmt.Sprintf("%d × %d", chips, mult))

		row := "  " + nameCell + valCell
		color := ColorWhite
		if level > 1 {
			color = ColorGold
			row += lipgloss.NewStyle().Foreground(ColorGold).
				Render(fmt.Sprintf(" Lv%d", level))
		}
		b.WriteString(lipgloss.NewStyle().Foreground(color).Render(row) + "\n")
	}

	style := SidebarStyle.Height(height)
	return style.Render(b.String())
}

// renderSidebarScore renders the "last play" Chips × Mult summary for the
// sidebar as a readable two-line text block.
//
// The full cinematic "big digits" display lives in the table area (see
// RenderScoreDisplay in bignum.go); in the narrow 22-col sidebar a plain
// colored text line reads much better than a cramped ASCII-art font.
func renderSidebarScore(chips int, mult float64, total int) string {
	chipsStr := lipgloss.NewStyle().Foreground(ColorBlue).Bold(true).
		Render(FormatNumber(chips))
	multStr := lipgloss.NewStyle().Foreground(ColorRed).Bold(true).
		Render(FormatNumber(int(mult)))
	totalStr := lipgloss.NewStyle().Foreground(ColorGold).Bold(true).
		Render(FormatNumber(total))
	sep := lipgloss.NewStyle().Foreground(ColorWhite).Render(" ✕ ")

	var b strings.Builder
	b.WriteString("  " + chipsStr + sep + multStr + "\n")
	b.WriteString("  = " + totalStr + "\n")
	return b.String()
}

func writeLabel(b *strings.Builder, label, value string) {
	line := fmt.Sprintf("  %s  %s\n",
		lipgloss.NewStyle().Foreground(ColorDim).Render(label),
		lipgloss.NewStyle().Foreground(ColorWhite).Bold(true).Render(value),
	)
	b.WriteString(line)
}

func handTypeKey(ht engine.HandType) string {
	keys := map[engine.HandType]string{
		engine.HighCard:      "hand.high_card",
		engine.Pair:          "hand.pair",
		engine.TwoPair:       "hand.two_pair",
		engine.ThreeOfAKind:  "hand.three_of_kind",
		engine.Straight:      "hand.straight",
		engine.Flush:         "hand.flush",
		engine.FullHouse:     "hand.full_house",
		engine.FourOfAKind:   "hand.four_of_kind",
		engine.StraightFlush: "hand.straight_flush",
		engine.RoyalFlush:    "hand.royal_flush",
	}
	return keys[ht]
}
