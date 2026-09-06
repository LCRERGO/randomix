package randomix

import "math"

// core turns an arbitrary raw 64-bit word stream into the full baseRandom
// method set (uniform sampling, floats, permutations, shuffle). Concrete PRNG
// types embed a *core wired to their own raw word source.
type core struct {
	next func() uint64
}

func newCore(next func() uint64) *core {
	return &core{next: next}
}

// Uint64 returns the next raw 64-bit word.
func (c *core) Uint64() uint64 { return c.next() }

func (c *core) Uint32() uint32 { return uint32(c.Uint64() >> 32) }

func (c *core) Uint() uint { return uint(c.Uint64()) }

func (c *core) UintN(n uint) uint {
	if n == 0 {
		panic("randomix: invalid argument to UintN")
	}
	return uint(uint64N(c.next, uint64(n)))
}

func (c *core) Uint32N(n uint32) uint32 {
	if n == 0 {
		panic("randomix: invalid argument to Uint32N")
	}
	return uint32(uint64N(c.next, uint64(n)))
}

func (c *core) Uint64N(n uint64) uint64 {
	if n == 0 {
		panic("randomix: invalid argument to Uint64N")
	}
	return uint64N(c.next, n)
}

func (c *core) Int() int { return int(uint(c.Uint64()) >> 1) }

func (c *core) Int32() int32 { return int32(c.Uint32() >> 1) }

func (c *core) Int64() int64 { return int64(c.Uint64() >> 1) }

func (c *core) Int32N(n int32) int32 {
	if n <= 0 {
		panic("randomix: invalid argument to Int32N")
	}
	return int32(uint64N(c.next, uint64(n)))
}

func (c *core) Int64N(n int64) int64 {
	if n <= 0 {
		panic("randomix: invalid argument to Int64N")
	}
	return int64(uint64N(c.next, uint64(n)))
}

func (c *core) IntN(n int) int {
	if n <= 0 {
		panic("randomix: invalid argument to IntN")
	}
	return int(uint64N(c.next, uint64(n)))
}

// Float64 returns a uniform float in [0, 1) with 53 random bits.
func (c *core) Float64() float64 {
	return float64(c.Uint64()>>11) * (1.0 / (1 << 53))
}

// Float32 returns a uniform float in [0, 1) with 24 random bits.
func (c *core) Float32() float32 {
	return float32(c.Uint64()>>40) * (1.0 / (1 << 24))
}

// NormFloat64 returns a standard normal value via the polar method.
func (c *core) NormFloat64() float64 {
	for {
		a := c.Float64()*2 - 1
		b := c.Float64()*2 - 1
		s := a*a + b*b
		if s >= 1 || s == 0 {
			continue
		}
		f := math.Sqrt(-2 * math.Log(s) / s)
		return a * f
	}
}

// ExpFloat64 returns an exponentially distributed value with mean 1.
func (c *core) ExpFloat64() float64 {
	return -math.Log(1 - c.Float64())
}

func (c *core) Shuffle(n int, swap func(i, j int)) {
	if n < 0 {
		panic("randomix: invalid argument to Shuffle")
	}
	for i := n - 1; i > 0; i-- {
		j := c.IntN(i + 1)
		swap(i, j)
	}
}

func (c *core) Perm(n int) []int {
	m := make([]int, n)
	for i := range m {
		m[i] = i
	}
	c.Shuffle(n, func(i, j int) { m[i], m[j] = m[j], m[i] })
	return m
}

// uint64N returns a value uniformly in [0, n) without modulo bias.
func uint64N(next func() uint64, n uint64) uint64 {
	if n&(n-1) == 0 { // power of two: mask
		return next() & (n - 1)
	}
	// Reject draws above the largest multiple of n that fits in a uint64 so
	// that x % n is perfectly uniform over the accepted range.
	limit := ^uint64(0) - ^uint64(0)%n
	for {
		if x := next(); x < limit {
			return x % n
		}
	}
}

// splitMixNext advances *state and returns one SplitMix64 word. Used both as a
// generator and as a deterministic seed expander for the other algorithms.
func splitMixNext(state *uint64) uint64 {
	*state += 0x9e3779b97f4a7c15
	z := *state
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

func rotl64(x uint64, k uint) uint64 { return x<<k | x>>(64-k) }

func rotr64(x uint64, k uint) uint64 {
	if k == 0 {
		return x
	}
	return x>>k | x<<(64-k)
}
