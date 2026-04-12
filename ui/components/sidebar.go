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

	// Hand type levels
	b.WriteString("──────────────────────\n")
	for _, ht := range engine.AllHandTypes() {
		level := game.HandLevels[ht]
		htKey := handTypeKey(ht)
		name := i18n.T(htKey)
		if level > 1 {
			b.WriteString(fmt.Sprintf("  %-8s Lv.%d\n", name, level))
		}
	}

	style := SidebarStyle.Height(height)
	return style.Render(b.String())
}

// renderSidebarScore renders the Chips × Mult block for the sidebar.
// Uses 2-line compact big digits that fit within ~22 char width.
func renderSidebarScore(chips int, mult float64, total int) string {
	chipsStr := FormatNumber(chips)
	multStr := FormatNumber(int(mult))
	totalStr := FormatNumber(total)

	chipsRendered := renderCompactBig(chipsStr, ColorBlue)
	multRendered := renderCompactBig(multStr, ColorRed)

	// Build 2-line display:
	// line 0:  <chips>  ✕  <mult>
	// line 1:  <chips>  ✕  <mult>
	chipsLines := strings.Split(chipsRendered, "\n")
	multLines := strings.Split(multRendered, "\n")

	xTop := lipgloss.NewStyle().Foreground(ColorWhite).Bold(true).Render("   ")
	xBot := lipgloss.NewStyle().Foreground(ColorWhite).Bold(true).Render(" ✕ ")

	var result strings.Builder
	if len(chipsLines) >= 2 && len(multLines) >= 2 {
		result.WriteString("  " + chipsLines[0] + xTop + multLines[0] + "\n")
		result.WriteString("  " + chipsLines[1] + xBot + multLines[1] + "\n")
	}

	// Total line
	result.WriteString("  = ")
	result.WriteString(lipgloss.NewStyle().Foreground(ColorGold).Bold(true).Render(totalStr))
	result.WriteString("\n")

	return result.String()
}

// 2-line compact big digits (each digit is 2 chars wide × 2 lines tall)
var compactDigits = map[rune][2]string{
	'0': {"▄▀▄", "▀▄▀"},
	'1': {" ▄ ", " █ "},
	'2': {"▄▀▄", "▀▄▀"},
	'3': {"▀▀▄", "▀▀▀"},
	'4': {"▄ ▄", "▀▀█"},
	'5': {"█▀▀", "▀▀█"},
	'6': {"▄▀ ", "▀▀█"},
	'7': {"▀▀█", "  █"},
	'8': {"▄▀▄", "▀▀▀"},
	'9': {"█▀▀", " ▀█"},
	',': {"   ", " . "},
}

func renderCompactBig(s string, color lipgloss.Color) string {
	lines := [2]strings.Builder{}
	for _, ch := range s {
		d, ok := compactDigits[ch]
		if !ok {
			d = [2]string{"   ", "   "}
		}
		lines[0].WriteString(d[0])
		lines[1].WriteString(d[1])
	}
	style := lipgloss.NewStyle().Foreground(color).Bold(true)
	return style.Render(lines[0].String()) + "\n" + style.Render(lines[1].String())
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
