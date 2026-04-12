package jokers

import (
	"testing"

	"github.com/yinuoli/balatro-cli/engine"
	"github.com/yinuoli/balatro-cli/testutil"
)

func TestFourFingersPassive(t *testing.T) {
	if !FourFingersDef.Passive.FourFingers {
		t.Error("FourFingers passive should be true")
	}
	if FourFingersDef.EffectFactory != nil {
		t.Error("FourFingers should have no scoring effect")
	}
}

func TestBlackboard(t *testing.T) {
	j := engine.NewJokerInstance(BlackboardDef)

	// All black held cards → ×3
	ctx := &engine.ScoringContext{
		HeldCards: testutil.MakeCards("2s", "5c", "8s"),
	}
	mod := j.Effect.OnScore(ctx)
	if mod.MultMult != 3.0 {
		t.Errorf("Blackboard all black: ×mult = %f, want 3.0", mod.MultMult)
	}

	// Mixed held cards → no effect
	ctx.HeldCards = testutil.MakeCards("2s", "5h", "8c")
	mod = j.Effect.OnScore(ctx)
	if mod.MultMult != 0 && mod.MultMult != 1.0 {
		t.Errorf("Blackboard mixed: ×mult = %f, want 1.0/0", mod.MultMult)
	}

	// Empty held cards → no effect
	ctx.HeldCards = nil
	mod = j.Effect.OnScore(ctx)
	if mod.MultMult != 0 && mod.MultMult != 1.0 {
		t.Errorf("Blackboard empty: ×mult = %f, want 1.0/0", mod.MultMult)
	}
}

func TestAcrobat(t *testing.T) {
	j := engine.NewJokerInstance(AcrobatDef)

	// Last hand → ×3
	ctx := &engine.ScoringContext{IsLastHand: true}
	mod := j.Effect.OnScore(ctx)
	if mod.MultMult != 3.0 {
		t.Errorf("Acrobat last hand: ×mult = %f, want 3.0", mod.MultMult)
	}

	// Not last hand → no effect
	ctx.IsLastHand = false
	mod = j.Effect.OnScore(ctx)
	if mod.MultMult != 0 && mod.MultMult != 1.0 {
		t.Errorf("Acrobat not last: ×mult = %f, want 1.0/0", mod.MultMult)
	}
}
