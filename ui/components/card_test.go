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
	if len(lines) != 7 {
		t.Errorf("card has %d lines, want 7", len(lines))
	}
}

func TestRenderCardSelected(t *testing.T) {
	card := engine.NewCard(engine.Ace, engine.Hearts)
	rendered := RenderCard(card, true, false, false)
	lines := strings.Split(rendered, "\n")
	// Selected card: first line should be empty (lifted)
	if strings.TrimSpace(lines[0]) != "" {
		t.Errorf("selected card first line should be empty, got %q", lines[0])
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
