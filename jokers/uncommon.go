package jokers

import "github.com/yinuoli/balatro-cli/engine"

// --- Four Fingers: passive only, 4 cards form straights/flushes ---
var FourFingersDef = &engine.JokerDef{
	ID: "four_fingers", NameKey: "joker.four_fingers.name", DescKey: "joker.four_fingers.desc",
	Rarity: engine.Uncommon, Cost: 7,
	Passive: engine.PassiveEffect{FourFingers: true},
}

// --- Blackboard: ×3 mult if all held cards are spades or clubs ---
type blackboardEffect struct{ engine.BaseJokerEffect }

func (j *blackboardEffect) OnScore(ctx *engine.ScoringContext) engine.ScoreModification {
	if len(ctx.HeldCards) == 0 {
		return engine.NoMod()
	}
	allBlack := true
	for _, card := range ctx.HeldCards {
		if !card.Suit.IsBlack() {
			allBlack = false
			break
		}
	}
	if allBlack {
		return engine.ScoreModification{MultMult: 3.0}
	}
	return engine.NoMod()
}

var BlackboardDef = &engine.JokerDef{
	ID: "blackboard", NameKey: "joker.blackboard.name", DescKey: "joker.blackboard.desc",
	Rarity: engine.Uncommon, Cost: 6,
	EffectFactory: func() engine.JokerEffect { return &blackboardEffect{} },
}

// --- Acrobat: ×3 mult on the last hand of a round ---
type acrobatEffect struct{ engine.BaseJokerEffect }

func (j *acrobatEffect) OnScore(ctx *engine.ScoringContext) engine.ScoreModification {
	if ctx.IsLastHand {
		return engine.ScoreModification{MultMult: 3.0}
	}
	return engine.NoMod()
}

var AcrobatDef = &engine.JokerDef{
	ID: "acrobat", NameKey: "joker.acrobat.name", DescKey: "joker.acrobat.desc",
	Rarity: engine.Uncommon, Cost: 6,
	EffectFactory: func() engine.JokerEffect { return &acrobatEffect{} },
}
