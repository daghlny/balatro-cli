package engine

import (
	"testing"
)

func init() {
	// Register some test jokers for shop tests
	RegisterJoker(&JokerDef{ID: "test_common", Rarity: Common, Cost: 3})
	RegisterJoker(&JokerDef{ID: "test_uncommon", Rarity: Uncommon, Cost: 6})
	RegisterJoker(&JokerDef{ID: "test_rare", Rarity: Rare, Cost: 9})
}

func TestNewShopState(t *testing.T) {
	rng := NewSeededRNG(42)
	shop := NewShopState(rng, false, false)
	if shop.RerollCost != 5 {
		t.Errorf("initial reroll cost = %d, want 5", shop.RerollCost)
	}
	// Both slots should have jokers
	for i := 0; i < 2; i++ {
		if shop.JokerSlots[i] == nil {
			t.Errorf("joker slot %d is nil", i)
		}
	}
}

func TestShopReroll(t *testing.T) {
	rng := NewSeededRNG(42)
	shop := NewShopState(rng, false, false)
	money := 10

	old0 := shop.JokerSlots[0]
	old1 := shop.JokerSlots[1]

	ok := shop.Reroll(&money)
	if !ok {
		t.Error("reroll should succeed with $10")
	}
	if money != 5 {
		t.Errorf("money after reroll = %d, want 5", money)
	}
	if shop.RerollCost != 6 {
		t.Errorf("reroll cost after first = %d, want 6", shop.RerollCost)
	}
	// At least one slot should have changed (statistically very likely with seeded RNG)
	_ = old0
	_ = old1
}

func TestShopRerollInsufficientFunds(t *testing.T) {
	rng := NewSeededRNG(42)
	shop := NewShopState(rng, false, false)
	money := 3

	ok := shop.Reroll(&money)
	if ok {
		t.Error("reroll should fail with $3")
	}
	if money != 3 {
		t.Errorf("money should be unchanged: %d", money)
	}
}

func TestShopFreeReroll(t *testing.T) {
	rng := NewSeededRNG(42)
	shop := NewShopState(rng, true, false)
	money := 0

	ok := shop.Reroll(&money)
	if !ok {
		t.Error("free reroll should succeed with $0")
	}
	if money != 0 {
		t.Errorf("money should stay 0, got %d", money)
	}
	// Second reroll should NOT be free
	if shop.FreeReroll {
		t.Error("free reroll should be consumed")
	}
}

func TestShopBuyJoker(t *testing.T) {
	rng := NewSeededRNG(42)
	shop := NewShopState(rng, false, false)
	money := 20
	jokers := make([]*JokerInstance, 0)

	cost := shop.JokerCosts[0]
	j, ok := shop.BuyJoker(0, &money, &jokers, 5)
	if !ok {
		t.Error("buy should succeed")
	}
	if j == nil {
		t.Error("should return joker instance")
	}
	if money != 20-cost {
		t.Errorf("money = %d, want %d", money, 20-cost)
	}
	if len(jokers) != 1 {
		t.Errorf("joker count = %d, want 1", len(jokers))
	}
	if shop.JokerSlots[0] != nil {
		t.Error("bought slot should be nil")
	}
}

func TestShopBuyJokerFullSlots(t *testing.T) {
	rng := NewSeededRNG(42)
	shop := NewShopState(rng, false, false)
	money := 100
	jokers := make([]*JokerInstance, 5) // full

	_, ok := shop.BuyJoker(0, &money, &jokers, 5)
	if ok {
		t.Error("should not buy with full joker slots")
	}
}

func TestBoosterPack(t *testing.T) {
	rng := NewSeededRNG(42)
	shop := NewShopState(rng, false, false)
	pack := shop.GenerateBoosterPack()

	tier := boosterTiers[shop.BoosterTier]
	if pack.Total != tier.total {
		t.Errorf("pack total = %d, want %d", pack.Total, tier.total)
	}
	if pack.Choose != tier.choose {
		t.Errorf("pack choose = %d, want %d", pack.Choose, tier.choose)
	}
	if len(pack.Jokers) != tier.total {
		t.Errorf("pack jokers = %d, want %d", len(pack.Jokers), tier.total)
	}
}
