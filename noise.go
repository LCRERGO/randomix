package randomix

import "math"

// PerlinNoise generates 2D Perlin noise for the given (x, y) coordinates.
// The result is a value in the range [-1, 1] and is deterministic for a given
// Randomix seed.
func (r *randomix) PerlinNoise(x, y float64) float64 {
	r.permOnce.Do(func() {
		r.perm = make([]int, 512)
		p := r.Perm(256)
		for i := range 512 {
			r.perm[i] = p[i&255]
		}
	})

	X := int(math.Floor(x)) & 255
	Y := int(math.Floor(y)) & 255
	xf := x - math.Floor(x)
	yf := y - math.Floor(y)
	u := xf * xf * (3 - 2*xf)
	v := yf * yf * (3 - 2*yf)

	aa := r.perm[r.perm[X]+Y]
	ab := r.perm[r.perm[X]+Y+1]
	ba := r.perm[r.perm[X+1]+Y]
	bb := r.perm[r.perm[X+1]+Y+1]

	x1 := lerp(grad(aa, xf, yf), grad(ba, xf-1, yf), u)
	x2 := lerp(grad(ab, xf, yf-1), grad(bb, xf-1, yf-1), u)

	return lerp(x1, x2, v)
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

func grad(hash int, x, y float64) float64 {
	h := hash & 3
	u := y
	if h&1 == 0 {
		u = x
	}
	v := 2 * x
	if (h & 1) == 0 {
		v = 2 * y
	}
	return u + v
}
