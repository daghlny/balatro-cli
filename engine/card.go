package engine

import "fmt"

// Suit represents a playing card suit.
type Suit int

const (
	Spades Suit = iota
	Hearts
	Clubs
	Diamonds
)

var suitNames = [4]string{"Spades", "Hearts", "Clubs", "Diamonds"}
var suitSymbols = [4]string{"♠", "♥", "♣", "♦"}

func (s Suit) String() string  { return suitNames[s] }
func (s Suit) Symbol() string  { return suitSymbols[s] }
func (s Suit) IsRed() bool     { return s == Hearts || s == Diamonds }
func (s Suit) IsBlack() bool   { return s == Spades || s == Clubs }
func AllSuits() []Suit          { return []Suit{Spades, Hearts, Clubs, Diamonds} }

// Rank represents a playing card rank.
type Rank int

const (
	Two   Rank = 2
	Three Rank = 3
	Four  Rank = 4
	Five  Rank = 5
	Six   Rank = 6
	Seven Rank = 7
	Eight Rank = 8
	Nine  Rank = 9
	Ten   Rank = 10
	Jack  Rank = 11
	Queen Rank = 12
	King  Rank = 13
	Ace   Rank = 14
)

var rankNames = map[Rank]string{
	Two: "2", Three: "3", Four: "4", Five: "5", Six: "6",
	Seven: "7", Eight: "8", Nine: "9", Ten: "10",
	Jack: "J", Queen: "Q", King: "K", Ace: "A",
}

func (r Rank) String() string { return rankNames[r] }

// Chips returns the chip value for this rank when scoring.
func (r Rank) Chips() int {
	switch {
	case r >= Two && r <= Ten:
		return int(r)
	case r == Jack || r == Queen || r == King:
		return 10
	case r == Ace:
		return 11
	default:
		return 0
	}
}

// IsFace returns true for J, Q, K.
func (r Rank) IsFace() bool {
	return r == Jack || r == Queen || r == King
}

// IsEven returns true for 2, 4, 6, 8, 10.
func (r Rank) IsEven() bool {
	return r >= Two && r <= Ten && int(r)%2 == 0
}

// IsOdd returns true for 3, 5, 7, 9, A.
func (r Rank) IsOdd() bool {
	return r == Ace || (r >= Three && r <= Nine && int(r)%2 == 1)
}

func AllRanks() []Rank {
	return []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
}

// Card represents a single playing card.
type Card struct {
	Suit     Suit
	Rank     Rank
	Debuffed bool
}

func NewCard(rank Rank, suit Suit) Card {
	return Card{Suit: suit, Rank: rank}
}

func (c Card) IsFace() bool { return c.Rank.IsFace() }

// Chips returns the chip value, or 0 if debuffed.
func (c Card) Chips() int {
	if c.Debuffed {
		return 0
	}
	return c.Rank.Chips()
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank.String(), c.Suit.Symbol())
}

// ShortCode returns a compact code like "Ks", "Ah", "10d".
func (c Card) ShortCode() string {
	suitChar := map[Suit]string{Spades: "s", Hearts: "h", Clubs: "c", Diamonds: "d"}
	return fmt.Sprintf("%s%s", c.Rank.String(), suitChar[c.Suit])
}
