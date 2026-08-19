package linalg

import (
	"math"
	"testing"
)

func TestSolveDiagonal(t *testing.T) {
	a := NewMat(3, 3)
	a[0][0] = 2
	a[1][1] = 3
	a[2][2] = 4
	b := Vec{4, 9, 8}
	x, err := Solve(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(x[0]-2) > 1e-9 || math.Abs(x[1]-3) > 1e-9 || math.Abs(x[2]-2) > 1e-9 {
		t.Fatalf("Solve = %v", x)
	}
}

func TestSolveGeneral(t *testing.T) {
	a := NewMat(3, 3)
	a[0] = []float64{3, 2, -1}
	a[1] = []float64{2, -2, 4}
	a[2] = []float64{-1, 0.5, -1}
	b := Vec{1, -2, 0}
	x, err := Solve(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Known solution: x=1, y=-2, z=-2.
	want := Vec{1, -2, -2}
	for i := range want {
		if math.Abs(x[i]-want[i]) > 1e-9 {
			t.Fatalf("Solve[%d] = %v, want %v", i, x[i], want[i])
		}
	}
}

func TestSolveSingular(t *testing.T) {
	a := NewMat(2, 2)
	a[0] = []float64{1, 2}
	a[1] = []float64{2, 4} // row 2 = 2 * row 1 -> singular
	b := Vec{3, 6}
	if _, err := Solve(a, b); err != ErrSingular {
		t.Fatalf("expected ErrSingular, got %v", err)
	}
}

func TestSolveResidual(t *testing.T) {
	a := NewMat(4, 4)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			a[i][j] = math.Sin(float64(i+1)) + float64(j) + 1
		}
		a[i][i] += 10
	}
	b := Vec{1, -2, 3, -4}
	x, err := Solve(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	res := a.MulVec(x).Sub(b)
	for i := range res {
		if math.Abs(res[i]) > 1e-9 {
			t.Fatalf("residual[%d] = %v", i, res[i])
		}
	}
}
