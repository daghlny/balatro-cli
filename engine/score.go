package engine

// ScoreModification represents a change to chips/mult from a joker or card effect.
type ScoreModification struct {
	PlusChips int
	PlusMult  int
	MultMult  float64 // 1.0 = no change, 1.5 = ×1.5, etc.
}

// NoMod returns a score modification with no effect.
func NoMod() ScoreModification {
	return ScoreModification{MultMult: 1.0}
}

// CardScoreDetail records how a single card contributed to the score.
type CardScoreDetail struct {
	Card      Card
	ChipsAdd  int
	Debuffed  bool
}

// JokerScoreDetail records how a single joker contributed to the score.
type JokerScoreDetail struct {
	JokerID   string
	JokerName string
	Mod       ScoreModification
}

// ScoreBreakdown provides a detailed view of how the score was calculated.
type ScoreBreakdown struct {
	HandType     HandType
	HandLevel    int
	BaseChips    int
	BaseMult     int
	CardDetails  []CardScoreDetail
	JokerDetails []JokerScoreDetail
	FinalChips   int
	FinalMult    float64
	TotalScore   int
}

// ScoringContext provides all the information needed for scoring.
type ScoringContext struct {
	HandType       HandType
	HandLevel      int
	ScoringCards   []Card
	HeldCards      []Card
	PlayedCards    []Card
	Jokers         []*JokerInstance
	RemainingHands int
	RemainingDisc  int
	Money          int
	DeckSize       int
	HandPlayCounts map[HandType]int
	IsLastHand     bool
	RoundHandTypes []HandType
	RNG            *RNG
}

// CalculateScore runs the 3-phase scoring pipeline.
func CalculateScore(ctx *ScoringContext) ScoreBreakdown {
	data := HandTypeDataTable[ctx.HandType]
	level := ctx.HandLevel
	if level < 1 {
		level = 1
	}

	// Phase 1: Base hand type values
	chips := data.BaseChips + (level-1)*data.LevelChips
	mult := float64(data.BaseMult + (level-1)*data.LevelMult)

	breakdown := ScoreBreakdown{
		HandType:  ctx.HandType,
		HandLevel: level,
		BaseChips: chips,
		BaseMult:  int(mult),
	}

	// Phase 2: Scoring cards (left to right)
	for _, card := range ctx.ScoringCards {
		detail := CardScoreDetail{Card: card, Debuffed: card.Debuffed}
		if !card.Debuffed {
			chipValue := card.Rank.Chips()
			chips += chipValue
			detail.ChipsAdd = chipValue
		}
		breakdown.CardDetails = append(breakdown.CardDetails, detail)
	}

	// Phase 3: Joker effects (left to right)
	for _, joker := range ctx.Jokers {
		if joker == nil || joker.Effect == nil {
			continue
		}

		totalMod := NoMod()

		// OnCardScore: called for each scoring card
		for _, card := range ctx.ScoringCards {
			if card.Debuffed {
				continue
			}
			mod := joker.Effect.OnCardScore(card, ctx)
			totalMod.PlusChips += mod.PlusChips
			totalMod.PlusMult += mod.PlusMult
			if mod.MultMult > 0 && mod.MultMult != 1.0 {
				if totalMod.MultMult == 0 || totalMod.MultMult == 1.0 {
					totalMod.MultMult = mod.MultMult
				} else {
					totalMod.MultMult *= mod.MultMult
				}
			}
		}

		// OnScore: called once per hand
		scoreMod := joker.Effect.OnScore(ctx)
		totalMod.PlusChips += scoreMod.PlusChips
		totalMod.PlusMult += scoreMod.PlusMult
		if scoreMod.MultMult > 0 && scoreMod.MultMult != 1.0 {
			if totalMod.MultMult == 0 || totalMod.MultMult == 1.0 {
				totalMod.MultMult = scoreMod.MultMult
			} else {
				totalMod.MultMult *= scoreMod.MultMult
			}
		}

		// Apply accumulated modification
		chips += totalMod.PlusChips
		mult += float64(totalMod.PlusMult)
		if totalMod.MultMult > 0 && totalMod.MultMult != 1.0 {
			mult *= totalMod.MultMult
		}

		breakdown.JokerDetails = append(breakdown.JokerDetails, JokerScoreDetail{
			JokerID:   joker.Def.ID,
			JokerName: joker.Def.ID,
			Mod:       totalMod,
		})
	}

	breakdown.FinalChips = chips
	breakdown.FinalMult = mult
	breakdown.TotalScore = int(float64(chips) * mult)

	return breakdown
}
