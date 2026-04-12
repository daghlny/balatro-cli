package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yinuoli/balatro-cli/engine"
	"github.com/yinuoli/balatro-cli/ui/components"
)

// Screen defines the interface for each game phase screen.
type Screen interface {
	Update(msg tea.Msg, game *engine.GameState) tea.Cmd
	View(game *engine.GameState, width, height int) string
	Focus() components.FocusTarget
}
