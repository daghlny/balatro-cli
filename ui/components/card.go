package components

import (
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

// rankGlyphs is a 3×3 block font used to draw the rank in each corner of the
// card. Digits mirror the style of bignum.go so the visual language matches
// the sidebar score display. Letters (A/J/Q/K) are added for face cards.
var rankGlyphs = map[rune][3]string{
	'0': {"█▀█", "█ █", "▀▀▀"},
	'1': {" ▄█", "  █", "  ▀"},
	'2': {"▀▀█", "█▀▀", "▀▀▀"},
	'3': {"▀▀█", " ▀█", "▀▀▀"},
	'4': {"█ █", "▀▀█", "  ▀"},
	'5': {"█▀▀", "▀▀█", "▀▀▀"},
	'6': {"█▀▀", "█▀█", "▀▀▀"},
	'7': {"▀▀█", "  █", "  ▀"},
	'8': {"█▀█", "█▀█", "▀▀▀"},
	'9': {"█▀█", "▀▀█", "▀▀▀"},
	'A': {"▄▀▄", "█▀█", "▀ ▀"},
	'J': {"▀▀█", "  █", "▄█▀"},
	'Q': {"█▀█", "█ █", "▀▀▄"},
	'K': {"█ █", "█▀▄", "▀ ▀"},
}

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
// Layout: large 3-line block rank in the upper-left corner, suit symbol
// centered, rank mirrored in the lower-right corner.
func buildCardLines(rank, suit string) []string {
	const w = innerCardWidth
	blank := strings.Repeat(" ", w)

	rankBlock := renderRankGlyph(rank) // 3 lines

	// Upper-left: one column of left padding, then the glyph, pad right to w.
	topLines := [3]string{
		padRight(" "+rankBlock[0], w),
		padRight(" "+rankBlock[1], w),
		padRight(" "+rankBlock[2], w),
	}
	// Lower-right: mirror — pad left to w, with one trailing space column.
	bottomLines := [3]string{
		padLeft(rankBlock[0]+" ", w),
		padLeft(rankBlock[1]+" ", w),
		padLeft(rankBlock[2]+" ", w),
	}
	// Centered suit row.
	suitRow := padCenter(suit, w)

	lines := make([]string, innerCardHeight)
	for i := range lines {
		lines[i] = blank
	}
	// Row layout (12 rows): pad, rank×3, pad×2, suit, pad×2, rank×3.
	lines[1] = topLines[0]
	lines[2] = topLines[1]
	lines[3] = topLines[2]
	lines[6] = suitRow
	lines[9] = bottomLines[0]
	lines[10] = bottomLines[1]
	lines[11] = bottomLines[2]
	return lines
}

// renderRankGlyph renders a rank string ("K", "10", "A", ...) as a 3-line
// block-font. Multi-character ranks are concatenated with a single space
// column between glyphs so "10" still fits comfortably inside the card.
func renderRankGlyph(rank string) [3]string {
	var b [3]strings.Builder
	for i, ch := range rank {
		if i > 0 {
			for j := 0; j < 3; j++ {
				b[j].WriteByte(' ')
			}
		}
		g, ok := rankGlyphs[ch]
		if !ok {
			g = [3]string{"   ", "   ", "   "}
		}
		for j := 0; j < 3; j++ {
			b[j].WriteString(g[j])
		}
	}
	return [3]string{b[0].String(), b[1].String(), b[2].String()}
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
