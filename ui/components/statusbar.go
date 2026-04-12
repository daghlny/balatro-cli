package components

import (
	"fmt"

	"github.com/yinuoli/balatro-cli/engine"
	"github.com/yinuoli/balatro-cli/i18n"
	
)

// FocusType indicates what kind of element the cursor is on.
type FocusType int

const (
	FocusNone FocusType = iota
	FocusHandCard
	FocusJokerSlot
	FocusShopJoker
	FocusShopBooster
)

// FocusTarget describes the current cursor focus for the status bar.
type FocusTarget struct {
	Type  FocusType
	Index int
}

// RenderStatusBar renders the bottom status bar based on current focus.
func RenderStatusBar(focus FocusTarget, game *engine.GameState, width int) string {
	var content string

	switch focus.Type {
	case FocusHandCard:
		if focus.Index >= 0 && focus.Index < len(game.Hand) {
			card := game.Hand[focus.Index]
			suitName := suitKey(card.Suit)
			content = fmt.Sprintf("%s%s  %s  %s: %d",
				card.Rank.String(), card.Suit.Symbol(),
				i18n.T(suitName),
				i18n.T("status.chips_value"), card.Rank.Chips())
			if card.Debuffed {
				content += "  " + i18n.T("status.debuffed")
			}
		}

	case FocusJokerSlot:
		if focus.Index >= 0 && focus.Index < len(game.Jokers) {
			j := game.Jokers[focus.Index]
			name := i18n.T(j.Def.NameKey)
			desc := i18n.T(j.Def.DescKey)
			rarity := rarityI18nKey(j.Def.Rarity)
			content = fmt.Sprintf("%s — %s | %s: $%d | %s: $%d | %s: %s",
				name, desc,
				i18n.T("status.cost"), j.Cost,
				i18n.T("status.sell_value"), j.SellValue(),
				i18n.T("status.rarity"), i18n.T(rarity))
		} else {
			content = i18n.T("status.empty_slot")
		}

	case FocusShopJoker:
		if game.Shop != nil && focus.Index >= 0 && focus.Index < 2 {
			def := game.Shop.JokerSlots[focus.Index]
			if def != nil {
				name := i18n.T(def.NameKey)
				desc := i18n.T(def.DescKey)
				cost := game.Shop.JokerCosts[focus.Index]
				content = fmt.Sprintf("%s — %s | %s: $%d | %s: %s",
					name, desc,
					i18n.T("status.cost"), cost,
					i18n.T("status.rarity"), i18n.T(rarityI18nKey(def.Rarity)))
			} else {
				content = i18n.T("shop.sold_out")
			}
		}

	default:
		content = "(←→)移动  (空格)" + i18n.T("action.select") + "  " +
			"(Enter)" + i18n.T("action.play") + "  " +
			"(D)" + i18n.T("action.discard") + "  " +
			"(S)" + i18n.T("action.sort") + "  " +
			"(Tab)切换  " +
			"(Q)" + i18n.T("action.quit")
	}

	style := StatusBarStyle.Width(width - 2)
	return style.Render(content)
}

func suitKey(s engine.Suit) string {
	keys := map[engine.Suit]string{
		engine.Spades: "suit.spades", engine.Hearts: "suit.hearts",
		engine.Clubs: "suit.clubs", engine.Diamonds: "suit.diamonds",
	}
	return keys[s]
}

func rarityI18nKey(r engine.Rarity) string {
	keys := map[engine.Rarity]string{
		engine.Common: "rarity.common", engine.Uncommon: "rarity.uncommon",
		engine.Rare: "rarity.rare", engine.Legendary: "rarity.legendary",
	}
	return keys[r]
}
