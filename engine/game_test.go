package engine

import "testing"

func TestNewGame(t *testing.T) {
	g := NewGame(NewSeededRNG(42))
	if g.Phase != PhaseTitle {
		t.Errorf("initial phase = %d, want PhaseTitle", g.Phase)
	}
	if g.Money != 4 {
		t.Errorf("initial money = %d, want 4", g.Money)
	}
	if g.Ante != 1 {
		t.Errorf("initial ante = %d, want 1", g.Ante)
	}
	if g.Deck.Remaining() != 52 {
		t.Errorf("deck size = %d, want 52", g.Deck.Remaining())
	}
	for _, ht := range AllHandTypes() {
		if g.HandLevels[ht] != 1 {
			t.Errorf("hand level for %s = %d, want 1", ht, g.HandLevels[ht])
		}
	}
}

func TestStartGame(t *testing.T) {
	g := NewGame(NewSeededRNG(42))
	g.StartGame()
	if g.Phase != PhaseBlindSelect {
		t.Errorf("phase after start = %d, want PhaseBlindSelect", g.Phase)
	}
}

func TestEnterBlindSmall(t *testing.T) {
	g := NewGame(NewSeededRNG(42))
	g.StartGame()
	g.EnterBlind() // Small blind, no boss
	if g.Phase != PhasePlay {
		t.Errorf("phase = %d, want PhasePlay", g.Phase)
	}
	if len(g.Hand) != 8 {
		t.Errorf("hand size = %d, want 8", len(g.Hand))
	}
	if g.RemainingHands != 4 {
		t.Errorf("remaining hands = %d, want 4", g.RemainingHands)
	}
	if g.TargetScore != 300 {
		t.Errorf("target = %d, want 300", g.TargetScore)
	}
}

func TestSkipBlind(t *testing.T) {
	g := NewGame(NewSeededRNG(42))
	g.StartGame()
	g.SkipBlind() // skip small blind
	if g.BlindIndex != 1 {
		t.Errorf("blind index = %d, want 1 (big blind)", g.BlindIndex)
	}
	if g.Phase != PhaseBlindSelect {
		t.Errorf("phase = %d, want PhaseBlindSelect", g.Phase)
	}
}

func TestPlayHandBasic(t *testing.T) {
	g := NewGame(NewSeededRNG(42))
	g.StartGame()
	g.EnterBlind()

	// Select first 2 cards
	g.ToggleSelect(0)
	g.ToggleSelect(1)

	err := g.PlayHand()
	if err != nil {
		t.Fatalf("PlayHand error: %v", err)
	}
	if g.Phase != PhasePlay {
		t.Errorf("phase = %d, want PhasePlay (stays in play)", g.Phase)
	}
	if g.RemainingHands != 3 {
		t.Errorf("remaining hands = %d, want 3", g.RemainingHands)
	}
	if g.LastBreakdown == nil {
		t.Error("LastBreakdown should not be nil")
	}
}

func TestPlayHandTooManyCards(t *testing.T) {
	g := NewGame(NewSeededRNG(42))
	g.StartGame()
	g.EnterBlind()

	// Select 6 cards — should fail
	for i := 0; i < 6; i++ {
		g.ToggleSelect(i)
	}
	err := g.PlayHand()
	if err == nil {
		t.Error("should fail with 6 cards selected")
	}
}

func TestDiscardSelected(t *testing.T) {
	g := NewGame(NewSeededRNG(42))
	g.StartGame()
	g.EnterBlind()

	g.ToggleSelect(0)
	g.ToggleSelect(1)
	err := g.DiscardSelected()
	if err != nil {
		t.Fatalf("DiscardSelected error: %v", err)
	}
	if g.RemainingDisc != 2 {
		t.Errorf("remaining discards = %d, want 2", g.RemainingDisc)
	}
	if len(g.Hand) != 8 {
		t.Errorf("hand size after discard = %d, want 8 (refilled)", len(g.Hand))
	}
}

func TestAdvanceBlindCycle(t *testing.T) {
	g := NewGame(NewSeededRNG(42))
	g.StartGame()

	// Advance through all 3 blinds of ante 1
	g.AdvanceBlind() // blind 0→1
	if g.BlindIndex != 1 || g.Ante != 1 {
		t.Errorf("after 1st advance: blind=%d ante=%d, want 1,1", g.BlindIndex, g.Ante)
	}
	g.AdvanceBlind() // blind 1→2
	if g.BlindIndex != 2 || g.Ante != 1 {
		t.Errorf("after 2nd advance: blind=%d ante=%d, want 2,1", g.BlindIndex, g.Ante)
	}
	g.AdvanceBlind() // blind 2→0, ante 1→2
	if g.BlindIndex != 0 || g.Ante != 2 {
		t.Errorf("after 3rd advance: blind=%d ante=%d, want 0,2", g.BlindIndex, g.Ante)
	}
}

func TestGameWinCondition(t *testing.T) {
	g := NewGame(NewSeededRNG(42))
	g.Ante = 8
	g.BlindIndex = 2 // boss blind
	g.AdvanceBlind() // should trigger win
	if !g.Won {
		t.Error("should have won after ante 8 boss blind")
	}
	if g.Phase != PhaseGameOver {
		t.Errorf("phase = %d, want PhaseGameOver", g.Phase)
	}
}

func TestSellJoker(t *testing.T) {
	g := NewGame(NewSeededRNG(42))
	def := &JokerDef{ID: "test", Rarity: Common, Cost: 6}
	j := NewJokerInstance(def)
	j.Cost = 6
	g.Jokers = append(g.Jokers, j)

	ok := g.SellJoker(0)
	if !ok {
		t.Error("sell should succeed")
	}
	if g.Money != 4+3 { // initial $4 + sell value floor(6/2)
		t.Errorf("money = %d, want 7", g.Money)
	}
	if len(g.Jokers) != 0 {
		t.Errorf("joker count = %d, want 0", len(g.Jokers))
	}
}
