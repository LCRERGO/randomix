package randomix

const (
	mtN          = 624
	mtM          = 397
	mtMatrixA    = 0x9908b0df
	mtUpperMask  = 0x80000000
	mtLowerMask  = 0x7fffffff
	mtInitFactor = 1812433253
)

// mt19937 is a Mersenne Twister (MT19937) with 32-bit output words.
type mt19937 struct {
	*core
	state [mtN]uint32
	idx   int
}

func newMT19937(seed uint64) baseRandom {
	g := &mt19937{}
	g.init(uint32(seed) ^ uint32(seed>>32))
	g.core = newCore(g.next)
	return g
}

func (g *mt19937) init(seed uint32) {
	g.state[0] = seed
	for i := 1; i < mtN; i++ {
		g.state[i] = mtInitFactor*(g.state[i-1]^(g.state[i-1]>>30)) + uint32(i)
	}
	g.idx = mtN
}

func (g *mt19937) twist() {
	for i := 0; i < mtN; i++ {
		y := (g.state[i] & mtUpperMask) | (g.state[(i+1)%mtN] & mtLowerMask)
		x := g.state[(i+mtM)%mtN] ^ (y >> 1)
		if y&1 != 0 {
			x ^= mtMatrixA
		}
		g.state[i] = x
	}
	g.idx = 0
}

func (g *mt19937) temper() uint32 {
	if g.idx >= mtN {
		g.twist()
	}
	y := g.state[g.idx]
	g.idx++
	y ^= y >> 11
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= y >> 18
	return y
}

func (g *mt19937) next() uint64 {
	return uint64(g.temper())<<32 | uint64(g.temper())
}
