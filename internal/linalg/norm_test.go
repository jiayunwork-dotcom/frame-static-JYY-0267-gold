package linalg

import (
	"math"
	"testing"
)

func TestNormInf(t *testing.T) {
	v := Vec{1, -4, 3}
	if v.NormInf() != 4 {
		t.Fatalf("NormInf = %v", v.NormInf())
	}
}

func TestNorm2(t *testing.T) {
	v := Vec{3, 4}
	if math.Abs(v.Norm2()-5) > 1e-12 {
		t.Fatalf("Norm2 = %v", v.Norm2())
	}
}

func TestAddScaled(t *testing.T) {
	a := Vec{1, 2}
	b := Vec{3, 4}
	out := a.AddScaled(b, 2)
	if out[0] != 7 || out[1] != 10 {
		t.Fatalf("AddScaled = %v", out)
	}
}

func TestFrobenius(t *testing.T) {
	m := NewMat(2, 2)
	m[0] = []float64{3, 4}
	m[1] = []float64{0, 0}
	if math.Abs(m.Frobenius()-5) > 1e-12 {
		t.Fatalf("Frobenius = %v", m.Frobenius())
	}
}

func TestMatCopyIndependent(t *testing.T) {
	m := NewMat(2, 2)
	m[0][0] = 7
	c := m.Copy()
	c[0][0] = 99
	if m[0][0] != 7 {
		t.Fatalf("copy shared backing array")
	}
}
