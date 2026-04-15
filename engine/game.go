package engine

// Phase represents the current game phase (screen).
type Phase int

const (
	PhaseTitle Phase = iota
	PhaseBlindSelect
	PhaseBossInfo
	PhasePlay
	PhaseScoring
	PhaseSettlement
	PhaseShop
	PhaseBooster
	PhaseGameOver
)

// GameState holds the complete state of a game in progress.
type GameState struct {
	Phase          Phase
	Ante           int // 1-8
	BlindIndex     int // 0=Small, 1=Big, 2=Boss
	Money          int
	Deck           *Deck
	Hand           []Card
	SelectedCards  []bool // parallel to Hand, true = selected
	Jokers         []*JokerInstance
	MaxJokerSlots  int
	HandSize       int
	BaseHandSize   int
	RemainingHands int
	RemainingDisc  int
	RoundScore     int
	TargetScore    int
	CurrentBoss    BossEffect
	HandLevels     map[HandType]int
	HandPlayCounts map[HandType]int
	RoundHandTypes []HandType // hand types played this round

	// Score breakdown from last play (for scoring screen)
	LastBreakdown *ScoreBreakdown

	// Skip rewards state
	NextShopFreeReroll    bool
	NextShopFreeJoker     bool
	NextRoundHandSizeBonus int
	LastSkipReward        *SkipReward

	// Shop
	Shop *ShopState

	// Win/loss
	Won  bool
	Lost bool

	RNG *RNG
}

// NewGame creates a new game with default settings.
func NewGame(rng *RNG) *GameState {
	g := &GameState{
		Phase:          PhaseTitle,
		Ante:           1,
		BlindIndex:     0,
		Money:          4,
		MaxJokerSlots:  5,
		HandSize:       8,
		BaseHandSize:   8,
		HandLevels:     make(map[HandType]int),
		HandPlayCounts: make(map[HandType]int),
		RNG:            rng,
	}
	for _, ht := range AllHandTypes() {
		g.HandLevels[ht] = 1
		g.HandPlayCounts[ht] = 0
	}
	g.Deck = NewStandardDeck(rng)
	return g
}

// StartGame transitions from title to the first blind select.
func (g *GameState) StartGame() {
	g.Phase = PhaseBlindSelect
}

// CurrentBlindType returns the current blind type.
func (g *GameState) CurrentBlindType() BlindType {
	switch g.BlindIndex {
	case 0:
		return SmallBlind
	case 1:
		return BigBlind
	default:
		return BossBlind
	}
}

// SkipBlind skips the current non-boss blind and awards a skip reward.
func (g *GameState) SkipBlind() {
	reward := RandomSkipReward(g.RNG)
	g.LastSkipReward = reward
	ApplySkipReward(reward, &g.Money, &g.NextShopFreeReroll, &g.NextShopFreeJoker, &g.NextRoundHandSizeBonus)
	g.AdvanceBlind()
}

// EnterBlind starts a blind round.
func (g *GameState) EnterBlind() {
	blind := g.CurrentBlindType()

	if blind == BossBlind {
		g.CurrentBoss = RandomBoss(g.Ante, g.RNG)
		g.Phase = PhaseBossInfo
	} else {
		g.CurrentBoss = nil
		g.startRound()
	}
}

// ConfirmBossInfo transitions from boss info to play phase.
func (g *GameState) ConfirmBossInfo() {
	g.startRound()
}

func (g *GameState) startRound() {
	// Reset the "last scoring breakdown" so neither the sidebar nor the
	// play screen's table area show stale display from the previous round.
	g.LastBreakdown = nil

	g.Deck.Shuffle()
	g.RemainingHands = 4
	g.RemainingDisc = 3
	g.RoundScore = 0
	g.RoundHandTypes = nil

	// Calculate hand size with bonuses
	g.HandSize = g.BaseHandSize + g.NextRoundHandSizeBonus
	g.NextRoundHandSizeBonus = 0

	// Apply passive joker effects
	passive := GetPassiveEffects(g.Jokers)
	g.HandSize += passive.HandSizeMod
	g.RemainingDisc += passive.DiscardMod

	// Apply boss effect
	if g.CurrentBoss != nil {
		roundState := g.makeRoundState()
		g.CurrentBoss.Apply(roundState)
		g.RemainingHands = roundState.RemainingHands
		g.RemainingDisc = roundState.RemainingDisc
		g.HandSize = roundState.HandSize
	}

	// Calculate target
	bossMult := 0.0
	if g.CurrentBoss != nil {
		bossMult = g.CurrentBoss.TargetMultiplier()
	}
	g.TargetScore = CalcTargetScore(g.Ante, g.CurrentBlindType(), bossMult)

	// Draw initial hand
	g.Hand = g.Deck.Draw(g.HandSize)
	g.applyDebuffs()
	g.SelectedCards = make([]bool, len(g.Hand))

	g.Phase = PhasePlay
}

// applyDebuffs marks cards as debuffed based on boss effect.
func (g *GameState) applyDebuffs() {
	if g.CurrentBoss == nil {
		return
	}
	for i := range g.Hand {
		g.Hand[i].Debuffed = g.CurrentBoss.IsDebuffed(g.Hand[i])
	}
}

// ToggleSelect toggles selection of a card in hand.
func (g *GameState) ToggleSelect(index int) {
	if index >= 0 && index < len(g.SelectedCards) {
		g.SelectedCards[index] = !g.SelectedCards[index]
	}
}

// SelectedCardsList returns the currently selected cards.
func (g *GameState) SelectedCardsList() []Card {
	var cards []Card
	for i, sel := range g.SelectedCards {
		if sel {
			cards = append(cards, g.Hand[i])
		}
	}
	return cards
}

// PlayHand plays the selected cards, scores them, and transitions to scoring.
func (g *GameState) PlayHand() error {
	selected := g.SelectedCardsList()
	if len(selected) < 1 || len(selected) > 5 {
		return &GameError{"must select 1-5 cards"}
	}
	if g.RemainingHands <= 0 {
		return &GameError{"no hands remaining"}
	}

	// Boss validation
	if g.CurrentBoss != nil {
		roundState := g.makeRoundState()
		if err := g.CurrentBoss.ValidatePlay(selected, roundState); err != nil {
			return err
		}
	}

	g.RemainingHands--

	// Detect hand type
	passive := GetPassiveEffects(g.Jokers)
	opts := HandDetectOpts{
		FourFingers:  passive.FourFingers,
		SmearedJoker: passive.SmearedJoker,
		Splash:       passive.Splash,
	}
	result := DetectHand(selected, opts)

	// Build held cards (cards not played)
	var heldCards []Card
	for i, sel := range g.SelectedCards {
		if !sel {
			heldCards = append(heldCards, g.Hand[i])
		}
	}

	// Score
	ctx := &ScoringContext{
		HandType:       result.Type,
		HandLevel:      g.HandLevels[result.Type],
		ScoringCards:   result.ScoringCards,
		HeldCards:      heldCards,
		PlayedCards:    selected,
		Jokers:         g.Jokers,
		RemainingHands: g.RemainingHands,
		RemainingDisc:  g.RemainingDisc,
		Money:          g.Money,
		DeckSize:       g.Deck.Remaining(),
		HandPlayCounts: g.HandPlayCounts,
		IsLastHand:     g.RemainingHands == 0,
		RoundHandTypes: g.RoundHandTypes,
		RNG:            g.RNG,
	}
	breakdown := CalculateScore(ctx)
	g.LastBreakdown = &breakdown
	g.RoundScore += breakdown.TotalScore

	// Update hand type tracking
	g.HandPlayCounts[result.Type]++
	if g.HandPlayCounts[result.Type]%3 == 0 {
		g.HandLevels[result.Type]++
	}
	g.RoundHandTypes = append(g.RoundHandTypes, result.Type)

	// Boss OnPlay hook
	if g.CurrentBoss != nil {
		roundState := g.makeRoundState()
		g.CurrentBoss.OnPlay(selected, roundState)
	}

	// Remove played cards from hand, discard them
	g.removeSelectedFromHand()

	// Draw new cards to refill hand (stay in PhasePlay, UI handles transitions)
	g.drawToHandSize()
	g.applyDebuffs()
	g.SelectedCards = make([]bool, len(g.Hand))

	return nil
}

// FinishScoring transitions from scoring display to next action.
func (g *GameState) FinishScoring() {
	if g.RoundScore >= g.TargetScore {
		// Won this blind
		g.Phase = PhaseSettlement
	} else if g.RemainingHands <= 0 {
		// Lost
		g.Lost = true
		g.Phase = PhaseGameOver
	} else {
		// Continue playing - draw new cards
		g.drawToHandSize()
		g.applyDebuffs()
		g.SelectedCards = make([]bool, len(g.Hand))
		g.Phase = PhasePlay
	}
}

// FinishSettlement calculates rewards and transitions to shop.
func (g *GameState) FinishSettlement() {
	reward := CalcReward(g.CurrentBlindType(), g.RemainingHands, g.Money)
	g.Money += reward

	// Joker OnRoundEnd hooks
	ctx := &ScoringContext{Money: g.Money, RNG: g.RNG}
	for _, j := range g.Jokers {
		if j != nil && j.Effect != nil {
			j.Effect.OnRoundEnd(ctx)
		}
	}

	g.Shop = NewShopState(g.RNG, g.NextShopFreeReroll, g.NextShopFreeJoker)
	g.NextShopFreeReroll = false
	g.NextShopFreeJoker = false
	g.Phase = PhaseShop
}

// LeaveShop advances to the next blind or ante.
func (g *GameState) LeaveShop() {
	g.AdvanceBlind()
}

// AdvanceBlind moves to the next blind, or next ante if all 3 done.
func (g *GameState) AdvanceBlind() {
	g.BlindIndex++
	if g.BlindIndex >= 3 {
		g.BlindIndex = 0
		g.Ante++
		if g.Ante > 8 {
			g.Won = true
			g.Phase = PhaseGameOver
			return
		}
	}
	g.Phase = PhaseBlindSelect
}

// DiscardSelected discards selected cards and draws replacements.
func (g *GameState) DiscardSelected() error {
	selected := g.SelectedCardsList()
	if len(selected) < 1 || len(selected) > 5 {
		return &GameError{"must select 1-5 cards to discard"}
	}
	if g.RemainingDisc <= 0 {
		return &GameError{"no discards remaining"}
	}

	g.RemainingDisc--
	g.Deck.Discard(selected)
	g.removeSelectedFromHand()
	g.drawToHandSize()
	g.applyDebuffs()
	g.SelectedCards = make([]bool, len(g.Hand))

	// Joker OnDiscard hooks
	ctx := &ScoringContext{RNG: g.RNG}
	for _, j := range g.Jokers {
		if j != nil && j.Effect != nil {
			j.Effect.OnDiscard(ctx)
		}
	}

	return nil
}

// SellJoker sells a joker at the given index.
func (g *GameState) SellJoker(index int) bool {
	if index < 0 || index >= len(g.Jokers) {
		return false
	}
	j := g.Jokers[index]
	g.Money += j.SellValue()

	// OnSell hook
	ctx := &ScoringContext{Money: g.Money, RNG: g.RNG}
	j.Effect.OnSell(ctx)

	// Remove from slice
	g.Jokers = append(g.Jokers[:index], g.Jokers[index+1:]...)
	return true
}

func (g *GameState) removeSelectedFromHand() {
	var remaining []Card
	for i, card := range g.Hand {
		if !g.SelectedCards[i] {
			remaining = append(remaining, card)
		} else {
			g.Deck.Discard([]Card{card})
		}
	}
	g.Hand = remaining
}

func (g *GameState) drawToHandSize() {
	need := g.HandSize - len(g.Hand)
	if need > 0 {
		drawn := g.Deck.Draw(need)
		g.Hand = append(g.Hand, drawn...)
	}
}

func (g *GameState) makeRoundState() *RoundState {
	return &RoundState{
		RemainingHands:  g.RemainingHands,
		RemainingDisc:   g.RemainingDisc,
		HandSize:        g.HandSize,
		Money:           &g.Money,
		PlayedHandTypes: g.RoundHandTypes,
		HandLevels:      g.HandLevels,
	}
}

// GameError represents a game rule violation.
type GameError struct {
	Message string
}

func (e *GameError) Error() string { return e.Message }
