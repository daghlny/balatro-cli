package testutil

import (
	"strings"

	"github.com/yinuoli/balatro-cli/engine"
)

// MakeCard creates a card from a shortcode like "Ks", "Ah", "10d", "2c".
// Format: <rank><suit> where suit is s/h/c/d.
func MakeCard(code string) engine.Card {
	if len(code) < 2 {
		panic("invalid card code: " + code)
	}

	suitChar := code[len(code)-1:]
	rankStr := code[:len(code)-1]

	suit := parseSuit(suitChar)
	rank := parseRank(rankStr)

	return engine.NewCard(rank, suit)
}

// MakeCards creates multiple cards from shortcodes.
func MakeCards(codes ...string) []engine.Card {
	cards := make([]engine.Card, len(codes))
	for i, code := range codes {
		cards[i] = MakeCard(code)
	}
	return cards
}

// NewTestRNG creates a deterministic RNG for testing.
func NewTestRNG(seed int64) *engine.RNG {
	return engine.NewSeededRNG(seed)
}

func parseSuit(s string) engine.Suit {
	switch strings.ToLower(s) {
	case "s":
		return engine.Spades
	case "h":
		return engine.Hearts
	case "c":
		return engine.Clubs
	case "d":
		return engine.Diamonds
	default:
		panic("invalid suit: " + s)
	}
}

func parseRank(s string) engine.Rank {
	switch strings.ToUpper(s) {
	case "2":
		return engine.Two
	case "3":
		return engine.Three
	case "4":
		return engine.Four
	case "5":
		return engine.Five
	case "6":
		return engine.Six
	case "7":
		return engine.Seven
	case "8":
		return engine.Eight
	case "9":
		return engine.Nine
	case "10":
		return engine.Ten
	case "J":
		return engine.Jack
	case "Q":
		return engine.Queen
	case "K":
		return engine.King
	case "A":
		return engine.Ace
	default:
		panic("invalid rank: " + s)
	}
}
