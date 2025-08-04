package randomix

import (
	"fmt"
	"math/rand/v2"
	"net"
	"strings"
)

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

type Randomix interface {
	baseRandom

	RandomIPv4() net.IP
	RandomIPv6() net.IP
	RandomMAC() net.HardwareAddr
	RandomEmail() string
	RandomPhoneNumber() string
	RandomCPF() string
	RandomCNPJ() string
	RandomUsername() string
	RandomPassword(int) string
}
type randomix struct {
	randomGenerator *rand.Rand
}

func NewRandomix(src rand.Source) Randomix {
	generator := rand.New(src)
	return &randomix{
		randomGenerator: generator,
	}
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

func (r *randomix) RandomIPv4() net.IP {
	return net.IPv4(
		byte(r.IntN(256)),
		byte(r.IntN(256)),
		byte(r.IntN(256)),
		byte(r.IntN(256)),
	)
}

func (r *randomix) RandomIPv6() net.IP {
	var ip [16]byte

	r.readBytes(ip[:])

	return net.IP(ip[:])
}

func (r *randomix) RandomMAC() net.HardwareAddr {
	mac := make([]byte, 6)
	for i := range mac {
		mac[i] = byte(r.Uint32())
	}

	mac[0] = (mac[0] | 0x02) & 0xFE

	hwAddr := net.HardwareAddr(mac)

	return hwAddr
}

func (r *randomix) RandomEmail() string {
	user := randomString(r, 8)
	domain := randomString(r, 5)
	tld := []string{"com", "net", "org", "io", "dev"}[r.IntN(5)]

	return fmt.Sprintf("%s@%s.%s", user, domain, tld)
}

func (r *randomix) RandomPhoneNumber() string {
	areaCode := r.IntN(90) + 10
	number := fmt.Sprintf("9%04d-%04d", r.IntN(10000), r.IntN(10000))

	return fmt.Sprintf("(%d) %s", areaCode, number)
}

func (r *randomix) RandomCPF() string {
	digits := make([]int, 9)
	for i := range digits {
		digits[i] = r.IntN(10)
	}

	sum := 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * (10 - i)
	}
	d1 := (sum * 10 % 11) % 10

	sum = 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * (11 - i)
	}
	sum += d1 * 2
	d2 := (sum * 10 % 11) % 10

	cpf := append(digits, d1, d2)
	return fmt.Sprintf("%d%d%d.%d%d%d.%d%d%d-%d%d",
		cpf[0], cpf[1], cpf[2], cpf[3], cpf[4], cpf[5],
		cpf[6], cpf[7], cpf[8], cpf[9], cpf[10],
	)
}

func (r *randomix) RandomCNPJ() string {
	base := make([]int, 12)
	for i := 0; i < 8; i++ {
		base[i] = r.IntN(10)
	}
	base[8], base[9], base[10], base[11] = 0, 0, 0, 1

	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 12; i++ {
		sum += base[i] * weights1[i]
	}
	d1 := 11 - sum%11
	if d1 >= 10 {
		d1 = 0
	}

	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum = 0
	for i := 0; i < 12; i++ {
		sum += base[i] * weights2[i]
	}
	sum += d1 * weights2[12]
	d2 := 11 - sum%11
	if d2 >= 10 {
		d2 = 0
	}

	cnpj := append(base, d1, d2)

	return fmt.Sprintf("%d%d.%d%d%d.%d%d%d/%d%d%d%d-%d%d",
		cnpj[0], cnpj[1], cnpj[2], cnpj[3], cnpj[4],
		cnpj[5], cnpj[6], cnpj[7], cnpj[8], cnpj[9],
		cnpj[10], cnpj[11], cnpj[12], cnpj[13],
	)
}

func (r *randomix) RandomUsername() string {
	adjs := []string{"cool", "fast", "lazy", "smart", "funny", "dark", "loud"}
	nouns := []string{"tiger", "wizard", "robot", "cat", "dog", "pirate", "ninja"}

	return fmt.Sprintf("%s_%s%d",
		adjs[r.IntN(len(adjs))],
		nouns[r.IntN(len(nouns))],
		r.IntN(1000),
	)
}

func (r *randomix) RandomPassword(length int) string {
	const charset = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*`
	pass := make([]byte, length)
	for i := range pass {
		pass[i] = charset[r.IntN(len(charset))]
	}

	return string(pass)
}

// utility methods

func (r *randomix) readBytes(buf []byte) {
	for i := range len(buf) {
		buf[i] = byte(r.Uint32())
	}
}

func randomString(r *randomix, length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[r.IntN(len(charset))])
	}

	return sb.String()
}
