# randomix

Go library for generating random data (CPF, CNPJ, email, UUID, etc.).

- **Module:** `github.com/LCRERGO/randomix`
- **Requires Go 1.27**
- **Entry point:** `NewRandomix(opts ...Option) Randomix` with `WithAlgorithm`/`WithSeed`; default is MT19937, crypto-seeded
- **Own PRNGs, no `math/rand`:** MT19937, LCG, SplitMix64, Xoshiro256, PCG64, ChaCha8 (see rngcore.go for the shared derived-ops core)
- **Secrets:** `RandomPassword` uses `crypto/rand`; UUIDs use `github.com/google/uuid` — never the configured PRNG
- **Single package** — all files are `package randomix`
- **Tests exist** in `rng_test.go` (includes reference-vector checks for MT19937/SplitMix64)

## Commands

```bash
go build ./# Build
go test ./# Run tests
go fmt ./# Format
go vet ./# Lint
```
