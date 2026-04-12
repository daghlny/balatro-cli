package engine

import "testing"

func TestCalcTargetScore(t *testing.T) {
	tests := []struct {
		ante int
		blind BlindType
		bossMult float64
		want int
	}{
		{1, SmallBlind, 0, 300},
		{1, BigBlind, 0, 450},
		{1, BossBlind, 0, 600},     // default 2x
		{1, BossBlind, 4.0, 1200},  // Wall 4x
		{4, SmallBlind, 0, 5000},
		{4, BigBlind, 0, 7500},
		{8, SmallBlind, 0, 50000},
		{8, BossBlind, 6.0, 300000}, // Violet Vessel 6x
	}
	for _, tt := range tests {
		got := CalcTargetScore(tt.ante, tt.blind, tt.bossMult)
		if got != tt.want {
			t.Errorf("CalcTargetScore(%d, %s, %f) = %d, want %d",
				tt.ante, tt.blind, tt.bossMult, got, tt.want)
		}
	}
}

func TestCalcReward(t *testing.T) {
	tests := []struct {
		blind    BlindType
		hands    int
		money    int
		want     int
	}{
		{SmallBlind, 2, 17, 8},  // 3 + 2 + floor(17/5)=3 = 8
		{BigBlind, 2, 17, 9},    // 4 + 2 + 3 = 9
		{BossBlind, 0, 25, 10},  // 5 + 0 + 5 = 10
		{SmallBlind, 0, 0, 3},   // 3 + 0 + 0 = 3
		{SmallBlind, 4, 30, 12}, // 3 + 4 + 5(cap) = 12
	}
	for _, tt := range tests {
		got := CalcReward(tt.blind, tt.hands, tt.money)
		if got != tt.want {
			t.Errorf("CalcReward(%s, %d, %d) = %d, want %d",
				tt.blind, tt.hands, tt.money, got, tt.want)
		}
	}
}

func TestCalcRewardSmallBlindExample(t *testing.T) {
	// From simple_game_rules: BigBlind, 2 hands left, $17 → $4+$2+$3 = $9
	got := CalcReward(BigBlind, 2, 17)
	if got != 9 {
		t.Errorf("example reward = %d, want 9", got)
	}
}

func TestSkipRewardMoney(t *testing.T) {
	money := 10
	freeReroll := false
	freeJoker := false
	handSizeBonus := 0
	reward := &SkipReward{Type: SkipRewardMoney}
	ApplySkipReward(reward, &money, &freeReroll, &freeJoker, &handSizeBonus)
	if money != 15 {
		t.Errorf("money after skip = %d, want 15", money)
	}
}

func TestSkipRewardDouble(t *testing.T) {
	money := 8
	freeReroll := false
	freeJoker := false
	handSizeBonus := 0
	reward := &SkipReward{Type: SkipRewardDouble}
	ApplySkipReward(reward, &money, &freeReroll, &freeJoker, &handSizeBonus)
	if money != 16 {
		t.Errorf("money after double = %d, want 16", money)
	}

	// Cap at $20
	money = 15
	ApplySkipReward(reward, &money, &freeReroll, &freeJoker, &handSizeBonus)
	if money != 20 {
		t.Errorf("money after double (capped) = %d, want 20", money)
	}
}
