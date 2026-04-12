package jokers

import "github.com/yinuoli/balatro-cli/engine"

func init() {
	// Common
	engine.RegisterJoker(JokerDef)
	engine.RegisterJoker(GreedyJokerDef)
	engine.RegisterJoker(JollyJokerDef)
	engine.RegisterJoker(SlyJokerDef)
	engine.RegisterJoker(HalfJokerDef)
	engine.RegisterJoker(BannerDef)

	// Uncommon
	engine.RegisterJoker(FourFingersDef)
	engine.RegisterJoker(BlackboardDef)
	engine.RegisterJoker(AcrobatDef)

	// Rare
	engine.RegisterJoker(TheDuoDef)
}
