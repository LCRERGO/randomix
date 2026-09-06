package randomix

// chacha8 is an 8-round ChaCha keystream generator (4 double rounds).
// It is deterministic and fast but here it is NOT cryptographically secure:
// it is seeded from a single 64-bit value and is only meant as a drop-in
// algorithm variant. Do not use it for secrets.
type chacha8 struct {
	*core
	key [8]uint32
	seq uint64
	buf [8]uint64
	pos int
}

func newChaCha8(seed uint64) baseRandom {
	g := &chacha8{}
	sm := seed
	for i := 0; i < len(g.key); i += 2 {
		g.key[i] = uint32(splitMixNext(&sm))
		g.key[i+1] = uint32(splitMixNext(&sm) >> 32)
	}
	g.core = newCore(g.next)
	g.pos = len(g.buf)
	return g
}

func (g *chacha8) next() uint64 {
	if g.pos >= len(g.buf) {
		g.refill()
	}
	v := g.buf[g.pos]
	g.pos++
	return v
}

func (g *chacha8) refill() {
	var s [16]uint32
	s[0] = 0x61707865
	s[1] = 0x3320646e
	s[2] = 0x79622d32
	s[3] = 0x6b206574
	copy(s[4:12], g.key[:])
	s[12] = uint32(g.seq)
	s[13] = uint32(g.seq >> 32)

	initial := s
	for range 4 { // ChaCha8 = 4 double rounds = 8 rounds
		// Column round.
		chachaQR(&s[0], &s[4], &s[8], &s[12])
		chachaQR(&s[1], &s[5], &s[9], &s[13])
		chachaQR(&s[2], &s[6], &s[10], &s[14])
		chachaQR(&s[3], &s[7], &s[11], &s[15])
		// Diagonal round.
		chachaQR(&s[0], &s[5], &s[10], &s[15])
		chachaQR(&s[1], &s[6], &s[11], &s[12])
		chachaQR(&s[2], &s[7], &s[8], &s[13])
		chachaQR(&s[3], &s[4], &s[9], &s[14])
	}
	for i := range s {
		s[i] += initial[i]
	}
	for i := range 8 {
		g.buf[i] = uint64(s[2*i]) | uint64(s[2*i+1])<<32
	}
	g.seq++
	g.pos = 0
}

// chachaQR is the ChaCha quarter round on words a, b, c, d.
func chachaQR(a, b, c, d *uint32) {
	*a += *b
	*d ^= *a
	*d = *d<<16 | *d>>16
	*c += *d
	*b ^= *c
	*b = *b<<12 | *b>>20
	*a += *b
	*d ^= *a
	*d = *d<<8 | *d>>24
	*c += *d
	*b ^= *c
	*b = *b<<7 | *b>>25
}
