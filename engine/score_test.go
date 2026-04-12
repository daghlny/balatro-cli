package engine

import "testing"

func TestBasicScoring(t *testing.T) {
	// Pair of Kings, level 1, no jokers
	ctx := &ScoringContext{
		HandType:     Pair,
		HandLevel:    1,
		ScoringCards: []Card{c(King, Spades), c(King, Hearts)},
		PlayedCards:  []Card{c(King, Spades), c(King, Hearts), c(Five, Clubs)},
	}
	b := CalculateScore(ctx)
	// Base: chips=10, mult=2
	// K♠: +10 → chips=20
	// K♥: +10 → chips=30
	// Total: 30 × 2 = 60
	if b.TotalScore != 60 {
		t.Errorf("basic pair score = %d, want 60", b.TotalScore)
	}
}

func TestScoringWithLevel(t *testing.T) {
	// Pair level 3: base_chips = 10 + 2*15 = 40, base_mult = 2 + 2*1 = 4
	ctx := &ScoringContext{
		HandType:     Pair,
		HandLevel:    3,
		ScoringCards: []Card{c(King, Spades), c(King, Hearts)},
	}
	b := CalculateScore(ctx)
	// chips: 40 + 10 + 10 = 60
	// mult: 4
	// total: 60 * 4 = 240
	if b.TotalScore != 240 {
		t.Errorf("pair level 3 score = %d, want 240", b.TotalScore)
	}
}

func TestScoringDebuffedCard(t *testing.T) {
	debuffedKing := Card{Suit: Spades, Rank: King, Debuffed: true}
	normalKing := c(King, Hearts)

	ctx := &ScoringContext{
		HandType:     Pair,
		HandLevel:    1,
		ScoringCards: []Card{debuffedKing, normalKing},
	}
	b := CalculateScore(ctx)
	// Base: chips=10, mult=2
	// Debuffed K♠: +0
	// K♥: +10
	// Total: 20 × 2 = 40
	if b.TotalScore != 40 {
		t.Errorf("debuffed pair score = %d, want 40", b.TotalScore)
	}
}

func TestScoringWithSimpleJoker(t *testing.T) {
	// Joker that gives +4 mult unconditionally
	jokerDef := &JokerDef{
		ID:     "joker",
		Rarity: Common,
		Cost:   2,
		EffectFactory: func() JokerEffect {
			return &testPlusMultJoker{mult: 4}
		},
	}
	joker := NewJokerInstance(jokerDef)

	ctx := &ScoringContext{
		HandType:     Pair,
		HandLevel:    1,
		ScoringCards: []Card{c(King, Spades), c(King, Hearts)},
		Jokers:       []*JokerInstance{joker},
	}
	b := CalculateScore(ctx)
	// chips: 10 + 10 + 10 = 30
	// mult: 2 + 4 = 6
	// total: 30 × 6 = 180
	if b.TotalScore != 180 {
		t.Errorf("score with joker = %d, want 180", b.TotalScore)
	}
}

func TestScoringJokerOrder(t *testing.T) {
	// +20 mult joker, then ×2 mult joker
	plusJoker := NewJokerInstance(&JokerDef{
		ID: "plus", Rarity: Common, Cost: 2,
		EffectFactory: func() JokerEffect { return &testPlusMultJoker{mult: 20} },
	})
	multJoker := NewJokerInstance(&JokerDef{
		ID: "mult", Rarity: Common, Cost: 2,
		EffectFactory: func() JokerEffect { return &testMultMultJoker{mult: 2.0} },
	})

	// Order: plus first, then mult
	ctx := &ScoringContext{
		HandType:     HighCard,
		HandLevel:    1,
		ScoringCards: []Card{c(Ace, Spades)},
		Jokers:       []*JokerInstance{plusJoker, multJoker},
	}
	b1 := CalculateScore(ctx)
	// chips: 5 + 11 = 16, mult: 1 + 20 = 21, × 2 = 42
	// total: 16 × 42 = 672
	if b1.TotalScore != 672 {
		t.Errorf("plus-then-mult score = %d, want 672", b1.TotalScore)
	}

	// Reverse order: mult first, then plus
	ctx.Jokers = []*JokerInstance{multJoker, plusJoker}
	b2 := CalculateScore(ctx)
	// chips: 16, mult: 1 × 2 = 2, + 20 = 22
	// total: 16 × 22 = 352
	if b2.TotalScore != 352 {
		t.Errorf("mult-then-plus score = %d, want 352", b2.TotalScore)
	}

	// Verify order matters
	if b1.TotalScore <= b2.TotalScore {
		t.Error("plus-before-mult should score higher than mult-before-plus")
	}
}

func TestScoringHighCard(t *testing.T) {
	ctx := &ScoringContext{
		HandType:     HighCard,
		HandLevel:    1,
		ScoringCards: []Card{c(Ace, Spades)},
	}
	b := CalculateScore(ctx)
	// chips: 5 + 11 = 16, mult: 1
	// total: 16
	if b.TotalScore != 16 {
		t.Errorf("high card Ace score = %d, want 16", b.TotalScore)
	}
}

func TestScoringFlush(t *testing.T) {
	ctx := &ScoringContext{
		HandType:  Flush,
		HandLevel: 1,
		ScoringCards: []Card{
			c(Two, Hearts), c(Five, Hearts), c(Eight, Hearts),
			c(Jack, Hearts), c(Ace, Hearts),
		},
	}
	b := CalculateScore(ctx)
	// Base: chips=35, mult=4
	// Cards: 2+5+8+10+11 = 36
	// Total chips: 35 + 36 = 71
	// Total: 71 × 4 = 284
	if b.TotalScore != 284 {
		t.Errorf("flush score = %d, want 284", b.TotalScore)
	}
}

// --- test helper joker implementations ---

type testPlusMultJoker struct {
	BaseJokerEffect
	mult int
}

func (j *testPlusMultJoker) OnScore(ctx *ScoringContext) ScoreModification {
	return ScoreModification{PlusMult: j.mult, MultMult: 1.0}
}

type testMultMultJoker struct {
	BaseJokerEffect
	mult float64
}

func (j *testMultMultJoker) OnScore(ctx *ScoringContext) ScoreModification {
	return ScoreModification{MultMult: j.mult}
}
