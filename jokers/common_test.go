package jokers

import (
	"testing"

	"github.com/yinuoli/balatro-cli/engine"
	"github.com/yinuoli/balatro-cli/testutil"
)

func TestJoker(t *testing.T) {
	j := engine.NewJokerInstance(JokerDef)
	ctx := &engine.ScoringContext{HandType: engine.HighCard}
	mod := j.Effect.OnScore(ctx)
	if mod.PlusMult != 4 {
		t.Errorf("Joker +mult = %d, want 4", mod.PlusMult)
	}
}

func TestGreedyJoker(t *testing.T) {
	j := engine.NewJokerInstance(GreedyJokerDef)
	ctx := &engine.ScoringContext{}

	// Diamond card: +3 mult
	diamond := testutil.MakeCard("Kd")
	mod := j.Effect.OnCardScore(diamond, ctx)
	if mod.PlusMult != 3 {
		t.Errorf("Greedy on diamond = %d, want 3", mod.PlusMult)
	}

	// Non-diamond: no effect
	spade := testutil.MakeCard("Ks")
	mod = j.Effect.OnCardScore(spade, ctx)
	if mod.PlusMult != 0 {
		t.Errorf("Greedy on spade = %d, want 0", mod.PlusMult)
	}
}

func TestJollyJoker(t *testing.T) {
	j := engine.NewJokerInstance(JollyJokerDef)

	tests := []struct {
		name     string
		handType engine.HandType
		wantMult int
	}{
		{"Pair", engine.Pair, 8},
		{"TwoPair (contains pair)", engine.TwoPair, 8},
		{"FullHouse (contains pair)", engine.FullHouse, 8},
		{"FourOfAKind (contains pair)", engine.FourOfAKind, 8},
		{"HighCard (no pair)", engine.HighCard, 0},
		{"Straight (no pair)", engine.Straight, 0},
		{"Flush (no pair)", engine.Flush, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &engine.ScoringContext{HandType: tt.handType}
			mod := j.Effect.OnScore(ctx)
			if mod.PlusMult != tt.wantMult {
				t.Errorf("Jolly on %s: +mult = %d, want %d", tt.handType, mod.PlusMult, tt.wantMult)
			}
		})
	}
}

func TestSlyJoker(t *testing.T) {
	j := engine.NewJokerInstance(SlyJokerDef)
	ctx := &engine.ScoringContext{HandType: engine.Pair}
	mod := j.Effect.OnScore(ctx)
	if mod.PlusChips != 50 {
		t.Errorf("Sly on pair: +chips = %d, want 50", mod.PlusChips)
	}
	ctx.HandType = engine.Flush
	mod = j.Effect.OnScore(ctx)
	if mod.PlusChips != 0 {
		t.Errorf("Sly on flush: +chips = %d, want 0", mod.PlusChips)
	}
}

func TestHalfJoker(t *testing.T) {
	j := engine.NewJokerInstance(HalfJokerDef)

	// 3 cards: triggers
	ctx := &engine.ScoringContext{PlayedCards: testutil.MakeCards("Ks", "Kh", "5c")}
	mod := j.Effect.OnScore(ctx)
	if mod.PlusMult != 20 {
		t.Errorf("Half on 3 cards: +mult = %d, want 20", mod.PlusMult)
	}

	// 4 cards: does not trigger
	ctx.PlayedCards = testutil.MakeCards("Ks", "Kh", "5c", "2d")
	mod = j.Effect.OnScore(ctx)
	if mod.PlusMult != 0 {
		t.Errorf("Half on 4 cards: +mult = %d, want 0", mod.PlusMult)
	}
}

func TestBanner(t *testing.T) {
	j := engine.NewJokerInstance(BannerDef)
	ctx := &engine.ScoringContext{RemainingDisc: 3}
	mod := j.Effect.OnScore(ctx)
	if mod.PlusChips != 90 { // 30 * 3
		t.Errorf("Banner with 3 discards: +chips = %d, want 90", mod.PlusChips)
	}
	ctx.RemainingDisc = 0
	mod = j.Effect.OnScore(ctx)
	if mod.PlusChips != 0 {
		t.Errorf("Banner with 0 discards: +chips = %d, want 0", mod.PlusChips)
	}
}
