package randomix

// lcg64 uses the full-period 64-bit linear congruential generator described by
// Knuth (MMIX): x' = 6364136223846793005*x + 1442695040888963407 (mod 2^64).
const (
	lcgMultiplier = 6364136223846793005
	lcgIncrement  = 1442695040888963407
)

// lcg is a 64-bit linear congruential generator.
type lcg struct {
	*core
	state uint64
}

func newLCG(seed uint64) baseRandom {
	g := &lcg{state: seed}
	g.core = newCore(g.next)
	return g
}

func (g *lcg) next() uint64 {
	g.state = g.state*lcgMultiplier + lcgIncrement
	return g.state
}
