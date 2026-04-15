package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/yinuoli/balatro-cli/engine"
	"github.com/yinuoli/balatro-cli/i18n"
)

// FocusType indicates what kind of element the cursor is on.
type FocusType int

const (
	FocusNone FocusType = iota
	FocusHandCard
	FocusJokerSlot
	FocusShopJoker
	FocusShopBooster
)

// FocusTarget describes the current cursor focus for the status bar.
type FocusTarget struct {
	Type  FocusType
	Index int
}

// RenderStatusBar renders the bottom status bar.
// Layout: <contextual info on the left>  <shortcut cheatsheet on the right>.
// Shortcuts are always visible so the player never has to guess a key.
func RenderStatusBar(focus FocusTarget, game *engine.GameState, width int) string {
	context := buildContextInfo(focus, game)
	shortcuts := buildShortcutsLine(game)

	// Inner width of the status bar (StatusBarStyle uses Padding(0,1) → 2 cells).
	barWidth := width - 2
	if barWidth < 0 {
		barWidth = 0
	}
	innerWidth := barWidth - 2
	if innerWidth < 0 {
		innerWidth = 0
	}

	shortcutsW := lipgloss.Width(shortcuts)
	contextW := lipgloss.Width(context)

	var body string
	switch {
	case shortcutsW >= innerWidth:
		// Not enough room to fit both — shortcuts take priority.
		body = shortcuts
	case context == "":
		body = lipgloss.NewStyle().Width(innerWidth).Align(lipgloss.Right).Render(shortcuts)
	case contextW+2+shortcutsW > innerWidth:
		// Truncate context to keep shortcuts fully visible.
		leftW := innerWidth - shortcutsW - 2
		left := lipgloss.NewStyle().Width(leftW).MaxWidth(leftW).Render(context)
		body = left + "  " + shortcuts
	default:
		leftW := innerWidth - shortcutsW
		left := lipgloss.NewStyle().Width(leftW).Render(context)
		body = left + shortcuts
	}

	return StatusBarStyle.Width(barWidth).Render(body)
}

// buildContextInfo returns the contextual info (card/joker details) based on focus.
// Returns an empty string when there is nothing specific to show.
func buildContextInfo(focus FocusTarget, game *engine.GameState) string {
	switch focus.Type {
	case FocusHandCard:
		if focus.Index >= 0 && focus.Index < len(game.Hand) {
			card := game.Hand[focus.Index]
			info := fmt.Sprintf("%s%s  %s  %s: %d",
				card.Rank.String(), card.Suit.Symbol(),
				i18n.T(suitKey(card.Suit)),
				i18n.T("status.chips_value"), card.Rank.Chips())
			if card.Debuffed {
				info += "  " + i18n.T("status.debuffed")
			}
			return info
		}
	case FocusJokerSlot:
		if focus.Index >= 0 && focus.Index < len(game.Jokers) {
			j := game.Jokers[focus.Index]
			return fmt.Sprintf("%s — %s | %s: $%d | %s: $%d | %s: %s",
				i18n.T(j.Def.NameKey), i18n.T(j.Def.DescKey),
				i18n.T("status.cost"), j.Cost,
				i18n.T("status.sell_value"), j.SellValue(),
				i18n.T("status.rarity"), i18n.T(rarityI18nKey(j.Def.Rarity)))
		}
		return i18n.T("status.empty_slot")
	case FocusShopJoker:
		if game.Shop != nil && focus.Index >= 0 && focus.Index < 2 {
			def := game.Shop.JokerSlots[focus.Index]
			if def != nil {
				return fmt.Sprintf("%s — %s | %s: $%d | %s: %s",
					i18n.T(def.NameKey), i18n.T(def.DescKey),
					i18n.T("status.cost"), game.Shop.JokerCosts[focus.Index],
					i18n.T("status.rarity"), i18n.T(rarityI18nKey(def.Rarity)))
			}
			return i18n.T("shop.sold_out")
		}
	}
	return ""
}

// buildShortcutsLine returns the keyboard cheatsheet for the current phase.
// The cheatsheet is rendered in a dim color so it doesn't compete visually
// with the contextual info on the left.
func buildShortcutsLine(game *engine.GameState) string {
	dim := lipgloss.NewStyle().Foreground(ColorDim)
	key := lipgloss.NewStyle().Foreground(ColorCyan).Bold(true)

	// item formats a "(key)label" pair with the key highlighted.
	item := func(k, label string) string {
		return key.Render("("+k+")") + dim.Render(label)
	}
	sep := dim.Render("  ")

	var parts []string
	switch game.Phase {
	case engine.PhaseShop:
		parts = []string{
			item("←→", i18n.T("action.select")),
			item("Enter", i18n.T("action.buy")),
			item("Q", i18n.T("action.quit")),
		}
	case engine.PhasePlay:
		parts = []string{
			item("←→", "移动"),
			item("Space", i18n.T("action.select")),
			item("Enter", i18n.T("action.play")),
			item("D", i18n.T("action.discard")),
			item("s", i18n.T("action.sort_rank")),
			item("S", i18n.T("action.sort_suit")),
			item("Tab", "切换"),
			item("Q", i18n.T("action.quit")),
		}
	default:
		parts = []string{
			item("Enter", i18n.T("action.continue")),
			item("Q", i18n.T("action.quit")),
		}
	}
	return joinWith(parts, sep)
}

func joinWith(parts []string, sep string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += sep
		}
		out += p
	}
	return out
}

func suitKey(s engine.Suit) string {
	keys := map[engine.Suit]string{
		engine.Spades: "suit.spades", engine.Hearts: "suit.hearts",
		engine.Clubs: "suit.clubs", engine.Diamonds: "suit.diamonds",
	}
	return keys[s]
}

func rarityI18nKey(r engine.Rarity) string {
	keys := map[engine.Rarity]string{
		engine.Common: "rarity.common", engine.Uncommon: "rarity.uncommon",
		engine.Rare: "rarity.rare", engine.Legendary: "rarity.legendary",
	}
	return keys[r]
}
