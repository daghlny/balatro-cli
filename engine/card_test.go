package engine

import "testing"

func TestRankChips(t *testing.T) {
	tests := []struct {
		rank Rank
		want int
	}{
		{Two, 2}, {Three, 3}, {Four, 4}, {Five, 5},
		{Six, 6}, {Seven, 7}, {Eight, 8}, {Nine, 9},
		{Ten, 10}, {Jack, 10}, {Queen, 10}, {King, 10}, {Ace, 11},
	}
	for _, tt := range tests {
		if got := tt.rank.Chips(); got != tt.want {
			t.Errorf("Rank(%d).Chips() = %d, want %d", tt.rank, got, tt.want)
		}
	}
}

func TestRankIsFace(t *testing.T) {
	faces := []Rank{Jack, Queen, King}
	nonFaces := []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Ace}
	for _, r := range faces {
		if !r.IsFace() {
			t.Errorf("Rank(%s).IsFace() = false, want true", r)
		}
	}
	for _, r := range nonFaces {
		if r.IsFace() {
			t.Errorf("Rank(%s).IsFace() = true, want false", r)
		}
	}
}

func TestRankEvenOdd(t *testing.T) {
	evens := []Rank{Two, Four, Six, Eight, Ten}
	odds := []Rank{Three, Five, Seven, Nine, Ace}

	for _, r := range evens {
		if !r.IsEven() {
			t.Errorf("Rank(%s).IsEven() = false, want true", r)
		}
	}
	for _, r := range odds {
		if !r.IsOdd() {
			t.Errorf("Rank(%s).IsOdd() = false, want true", r)
		}
	}
	// Face cards are neither even nor odd
	if Jack.IsEven() || Jack.IsOdd() {
		t.Error("Jack should be neither even nor odd")
	}
}

func TestCardDebuffedChips(t *testing.T) {
	c := Card{Suit: Spades, Rank: King, Debuffed: true}
	if c.Chips() != 0 {
		t.Errorf("debuffed card chips = %d, want 0", c.Chips())
	}
	c.Debuffed = false
	if c.Chips() != 10 {
		t.Errorf("normal King chips = %d, want 10", c.Chips())
	}
}

func TestSuitProperties(t *testing.T) {
	if !Spades.IsBlack() || Spades.IsRed() {
		t.Error("Spades should be black")
	}
	if !Hearts.IsRed() || Hearts.IsBlack() {
		t.Error("Hearts should be red")
	}
	if !Clubs.IsBlack() {
		t.Error("Clubs should be black")
	}
	if !Diamonds.IsRed() {
		t.Error("Diamonds should be red")
	}
}

func TestCardString(t *testing.T) {
	c := NewCard(King, Spades)
	if s := c.String(); s != "K♠" {
		t.Errorf("Card.String() = %q, want %q", s, "K♠")
	}
	c2 := NewCard(Ten, Hearts)
	if s := c2.String(); s != "10♥" {
		t.Errorf("Card.String() = %q, want %q", s, "10♥")
	}
}
