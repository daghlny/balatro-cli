package jokers

import "github.com/yinuoli/balatro-cli/engine"

// --- Joker: +4 mult unconditionally ---
type jokerEffect struct{ engine.BaseJokerEffect }

func (j *jokerEffect) OnScore(ctx *engine.ScoringContext) engine.ScoreModification {
	return engine.ScoreModification{PlusMult: 4, MultMult: 1.0}
}

var JokerDef = &engine.JokerDef{
	ID: "joker", NameKey: "joker.joker.name", DescKey: "joker.joker.desc",
	Rarity: engine.Common, Cost: 2,
	EffectFactory: func() engine.JokerEffect { return &jokerEffect{} },
}

// --- Greedy Joker: diamond cards +3 mult per scoring card ---
type greedyJokerEffect struct{ engine.BaseJokerEffect }

func (j *greedyJokerEffect) OnCardScore(card engine.Card, ctx *engine.ScoringContext) engine.ScoreModification {
	if card.Suit == engine.Diamonds {
		return engine.ScoreModification{PlusMult: 3, MultMult: 1.0}
	}
	return engine.NoMod()
}

var GreedyJokerDef = &engine.JokerDef{
	ID: "greedy_joker", NameKey: "joker.greedy_joker.name", DescKey: "joker.greedy_joker.desc",
	Rarity: engine.Common, Cost: 5,
	EffectFactory: func() engine.JokerEffect { return &greedyJokerEffect{} },
}

// --- Jolly Joker: +8 mult if hand contains a pair ---
type jollyJokerEffect struct{ engine.BaseJokerEffect }

func (j *jollyJokerEffect) OnScore(ctx *engine.ScoringContext) engine.ScoreModification {
	if handContains(ctx.HandType, engine.Pair) {
		return engine.ScoreModification{PlusMult: 8, MultMult: 1.0}
	}
	return engine.NoMod()
}

var JollyJokerDef = &engine.JokerDef{
	ID: "jolly_joker", NameKey: "joker.jolly_joker.name", DescKey: "joker.jolly_joker.desc",
	Rarity: engine.Common, Cost: 3,
	EffectFactory: func() engine.JokerEffect { return &jollyJokerEffect{} },
}

// --- Sly Joker: +50 chips if hand contains a pair ---
type slyJokerEffect struct{ engine.BaseJokerEffect }

func (j *slyJokerEffect) OnScore(ctx *engine.ScoringContext) engine.ScoreModification {
	if handContains(ctx.HandType, engine.Pair) {
		return engine.ScoreModification{PlusChips: 50, MultMult: 1.0}
	}
	return engine.NoMod()
}

var SlyJokerDef = &engine.JokerDef{
	ID: "sly_joker", NameKey: "joker.sly_joker.name", DescKey: "joker.sly_joker.desc",
	Rarity: engine.Common, Cost: 3,
	EffectFactory: func() engine.JokerEffect { return &slyJokerEffect{} },
}

// --- Half Joker: +20 mult if played ≤3 cards ---
type halfJokerEffect struct{ engine.BaseJokerEffect }

func (j *halfJokerEffect) OnScore(ctx *engine.ScoringContext) engine.ScoreModification {
	if len(ctx.PlayedCards) <= 3 {
		return engine.ScoreModification{PlusMult: 20, MultMult: 1.0}
	}
	return engine.NoMod()
}

var HalfJokerDef = &engine.JokerDef{
	ID: "half_joker", NameKey: "joker.half_joker.name", DescKey: "joker.half_joker.desc",
	Rarity: engine.Common, Cost: 5,
	EffectFactory: func() engine.JokerEffect { return &halfJokerEffect{} },
}

// --- Banner: +30 chips per remaining discard ---
type bannerEffect struct{ engine.BaseJokerEffect }

func (j *bannerEffect) OnScore(ctx *engine.ScoringContext) engine.ScoreModification {
	return engine.ScoreModification{PlusChips: 30 * ctx.RemainingDisc, MultMult: 1.0}
}

var BannerDef = &engine.JokerDef{
	ID: "banner", NameKey: "joker.banner.name", DescKey: "joker.banner.desc",
	Rarity: engine.Common, Cost: 5,
	EffectFactory: func() engine.JokerEffect { return &bannerEffect{} },
}

// handContains checks if a hand type "contains" a sub-type.
// e.g. FullHouse contains Pair and ThreeOfAKind, FourOfAKind contains Pair, etc.
func handContains(actual engine.HandType, query engine.HandType) bool {
	switch query {
	case engine.Pair:
		return actual == engine.Pair || actual == engine.TwoPair ||
			actual == engine.FullHouse || actual == engine.FourOfAKind ||
			actual == engine.ThreeOfAKind
	case engine.TwoPair:
		return actual == engine.TwoPair || actual == engine.FullHouse
	case engine.ThreeOfAKind:
		return actual == engine.ThreeOfAKind || actual == engine.FullHouse ||
			actual == engine.FourOfAKind
	case engine.Straight:
		return actual == engine.Straight || actual == engine.StraightFlush ||
			actual == engine.RoyalFlush
	case engine.Flush:
		return actual == engine.Flush || actual == engine.StraightFlush ||
			actual == engine.RoyalFlush
	case engine.FourOfAKind:
		return actual == engine.FourOfAKind
	default:
		return actual == query
	}
}
