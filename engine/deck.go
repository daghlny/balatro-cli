package engine

// Deck manages a draw pile and discard pile.
type Deck struct {
	drawPile    []Card
	discardPile []Card
	rng         *RNG
}

// NewStandardDeck creates a 52-card deck (unshuffled).
func NewStandardDeck(rng *RNG) *Deck {
	cards := make([]Card, 0, 52)
	for _, suit := range AllSuits() {
		for _, rank := range AllRanks() {
			cards = append(cards, NewCard(rank, suit))
		}
	}
	return &Deck{drawPile: cards, rng: rng}
}

// NewDeckFromCards creates a deck from a specific set of cards.
func NewDeckFromCards(cards []Card, rng *RNG) *Deck {
	cp := make([]Card, len(cards))
	copy(cp, cards)
	return &Deck{drawPile: cp, rng: rng}
}

// Shuffle randomizes the draw pile.
func (d *Deck) Shuffle() {
	d.rng.Shuffle(len(d.drawPile), func(i, j int) {
		d.drawPile[i], d.drawPile[j] = d.drawPile[j], d.drawPile[i]
	})
}

// Draw takes n cards from the top of the draw pile.
// If the draw pile is exhausted, the discard pile is shuffled back in.
// Returns fewer cards if both piles together have fewer than n.
func (d *Deck) Draw(n int) []Card {
	drawn := make([]Card, 0, n)
	for i := 0; i < n; i++ {
		if len(d.drawPile) == 0 {
			if len(d.discardPile) == 0 {
				break
			}
			d.drawPile = d.discardPile
			d.discardPile = nil
			d.Shuffle()
		}
		drawn = append(drawn, d.drawPile[0])
		d.drawPile = d.drawPile[1:]
	}
	return drawn
}

// Discard moves cards to the discard pile.
func (d *Deck) Discard(cards []Card) {
	d.discardPile = append(d.discardPile, cards...)
}

// Remaining returns the number of cards in the draw pile.
func (d *Deck) Remaining() int {
	return len(d.drawPile)
}

// TotalRemaining returns cards in both draw and discard piles.
func (d *Deck) TotalRemaining() int {
	return len(d.drawPile) + len(d.discardPile)
}

// Reclaim moves the discard pile back into the draw pile, so a new round starts with the
// full deck (as in real balatro). Shuffle afterwards.
func (d *Deck) Reclaim() {
	d.drawPile = append(d.drawPile, d.discardPile...)
	d.discardPile = nil
}

// AddCards adds cards to the draw pile (used by DNA joker).
func (d *Deck) AddCards(cards []Card) {
	d.drawPile = append(d.drawPile, cards...)
}
