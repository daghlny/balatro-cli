package jokers

import (
	"testing"

	"github.com/yinuoli/balatro-cli/engine"
)

func TestTheDuo(t *testing.T) {
	j := engine.NewJokerInstance(TheDuoDef)

	tests := []struct {
		name     string
		handType engine.HandType
		wantMult float64
	}{
		{"Pair", engine.Pair, 2.0},
		{"FullHouse", engine.FullHouse, 2.0},
		{"Straight (no pair)", engine.Straight, 1.0},
		{"Flush (no pair)", engine.Flush, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &engine.ScoringContext{HandType: tt.handType}
			mod := j.Effect.OnScore(ctx)
			if mod.MultMult != tt.wantMult {
				t.Errorf("TheDuo on %s: ×mult = %f, want %f", tt.handType, mod.MultMult, tt.wantMult)
			}
		})
	}
}
