package components

import (
	"github.com/charmbracelet/lipgloss"
)

// ComposePlayLayout assembles the main play screen layout:
// sidebar | [top: jokers / middle: table / bottom: hand]
// ─────────────────────────────────────────────────────
// status bar
func ComposePlayLayout(sidebar, jokerArea, tableArea, handArea, statusBar string, width, height int) string {
	// Right panel: stack top/middle/bottom vertically
	rightPanel := lipgloss.JoinVertical(lipgloss.Left, jokerArea, tableArea, handArea)

	// Main content: sidebar + right panel
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, rightPanel)

	// Full layout: main + status bar
	return lipgloss.JoinVertical(lipgloss.Left, mainContent, statusBar)
}
