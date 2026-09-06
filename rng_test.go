package randomix

import (
	"math"
	"testing"
)

var allAlgorithms = []struct {
	name string
	algo Algorithm
}{
	{"MT19937", MT19937},
	{"LCG", LCG},
	{"SplitMix64", SplitMix64},
	{"Xoshiro256", Xoshiro256},
	{"PCG64", PCG64},
	{"ChaCha8", ChaCha8},
}

func seeded(algo Algorithm, seed uint64) Randomix {
	return NewRandomix(WithAlgorithm(algo), WithSeed(seed))
}

func TestMT19937ReferenceVector(t *testing.T) {
	// First tempered word of init_genrand(5489) is 3499211612 (mt19937ar.c).
	r := seeded(MT19937, 5489)
	if got := r.Uint64() >> 32; got != 3499211612 {
		t.Errorf("MT19937(5489) first word = %d, want 3499211612", got)
	}
}

func TestSplitMix64ReferenceVector(t *testing.T) {
	r := seeded(SplitMix64, 0)
	if got := r.Uint64(); got != 0xe220a8397b1dcdaf {
		t.Errorf("SplitMix64(0) first = %#x, want 0xe220a8397b1dcdaf", got)
	}
}

func TestLCGReferenceVector(t *testing.T) {
	r := seeded(LCG, 42)
	s := uint64(42)
	want := s*lcgMultiplier + lcgIncrement
	if got := r.Uint64(); got != want {
		t.Errorf("LCG(42) first = %d, want %d", got, want)
	}
	want = want*lcgMultiplier + lcgIncrement
	if got := r.Uint64(); got != want {
		t.Errorf("LCG(42) second = %d, want %d", got, want)
	}
}

func TestDefaultIsMT19937Seeded(t *testing.T) {
	a := NewRandomix(WithSeed(1234))
	b := seeded(MT19937, 1234)
	for range 16 {
		if a.Uint64() != b.Uint64() {
			t.Fatal("default algorithm is not MT19937")
		}
	}
}

func TestDeterministicPerSeed(t *testing.T) {
	for _, tc := range allAlgorithms {
		a := seeded(tc.algo, 987654321)
		b := seeded(tc.algo, 987654321)
		ca := make([]uint64, 32)
		cb := make([]uint64, 32)
		for i := range ca {
			ca[i] = a.Uint64()
			cb[i] = b.Uint64()
		}
		for i := range ca {
			if ca[i] != cb[i] {
				t.Fatalf("%s: same seed diverged at draw %d", tc.name, i)
			}
		}
		if a.IntN(100) != b.IntN(100) || a.Float64() != b.Float64() {
			t.Fatalf("%s: derived draws diverged", tc.name)
		}
	}
}

func TestDifferentSeedsDiffer(t *testing.T) {
	for _, tc := range allAlgorithms {
		a := seeded(tc.algo, 1)
		b := seeded(tc.algo, 2)
		diff := false
		for range 8 {
			if a.Uint64() != b.Uint64() {
				diff = true
				break
			}
		}
		if !diff {
			t.Errorf("%s: different seeds produced identical streams", tc.name)
		}
	}
}

func TestAlgorithmsAreDistinct(t *testing.T) {
	first := make([]uint64, 0, len(allAlgorithms))
	for _, tc := range allAlgorithms {
		first = append(first, seeded(tc.algo, 42).Uint64())
	}
	seen := map[uint64]bool{}
	for _, v := range first {
		seen[v] = true
	}
	if len(seen) < len(allAlgorithms) {
		t.Fatalf("expected %d distinct first values, got %d", len(allAlgorithms), len(seen))
	}
}

func TestRandomixAutoSeeds(t *testing.T) {
	a := NewRandomix(WithAlgorithm(MT19937))
	b := NewRandomix(WithAlgorithm(MT19937))
	if a.Uint64() == b.Uint64() {
		t.Fatal("two auto-seeded instances produced the same stream")
	}
}

func TestUniformBounds(t *testing.T) {
	for _, tc := range allAlgorithms {
		t.Run(tc.name, func(t *testing.T) {
			r := seeded(tc.algo, 7)

			for range 2000 {
				if v := r.Float64(); v < 0 || v >= 1 {
					t.Fatalf("Float64 out of range: %v", v)
				}
				if v := r.Float32(); v < 0 || v >= 1 {
					t.Fatalf("Float32 out of range: %v", v)
				}
				if v := r.IntN(10); v < 0 || v >= 10 {
					t.Fatalf("IntN out of range: %d", v)
				}
				if v := r.Int32N(10); v < 0 || v >= 10 {
					t.Fatalf("Int32N out of range: %d", v)
				}
				if v := r.Int64N(10); v < 0 || v >= 10 {
					t.Fatalf("Int64N out of range: %d", v)
				}
				if v := r.UintN(10); v >= 10 {
					t.Fatalf("UintN out of range: %d", v)
				}
				if v := r.Uint32N(10); v >= 10 {
					t.Fatalf("Uint32N out of range: %d", v)
				}
				if v := r.Uint64N(10); v >= 10 {
					t.Fatalf("Uint64N out of range: %d", v)
				}
				if v := r.Int(); v < 0 {
					t.Fatalf("Int negative: %d", v)
				}
				if v := r.Int32(); v < 0 {
					t.Fatalf("Int32 negative: %d", v)
				}
				if v := r.Int64(); v < 0 {
					t.Fatalf("Int64 negative: %d", v)
				}
			}

			if v := r.IntN(1); v != 0 {
				t.Fatalf("IntN(1) = %d, want 0", v)
			}
			if v := r.Uint64N(1); v != 0 {
				t.Fatalf("Uint64N(1) = %d, want 0", v)
			}

			// Every residue of a small modulus should be reachable.
			seen := [7]bool{}
			for range 5000 {
				seen[r.Uint64N(7)] = true
			}
			for res, ok := range seen {
				if !ok {
					t.Fatalf("residue %d never produced by Uint64N(7)", res)
				}
			}

			if v := r.ExpFloat64(); v < 0 || math.IsInf(v, 1) {
				t.Fatalf("ExpFloat64 out of range: %v", v)
			}
			for range 100 {
				if v := r.NormFloat64(); math.IsNaN(v) || math.IsInf(v, 0) {
					t.Fatalf("NormFloat64 not finite: %v", v)
				}
			}

			perm := r.Perm(10)
			sum := 0
			for _, v := range perm {
				sum += v
			}
			if sum != 45 {
				t.Fatalf("Perm(10) is not a permutation: %v", perm)
			}

			arr := make([]int, 12)
			for i := range arr {
				arr[i] = i
			}
			r.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
			sum = 0
			for _, v := range arr {
				sum += v
			}
			if sum != 66 {
				t.Fatalf("Shuffle corrupted the slice: %v", arr)
			}
		})
	}
}

func TestFloatMeans(t *testing.T) {
	for _, tc := range allAlgorithms {
		r := seeded(tc.algo, 99)
		var sum float64
		const n = 50000
		for range n {
			sum += r.Float64()
		}
		mean := sum / n
		if mean < 0.47 || mean > 0.53 {
			t.Errorf("%s: Float64 mean = %.4f, want ~0.5", tc.name, mean)
		}
	}
}

func TestZeroArgumentPanics(t *testing.T) {
	for _, tc := range allAlgorithms {
		t.Run(tc.name, func(t *testing.T) {
			r := seeded(tc.algo, 3)
			for name, fn := range map[string]func(){
				"UintN":   func() { r.UintN(0) },
				"Uint32N": func() { r.Uint32N(0) },
				"Uint64N": func() { r.Uint64N(0) },
				"IntN":    func() { r.IntN(0) },
				"Int32N":  func() { r.Int32N(0) },
				"Int64N":  func() { r.Int64N(0) },
			} {
				func() {
					defer func() {
						if recover() == nil {
							t.Errorf("%s did not panic on n=0", name)
						}
					}()
					fn()
				}()
			}
		})
	}
}

func TestRandomPassword(t *testing.T) {
	if got := NewRandomix().RandomPassword(0); got != "" {
		t.Fatalf("RandomPassword(0) = %q, want empty", got)
	}

	r := seeded(MT19937, 5)
	pass := r.RandomPassword(64)
	if len(pass) != 64 {
		t.Fatalf("RandomPassword length = %d, want 64", len(pass))
	}
	valid := map[byte]bool{}
	for i := range passwordCharset {
		valid[passwordCharset[i]] = true
	}
	for i := 0; i < len(pass); i++ {
		if !valid[pass[i]] {
			t.Fatalf("password contains out-of-charset byte %q", pass[i])
		}
	}
}

func TestRandomDataBuilders(t *testing.T) {
	for _, tc := range allAlgorithms {
		r := seeded(tc.algo, 2024)
		if len(r.RandomIPv4()) != 16 {
			t.Errorf("%s: RandomIPv4 wrong length", tc.name)
		}
		if len(r.RandomIPv6()) != 16 {
			t.Errorf("%s: RandomIPv6 wrong length", tc.name)
		}
		if len(r.RandomMAC()) != 6 {
			t.Errorf("%s: RandomMAC wrong length", tc.name)
		}
		if len(r.RandomEmail()) == 0 {
			t.Errorf("%s: RandomEmail empty", tc.name)
		}
		if len(r.RandomPhoneNumber()) == 0 {
			t.Errorf("%s: RandomPhoneNumber empty", tc.name)
		}
		cpf := r.RandomCPF()
		if len(cpf) != 14 {
			t.Errorf("%s: RandomCPF %q not 14 chars", tc.name, cpf)
		}
		if len(r.RandomCNPJ()) != 18 {
			t.Errorf("%s: RandomCNPJ not 18 chars", tc.name)
		}
		if len(r.RandomUsername()) == 0 {
			t.Errorf("%s: RandomUsername empty", tc.name)
		}
		for i := 1; i <= 3; i++ {
			if v := r.PerlinNoise(float64(i), 2.5); v < -1.001 || v > 1.001 {
				t.Errorf("%s: PerlinNoise out of [-1,1]: %v", tc.name, v)
			}
		}
	}
}
