package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yinuoli/balatro-cli/engine"
)

// Card inner dimensions (exclusive of the rounded border).
// Outer size becomes innerCardWidth+2 × innerCardHeight+2 = 18 × 14.
const (
	innerCardWidth  = 16
	innerCardHeight = 12
)

// RenderCard renders a playing card.
// selected: chosen for play/discard (gold border)
// focused:  cursor is on this card (cyan border)
// debuffed: boss debuff applied (dimmed)
//
// The lift animation for selected cards is NOT done here; it happens at the
// row level inside RenderHand so that the card keeps both of its borders and
// the "rise above" effect stays consistent across neighbours.
func RenderCard(card engine.Card, selected, focused, debuffed bool) string {
	rankStr := card.Rank.String()
	suitSym := card.Suit.Symbol()

	fg := ColorWhite
	if card.Suit.IsRed() {
		fg = ColorRed
	}
	if debuffed {
		fg = ColorDim
	}

	borderColor := ColorBorder
	if selected {
		borderColor = ColorSelected
	}
	if focused {
		borderColor = ColorFocused
	}

	lines := buildCardLines(rankStr, suitSym)

	style := lipgloss.NewStyle().
		Foreground(fg).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(innerCardWidth)

	return style.Render(strings.Join(lines, "\n"))
}

// buildCardLines returns the innerCardHeight inner rows for a card.
// Layout: rank in the upper-left corner, suit symbol centered, rank mirrored
// in the lower-right corner. Lots of whitespace by design — cards breathe.
func buildCardLines(rank, suit string) []string {
	const w = innerCardWidth
	blank := strings.Repeat(" ", w)

	// Upper-left rank: " K" + padding (right-padded to w).
	topRank := padRight(fmt.Sprintf(" %s", rank), w)
	// Lower-right rank: padding + "K " (left-padded to w).
	bottomRank := padLeft(fmt.Sprintf("%s ", rank), w)
	// Centered suit row.
	suitRow := padCenter(suit, w)

	lines := make([]string, innerCardHeight)
	for i := range lines {
		lines[i] = blank
	}
	// rows (0-indexed of 12): top pad, rank, blank*4, suit, blank*4, bottomRank, bot pad
	lines[1] = topRank
	lines[6] = suitRow
	lines[10] = bottomRank
	return lines
}

// RenderCardBack renders a face-down card at the same outer size as RenderCard.
func RenderCardBack() string {
	style := lipgloss.NewStyle().
		Foreground(ColorDim).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorDim).
		Width(innerCardWidth)

	row := strings.Repeat("░", innerCardWidth)
	lines := make([]string, innerCardHeight)
	for i := range lines {
		lines[i] = row
	}
	return style.Render(strings.Join(lines, "\n"))
}

// RenderHand renders a row of cards horizontally.
//
// When any card is selected, unselected cards get a leading blank row prepended
// and the row is top-aligned. That makes selected cards visibly sit one row
// higher than their neighbours — a proper "lift" — while every card keeps its
// full border intact.
func RenderHand(cards []engine.Card, selected []bool, cursorIndex int, debuffed []bool) string {
	if len(cards) == 0 {
		return ""
	}

	anySelected := false
	for _, sel := range selected {
		if sel {
			anySelected = true
			break
		}
	}

	rendered := make([]string, len(cards))
	for i, card := range cards {
		isSel := i < len(selected) && selected[i]
		isFoc := i == cursorIndex
		isDbf := i < len(debuffed) && debuffed[i]
		rendered[i] = RenderCard(card, isSel, isFoc, isDbf)
		if anySelected && !isSel {
			// Push unselected cards down by one row so selected ones appear lifted.
			rendered[i] = "\n" + rendered[i]
		}
	}
	if anySelected {
		return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
	}
	return lipgloss.JoinHorizontal(lipgloss.Bottom, rendered...)
}

func padRight(s string, width int) string {
	if lipgloss.Width(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-lipgloss.Width(s))
}

func padLeft(s string, width int) string {
	if lipgloss.Width(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-lipgloss.Width(s)) + s
}

func padCenter(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return s
	}
	left := (width - w) / 2
	right := width - w - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}
