package engine

import "testing"

func TestSuitDebuffBosses(t *testing.T) {
	tests := []struct {
		name    string
		boss    BossEffect
		card    Card
		debuff  bool
	}{
		{"Club debuffs clubs", NewTheClub(), NewCard(King, Clubs), true},
		{"Club ignores spades", NewTheClub(), NewCard(King, Spades), false},
		{"Goad debuffs spades", NewTheGoad(), NewCard(Ace, Spades), true},
		{"Goad ignores hearts", NewTheGoad(), NewCard(Ace, Hearts), false},
		{"Head debuffs hearts", NewTheHead(), NewCard(Two, Hearts), true},
		{"Window debuffs diamonds", NewTheWindow(), NewCard(Queen, Diamonds), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.boss.IsDebuffed(tt.card); got != tt.debuff {
				t.Errorf("IsDebuffed = %v, want %v", got, tt.debuff)
			}
		})
	}
}

func TestPlantDebuffsFaceCards(t *testing.T) {
	boss := NewThePlant()
	if !boss.IsDebuffed(NewCard(Jack, Spades)) {
		t.Error("Plant should debuff Jack")
	}
	if !boss.IsDebuffed(NewCard(Queen, Hearts)) {
		t.Error("Plant should debuff Queen")
	}
	if boss.IsDebuffed(NewCard(Ten, Spades)) {
		t.Error("Plant should not debuff 10")
	}
	if boss.IsDebuffed(NewCard(Ace, Spades)) {
		t.Error("Plant should not debuff Ace")
	}
}

func TestPsychicValidation(t *testing.T) {
	boss := NewThePsychic()
	state := &RoundState{}

	// 5 cards: valid
	if err := boss.ValidatePlay(make([]Card, 5), state); err != nil {
		t.Errorf("5 cards should be valid: %v", err)
	}
	// 3 cards: invalid
	if err := boss.ValidatePlay(make([]Card, 3), state); err == nil {
		t.Error("3 cards should be invalid for Psychic")
	}
}

func TestManacleReducesHandSize(t *testing.T) {
	boss := NewTheManacle()
	state := &RoundState{HandSize: 8}
	boss.Apply(state)
	if state.HandSize != 7 {
		t.Errorf("hand size = %d, want 7", state.HandSize)
	}
}

func TestWallTargetMultiplier(t *testing.T) {
	boss := NewTheWall()
	if m := boss.TargetMultiplier(); m != 4.0 {
		t.Errorf("Wall mult = %f, want 4.0", m)
	}
}

func TestNeedleOneHand(t *testing.T) {
	boss := NewTheNeedle()
	state := &RoundState{RemainingHands: 4}
	boss.Apply(state)
	if state.RemainingHands != 1 {
		t.Errorf("hands = %d, want 1", state.RemainingHands)
	}
}

func TestWaterNoDiscards(t *testing.T) {
	boss := NewTheWater()
	state := &RoundState{RemainingDisc: 3}
	boss.Apply(state)
	if state.RemainingDisc != 0 {
		t.Errorf("discards = %d, want 0", state.RemainingDisc)
	}
}

func TestToothReducesMoney(t *testing.T) {
	boss := NewTheTooth()
	money := 10
	state := &RoundState{Money: &money}
	cards := make([]Card, 3)
	boss.OnPlay(cards, state)
	if money != 7 {
		t.Errorf("money = %d, want 7", money)
	}
}

func TestVioletVesselTarget(t *testing.T) {
	boss := NewVioletVessel()
	if m := boss.TargetMultiplier(); m != 6.0 {
		t.Errorf("VioletVessel mult = %f, want 6.0", m)
	}
}
