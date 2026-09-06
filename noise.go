package randomix

import "math"

// NoiseFlavor selects the base lattice algorithm a Noise instance uses.
type NoiseFlavor uint8

const (
	// NoisePerlin is classic Perlin gradient noise.
	NoisePerlin NoiseFlavor = iota
	// NoiseValue is value noise (scalar lattice interpolation).
	NoiseValue
)

// Noise generates procedural noise. The underlying flavor (Perlin or value)
// and any fractal (fBm) settings are fixed at construction time via the
// options passed to Randomix.Noise; the dimensionality is chosen per call.
type Noise interface {
	Noise1D(x float64) float64
	Noise2D(x, y float64) float64
	Noise3D(x, y, z float64) float64
	Noise4D(x, y, z, w float64) float64
}

// NoiseOption tunes a Noise instance.
type NoiseOption func(*noiseConfig)

type noiseConfig struct {
	flavor      NoiseFlavor
	octaves     int
	lacunarity  float64
	persistence float64
	frequency   float64
	hasRange    bool
	lo, hi      float64
}

// WithFlavor selects Perlin or value noise. Default: NoisePerlin.
func WithFlavor(f NoiseFlavor) NoiseOption {
	return func(c *noiseConfig) { c.flavor = f }
}

// WithOctaves enables fractal (fBm) summation. 1 disables it. Default: 1.
func WithOctaves(n int) NoiseOption {
	return func(c *noiseConfig) { c.octaves = n }
}

// WithLacunarity sets the frequency multiplier between octaves. Default: 2.
func WithLacunarity(l float64) NoiseOption {
	return func(c *noiseConfig) { c.lacunarity = l }
}

// WithPersistence sets the amplitude multiplier between octaves. Default: 0.5.
func WithPersistence(p float64) NoiseOption {
	return func(c *noiseConfig) { c.persistence = p }
}

// WithFrequency scales input coordinates before sampling. Default: 1.
func WithFrequency(f float64) NoiseOption {
	return func(c *noiseConfig) { c.frequency = f }
}

// WithRange remaps the roughly [-1, 1] output onto [lo, hi].
func WithRange(lo, hi float64) NoiseOption {
	return func(c *noiseConfig) {
		c.hasRange = true
		c.lo, c.hi = lo, hi
	}
}

func defaultNoiseConfig() noiseConfig {
	return noiseConfig{
		flavor:      NoisePerlin,
		octaves:     1,
		lacunarity:  2,
		persistence: 0.5,
		frequency:   1,
	}
}

func (c noiseConfig) validate() {
	if c.octaves < 1 {
		panic("randomix: WithOctaves must be >= 1")
	}
	if c.lacunarity <= 0 || math.IsNaN(c.lacunarity) {
		panic("randomix: WithLacunarity must be > 0")
	}
	if c.persistence <= 0 || math.IsNaN(c.persistence) {
		panic("randomix: WithPersistence must be > 0")
	}
	if c.frequency <= 0 || math.IsNaN(c.frequency) {
		panic("randomix: WithFrequency must be > 0")
	}
	if c.hasRange && c.lo >= c.hi {
		panic("randomix: WithRange requires lo < hi")
	}
	if c.flavor != NoisePerlin && c.flavor != NoiseValue {
		panic("randomix: unknown noise flavor")
	}
}

// Noise returns a new Noise whose permutation table is drawn from this
// Randomix's PRNG stream. For reproducible fields, create the Noise before
// drawing anything else from the Randomix (the table consumes 256 draws).
func (r *randomix) Noise(opts ...NoiseOption) Noise {
	cfg := defaultNoiseConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	cfg.validate()

	perm := r.Perm(256)
	tbl := make([]int, 512)
	for i := range 512 {
		tbl[i] = perm[i&255]
	}

	return &noise{tbl: tbl, cfg: cfg}
}

type noise struct {
	tbl []int
	cfg noiseConfig
}

// Noise1D samples the configured noise at x.
func (n *noise) Noise1D(x float64) float64 {
	c := [4]float64{x * n.cfg.frequency}
	return n.apply(n.sample(1, c))
}

// Noise2D samples the configured noise at (x, y).
func (n *noise) Noise2D(x, y float64) float64 {
	c := [4]float64{x * n.cfg.frequency, y * n.cfg.frequency}
	return n.apply(n.sample(2, c))
}

// Noise3D samples the configured noise at (x, y, z).
func (n *noise) Noise3D(x, y, z float64) float64 {
	c := [4]float64{x * n.cfg.frequency, y * n.cfg.frequency, z * n.cfg.frequency}
	return n.apply(n.sample(3, c))
}

// Noise4D samples the configured noise at (x, y, z, w).
func (n *noise) Noise4D(x, y, z, w float64) float64 {
	c := [4]float64{x * n.cfg.frequency, y * n.cfg.frequency, z * n.cfg.frequency, w * n.cfg.frequency}
	return n.apply(n.sample(4, c))
}

// PerlinNoise returns a classic 2D Perlin value in roughly [-1, 1]. It is a
// convenience over Noise() and stays deterministic for a given Randomix seed:
// the underlying default Noise is built lazily on first call and cached.
func (r *randomix) PerlinNoise(x, y float64) float64 {
	return r.defaultNoise().Noise2D(x, y)
}

func (r *randomix) defaultNoise() Noise {
	r.noiseOnce.Do(func() {
		r.noise = r.Noise()
	})
	return r.noise
}

// sample returns the raw (pre-range-remap) noise value. When octaves > 1 it
// sums the flavor over increasing lacunarity, normalizing by total amplitude
// so the output stays approximately within [-1, 1].
func (n *noise) sample(d int, c [4]float64) float64 {
	if n.cfg.octaves == 1 {
		return n.single(d, c)
	}

	total, amp, norm := 0.0, 1.0, 0.0
	for range n.cfg.octaves {
		total += amp * n.single(d, c)
		norm += amp
		amp *= n.cfg.persistence
		for i := range d {
			c[i] *= n.cfg.lacunarity
		}
	}
	return total / norm
}

func (n *noise) single(d int, c [4]float64) float64 {
	var cell [4]int64
	var f, u [4]float64
	for i := range d {
		cell[i] = int64(math.Floor(c[i]))
		f[i] = c[i] - float64(cell[i])
		u[i] = fade(f[i])
	}

	if n.cfg.flavor == NoiseValue {
		return n.valueLattice(d, cell, u, f)
	}
	return n.perlinLattice(d, cell, u, f)
}

// valueLattice interpolates scalar random lattice values with a fade-weighted
// hypercube sum.
func (n *noise) valueLattice(d int, cell [4]int64, u, f [4]float64) float64 {
	sum := 0.0
	for mask := 0; mask < 1<<d; mask++ {
		var corner [4]int64
		w := 1.0
		for i := range d {
			if mask>>i&1 == 0 {
				corner[i] = cell[i]
				w *= 1 - u[i]
			} else {
				corner[i] = cell[i] + 1
				w *= u[i]
			}
		}
		h := n.latticeHash(d, corner)
		sum += w * n.latticeValue(h)
	}
	return sum
}

// perlinLattice interpolates gradient noise with a fade-weighted hypercube sum.
func (n *noise) perlinLattice(d int, cell [4]int64, u, f [4]float64) float64 {
	sum := 0.0
	for mask := 0; mask < 1<<d; mask++ {
		var corner [4]int64
		var dist [4]float64
		w := 1.0
		for i := range d {
			if mask>>i&1 == 0 {
				corner[i] = cell[i]
				dist[i] = f[i]
				w *= 1 - u[i]
			} else {
				corner[i] = cell[i] + 1
				dist[i] = f[i] - 1
				w *= u[i]
			}
		}
		h := n.latticeHash(d, corner)
		sum += w * n.gradientDot(d, h, dist)
	}
	return sum
}

// latticeHash folds corner coordinates through the permutation table,
// returning a value in [0, 255].
func (n *noise) latticeHash(d int, corner [4]int64) int {
	v := 0
	for i := range d {
		v = n.tbl[v+int(corner[i]&255)]
	}
	return v
}

// latticeValue maps a permutation entry onto [-1, 1] for value noise.
func (n *noise) latticeValue(h int) float64 {
	return (2*float64(h) - 255) / 255
}

// gradientDot returns the dot product of a unit gradient direction derived
// from hash h with the corner distance vector dist.
func (n *noise) gradientDot(d int, h int, dist [4]float64) float64 {
	dir := unitDirection(d, h)
	dot := 0.0
	for i := range d {
		dot += dir[i] * dist[i]
	}
	return dot
}

// unitDirection derives a unit vector in d dimensions from a hash using a
// deterministic scramble over {-1, 0, 1} components.
func unitDirection(d, h int) [4]float64 {
	var dir [4]float64
	x := uint64(h)
	nonZero := 0
	for i := range d {
		x ^= x << 13
		x ^= x >> 7
		x ^= x << 17
		switch x & 3 {
		case 0:
			dir[i] = -1
		case 2:
			dir[i] = 1
		default:
			dir[i] = 0
		}
		if dir[i] != 0 {
			nonZero++
		}
	}
	if nonZero == 0 {
		dir[0] = 1
		nonZero = 1
	}
	scale := 1 / math.Sqrt(float64(nonZero))
	for i := range d {
		dir[i] *= scale
	}
	return dir
}

// apply remaps the raw sample to the configured output range, if any.
func (n *noise) apply(v float64) float64 {
	if !n.cfg.hasRange {
		return v
	}
	if v < -1 {
		v = -1
	}
	if v > 1 {
		v = 1
	}
	return n.cfg.lo + (v+1)/2*(n.cfg.hi-n.cfg.lo)
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}
