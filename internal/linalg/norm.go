package linalg

import "math"

// NormInf returns the maximum absolute component of the vector.
func (a Vec) NormInf() float64 {
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

// Norm2 returns the Euclidean norm of the vector.
func (a Vec) Norm2() float64 {
	var s float64
	for _, v := range a {
		s += v * v
	}
	return math.Sqrt(s)
}

// AddScaled returns a + s*b.
func (a Vec) AddScaled(b Vec, s float64) Vec {
	if len(a) != len(b) {
		panic("linalg: vec addscaled length mismatch")
	}
	out := make(Vec, len(a))
	for i := range a {
		out[i] = a[i] + s*b[i]
	}
	return out
}

// Frobenius returns the Frobenius norm of the matrix (root of the sum of
// squares of all entries).
func (a Mat) Frobenius() float64 {
	var s float64
	for _, row := range a {
		for _, v := range row {
			s += v * v
		}
	}
	return math.Sqrt(s)
}

// MaxAbs returns the largest absolute entry of the matrix.
func (a Mat) MaxAbs() float64 {
	m := 0.0
	for _, row := range a {
		for _, v := range row {
			if v < 0 {
				v = -v
			}
			if v > m {
				m = v
			}
		}
	}
	return m
}

// RowScale returns a copy of A with row r scaled by s.
func (a Mat) RowScale(r int, s float64) Mat {
	out := NewMat(a.Rows(), a.Cols())
	for i := range a {
		copy(out[i], a[i])
	}
	if r < 0 || r >= len(out) {
		return out
	}
	for j := range out[r] {
		out[r][j] *= s
	}
	return out
}
