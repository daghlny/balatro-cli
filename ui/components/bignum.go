package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// 3-line tall block digit font
var bigDigits = map[rune][3]string{
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
	',': {"   ", "   ", " ▄ "},
	'.': {"   ", "   ", " ▄ "},
	' ': {"   ", "   ", "   "},
}

// RenderBigNumber renders a number string in 3-line tall block font.
func RenderBigNumber(s string, color lipgloss.Color) string {
	lines := [3]strings.Builder{}
	for _, ch := range s {
		d, ok := bigDigits[ch]
		if !ok {
			d = [3]string{"   ", "   ", "   "}
		}
		for i := 0; i < 3; i++ {
			lines[i].WriteString(d[i])
			lines[i].WriteString(" ")
		}
	}

	style := lipgloss.NewStyle().Foreground(color).Bold(true)
	result := make([]string, 3)
	for i := 0; i < 3; i++ {
		result[i] = style.Render(lines[i].String())
	}
	return strings.Join(result, "\n")
}

// RenderScoreDisplay renders "Chips × Mult" in big stylized format.
func RenderScoreDisplay(chips int, mult float64) string {
	chipsStr := FormatNumber(chips)
	multStr := FormatNumber(int(mult))

	chipsRendered := RenderBigNumber(chipsStr, ColorBlue)
	multRendered := RenderBigNumber(multStr, ColorRed)

	// The × symbol between them, vertically centered (3 lines)
	xSymbol := lipgloss.NewStyle().Foreground(ColorWhite).Bold(true).Render(
		"   \n ✕ \n   ")

	return lipgloss.JoinHorizontal(lipgloss.Center, chipsRendered, xSymbol, multRendered)
}
