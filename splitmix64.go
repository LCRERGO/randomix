package randomix

// splitmix64 is the SplitMix64 generator.
type splitmix64 struct {
	*core
	state uint64
}

func newSplitMix64(seed uint64) baseRandom {
	g := &splitmix64{state: seed}
	g.core = newCore(g.next)
	return g
}

func (g *splitmix64) next() uint64 {
	return splitMixNext(&g.state)
}
