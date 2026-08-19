// Package linalg provides a tiny dense linear-algebra toolkit used by the
// frame solver: vectors, row-major matrices, and a Gaussian-elimination
// solver with partial pivoting. It is intentionally dependency-free so the
// solver stays in pure Go.
package linalg

// Vec is a column vector of float64 values.
type Vec []float64

// NewVec allocates a zero vector of length n.
func NewVec(n int) Vec { return make(Vec, n) }

// Len returns the vector length.
func (a Vec) Len() int { return len(a) }

// Sum returns the sum of all components.
func (a Vec) Sum() float64 {
	var s float64
	for _, v := range a {
		s += v
	}
	return s
}

// MaxAbs returns the largest absolute component.
func (a Vec) MaxAbs() float64 {
	m := 0.0
	for _, v := range a {
		if v < 0 {
			v = -v
		}
		if v > m {
			m = v
		}
	}
	return m
}

// Add returns the elementwise sum a+b.
func (a Vec) Add(b Vec) Vec {
	if len(a) != len(b) {
		panic("linalg: vec add length mismatch")
	}
	out := make(Vec, len(a))
	for i := range a {
		out[i] = a[i] + b[i]
	}
	return out
}

// Sub returns the elementwise difference a-b.
func (a Vec) Sub(b Vec) Vec {
	if len(a) != len(b) {
		panic("linalg: vec sub length mismatch")
	}
	out := make(Vec, len(a))
	for i := range a {
		out[i] = a[i] - b[i]
	}
	return out
}

// Scale returns s*a.
func (a Vec) Scale(s float64) Vec {
	out := make(Vec, len(a))
	for i := range a {
		out[i] = a[i] * s
	}
	return out
}

// Dot returns the inner product a·b.
func (a Vec) Dot(b Vec) float64 {
	if len(a) != len(b) {
		panic("linalg: vec dot length mismatch")
	}
	var s float64
	for i := range a {
		s += a[i] * b[i]
	}
	return s
}

// Copy returns an independent copy of the vector.
func (a Vec) Copy() Vec {
	out := make(Vec, len(a))
	copy(out, a)
	return out
}
