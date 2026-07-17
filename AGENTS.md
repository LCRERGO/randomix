# randomix

Go library for generating random data (CPF, CNPJ, email, UUID, etc.).

- **Module:** `github.com/LCRERGO/randomix`
- **Requires Go 1.24.5** (uses `math/rand/v2`)
- **Entry point:** `randomix.NewRandomix(src)` returns a `Randomix` interface
- **Single package** — all files are `package randomix`
- **`noise.go` is intentionally empty** (placeholder)
- **No test files exist yet**

## Commands

```bash
go build ./# Build
go test ./# Run tests (will pass — zero tests)
go fmt ./# Format
go vet ./# Lint
```
