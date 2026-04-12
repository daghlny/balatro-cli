package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yinuoli/balatro-cli/engine"
	"github.com/yinuoli/balatro-cli/i18n"
	
)

// RenderJokerCard renders a 10×7 joker card.
func RenderJokerCard(joker *engine.JokerInstance, focused bool) string {
	name := i18n.T(joker.Def.NameKey)
	desc := i18n.T(joker.Def.DescKey)
	cost := fmt.Sprintf("$%d", joker.Cost)

	borderColor := rarityColor(joker.Def.Rarity)
	if focused {
		borderColor = ColorFocused
	}

	// Split name into 2 lines if needed (max 8 chars per line)
	nameLine1, nameLine2 := splitName(name, 8)

	// Truncate desc to 8 chars
	if len(desc) > 8 {
		desc = desc[:8]
	}

	lines := []string{
		padCenter(nameLine1, 8),
		padCenter(nameLine2, 8),
		"        ",
		padCenter(desc, 8),
		padCenter(cost, 8),
	}

	style := lipgloss.NewStyle().
		Foreground(ColorWhite).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(8)

	return style.Render(strings.Join(lines, "\n"))
}

// RenderEmptyJokerSlot renders an empty joker slot with dashed border.
func RenderEmptyJokerSlot() string {
	lines := []string{
		"        ",
		"  空位  ",
		"        ",
		"        ",
		"        ",
	}
	style := lipgloss.NewStyle().
		Foreground(ColorDim).
		BorderStyle(lipgloss.HiddenBorder()).
		BorderForeground(ColorDim).
		Width(8)

	return style.Render(strings.Join(lines, "\n"))
}

// RenderJokerSlots renders all joker slots (filled + empty).
func RenderJokerSlots(jokers []*engine.JokerInstance, maxSlots int, focusIndex int) string {
	parts := make([]string, maxSlots)
	for i := 0; i < maxSlots; i++ {
		if i < len(jokers) {
			parts[i] = RenderJokerCard(jokers[i], i == focusIndex)
		} else {
			parts[i] = RenderEmptyJokerSlot()
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

func rarityColor(r engine.Rarity) lipgloss.Color {
	switch r {
	case engine.Common:
		return ColorBlue
	case engine.Uncommon:
		return ColorGreen
	case engine.Rare:
		return ColorPurple
	case engine.Legendary:
		return ColorGold
	default:
		return ColorBorder
	}
}

func splitName(name string, maxWidth int) (string, string) {
	runes := []rune(name)
	if len(runes) <= maxWidth {
		return name, ""
	}
	return string(runes[:maxWidth]), string(runes[maxWidth:])
}
