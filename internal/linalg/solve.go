package linalg

import "errors"

// ErrSingular is returned when the linear system has no unique solution
// (the matrix is numerically singular after pivoting).
var ErrSingular = errors.New("linalg: singular matrix")

// Solve returns x such that A·x = b, using Gaussian elimination with partial
// pivoting. The input matrix is copied, so callers may reuse A afterwards.
// When A is (numerically) singular ErrSingular is returned.
func Solve(a Mat, b Vec) (Vec, error) {
	n := len(b)
	if len(a) != n {
		panic("linalg: solve dim mismatch")
	}
	// Augmented matrix copy.
	m := NewMat(n, n+1)
	for i := 0; i < n; i++ {
		copy(m[i], a[i])
		m[i][n] = b[i]
	}
	const eps = 1e-12
	for col := 0; col < n; col++ {
		pivot := col
		max := abs(m[col][col])
		for r := col + 1; r < n; r++ {
			if v := abs(m[r][col]); v > max {
				max = v
				pivot = r
			}
		}
		if max < eps {
			return nil, ErrSingular
		}
		m[col], m[pivot] = m[pivot], m[col]
		diag := m[col][col]
		for r := col + 1; r < n; r++ {
			f := m[r][col] / diag
			if f == 0 {
				continue
			}
			row := m[r]
			prow := m[col]
			for c := col; c <= n; c++ {
				row[c] -= f * prow[c]
			}
		}
	}
	x := NewVec(n)
	for i := n - 1; i >= 0; i-- {
		sum := m[i][n]
		for j := i + 1; j < n; j++ {
			sum -= m[i][j] * x[j]
		}
		x[i] = sum / m[i][i]
	}
	return x, nil
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
