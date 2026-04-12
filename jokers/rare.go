package jokers

import "github.com/yinuoli/balatro-cli/engine"

// --- The Duo: ×2 mult if hand contains a pair ---
type theDuoEffect struct{ engine.BaseJokerEffect }

func (j *theDuoEffect) OnScore(ctx *engine.ScoringContext) engine.ScoreModification {
	if handContains(ctx.HandType, engine.Pair) {
		return engine.ScoreModification{MultMult: 2.0}
	}
	return engine.NoMod()
}

var TheDuoDef = &engine.JokerDef{
	ID: "the_duo", NameKey: "joker.the_duo.name", DescKey: "joker.the_duo.desc",
	Rarity: engine.Rare, Cost: 8,
	EffectFactory: func() engine.JokerEffect { return &theDuoEffect{} },
}
