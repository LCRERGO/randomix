package randomix

type baseRandom interface {
	ExpFloat64() float64
	Float32() float32
	Float64() float64
	Int() int
	Int32() int32
	Int32N(n int32) int32
	Int64() int64
	Int64N(n int64) int64
	IntN(n int) int
	NormFloat64() float64
	Perm(n int) []int
	Shuffle(n int, swap func(i, j int))
	Uint() uint
	Uint32() uint32
	Uint32N(n uint32) uint32
	Uint64() uint64
	Uint64N(n uint64) uint64
	UintN(n uint) uint
}

func (r *randomix) ExpFloat64() float64 { return r.randomGenerator.ExpFloat64() }

func (r *randomix) Float32() float32 { return r.randomGenerator.Float32() }

func (r *randomix) Float64() float64 { return r.randomGenerator.Float64() }

func (r *randomix) Int() int { return r.randomGenerator.Int() }

func (r *randomix) Int32() int32 { return r.randomGenerator.Int32() }

func (r *randomix) Int32N(n int32) int32 { return r.randomGenerator.Int32N(n) }

func (r *randomix) Int64() int64 { return r.randomGenerator.Int64() }

func (r *randomix) Int64N(n int64) int64 { return r.randomGenerator.Int64N(n) }

func (r *randomix) IntN(n int) int { return r.randomGenerator.IntN(n) }

func (r *randomix) NormFloat64() float64 { return r.randomGenerator.NormFloat64() }

func (r *randomix) Perm(n int) []int { return r.randomGenerator.Perm(n) }

func (r *randomix) Shuffle(n int, swap func(i int, j int)) { r.randomGenerator.Shuffle(n, swap) }

func (r *randomix) Uint() uint { return r.randomGenerator.Uint() }

func (r *randomix) Uint32() uint32 { return r.randomGenerator.Uint32() }

func (r *randomix) Uint32N(n uint32) uint32 { return r.randomGenerator.Uint32N(n) }

func (r *randomix) Uint64() uint64 { return r.randomGenerator.Uint64() }

func (r *randomix) Uint64N(n uint64) uint64 { return r.randomGenerator.Uint64N(n) }

func (r *randomix) UintN(n uint) uint { return r.randomGenerator.UintN(n) }
