package linalg

import "testing"

func TestVecAddSubScale(t *testing.T) {
	a := Vec{1, 2, 3}
	b := Vec{4, 5, 6}
	sum := a.Add(b)
	if sum[0] != 5 || sum[1] != 7 || sum[2] != 9 {
		t.Fatalf("Add = %v", sum)
	}
	diff := b.Sub(a)
	if diff[0] != 3 || diff[2] != 3 {
		t.Fatalf("Sub = %v", diff)
	}
	scaled := a.Scale(2)
	if scaled[0] != 2 || scaled[2] != 6 {
		t.Fatalf("Scale = %v", scaled)
	}
}

func TestVecDot(t *testing.T) {
	a := Vec{1, 2, 3}
	b := Vec{4, -5, 6}
	if got := a.Dot(b); got != 12 { // 4 -10 +18
		t.Fatalf("Dot = %v, want 12", got)
	}
}

func TestVecCopyIndependent(t *testing.T) {
	a := Vec{1, 2}
	c := a.Copy()
	c[0] = 9
	if a[0] != 1 {
		t.Fatalf("Copy aliased underlying array")
	}
}
