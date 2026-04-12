package engine

import "testing"

func c(rank Rank, suit Suit) Card { return NewCard(rank, suit) }

var noOpts = HandDetectOpts{}

func TestDetectHighCard(t *testing.T) {
	cards := []Card{c(Two, Spades), c(Five, Hearts), c(Nine, Clubs)}
	r := DetectHand(cards, noOpts)
	if r.Type != HighCard {
		t.Errorf("got %s, want HighCard", r.Type)
	}
	if len(r.ScoringCards) != 1 || r.ScoringCards[0].Rank != Nine {
		t.Errorf("scoring card should be 9, got %v", r.ScoringCards)
	}
}

func TestDetectPair(t *testing.T) {
	cards := []Card{c(King, Spades), c(King, Hearts), c(Five, Clubs)}
	r := DetectHand(cards, noOpts)
	if r.Type != Pair {
		t.Errorf("got %s, want Pair", r.Type)
	}
	if len(r.ScoringCards) != 2 {
		t.Errorf("scoring cards = %d, want 2", len(r.ScoringCards))
	}
}

func TestDetectTwoPair(t *testing.T) {
	cards := []Card{c(King, Spades), c(King, Hearts), c(Five, Clubs), c(Five, Diamonds), c(Two, Spades)}
	r := DetectHand(cards, noOpts)
	if r.Type != TwoPair {
		t.Errorf("got %s, want TwoPair", r.Type)
	}
	if len(r.ScoringCards) != 4 {
		t.Errorf("scoring cards = %d, want 4", len(r.ScoringCards))
	}
}

func TestDetectThreeOfAKind(t *testing.T) {
	cards := []Card{c(Seven, Spades), c(Seven, Hearts), c(Seven, Clubs), c(Two, Diamonds)}
	r := DetectHand(cards, noOpts)
	if r.Type != ThreeOfAKind {
		t.Errorf("got %s, want ThreeOfAKind", r.Type)
	}
	if len(r.ScoringCards) != 3 {
		t.Errorf("scoring cards = %d, want 3", len(r.ScoringCards))
	}
}

func TestDetectStraight(t *testing.T) {
	cards := []Card{c(Three, Spades), c(Four, Hearts), c(Five, Clubs), c(Six, Diamonds), c(Seven, Spades)}
	r := DetectHand(cards, noOpts)
	if r.Type != Straight {
		t.Errorf("got %s, want Straight", r.Type)
	}
	if len(r.ScoringCards) != 5 {
		t.Errorf("scoring cards = %d, want 5", len(r.ScoringCards))
	}
}

func TestDetectStraightAceLow(t *testing.T) {
	cards := []Card{c(Ace, Spades), c(Two, Hearts), c(Three, Clubs), c(Four, Diamonds), c(Five, Spades)}
	r := DetectHand(cards, noOpts)
	if r.Type != Straight {
		t.Errorf("got %s, want Straight (A-low)", r.Type)
	}
}

func TestDetectStraightAceHigh(t *testing.T) {
	cards := []Card{c(Ten, Spades), c(Jack, Hearts), c(Queen, Clubs), c(King, Diamonds), c(Ace, Spades)}
	r := DetectHand(cards, noOpts)
	// This should be a Straight (not flush since different suits)
	if r.Type != Straight {
		t.Errorf("got %s, want Straight (A-high)", r.Type)
	}
}

func TestDetectFlush(t *testing.T) {
	cards := []Card{c(Two, Hearts), c(Five, Hearts), c(Eight, Hearts), c(Jack, Hearts), c(Ace, Hearts)}
	r := DetectHand(cards, noOpts)
	if r.Type != Flush {
		t.Errorf("got %s, want Flush", r.Type)
	}
}

func TestDetectFullHouse(t *testing.T) {
	cards := []Card{c(King, Spades), c(King, Hearts), c(King, Clubs), c(Two, Diamonds), c(Two, Spades)}
	r := DetectHand(cards, noOpts)
	if r.Type != FullHouse {
		t.Errorf("got %s, want FullHouse", r.Type)
	}
	if len(r.ScoringCards) != 5 {
		t.Errorf("scoring cards = %d, want 5", len(r.ScoringCards))
	}
}

func TestDetectFourOfAKind(t *testing.T) {
	cards := []Card{c(Ace, Spades), c(Ace, Hearts), c(Ace, Clubs), c(Ace, Diamonds), c(Two, Spades)}
	r := DetectHand(cards, noOpts)
	if r.Type != FourOfAKind {
		t.Errorf("got %s, want FourOfAKind", r.Type)
	}
	if len(r.ScoringCards) != 4 {
		t.Errorf("scoring cards = %d, want 4", len(r.ScoringCards))
	}
}

func TestDetectStraightFlush(t *testing.T) {
	cards := []Card{c(Five, Clubs), c(Six, Clubs), c(Seven, Clubs), c(Eight, Clubs), c(Nine, Clubs)}
	r := DetectHand(cards, noOpts)
	if r.Type != StraightFlush {
		t.Errorf("got %s, want StraightFlush", r.Type)
	}
}

func TestDetectRoyalFlush(t *testing.T) {
	cards := []Card{c(Ten, Spades), c(Jack, Spades), c(Queen, Spades), c(King, Spades), c(Ace, Spades)}
	r := DetectHand(cards, noOpts)
	if r.Type != RoyalFlush {
		t.Errorf("got %s, want RoyalFlush", r.Type)
	}
}

// Priority: FourOfAKind beats Flush even when all same suit
func TestPriorityFourOverFlush(t *testing.T) {
	cards := []Card{c(King, Spades), c(King, Hearts), c(King, Clubs), c(King, Diamonds), c(Two, Spades)}
	r := DetectHand(cards, noOpts)
	if r.Type != FourOfAKind {
		t.Errorf("got %s, want FourOfAKind (over flush)", r.Type)
	}
}

// 4 cards of same suit + 1 different = NOT a flush (need 5)
func TestNotFlushWith4(t *testing.T) {
	cards := []Card{c(Two, Hearts), c(Five, Hearts), c(Eight, Hearts), c(Jack, Hearts), c(Ace, Spades)}
	r := DetectHand(cards, noOpts)
	if r.Type == Flush {
		t.Error("4 of same suit should NOT be a flush without FourFingers")
	}
}

// FourFingers: 4 cards can form flush
func TestFourFingersFlush(t *testing.T) {
	cards := []Card{c(Two, Hearts), c(Five, Hearts), c(Eight, Hearts), c(Jack, Hearts), c(Ace, Spades)}
	r := DetectHand(cards, HandDetectOpts{FourFingers: true})
	if r.Type != Flush {
		t.Errorf("got %s, want Flush with FourFingers", r.Type)
	}
}

// FourFingers: 4 cards can form straight
func TestFourFingersStraight(t *testing.T) {
	cards := []Card{c(Three, Spades), c(Four, Hearts), c(Five, Clubs), c(Six, Diamonds), c(King, Spades)}
	r := DetectHand(cards, HandDetectOpts{FourFingers: true})
	if r.Type != Straight {
		t.Errorf("got %s, want Straight with FourFingers", r.Type)
	}
}

// SmearedJoker: red suits interchangeable
func TestSmearedJokerFlush(t *testing.T) {
	cards := []Card{c(Two, Hearts), c(Five, Diamonds), c(Eight, Hearts), c(Jack, Diamonds), c(Ace, Hearts)}
	r := DetectHand(cards, HandDetectOpts{SmearedJoker: true})
	if r.Type != Flush {
		t.Errorf("got %s, want Flush with SmearedJoker (mixed red)", r.Type)
	}
}

// Splash: all played cards score
func TestSplashAllCardsScore(t *testing.T) {
	cards := []Card{c(King, Spades), c(King, Hearts), c(Five, Clubs)}
	r := DetectHand(cards, HandDetectOpts{Splash: true})
	if r.Type != Pair {
		t.Errorf("got %s, want Pair", r.Type)
	}
	if len(r.ScoringCards) != 3 {
		t.Errorf("with Splash, scoring cards = %d, want 3 (all played)", len(r.ScoringCards))
	}
}

// Empty hand
func TestDetectEmptyHand(t *testing.T) {
	r := DetectHand([]Card{}, noOpts)
	if r.Type != HighCard {
		t.Errorf("got %s, want HighCard for empty hand", r.Type)
	}
}

// Single card
func TestDetectSingleCard(t *testing.T) {
	cards := []Card{c(Ace, Spades)}
	r := DetectHand(cards, noOpts)
	if r.Type != HighCard {
		t.Errorf("got %s, want HighCard for single card", r.Type)
	}
	if len(r.ScoringCards) != 1 {
		t.Errorf("scoring cards = %d, want 1", len(r.ScoringCards))
	}
}
