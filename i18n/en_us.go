package i18n

func init() {
	Register("en-US", map[string]string{
		"hand.high_card":      "High Card",
		"hand.pair":           "Pair",
		"hand.two_pair":       "Two Pair",
		"hand.three_of_kind":  "Three of a Kind",
		"hand.straight":       "Straight",
		"hand.flush":          "Flush",
		"hand.full_house":     "Full House",
		"hand.four_of_kind":   "Four of a Kind",
		"hand.straight_flush": "Straight Flush",
		"hand.royal_flush":    "Royal Flush",

		"ui.ante": "Ante", "ui.score": "Score", "ui.target": "Target",
		"ui.money": "Money", "ui.hands": "Hands", "ui.discards": "Discards",
		"ui.hand_reference": "Poker Hands",

		"action.play": "Play", "action.discard": "Discard",
		"action.sort_rank": "Sort Rank", "action.sort_suit": "Sort Suit",
		"action.start": "Start Game", "action.quit": "Quit",

		"title.name":        "BALATRO",
		"title.subtitle":    "Terminal Edition",
		"title.press_start": "Press Enter to start",
	})
}
