package randomix

import "math/bits"

// pcg64 is the PCG64 generator: a 128-bit LCG combined with the XSL-RR output
// function (same construction used by the rand_pcg Pcg64), constants from the
// PCG reference design by Melissa E. O'Neill.
const (
	pcg64MultHi = 0x2360ed051fc65da4
	pcg64MultLo = 0x4385df649fccf645
	pcg64IncHi  = 0x5851f42d4c957f2d
	pcg64IncLo  = 0x14057b7ef767814f
)

// pcg64 holds its 128-bit state as two 64-bit limbs.
type pcg64 struct {
	*core
	hi, lo uint64
}

func newPCG64(seed uint64) baseRandom {
	g := &pcg64{}
	sm := seed
	g.hi = splitMixNext(&sm)
	g.lo = splitMixNext(&sm)
	g.core = newCore(g.next)
	return g
}

func (g *pcg64) next() uint64 {
	oldHi, oldLo := g.hi, g.lo
	g.hi, g.lo = pcg64Step(oldHi, oldLo)

	// XSL-RR output over the old state: rotate hi^lo right by the top 6 bits.
	return rotr64(oldHi^oldLo, uint(oldHi>>58))
}

// pcg64Step advances the 128-bit LCG state: state' = state*mult + inc (mod 2^128).
func pcg64Step(hi, lo uint64) (uint64, uint64) {
	// Schoolbook 128x128 multiplication, keeping only the low 128 bits.
	// The hi*multHi term only lands at bit 128+, so it is dropped.
	p1Hi, p1Lo := bits.Mul64(lo, pcg64MultLo)
	p2Hi, p2Lo := bits.Mul64(hi, pcg64MultLo)
	p3Hi, p3Lo := bits.Mul64(lo, pcg64MultHi)

	resLo := p1Lo
	mid, carry := bits.Add64(p1Hi, p2Lo, 0)
	mid, carry = bits.Add64(mid, p3Lo, carry)
	resHi, _ := bits.Add64(p2Hi, p3Hi, carry)

	// Add the 128-bit increment.
	resLo, carry = bits.Add64(resLo, pcg64IncLo, 0)
	resHi, _ = bits.Add64(resHi, pcg64IncHi, carry)

	return resHi, resLo
}
