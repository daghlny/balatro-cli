package engine

// Rarity represents joker rarity tiers.
type Rarity int

const (
	Common Rarity = iota
	Uncommon
	Rare
	Legendary
)

var rarityNames = map[Rarity]string{
	Common: "Common", Uncommon: "Uncommon", Rare: "Rare", Legendary: "Legendary",
}

func (r Rarity) String() string { return rarityNames[r] }

// PassiveEffect represents joker effects that modify game rules (not scoring).
type PassiveEffect struct {
	FourFingers  bool
	Splash       bool
	SmearedJoker bool
	HandSizeMod  int // +1 for Juggler
	DiscardMod   int // +1 for Drunkard
}

// JokerDef is the static definition of a joker type, registered in the global registry.
type JokerDef struct {
	ID            string
	NameKey       string // i18n key
	DescKey       string // i18n key
	Rarity        Rarity
	Cost          int
	Passive       PassiveEffect
	EffectFactory func() JokerEffect // nil = no scoring effect
}

// JokerEffect is the interface for joker scoring behavior.
type JokerEffect interface {
	// OnCardScore is called for each non-debuffed scoring card (Phase 2 extension).
	OnCardScore(card Card, ctx *ScoringContext) ScoreModification
	// OnScore is called once per hand after all cards scored (Phase 3).
	OnScore(ctx *ScoringContext) ScoreModification
	// Lifecycle hooks for stateful jokers.
	OnHandPlayed(ctx *ScoringContext)
	OnDiscard(ctx *ScoringContext)
	OnRoundEnd(ctx *ScoringContext)
	OnSell(ctx *ScoringContext)
}

// BaseJokerEffect provides default no-op implementations.
type BaseJokerEffect struct{}

func (b BaseJokerEffect) OnCardScore(card Card, ctx *ScoringContext) ScoreModification {
	return NoMod()
}
func (b BaseJokerEffect) OnScore(ctx *ScoringContext) ScoreModification {
	return NoMod()
}
func (b BaseJokerEffect) OnHandPlayed(ctx *ScoringContext) {}
func (b BaseJokerEffect) OnDiscard(ctx *ScoringContext)    {}
func (b BaseJokerEffect) OnRoundEnd(ctx *ScoringContext)   {}
func (b BaseJokerEffect) OnSell(ctx *ScoringContext)       {}

// JokerInstance represents a joker owned by the player during a game.
type JokerInstance struct {
	Def    *JokerDef
	Effect JokerEffect
	Cost   int  // Purchase price (sell value = floor(Cost/2))
}

// SellValue returns the money gained from selling this joker.
func (j *JokerInstance) SellValue() int {
	return j.Cost / 2
}

// NewJokerInstance creates an instance from a definition.
func NewJokerInstance(def *JokerDef) *JokerInstance {
	var effect JokerEffect
	if def.EffectFactory != nil {
		effect = def.EffectFactory()
	} else {
		effect = BaseJokerEffect{}
	}
	return &JokerInstance{
		Def:    def,
		Effect: effect,
		Cost:   def.Cost,
	}
}

// JokerRegistry is the global registry of all joker definitions.
var JokerRegistry = map[string]*JokerDef{}

// RegisterJoker adds a joker definition to the registry.
func RegisterJoker(def *JokerDef) {
	JokerRegistry[def.ID] = def
}

// GetPassiveEffects aggregates passive effects from all jokers.
func GetPassiveEffects(jokers []*JokerInstance) PassiveEffect {
	result := PassiveEffect{}
	for _, j := range jokers {
		if j == nil {
			continue
		}
		p := j.Def.Passive
		result.FourFingers = result.FourFingers || p.FourFingers
		result.Splash = result.Splash || p.Splash
		result.SmearedJoker = result.SmearedJoker || p.SmearedJoker
		result.HandSizeMod += p.HandSizeMod
		result.DiscardMod += p.DiscardMod
	}
	return result
}
