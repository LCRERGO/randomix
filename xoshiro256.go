package randomix

// xoshiro256 is the xoshiro256** generator (Blackman & Vigna).
type xoshiro256 struct {
	*core
	s [4]uint64
}

func newXoshiro256(seed uint64) baseRandom {
	g := &xoshiro256{}
	sm := seed
	for i := range g.s {
		g.s[i] = splitMixNext(&sm)
	}
	if g.s[0] == 0 && g.s[1] == 0 && g.s[2] == 0 && g.s[3] == 0 {
		g.s[0] = 1
	}
	g.core = newCore(g.next)
	return g
}

func (g *xoshiro256) next() uint64 {
	result := rotl64(g.s[1]*5, 7) * 9
	t := g.s[1] << 17

	g.s[2] ^= g.s[0]
	g.s[3] ^= g.s[1]
	g.s[1] ^= g.s[2]
	g.s[0] ^= g.s[3]
	g.s[2] ^= t
	g.s[3] = rotl64(g.s[3], 45)

	return result
}
