package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yinuoli/balatro-cli/engine"
	
)

// RenderCard renders a 9×7 playing card.
// selected: card is chosen for play/discard (gold border)
// focused: cursor is on this card (cyan border)
// debuffed: card is debuffed by boss (dim)
func RenderCard(card engine.Card, selected, focused, debuffed bool) string {
	rankStr := card.Rank.String()
	suitSym := card.Suit.Symbol()

	// Determine foreground color
	fg := ColorWhite
	if card.Suit.IsRed() {
		fg = ColorRed
	}
	if debuffed {
		fg = ColorDim
	}

	// Determine border color
	borderColor := ColorBorder
	if selected {
		borderColor = ColorSelected
	}
	if focused {
		borderColor = ColorFocused
	}

	// Build the 5 inner lines (7×1 each, padded to width 7)
	// Line 1: rank top-left
	// Line 2: empty
	// Line 3: suit center
	// Line 4: empty
	// Line 5: rank bottom-right
	topLeft := fmt.Sprintf(" %-6s", rankStr)
	center := fmt.Sprintf("   %s   ", suitSym)
	bottomRight := fmt.Sprintf(" %6s", rankStr)

	// Trim/pad to exactly 7 chars
	topLeft = padRight(topLeft, 7)
	center = padCenter(suitSym, 7)
	bottomRight = padLeft(rankStr, 7)

	lines := []string{topLeft, "       ", center, "       ", bottomRight}

	style := lipgloss.NewStyle().
		Foreground(fg).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(7)

	content := strings.Join(lines, "\n")
	rendered := style.Render(content)

	// If selected, add visual "lift" (1 empty line above, remove 1 below)
	if selected {
		renderedLines := strings.Split(rendered, "\n")
		// Prepend empty line, remove last line
		lifted := append([]string{""}, renderedLines[:len(renderedLines)-1]...)
		return strings.Join(lifted, "\n")
	}

	return rendered
}

// RenderCardBack renders a face-down card.
func RenderCardBack() string {
	style := lipgloss.NewStyle().
		Foreground(ColorDim).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorDim).
		Width(7)

	lines := []string{
		"░░░░░░░",
		"░░░░░░░",
		"░░░░░░░",
		"░░░░░░░",
		"░░░░░░░",
	}
	return style.Render(strings.Join(lines, "\n"))
}

// RenderHand renders a row of cards horizontally.
func RenderHand(cards []engine.Card, selected []bool, cursorIndex int, debuffed []bool) string {
	if len(cards) == 0 {
		return ""
	}
	rendered := make([]string, len(cards))
	for i, card := range cards {
		isSel := i < len(selected) && selected[i]
		isFoc := i == cursorIndex
		isDbf := i < len(debuffed) && debuffed[i]
		rendered[i] = RenderCard(card, isSel, isFoc, isDbf)
	}
	return lipgloss.JoinHorizontal(lipgloss.Bottom, rendered...)
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

func padLeft(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return strings.Repeat(" ", width-len(s)) + s
}

func padCenter(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	left := (width - len(s)) / 2
	right := width - len(s) - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}
