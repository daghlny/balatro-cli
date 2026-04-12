package engine

// ShopState holds the current state of the shop.
type ShopState struct {
	JokerSlots   [2]*JokerDef  // 2 joker cards for sale
	JokerCosts   [2]int        // actual prices
	BoosterTier  int           // 0=normal($4), 1=large($6), 2=super($8)
	RerollCost   int
	FreeReroll   bool          // from skip reward
	FreeJoker    bool          // from skip reward (index 0 is free)
	RNG          *RNG
}

// BoosterPack represents a buffoon booster pack.
type BoosterPack struct {
	Tier    int    // 0=normal, 1=large, 2=super
	Cost    int
	Total   int    // total jokers to show
	Choose  int    // how many to pick
	Jokers  []*JokerDef
}

var boosterTiers = []struct {
	cost   int
	total  int
	choose int
}{
	{4, 2, 1},  // normal
	{6, 4, 1},  // large
	{8, 4, 2},  // super
}

// NewShopState creates a fresh shop with random inventory.
func NewShopState(rng *RNG, freeReroll bool, freeJoker bool) *ShopState {
	s := &ShopState{
		RerollCost: 5,
		FreeReroll: freeReroll,
		FreeJoker:  freeJoker,
		RNG:        rng,
	}
	s.generateJokers()
	s.BoosterTier = rng.Intn(3)
	return s
}

func (s *ShopState) generateJokers() {
	for i := 0; i < 2; i++ {
		def := randomJokerByRarity(s.RNG)
		s.JokerSlots[i] = def
		s.JokerCosts[i] = randomCostForRarity(def.Rarity, s.RNG)
	}
	if s.FreeJoker && s.JokerSlots[0] != nil {
		s.JokerCosts[0] = 0
	}
}

// Reroll regenerates the 2 joker slots.
func (s *ShopState) Reroll(money *int) bool {
	cost := s.RerollCost
	if s.FreeReroll {
		cost = 0
		s.FreeReroll = false
	}
	if *money < cost {
		return false
	}
	*money -= cost
	s.RerollCost++
	s.generateJokers()
	return true
}

// BuyJoker purchases a joker from the shop.
func (s *ShopState) BuyJoker(index int, money *int, jokers *[]*JokerInstance, maxSlots int) (*JokerInstance, bool) {
	if index < 0 || index >= 2 || s.JokerSlots[index] == nil {
		return nil, false
	}
	if len(*jokers) >= maxSlots {
		return nil, false
	}
	cost := s.JokerCosts[index]
	if *money < cost {
		return nil, false
	}
	*money -= cost
	inst := NewJokerInstance(s.JokerSlots[index])
	inst.Cost = cost
	*jokers = append(*jokers, inst)
	s.JokerSlots[index] = nil
	return inst, true
}

// GenerateBoosterPack creates a booster pack for the current tier.
func (s *ShopState) GenerateBoosterPack() *BoosterPack {
	tier := boosterTiers[s.BoosterTier]
	pack := &BoosterPack{
		Tier:   s.BoosterTier,
		Cost:   tier.cost,
		Total:  tier.total,
		Choose: tier.choose,
	}
	for i := 0; i < tier.total; i++ {
		pack.Jokers = append(pack.Jokers, randomJokerByRarity(s.RNG))
	}
	return pack
}

// randomJokerByRarity picks a joker from the registry by rarity weights.
func randomJokerByRarity(rng *RNG) *JokerDef {
	// Weights: Common 70%, Uncommon 25%, Rare 5%
	roll := rng.Intn(100)
	var targetRarity Rarity
	switch {
	case roll < 70:
		targetRarity = Common
	case roll < 95:
		targetRarity = Uncommon
	default:
		targetRarity = Rare
	}

	// Collect all jokers of the target rarity
	var candidates []*JokerDef
	for _, def := range JokerRegistry {
		if def.Rarity == targetRarity {
			candidates = append(candidates, def)
		}
	}
	if len(candidates) == 0 {
		// Fallback: pick any joker
		for _, def := range JokerRegistry {
			candidates = append(candidates, def)
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	return candidates[rng.Intn(len(candidates))]
}

// randomCostForRarity generates a price within the range for a rarity.
func randomCostForRarity(r Rarity, rng *RNG) int {
	switch r {
	case Common:
		return 2 + rng.Intn(5) // $2-$6
	case Uncommon:
		return 5 + rng.Intn(4) // $5-$8
	case Rare:
		return 8 + rng.Intn(3) // $8-$10
	default:
		return 10
	}
}
