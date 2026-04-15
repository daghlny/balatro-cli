package components

import (
	"github.com/charmbracelet/lipgloss"
)

// ComposePlayLayout assembles the main play screen layout:
// sidebar | [top: jokers / middle: table / bottom: hand]
//
// The status bar is rendered by the root AppModel (see ui/app.go) and
// joined below this content — don't append one here, or it will render
// twice.
func ComposePlayLayout(sidebar, jokerArea, tableArea, handArea string, width, height int) string {
	rightPanel := lipgloss.JoinVertical(lipgloss.Left, jokerArea, tableArea, handArea)
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, rightPanel)
}
