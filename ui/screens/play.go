package screens

import (
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yinuoli/balatro-cli/engine"
	"github.com/yinuoli/balatro-cli/i18n"
	"github.com/yinuoli/balatro-cli/ui/components"
)

type PlayScreen struct {
	cursor    int
	focusArea int // 0=hand, 1=jokers
	errMsg    string
	// Table display
	tableCards []engine.Card
	tableBD    *engine.ScoreBreakdown // last play breakdown
}

func NewPlayScreen() *PlayScreen {
	return &PlayScreen{}
}

func (s *PlayScreen) Update(msg tea.Msg, game *engine.GameState) tea.Cmd {
	s.errMsg = ""

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "left", "h":
			if s.cursor > 0 {
				s.cursor--
			}
		case "right", "l":
			maxIdx := len(game.Hand) - 1
			if s.focusArea == 1 {
				maxIdx = game.MaxJokerSlots - 1
			}
			if s.cursor < maxIdx {
				s.cursor++
			}
		case "tab":
			s.focusArea = 1 - s.focusArea
			s.cursor = 0
		case " ":
			if s.focusArea == 0 {
				game.ToggleSelect(s.cursor)
			}
		case "enter":
			if s.focusArea == 0 {
				s.doPlay(game)
			}
		case "d":
			if s.focusArea == 0 {
				err := game.DiscardSelected()
				s.tableCards = nil
				s.tableBD = nil
				if err != nil {
					s.errMsg = err.Error()
				}
				if s.cursor >= len(game.Hand) {
					s.cursor = maxInt(0, len(game.Hand)-1)
				}
			}
		case "s":
			sortHandByRank(game)
			s.cursor = clampCursor(s.cursor, game, s.focusArea)
		case "S":
			sortHandBySuit(game)
			s.cursor = clampCursor(s.cursor, game, s.focusArea)
		case "q", "ctrl+c":
			return tea.Quit
		}
	}
	return nil
}

// doPlay plays selected cards, scores inline, and handles round outcome.
func (s *PlayScreen) doPlay(game *engine.GameState) {
	selected := game.SelectedCardsList()
	s.tableCards = make([]engine.Card, len(selected))
	copy(s.tableCards, selected)

	err := game.PlayHand()
	if err != nil {
		s.errMsg = err.Error()
		s.tableCards = nil
		s.tableBD = nil
		return
	}

	s.tableBD = game.LastBreakdown

	// Check round outcome inline
	if game.RoundScore >= game.TargetScore {
		// Won this blind — transition to settlement
		game.Phase = engine.PhaseSettlement
	} else if game.RemainingHands <= 0 {
		// Lost
		game.Lost = true
		game.Phase = engine.PhaseGameOver
	}
	// Otherwise stay in PhasePlay — player continues
}

func (s *PlayScreen) View(game *engine.GameState, width, height int) string {
	// The PlayScreen instance is reused across rounds (see ui/app.go). When a
	// new round starts, the engine clears LastBreakdown — mirror that here so
	// the middle table doesn't keep showing the previous round's cards/score.
	if game.LastBreakdown == nil {
		s.tableCards = nil
		s.tableBD = nil
	}

	sidebar := components.RenderSidebar(game, height-2)
	sidebarWidth := 28 // sidebar style width(26) + border
	rightWidth := width - sidebarWidth

	// Top: Joker slots
	jokerFocusIdx := -1
	if s.focusArea == 1 {
		jokerFocusIdx = s.cursor
	}
	jokerArea := components.RenderJokerSlots(game.Jokers, game.MaxJokerSlots, jokerFocusIdx)
	jokerArea = lipgloss.NewStyle().Width(rightWidth).Align(lipgloss.Center).Render(jokerArea)

	// Middle: Table area with score display
	tableArea := s.renderTable(game, rightWidth)

	// Bottom: Hand
	handCursorIdx := -1
	if s.focusArea == 0 {
		handCursorIdx = s.cursor
	}
	debuffed := make([]bool, len(game.Hand))
	for i, c := range game.Hand {
		debuffed[i] = c.Debuffed
	}
	handArea := components.RenderHand(game.Hand, game.SelectedCards, handCursorIdx, debuffed)
	handArea = lipgloss.NewStyle().Width(rightWidth).Align(lipgloss.Center).Render(handArea)

	// Error message
	if s.errMsg != "" {
		errLine := lipgloss.NewStyle().Foreground(components.ColorRed).Render("  ⚠ " + s.errMsg)
		handArea = lipgloss.JoinVertical(lipgloss.Left, handArea, errLine)
	}

	// Selection count
	selCount := 0
	for _, sel := range game.SelectedCards {
		if sel {
			selCount++
		}
	}
	if selCount > 0 {
		selInfo := lipgloss.NewStyle().Foreground(components.ColorGold).
			Render(fmt.Sprintf("  已选 %d 张  (Enter 出牌 | D 弃牌)", selCount))
		handArea = lipgloss.JoinVertical(lipgloss.Left, handArea, selInfo)
	}

	return components.ComposePlayLayout(sidebar, jokerArea, tableArea, handArea, width, height)
}

// renderTable renders the middle table area: played cards + Chips × Mult display.
func (s *PlayScreen) renderTable(game *engine.GameState, rightWidth int) string {
	if len(s.tableCards) == 0 && s.tableBD == nil {
		return lipgloss.NewStyle().Width(rightWidth).Align(lipgloss.Center).
			Height(16).Foreground(components.ColorDim).
			Render("\n\n\n\n    空格选牌  Enter 出牌")
	}

	// Played cards row
	debuffed := make([]bool, len(s.tableCards))
	for i, c := range s.tableCards {
		debuffed[i] = c.Debuffed
	}
	cardsRendered := components.RenderHand(s.tableCards, nil, -1, debuffed)

	if s.tableBD == nil {
		return lipgloss.NewStyle().Width(rightWidth).Align(lipgloss.Center).
			Height(16).Render(cardsRendered)
	}

	bd := s.tableBD
	htKey := handTypeKeyFromType(bd.HandType)
	htName := i18n.T(htKey)

	// Hand type label
	handLabel := lipgloss.NewStyle().Foreground(components.ColorCyan).Bold(true).
		Render(fmt.Sprintf("%s Lv.%d", htName, bd.HandLevel))

	// Big Chips × Mult display
	scoreDisplay := components.RenderScoreDisplay(bd.FinalChips, bd.FinalMult)

	// Total score
	totalStr := fmt.Sprintf("= %s",
		lipgloss.NewStyle().Foreground(components.ColorGold).Bold(true).
			Render(components.FormatNumber(bd.TotalScore)))

	// Compose: cards on left, score info on right
	scoreInfo := lipgloss.JoinVertical(lipgloss.Center, handLabel, "", scoreDisplay, totalStr)
	tableContent := lipgloss.JoinHorizontal(lipgloss.Center, cardsRendered, "   ", scoreInfo)

	return lipgloss.NewStyle().Width(rightWidth).Align(lipgloss.Center).
		Height(16).Render(tableContent)
}

func (s *PlayScreen) Focus() components.FocusTarget {
	switch s.focusArea {
	case 1:
		return components.FocusTarget{Type: components.FocusJokerSlot, Index: s.cursor}
	default:
		return components.FocusTarget{Type: components.FocusHandCard, Index: s.cursor}
	}
}

func handTypeKeyFromType(ht engine.HandType) string {
	keys := map[engine.HandType]string{
		engine.HighCard: "hand.high_card", engine.Pair: "hand.pair",
		engine.TwoPair: "hand.two_pair", engine.ThreeOfAKind: "hand.three_of_kind",
		engine.Straight: "hand.straight", engine.Flush: "hand.flush",
		engine.FullHouse: "hand.full_house", engine.FourOfAKind: "hand.four_of_kind",
		engine.StraightFlush: "hand.straight_flush", engine.RoyalFlush: "hand.royal_flush",
	}
	return keys[ht]
}

// sortHandByRank sorts the hand by rank descending (ties broken by suit asc).
// Selections travel with their cards so the player doesn't lose in-progress picks.
func sortHandByRank(game *engine.GameState) {
	sortHandWith(game, func(a, b engine.Card) bool {
		if a.Rank != b.Rank {
			return a.Rank > b.Rank
		}
		return a.Suit < b.Suit
	})
}

// sortHandBySuit groups by suit (Spades → Hearts → Clubs → Diamonds),
// then orders each group by rank descending.
func sortHandBySuit(game *engine.GameState) {
	sortHandWith(game, func(a, b engine.Card) bool {
		if a.Suit != b.Suit {
			return a.Suit < b.Suit
		}
		return a.Rank > b.Rank
	})
}

// sortHandWith reorders the hand using the given comparator while keeping
// each card's selected flag attached to the card (not the slot).
func sortHandWith(game *engine.GameState, less func(a, b engine.Card) bool) {
	n := len(game.Hand)
	if n == 0 {
		return
	}
	type entry struct {
		card engine.Card
		sel  bool
	}
	entries := make([]entry, n)
	for i, c := range game.Hand {
		var sel bool
		if i < len(game.SelectedCards) {
			sel = game.SelectedCards[i]
		}
		entries[i] = entry{c, sel}
	}
	sort.SliceStable(entries, func(i, j int) bool {
		return less(entries[i].card, entries[j].card)
	})
	if len(game.SelectedCards) != n {
		game.SelectedCards = make([]bool, n)
	}
	for i, e := range entries {
		game.Hand[i] = e.card
		game.SelectedCards[i] = e.sel
	}
}

// clampCursor keeps the cursor within bounds of whichever area is focused.
func clampCursor(cursor int, game *engine.GameState, focusArea int) int {
	maxIdx := len(game.Hand) - 1
	if focusArea == 1 {
		maxIdx = game.MaxJokerSlots - 1
	}
	if cursor < 0 {
		return 0
	}
	if cursor > maxIdx {
		if maxIdx < 0 {
			return 0
		}
		return maxIdx
	}
	return cursor
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
