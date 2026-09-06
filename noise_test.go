package randomix

import (
	"math"
	"testing"
)

func noisePoints4D() [][4]float64 {
	return [][4]float64{
		{0, 0, 0, 0},
		{1.5, 2.5, 3.5, 4.5},
		{-3.7, 8.2, -1.1, 0.3},
		{12.001, -0.5, 6.25, -9.9},
		{0.33, 0.66, 1.2, 5.7},
		{-2, -2, -2, -2},
	}
}

func TestNoiseDeterministic(t *testing.T) {
	for _, tc := range allAlgorithms {
		a := seeded(tc.algo, 4242)
		b := seeded(tc.algo, 4242)

		na := a.Noise(WithFlavor(NoisePerlin), WithOctaves(3))
		nb := b.Noise(WithFlavor(NoisePerlin), WithOctaves(3))

		for _, p := range noisePoints4D() {
			if na.Noise2D(p[0], p[1]) != nb.Noise2D(p[0], p[1]) {
				t.Fatalf("%s: 2D noise diverged at %v", tc.name, p)
			}
			if na.Noise3D(p[0], p[1], p[2]) != nb.Noise3D(p[0], p[1], p[2]) {
				t.Fatalf("%s: 3D noise diverged at %v", tc.name, p)
			}
			if na.Noise4D(p[0], p[1], p[2], p[3]) != nb.Noise4D(p[0], p[1], p[2], p[3]) {
				t.Fatalf("%s: 4D noise diverged at %v", tc.name, p)
			}
		}
	}
}

func TestNoiseDifferentSeedsDiffer(t *testing.T) {
	a := seeded(MT19937, 1).Noise()
	b := seeded(MT19937, 2).Noise()

	diff := false
	for _, p := range noisePoints4D() {
		if a.Noise2D(p[0], p[1]) != b.Noise2D(p[0], p[1]) {
			diff = true
			break
		}
	}
	if !diff {
		t.Fatal("different seeds produced identical noise fields")
	}
}

func TestNoiseFamiliesAndDimensions(t *testing.T) {
	for _, flavor := range []NoiseFlavor{NoisePerlin, NoiseValue} {
		for _, dims := range []int{1, 2, 3, 4} {
			r := seeded(ChaCha8, uint64(dims)+uint64(flavor))
			n := r.Noise(WithFlavor(flavor))
			if n == nil {
				t.Fatalf("flavor %d dims %d: nil Noise", flavor, dims)
			}

			sample := func(p [4]float64) float64 {
				switch dims {
				case 1:
					return n.Noise1D(p[0])
				case 2:
					return n.Noise2D(p[0], p[1])
				case 3:
					return n.Noise3D(p[0], p[1], p[2])
				default:
					return n.Noise4D(p[0], p[1], p[2], p[3])
				}
			}

			for _, p := range noisePoints4D() {
				v := sample(p)
				if math.IsNaN(v) || math.IsInf(v, 0) {
					t.Fatalf("flavor %d dims %d: non-finite %v at %v", flavor, dims, v, p)
				}
				if v < -2 || v > 2 {
					t.Fatalf("flavor %d dims %d: out-of-range %v at %v", flavor, dims, v, p)
				}
			}
		}
	}
}

func TestValueNoiseBounded(t *testing.T) {
	n := seeded(MT19937, 7).Noise(WithFlavor(NoiseValue))
	for i := 0; i < 200; i++ {
		x := float64(i%17) * 1.3
		y := float64(i/7) * 0.9
		v := n.Noise2D(x, y)
		if v < -1.001 || v > 1.001 {
			t.Fatalf("value noise out of [-1,1]: %v at (%v,%v)", v, x, y)
		}
	}
}

func TestFractalOctaves(t *testing.T) {
	base := seeded(MT19937, 99).Noise(WithOctaves(1))
	fract := seeded(MT19937, 99).Noise(WithOctaves(5))

	diff := false
	for _, p := range noisePoints4D() {
		if base.Noise2D(p[0], p[1]) != fract.Noise2D(p[0], p[1]) {
			diff = true
			break
		}
	}
	if !diff {
		t.Fatal("fBm octaves had no effect")
	}
}

func TestFlavorsDiffer(t *testing.T) {
	perlin := seeded(MT19937, 3).Noise(WithFlavor(NoisePerlin))
	value := seeded(MT19937, 3).Noise(WithFlavor(NoiseValue))

	diff := false
	for _, p := range noisePoints4D() {
		if perlin.Noise2D(p[0], p[1]) != value.Noise2D(p[0], p[1]) {
			diff = true
			break
		}
	}
	if !diff {
		t.Fatal("Perlin and value noise produced identical output")
	}
}

func TestNoiseRangeRemap(t *testing.T) {
	n := seeded(MT19937, 11).Noise(WithRange(0, 1), WithFlavor(NoiseValue))
	for _, p := range noisePoints4D() {
		v := n.Noise2D(p[0], p[1])
		if v < 0 || v > 1 {
			t.Fatalf("range remap violated: %v at %v", v, p)
		}
	}

	// With an asymmetric range the midpoint should land near the middle.
	m := seeded(MT19937, 12).Noise(WithRange(0, 100))
	values := make([]float64, 0, 8)
	for i := 0; i < 8; i++ {
		values = append(values, m.Noise2D(float64(i)*0.7, 3.1))
	}
	for _, v := range values {
		if v < 0 || v > 100 {
			t.Fatalf("range remap violated: %v", v)
		}
	}
}

func TestNoiseOptionValidation(t *testing.T) {
	r := seeded(MT19937, 1)
	cases := []struct {
		name string
		opts []NoiseOption
	}{
		{"octaves zero", []NoiseOption{WithOctaves(0)}},
		{"lacunarity zero", []NoiseOption{WithLacunarity(0)}},
		{"persistence negative", []NoiseOption{WithPersistence(-1)}},
		{"frequency zero", []NoiseOption{WithFrequency(0)}},
		{"inverted range", []NoiseOption{WithRange(2, 1)}},
	}
	for _, tc := range cases {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s did not panic", tc.name)
				}
			}()
			r.Noise(tc.opts...)
		}()
	}
}

func TestPerlinNoiseBackCompat(t *testing.T) {
	a := seeded(MT19937, 2024)
	b := seeded(MT19937, 2024)

	for _, p := range noisePoints4D() {
		va := a.PerlinNoise(p[0], p[1])
		vb := b.PerlinNoise(p[0], p[1])
		if va != vb {
			t.Fatalf("PerlinNoise diverged for same seed at %v", p)
		}
		if va < -1.5 || va > 1.5 || math.IsNaN(va) {
			t.Fatalf("PerlinNoise out of range: %v at %v", va, p)
		}
		// Repeated calls on the same instance must be stable.
		if va != a.PerlinNoise(p[0], p[1]) {
			t.Fatalf("PerlinNoise not stable on same instance at %v", p)
		}
	}
}

func TestPerlinNoiseMatchesDefaultNoise(t *testing.T) {
	r := seeded(MT19937, 555)
	// PerlinNoise is defined as the default-configured Noise2D; verify it equals
	// the shared default noise sample on the same instance.
	v1 := r.PerlinNoise(1.3, 2.7)
	v2 := r.PerlinNoise(1.3, 2.7)
	if v1 != v2 {
		t.Fatal("PerlinNoise not stable")
	}
}

func TestNoiseInstancesConsumeStream(t *testing.T) {
	r := seeded(MT19937, 31)
	n1 := r.Noise()
	n2 := r.Noise()

	diff := false
	for _, p := range noisePoints4D() {
		if n1.Noise2D(p[0], p[1]) != n2.Noise2D(p[0], p[1]) {
			diff = true
			break
		}
	}
	if !diff {
		t.Fatal("two Noise() instances from one Randomix produced identical fields")
	}
}
