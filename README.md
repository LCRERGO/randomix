# randomix

A small Go library for generating random data: Brazilian IDs (CPF/CNPJ), email
addresses, phone numbers, MAC/IP addresses, passwords, UUIDs, procedural noise,
and more — backed by its **own PRNG implementations** (no `math/rand`).

```go
import "github.com/LCRERGO/randomix"
```

## Features

- **Six built-in PRNGs, all implemented in-repo**: MT19937, LCG, SplitMix64,
  xoshiro256**, PCG64, and ChaCha8. No dependency on the stdlib RNG.
- **Pluggable algorithms**: pick the generator per instance with
  `WithAlgorithm`.
- **Deterministic or random**: seed explicitly with `WithSeed`, or let the
  instance auto-seed from `crypto/rand` for a fresh stream every time.
- **Crypto-safe where it matters**: `RandomPassword` draws from `crypto/rand`
  and UUIDs use `github.com/google/uuid` — never the (predictable) PRNG.
- **Procedural noise**: classic Perlin and value noise in 1D–4D, with fractal
  (fBm) octaves, configurable frequency and output range.
- **Familiar API**: the generator surface mirrors `math/rand/v2`
  (`IntN`, `Float64`, `Perm`, `Shuffle`, `NormFloat64`, …).

## Requirements

- Go 1.27+

## Usage

### Quick start

```go
rng := randomix.NewRandomix()

fmt.Println(rng.RandomCPF())       // 529.982.247-25
fmt.Println(rng.RandomCNPJ())      // 11.222.333/0001-81
fmt.Println(rng.RandomEmail())     // x7kd9p2q@m2klo.com
fmt.Println(rng.RandomUsername())  // cool_tiger42
fmt.Println(rng.RandomPhoneNumber()) // (11) 91234-5678
fmt.Println(rng.RandomPassword(16))  // cV5#kL9@pQ2!rT7$x
fmt.Println(rng.RandomUUIDv4())    // 7c9aef1e-... (RFC 4122)
```

### Choosing an algorithm and seed

```go
// Deterministic MT19937.
rng := randomix.NewRandomix(randomix.WithSeed(42))

// The same seed reproduces the exact same stream on every run.
again := randomix.NewRandomix(randomix.WithSeed(42))

// A different algorithm, also deterministic.
lcg := randomix.NewRandomix(randomix.WithAlgorithm(randomix.LCG),
                            randomix.WithSeed(42))
```

| Constant        | Generator                         | Default |
|-----------------|-----------------------------------|---------|
| `MT19937`       | Mersenne Twister (32-bit, 624-word state) | yes      |
| `LCG`           | 64-bit linear congruential (Knuth) |         |
| `SplitMix64`    | SplitMix64                        |         |
| `Xoshiro256`    | xoshiro256**                      |         |
| `PCG64`         | PCG64 XSL-RR (128-bit state)      |         |
| `ChaCha8`       | 8-round ChaCha keystream (not crypto) |         |

`NewRandomix()` with no options returns a crypto-seeded MT19937 instance.

### Generator primitives

Every `Randomix` exposes the `math/rand/v2`-style method set:

```go
n  := rng.IntN(100)          // uniform int in [0, 100)
f  := rng.Float64()          // uniform float in [0, 1)
u  := rng.Uint64N(1 << 40)   // uniform uint64 in [0, 2^40)
g  := rng.NormFloat64()      // standard normal
p  := rng.Perm(10)           // random permutation of 0..9
rng.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
```

`IntN`/`UintN`/etc. use rejection sampling, so ranges have no modulo bias.

### Seeding and determinism

- With `WithSeed`, output is fully reproducible for a given sequence of calls.
- Without it, a seed is pulled from `crypto/rand` each construction.
- Secrecy note: because MT/LCG-class generators are predictable,
  `RandomPassword` and UUIDs intentionally ignore the configured PRNG.

### Procedural noise

Create a `Noise` from a `Randomix`. The permutation table is drawn from the
instance's PRNG stream at creation time, so for reproducible fields create the
`Noise` before drawing anything else from that `Randomix`.

```go
n := rng.Noise(randomix.WithFlavor(randomix.NoisePerlin),
               randomix.WithOctaves(4),
               randomix.WithLacunarity(2),
               randomix.WithPersistence(0.5))

v1 := n.Noise2D(12.3, -4.5)          // in roughly [-1, 1]
v2 := n.Noise3D(1.0, 2.0, 3.0)
v3 := n.Noise4D(0.5, 1.5, 2.5, 3.5)  // great for hash-free smooth noise fields
```

Options:

| Option                 | Default  | Purpose                                  |
|------------------------|----------|------------------------------------------|
| `WithFlavor(flavor)`   | `NoisePerlin` | `NoisePerlin` (gradient) or `NoiseValue` (value noise) |
| `WithOctaves(n)`       | `1`      | fBm octave count; `1` disables fractals  |
| `WithLacunarity(l)`    | `2`      | frequency multiplier between octaves     |
| `WithPersistence(p)`   | `0.5`    | amplitude multiplier between octaves     |
| `WithFrequency(f)`     | `1`      | scales input coordinates                 |
| `WithRange(lo, hi)`    | none     | remaps the roughly [-1, 1] output onto `[lo, hi]` |

The convenience method `PerlinNoise(x, y float64)` still exists on `Randomix`
and is equivalent to a default-configured 2D Perlin `Noise`.

## Installation

```bash
go get github.com/LCRERGO/randomix
```

## License

[MIT](LICENSE)
