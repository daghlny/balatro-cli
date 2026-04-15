package components

import (
	"strings"
	"testing"

	"github.com/yinuoli/balatro-cli/engine"
)

func TestRenderCardHasCorrectLineCount(t *testing.T) {
	card := engine.NewCard(engine.King, engine.Spades)
	rendered := RenderCard(card, false, false, false)
	lines := strings.Split(rendered, "\n")
	want := CardHeight
	if len(lines) != want {
		t.Errorf("card has %d lines, want %d", len(lines), want)
	}
}

// Selection no longer changes the card's rendered shape — the "lift" is
// applied row-wise in RenderHand so selected cards keep every border.
func TestRenderCardSelectedKeepsFullBorders(t *testing.T) {
	card := engine.NewCard(engine.Ace, engine.Hearts)
	rendered := RenderCard(card, true, false, false)
	lines := strings.Split(rendered, "\n")
	if len(lines) != CardHeight {
		t.Fatalf("selected card has %d lines, want %d", len(lines), CardHeight)
	}
	if strings.TrimSpace(lines[0]) == "" {
		t.Error("selected card should still have a top border on line 0")
	}
	if strings.TrimSpace(lines[len(lines)-1]) == "" {
		t.Error("selected card should still have a bottom border on the last line")
	}
}

// When any card in the hand is selected, the row is top-aligned and
// unselected cards get one extra row prepended — so the overall row is one
// row taller than a single card. Selected cards visibly rise above.
func TestRenderHandLiftsSelectedCards(t *testing.T) {
	cards := []engine.Card{
		engine.NewCard(engine.King, engine.Spades),
		engine.NewCard(engine.Queen, engine.Hearts),
	}
	selected := []bool{false, true}
	debuffed := []bool{false, false}
	rendered := RenderHand(cards, selected, -1, debuffed)
	lines := strings.Split(rendered, "\n")
	if len(lines) != CardHeight+1 {
		t.Errorf("hand with a selection has %d lines, want %d", len(lines), CardHeight+1)
	}
}

func TestRenderHand(t *testing.T) {
	cards := []engine.Card{
		engine.NewCard(engine.King, engine.Spades),
		engine.NewCard(engine.Queen, engine.Hearts),
		engine.NewCard(engine.Five, engine.Clubs),
	}
	selected := []bool{false, true, false}
	debuffed := []bool{false, false, false}
	rendered := RenderHand(cards, selected, 0, debuffed)
	if rendered == "" {
		t.Error("RenderHand returned empty string")
	}
}

func TestRenderCardBack(t *testing.T) {
	rendered := RenderCardBack()
	if rendered == "" {
		t.Error("RenderCardBack returned empty string")
	}
	if !strings.Contains(rendered, "░") {
		t.Error("card back should contain ░ characters")
	}
}
