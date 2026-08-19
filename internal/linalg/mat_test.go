package linalg

import "testing"

func TestMatMul(t *testing.T) {
	a := NewMat(2, 3)
	a[0] = []float64{1, 2, 3}
	a[1] = []float64{4, 5, 6}
	b := NewMat(3, 2)
	b[0] = []float64{7, 8}
	b[1] = []float64{9, 10}
	b[2] = []float64{11, 12}
	got := a.Mul(b)
	want := NewMat(2, 2)
	want[0] = []float64{58, 64}  // 1*7+2*9+3*11, 1*8+2*10+3*12
	want[1] = []float64{139, 154} // 4*7+5*9+6*11, 4*8+5*10+6*12
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if got[i][j] != want[i][j] {
				t.Fatalf("Mul[%d][%d] = %v, want %v", i, j, got[i][j], want[i][j])
			}
		}
	}
}

func TestMatMulVec(t *testing.T) {
	a := NewMat(2, 2)
	a[0] = []float64{2, 0}
	a[1] = []float64{0, 3}
	v := Vec{5, 7}
	got := a.MulVec(v)
	if got[0] != 10 || got[1] != 21 {
		t.Fatalf("MulVec = %v", got)
	}
}

func TestMatTranspose(t *testing.T) {
	a := NewMat(2, 3)
	a[0] = []float64{1, 2, 3}
	a[1] = []float64{4, 5, 6}
	tr := a.T()
	if tr.Rows() != 3 || tr.Cols() != 2 {
		t.Fatalf("T dims %dx%d", tr.Rows(), tr.Cols())
	}
	if tr[0][1] != 4 || tr[2][0] != 3 {
		t.Fatalf("T = %v", tr)
	}
}

func TestIdentity(t *testing.T) {
	id := Identity(3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			want := 0.0
			if i == j {
				want = 1
			}
			if id[i][j] != want {
				t.Fatalf("Identity[%d][%d] = %v", i, j, id[i][j])
			}
		}
	}
}
