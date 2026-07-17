package randomix

import (
	"fmt"
	"math/rand/v2"
	"net"
	"sync"

	"github.com/google/uuid"
)

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
	RandomPassword(length int) string
	RandomUUIDv4() uuid.UUID
	RandomUUIDv7() (uuid.UUID, error)
	PerlinNoise(x, y float64) float64
}
type randomix struct {
	randomGenerator *rand.Rand
	perm            []int
	permOnce        sync.Once
}

// NewRandomix creates a new Randomix instance seeded from src.
func NewRandomix(src rand.Source) Randomix {
	generator := rand.New(src)
	return &randomix{
		randomGenerator: generator,
	}
}

// RandomIPv4 returns a random IPv4 address.
func (r *randomix) RandomIPv4() net.IP {
	return net.IPv4(
		byte(r.IntN(256)),
		byte(r.IntN(256)),
		byte(r.IntN(256)),
		byte(r.IntN(256)),
	)
}

// RandomIPv6 returns a random IPv6 address.
func (r *randomix) RandomIPv6() net.IP {
	var ip [16]byte

	r.readBytes(ip[:])

	return net.IP(ip[:])
}

// RandomMAC returns a random MAC address (locally administered, unicast).
func (r *randomix) RandomMAC() net.HardwareAddr {
	mac := make([]byte, 6)
	for i := range mac {
		mac[i] = byte(r.Uint32())
	}

	mac[0] = (mac[0] | 0x02) & 0xFE

	hwAddr := net.HardwareAddr(mac)

	return hwAddr
}

// RandomEmail returns a random email address.
func (r *randomix) RandomEmail() string {
	user := randomString(r, 8)
	domain := randomString(r, 5)
	tld := []string{"com", "net", "org", "io", "dev"}[r.IntN(5)]

	return fmt.Sprintf("%s@%s.%s", user, domain, tld)
}

// RandomPhoneNumber returns a random Brazilian-style phone number.
func (r *randomix) RandomPhoneNumber() string {
	areaCode := r.IntN(90) + 10
	number := fmt.Sprintf("9%04d-%04d", r.IntN(10000), r.IntN(10000))

	return fmt.Sprintf("(%d) %s", areaCode, number)
}

// RandomCPF returns a random valid CPF (Brazilian individual taxpayer ID).
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

// RandomCNPJ returns a random valid CNPJ (Brazilian company ID).
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

// RandomUsername returns a random username like "cool_tiger42".
func (r *randomix) RandomUsername() string {
	adjs := []string{"cool", "fast", "lazy", "smart", "funny", "dark", "loud"}
	nouns := []string{"tiger", "wizard", "robot", "cat", "dog", "pirate", "ninja"}

	return fmt.Sprintf("%s_%s%d",
		adjs[r.IntN(len(adjs))],
		nouns[r.IntN(len(nouns))],
		r.IntN(1000),
	)
}

// RandomPassword returns a random password of the given length.
func (r *randomix) RandomPassword(length int) string {
	const charset = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*`
	pass := make([]byte, length)
	for i := range pass {
		pass[i] = charset[r.IntN(len(charset))]
	}

	return string(pass)
}

// RandomUUIDv4 returns a random UUIDv4.
func (r *randomix) RandomUUIDv4() uuid.UUID {
	return uuid.New()
}

// RandomUUIDv7 returns a random UUIDv7 (time-ordered).
func (r *randomix) RandomUUIDv7() (uuid.UUID, error) {
	return uuid.NewV7()
}
