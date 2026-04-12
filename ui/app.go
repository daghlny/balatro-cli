package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yinuoli/balatro-cli/engine"
	"github.com/yinuoli/balatro-cli/ui/components"
	"github.com/yinuoli/balatro-cli/ui/screens"
)

// AppModel is the root Bubble Tea model that routes to phase-specific screens.
type AppModel struct {
	game    *engine.GameState
	screens map[engine.Phase]screens.Screen
	width   int
	height  int
}

func NewAppModel(game *engine.GameState) AppModel {
	return AppModel{
		game: game,
		screens: map[engine.Phase]screens.Screen{
			engine.PhaseTitle:       &screens.TitleScreen{},
			engine.PhaseBlindSelect: &screens.BlindSelectScreen{},
			engine.PhaseBossInfo:    &screens.BossInfoScreen{},
			engine.PhasePlay:        screens.NewPlayScreen(),
			engine.PhaseScoring:     &screens.ScoringScreen{},
			engine.PhaseSettlement:  &screens.SettlementScreen{},
			engine.PhaseShop:        &screens.ShopScreen{},
			engine.PhaseGameOver:    &screens.GameOverScreen{},
		},
		width:  120,
		height: 40,
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	screen, ok := m.screens[m.game.Phase]
	if !ok {
		return m, nil
	}

	cmd := screen.Update(msg, m.game)
	return m, cmd
}

func (m AppModel) View() string {
	screen, ok := m.screens[m.game.Phase]
	if !ok {
		return "Unknown phase"
	}

	// Reserve 2 lines for status bar
	mainHeight := m.height - 2
	mainContent := screen.View(m.game, m.width, mainHeight)

	statusBar := components.RenderStatusBar(screen.Focus(), m.game, m.width)

	return mainContent + "\n" + statusBar
}
