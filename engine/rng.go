package engine

import "math/rand"

// RNG wraps a random source for reproducible game behavior.
// Use NewRNG for production, NewSeededRNG for tests.
type RNG struct {
	r *rand.Rand
}

func NewRNG() *RNG {
	return &RNG{r: rand.New(rand.NewSource(rand.Int63()))}
}

func NewSeededRNG(seed int64) *RNG {
	return &RNG{r: rand.New(rand.NewSource(seed))}
}

// Intn returns a random int in [0, n).
func (rng *RNG) Intn(n int) int {
	return rng.r.Intn(n)
}

// Float64 returns a random float64 in [0.0, 1.0).
func (rng *RNG) Float64() float64 {
	return rng.r.Float64()
}

// Shuffle randomizes the order of elements using Fisher-Yates.
func (rng *RNG) Shuffle(n int, swap func(i, j int)) {
	rng.r.Shuffle(n, swap)
}
