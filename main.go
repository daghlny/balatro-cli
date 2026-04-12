package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yinuoli/balatro-cli/engine"
	_ "github.com/yinuoli/balatro-cli/i18n"   // register locales
	_ "github.com/yinuoli/balatro-cli/jokers" // register jokers
	"github.com/yinuoli/balatro-cli/ui"
)

func main() {
	game := engine.NewGame(engine.NewRNG())
	model := ui.NewAppModel(game)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
