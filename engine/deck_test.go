package engine

import "testing"

func TestNewStandardDeck(t *testing.T) {
	rng := NewSeededRNG(42)
	d := NewStandardDeck(rng)
	if d.Remaining() != 52 {
		t.Errorf("standard deck size = %d, want 52", d.Remaining())
	}
}

func TestDeckDraw(t *testing.T) {
	rng := NewSeededRNG(42)
	d := NewStandardDeck(rng)
	cards := d.Draw(5)
	if len(cards) != 5 {
		t.Errorf("drew %d cards, want 5", len(cards))
	}
	if d.Remaining() != 47 {
		t.Errorf("remaining = %d, want 47", d.Remaining())
	}
}

func TestDeckDrawAll(t *testing.T) {
	rng := NewSeededRNG(42)
	d := NewStandardDeck(rng)
	cards := d.Draw(52)
	if len(cards) != 52 {
		t.Errorf("drew %d cards, want 52", len(cards))
	}
	if d.Remaining() != 0 {
		t.Errorf("remaining = %d, want 0", d.Remaining())
	}
	// Drawing from empty deck should return nothing
	cards = d.Draw(1)
	if len(cards) != 0 {
		t.Errorf("drew %d from empty deck, want 0", len(cards))
	}
}

func TestDeckDiscardAndRedraw(t *testing.T) {
	rng := NewSeededRNG(42)
	d := NewStandardDeck(rng)
	hand := d.Draw(52) // draw all
	d.Discard(hand)    // discard all back

	// Now draw again — should reshuffle discard pile
	redrawn := d.Draw(5)
	if len(redrawn) != 5 {
		t.Errorf("redrew %d cards, want 5", len(redrawn))
	}
}

func TestDeckShuffle(t *testing.T) {
	rng1 := NewSeededRNG(42)
	rng2 := NewSeededRNG(99)
	d1 := NewStandardDeck(rng1)
	d2 := NewStandardDeck(rng2)
	d1.Shuffle()
	d2.Shuffle()

	// With different seeds, shuffled order should differ
	cards1 := d1.Draw(5)
	cards2 := d2.Draw(5)
	same := true
	for i := range cards1 {
		if cards1[i] != cards2[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("different seeds produced same shuffle order")
	}
}

func TestDeckTotalRemaining(t *testing.T) {
	rng := NewSeededRNG(42)
	d := NewStandardDeck(rng)
	hand := d.Draw(8)
	d.Discard(hand[:3])
	if d.TotalRemaining() != 47 { // 44 in draw + 3 in discard
		t.Errorf("total remaining = %d, want 47", d.TotalRemaining())
	}
}
