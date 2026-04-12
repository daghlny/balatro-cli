package engine

import "math"

// BlindType represents the three blind types within an ante.
type BlindType int

const (
	SmallBlind BlindType = iota
	BigBlind
	BossBlind
)

var blindTypeNames = map[BlindType]string{
	SmallBlind: "Small Blind",
	BigBlind:   "Big Blind",
	BossBlind:  "Boss Blind",
}

func (b BlindType) String() string { return blindTypeNames[b] }

// BaseReward returns the money reward for beating this blind type.
func (b BlindType) BaseReward() int {
	switch b {
	case SmallBlind:
		return 3
	case BigBlind:
		return 4
	case BossBlind:
		return 5
	default:
		return 0
	}
}

// TargetMultiplier returns the score multiplier for this blind type.
func (b BlindType) TargetMultiplier() float64 {
	switch b {
	case SmallBlind:
		return 1.0
	case BigBlind:
		return 1.5
	case BossBlind:
		return 2.0
	default:
		return 1.0
	}
}

// AnteBaseScore returns the base score for the given ante (1-8).
var AnteBaseScore = [9]int{
	0,      // ante 0 (unused)
	300,    // ante 1
	800,    // ante 2
	2000,   // ante 3
	5000,   // ante 4
	11000,  // ante 5
	20000,  // ante 6
	35000,  // ante 7
	50000,  // ante 8
}

// CalcTargetScore calculates the target score for a given ante and blind.
// bossMultiplier overrides the default 2x for boss blinds (e.g. Wall=4x).
func CalcTargetScore(ante int, blind BlindType, bossMultiplier float64) int {
	if ante < 1 || ante > 8 {
		return 0
	}
	base := AnteBaseScore[ante]
	mult := blind.TargetMultiplier()
	if blind == BossBlind && bossMultiplier > 0 {
		mult = bossMultiplier
	}
	return int(math.Round(float64(base) * mult))
}

// CalcReward calculates money earned after beating a blind.
func CalcReward(blind BlindType, remainingHands int, currentMoney int) int {
	base := blind.BaseReward()
	handBonus := remainingHands
	interest := currentMoney / 5
	if interest > 5 {
		interest = 5
	}
	return base + handBonus + interest
}

// SkipRewardType represents the reward for skipping a blind.
type SkipRewardType int

const (
	SkipRewardMoney    SkipRewardType = iota // +$5
	SkipRewardReroll                         // Next shop free reroll
	SkipRewardFreeJoker                      // Next shop has 1 free joker
	SkipRewardDouble                         // Double money (cap $20)
	SkipRewardHandSize                       // Next round +3 hand size
)

// SkipReward holds a skip reward and its value.
type SkipReward struct {
	Type SkipRewardType
}

// RandomSkipReward generates a random skip reward.
func RandomSkipReward(rng *RNG) *SkipReward {
	return &SkipReward{Type: SkipRewardType(rng.Intn(5))}
}

// ApplySkipReward applies the skip reward to game state.
func ApplySkipReward(reward *SkipReward, money *int, nextShopFreeReroll *bool, nextShopFreeJoker *bool, nextRoundHandSizeBonus *int) {
	switch reward.Type {
	case SkipRewardMoney:
		*money += 5
	case SkipRewardReroll:
		*nextShopFreeReroll = true
	case SkipRewardFreeJoker:
		*nextShopFreeJoker = true
	case SkipRewardDouble:
		doubled := *money * 2
		if doubled > 20 {
			doubled = 20
		}
		if doubled > *money {
			*money = doubled
		}
	case SkipRewardHandSize:
		*nextRoundHandSizeBonus += 3
	}
}
