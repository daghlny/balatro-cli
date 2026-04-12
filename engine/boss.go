package engine

import "fmt"

// BossEffect defines the interface for boss blind special effects.
type BossEffect interface {
	ID() string
	NameKey() string
	DescKey() string
	Apply(state *RoundState)
	TargetMultiplier() float64
	IsDebuffed(card Card) bool
	ValidatePlay(cards []Card, state *RoundState) error
	OnPlay(cards []Card, state *RoundState)
}

// RoundState holds mutable state during a blind round, modified by boss effects.
type RoundState struct {
	RemainingHands   int
	RemainingDisc    int
	HandSize         int
	Money            *int
	PlayedHandTypes  []HandType // hand types played this round (for Mouth/Eye)
	HandLevels       map[HandType]int
}

// BaseBossEffect provides default no-op implementations.
type BaseBossEffect struct {
	BossID      string
	BossNameKey string
	BossDescKey string
}

func (b BaseBossEffect) ID() string                                     { return b.BossID }
func (b BaseBossEffect) NameKey() string                                { return b.BossNameKey }
func (b BaseBossEffect) DescKey() string                                { return b.BossDescKey }
func (b BaseBossEffect) Apply(s *RoundState)                            {}
func (b BaseBossEffect) TargetMultiplier() float64                      { return 2.0 }
func (b BaseBossEffect) IsDebuffed(c Card) bool                         { return false }
func (b BaseBossEffect) ValidatePlay(c []Card, s *RoundState) error     { return nil }
func (b BaseBossEffect) OnPlay(c []Card, s *RoundState)                 {}

// --- Suit debuff bosses ---

type suitDebuffBoss struct {
	BaseBossEffect
	suit Suit
}

func (b *suitDebuffBoss) IsDebuffed(c Card) bool { return c.Suit == b.suit }

func NewTheClub() BossEffect {
	return &suitDebuffBoss{BaseBossEffect{"the_club", "boss.the_club.name", "boss.the_club.desc"}, Clubs}
}
func NewTheGoad() BossEffect {
	return &suitDebuffBoss{BaseBossEffect{"the_goad", "boss.the_goad.name", "boss.the_goad.desc"}, Spades}
}
func NewTheHead() BossEffect {
	return &suitDebuffBoss{BaseBossEffect{"the_head", "boss.the_head.name", "boss.the_head.desc"}, Hearts}
}
func NewTheWindow() BossEffect {
	return &suitDebuffBoss{BaseBossEffect{"the_window", "boss.the_window.name", "boss.the_window.desc"}, Diamonds}
}

// --- The Plant: face cards debuffed ---
type thePlant struct{ BaseBossEffect }

func (b *thePlant) IsDebuffed(c Card) bool { return c.IsFace() }
func NewThePlant() BossEffect {
	return &thePlant{BaseBossEffect{"the_plant", "boss.the_plant.name", "boss.the_plant.desc"}}
}

// --- The Psychic: must play exactly 5 cards ---
type thePsychic struct{ BaseBossEffect }

func (b *thePsychic) ValidatePlay(cards []Card, s *RoundState) error {
	if len(cards) != 5 {
		return fmt.Errorf("must play exactly 5 cards")
	}
	return nil
}
func NewThePsychic() BossEffect {
	return &thePsychic{BaseBossEffect{"the_psychic", "boss.the_psychic.name", "boss.the_psychic.desc"}}
}

// --- The Manacle: hand size -1 ---
type theManacle struct{ BaseBossEffect }

func (b *theManacle) Apply(s *RoundState) { s.HandSize-- }
func NewTheManacle() BossEffect {
	return &theManacle{BaseBossEffect{"the_manacle", "boss.the_manacle.name", "boss.the_manacle.desc"}}
}

// --- The Flint: base chips and mult halved (handled in scoring) ---
type theFlint struct{ BaseBossEffect }

func NewTheFlint() BossEffect {
	return &theFlint{BaseBossEffect{"the_flint", "boss.the_flint.name", "boss.the_flint.desc"}}
}

// --- The Wall: 4x target ---
type theWall struct{ BaseBossEffect }

func (b *theWall) TargetMultiplier() float64 { return 4.0 }
func NewTheWall() BossEffect {
	return &theWall{BaseBossEffect{"the_wall", "boss.the_wall.name", "boss.the_wall.desc"}}
}

// --- The Needle: only 1 hand ---
type theNeedle struct{ BaseBossEffect }

func (b *theNeedle) Apply(s *RoundState) { s.RemainingHands = 1 }
func NewTheNeedle() BossEffect {
	return &theNeedle{BaseBossEffect{"the_needle", "boss.the_needle.name", "boss.the_needle.desc"}}
}

// --- The Water: no discards ---
type theWater struct{ BaseBossEffect }

func (b *theWater) Apply(s *RoundState) { s.RemainingDisc = 0 }
func NewTheWater() BossEffect {
	return &theWater{BaseBossEffect{"the_water", "boss.the_water.name", "boss.the_water.desc"}}
}

// --- The Mouth: can only play one hand type per round ---
type theMouth struct{ BaseBossEffect }

func (b *theMouth) ValidatePlay(cards []Card, s *RoundState) error {
	if len(s.PlayedHandTypes) > 0 {
		// Any hand type check happens after detection, so we use the state to track
		// The actual validation needs the detected hand type, handled in game.go
	}
	return nil
}
func NewTheMouth() BossEffect {
	return &theMouth{BaseBossEffect{"the_mouth", "boss.the_mouth.name", "boss.the_mouth.desc"}}
}

// --- The Eye: each hand must be a different type ---
type theEye struct{ BaseBossEffect }

func NewTheEye() BossEffect {
	return &theEye{BaseBossEffect{"the_eye", "boss.the_eye.name", "boss.the_eye.desc"}}
}

// --- The Arm: played hand type loses 1 level ---
type theArm struct{ BaseBossEffect }

func (b *theArm) OnPlay(cards []Card, s *RoundState) {
	// Hand type level reduction handled in game.go after detection
}
func NewTheArm() BossEffect {
	return &theArm{BaseBossEffect{"the_arm", "boss.the_arm.name", "boss.the_arm.desc"}}
}

// --- The Tooth: -$1 per card played ---
type theTooth struct{ BaseBossEffect }

func (b *theTooth) OnPlay(cards []Card, s *RoundState) {
	*s.Money -= len(cards)
	if *s.Money < 0 {
		*s.Money = 0
	}
}
func NewTheTooth() BossEffect {
	return &theTooth{BaseBossEffect{"the_tooth", "boss.the_tooth.name", "boss.the_tooth.desc"}}
}

// --- Finale bosses (Ante 8) ---

type violetVessel struct{ BaseBossEffect }

func (b *violetVessel) TargetMultiplier() float64 { return 6.0 }
func NewVioletVessel() BossEffect {
	return &violetVessel{BaseBossEffect{"violet_vessel", "boss.violet_vessel.name", "boss.violet_vessel.desc"}}
}

type crimsonHeart struct{ BaseBossEffect }

func NewCrimsonHeart() BossEffect {
	return &crimsonHeart{BaseBossEffect{"crimson_heart", "boss.crimson_heart.name", "boss.crimson_heart.desc"}}
}

type ceruleanBell struct{ BaseBossEffect }

func NewCeruleanBell() BossEffect {
	return &ceruleanBell{BaseBossEffect{"cerulean_bell", "boss.cerulean_bell.name", "boss.cerulean_bell.desc"}}
}

// AllRegularBosses returns constructors for the 15 regular boss blinds.
var AllRegularBosses = []func() BossEffect{
	NewTheClub, NewTheGoad, NewTheHead, NewTheWindow, NewThePlant,
	NewThePsychic, NewTheManacle, NewTheFlint, NewTheWall, NewTheNeedle,
	NewTheWater, NewTheMouth, NewTheEye, NewTheArm, NewTheTooth,
}

// AllFinaleBosses returns constructors for the 3 finale boss blinds.
var AllFinaleBosses = []func() BossEffect{
	NewVioletVessel, NewCrimsonHeart, NewCeruleanBell,
}

// RandomBoss picks a random boss for the given ante.
func RandomBoss(ante int, rng *RNG) BossEffect {
	if ante >= 8 {
		return AllFinaleBosses[rng.Intn(len(AllFinaleBosses))]()
	}
	return AllRegularBosses[rng.Intn(len(AllRegularBosses))]()
}
