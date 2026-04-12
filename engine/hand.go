package engine

import "sort"

// HandType represents a poker hand type, ordered by strength.
type HandType int

const (
	HighCard HandType = iota
	Pair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

var handTypeNames = map[HandType]string{
	HighCard:      "High Card",
	Pair:          "Pair",
	TwoPair:       "Two Pair",
	ThreeOfAKind:  "Three of a Kind",
	Straight:      "Straight",
	Flush:         "Flush",
	FullHouse:     "Full House",
	FourOfAKind:   "Four of a Kind",
	StraightFlush: "Straight Flush",
	RoyalFlush:    "Royal Flush",
}

func (h HandType) String() string { return handTypeNames[h] }

// AllHandTypes returns all hand types ordered by strength (lowest first).
func AllHandTypes() []HandType {
	return []HandType{
		HighCard, Pair, TwoPair, ThreeOfAKind, Straight, Flush,
		FullHouse, FourOfAKind, StraightFlush, RoyalFlush,
	}
}

// HandTypeData holds the base values and level-up increments for a hand type.
type HandTypeData struct {
	BaseChips  int
	BaseMult   int
	LevelChips int
	LevelMult  int
}

var HandTypeDataTable = map[HandType]HandTypeData{
	HighCard:      {5, 1, 10, 1},
	Pair:          {10, 2, 15, 1},
	TwoPair:       {20, 2, 20, 1},
	ThreeOfAKind:  {30, 3, 20, 2},
	Straight:      {30, 4, 30, 3},
	Flush:         {35, 4, 15, 2},
	FullHouse:     {40, 4, 25, 2},
	FourOfAKind:   {60, 7, 30, 3},
	StraightFlush: {100, 8, 40, 4},
	RoyalFlush:    {100, 8, 40, 4},
}

// HandResult contains the result of hand detection.
type HandResult struct {
	Type         HandType
	ScoringCards []Card // Cards that participate in scoring
	AllPlayed    []Card // All played cards
}

// HandDetectOpts contains passive joker effects that alter hand detection.
type HandDetectOpts struct {
	FourFingers  bool // 4 cards can form straights/flushes
	SmearedJoker bool // Red suits are interchangeable, black suits are interchangeable
	Splash       bool // All played cards participate in scoring
}

// DetectHand identifies the best poker hand from the played cards.
func DetectHand(cards []Card, opts HandDetectOpts) HandResult {
	result := HandResult{AllPlayed: cards}

	if len(cards) == 0 {
		result.Type = HighCard
		return result
	}

	flushReq := 5
	straightReq := 5
	if opts.FourFingers {
		flushReq = 4
		straightReq = 4
	}

	rankGroups := groupByRank(cards)
	isFlush, flushCards := checkFlush(cards, flushReq, opts.SmearedJoker)
	isStraight, straightCards := checkStraight(cards, straightReq)

	// Check from highest to lowest priority
	if isFlush && isStraight {
		sfCards := findStraightFlush(cards, straightReq, opts.SmearedJoker)
		if sfCards != nil {
			if isRoyal(sfCards) {
				result.Type = RoyalFlush
			} else {
				result.Type = StraightFlush
			}
			result.ScoringCards = sfCards
			if opts.Splash {
				result.ScoringCards = cards
			}
			return result
		}
	}

	if fourCards := findNOfAKind(rankGroups, 4); fourCards != nil {
		result.Type = FourOfAKind
		result.ScoringCards = fourCards
		if opts.Splash {
			result.ScoringCards = cards
		}
		return result
	}

	threeCards := findNOfAKind(rankGroups, 3)
	pairs := findAllPairs(rankGroups)

	if threeCards != nil && len(pairs) > 0 {
		// Make sure the pair is different from the three of a kind
		for _, pair := range pairs {
			if pair[0].Rank != threeCards[0].Rank {
				result.Type = FullHouse
				result.ScoringCards = append(threeCards, pair...)
				if opts.Splash {
					result.ScoringCards = cards
				}
				return result
			}
		}
	}

	if isFlush {
		result.Type = Flush
		result.ScoringCards = flushCards
		if opts.Splash {
			result.ScoringCards = cards
		}
		return result
	}

	if isStraight {
		result.Type = Straight
		result.ScoringCards = straightCards
		if opts.Splash {
			result.ScoringCards = cards
		}
		return result
	}

	if threeCards != nil {
		result.Type = ThreeOfAKind
		result.ScoringCards = threeCards
		if opts.Splash {
			result.ScoringCards = cards
		}
		return result
	}

	if len(pairs) >= 2 {
		result.Type = TwoPair
		result.ScoringCards = append(pairs[0], pairs[1]...)
		if opts.Splash {
			result.ScoringCards = cards
		}
		return result
	}

	if len(pairs) == 1 {
		result.Type = Pair
		result.ScoringCards = pairs[0]
		if opts.Splash {
			result.ScoringCards = cards
		}
		return result
	}

	result.Type = HighCard
	// Highest card scores
	best := cards[0]
	for _, c := range cards[1:] {
		if c.Rank > best.Rank {
			best = c
		}
	}
	result.ScoringCards = []Card{best}
	if opts.Splash {
		result.ScoringCards = cards
	}
	return result
}

// groupByRank groups cards by their rank.
func groupByRank(cards []Card) map[Rank][]Card {
	groups := make(map[Rank][]Card)
	for _, c := range cards {
		groups[c.Rank] = append(groups[c.Rank], c)
	}
	return groups
}

// checkFlush checks if there are enough cards of the same suit.
func checkFlush(cards []Card, required int, smeared bool) (bool, []Card) {
	if smeared {
		// Red suits (Hearts, Diamonds) are interchangeable
		// Black suits (Spades, Clubs) are interchangeable
		var red, black []Card
		for _, c := range cards {
			if c.Suit.IsRed() {
				red = append(red, c)
			} else {
				black = append(black, c)
			}
		}
		if len(red) >= required {
			return true, red[:required]
		}
		if len(black) >= required {
			return true, black[:required]
		}
		return false, nil
	}

	suitGroups := make(map[Suit][]Card)
	for _, c := range cards {
		suitGroups[c.Suit] = append(suitGroups[c.Suit], c)
	}
	for _, group := range suitGroups {
		if len(group) >= required {
			return true, group[:required]
		}
	}
	return false, nil
}

// checkStraight checks if cards form a straight.
func checkStraight(cards []Card, required int) (bool, []Card) {
	if len(cards) < required {
		return false, nil
	}

	// Get unique ranks sorted
	rankSet := make(map[Rank]Card)
	for _, c := range cards {
		if _, exists := rankSet[c.Rank]; !exists {
			rankSet[c.Rank] = c
		}
	}

	ranks := make([]Rank, 0, len(rankSet))
	for r := range rankSet {
		ranks = append(ranks, r)
	}
	sort.Slice(ranks, func(i, j int) bool { return ranks[i] < ranks[j] })

	// Check for A-low straight (A-2-3-4-5)
	if required <= 5 {
		aLow := []Rank{Ace, Two, Three, Four, Five}
		aLowCards := findStraightRun(aLow[:required], rankSet)
		if aLowCards != nil {
			return true, aLowCards
		}
	}

	// Check for consecutive runs
	for i := len(ranks) - 1; i >= required-1; i-- {
		consecutive := true
		for j := 0; j < required-1; j++ {
			if ranks[i-j]-ranks[i-j-1] != 1 {
				consecutive = false
				break
			}
		}
		if consecutive {
			result := make([]Card, required)
			for j := 0; j < required; j++ {
				result[j] = rankSet[ranks[i-j]]
			}
			// Reverse so it's low to high
			for l, r := 0, len(result)-1; l < r; l, r = l+1, r-1 {
				result[l], result[r] = result[r], result[l]
			}
			return true, result
		}
	}

	return false, nil
}

// findStraightRun checks if all ranks in the run exist in rankSet.
func findStraightRun(run []Rank, rankSet map[Rank]Card) []Card {
	result := make([]Card, 0, len(run))
	for _, r := range run {
		c, ok := rankSet[r]
		if !ok {
			return nil
		}
		result = append(result, c)
	}
	return result
}

// findStraightFlush finds cards that form both a straight and flush.
func findStraightFlush(cards []Card, straightReq int, smeared bool) []Card {
	// Group by effective suit
	type suitGroup struct {
		cards []Card
	}
	groups := make(map[int][]Card) // key: effective suit group ID

	for _, c := range cards {
		key := int(c.Suit)
		if smeared {
			if c.Suit.IsRed() {
				key = 100 // all red
			} else {
				key = 200 // all black
			}
		}
		groups[key] = append(groups[key], c)
	}

	for _, group := range groups {
		if len(group) >= straightReq {
			ok, sfCards := checkStraight(group, straightReq)
			if ok {
				return sfCards
			}
		}
	}
	return nil
}

// isRoyal checks if the straight flush is a royal flush (10-J-Q-K-A).
func isRoyal(cards []Card) bool {
	ranks := make(map[Rank]bool)
	for _, c := range cards {
		ranks[c.Rank] = true
	}
	return ranks[Ten] && ranks[Jack] && ranks[Queen] && ranks[King] && ranks[Ace]
}

// findNOfAKind finds n cards of the same rank (returns highest).
func findNOfAKind(groups map[Rank][]Card, n int) []Card {
	var best []Card
	var bestRank Rank
	for rank, group := range groups {
		if len(group) >= n && (best == nil || rank > bestRank) {
			best = group[:n]
			bestRank = rank
		}
	}
	return best
}

// findAllPairs returns all pairs found, sorted by rank descending.
func findAllPairs(groups map[Rank][]Card) [][]Card {
	var pairs [][]Card
	for _, group := range groups {
		if len(group) >= 2 {
			pairs = append(pairs, group[:2])
		}
	}
	// Sort by rank descending
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i][0].Rank > pairs[j][0].Rank
	})
	return pairs
}
