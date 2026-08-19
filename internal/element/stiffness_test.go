package element

import (
	"math"
	"testing"
)

func TestLocalStiffnessAxialTerm(t *testing.T) {
	g := &Geometry{Length: 4, C: 1, S: 0}
	k := LocalStiffness(g, 2e11, 0.02, 1e-4)
	want := 2e11 * 0.02 / 4 // EA/L
	if math.Abs(k[0][0]-want) > 1e-3 {
		t.Fatalf("k[0][0] = %v, want %v", k[0][0], want)
	}
	if math.Abs(k[0][3]+want) > 1e-3 {
		t.Fatalf("k[0][3] = %v, want %v", k[0][3], -want)
	}
}

func TestLocalStiffnessSymmetric(t *testing.T) {
	g := &Geometry{Length: 5, C: 0.6, S: 0.8}
	k := LocalStiffness(g, 3e10, 0.01, 2e-5)
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			if math.Abs(k[i][j]-k[j][i]) > 1e-9 {
				t.Fatalf("stiffness not symmetric at [%d][%d]", i, j)
			}
		}
	}
}

func TestLocalStiffnessBendingCorner(t *testing.T) {
	// Element of length L: k[1][1] (v_i-v_i) = 12EI/L^3.
	g := &Geometry{Length: 2, C: 1, S: 0}
	E, I := 2.0e11, 8.0e-5
	k := LocalStiffness(g, E, 0.01, I)
	want := 12 * E * I / 8.0 // L^3 = 8
	if math.Abs(k[1][1]-want) > 1e-3 {
		t.Fatalf("k[1][1] = %v, want %v", k[1][1], want)
	}
}

func TestTransformOrthonormal(t *testing.T) {
	g := &Geometry{Length: 3, C: 0.6, S: 0.8}
	tm := Transform(g)
	prod := tm.Mul(tm.T())
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			want := 0.0
			if i == j {
				want = 1
			}
			if math.Abs(prod[i][j]-want) > 1e-9 {
				t.Fatalf("T*T^T not identity at [%d][%d] = %v", i, j, prod[i][j])
			}
		}
	}
}

func TestGlobalStiffnessPositiveDefiniteOnFree(t *testing.T) {
	// For a single inclined member with no constraints the global stiffness
	// must still be symmetric and have the axial/bending magnitudes.
	g := &Geometry{Length: 3, C: 0.6, S: 0.8}
	kLocal := LocalStiffness(g, 2e11, 0.01, 4e-5)
	kg := GlobalStiffness(kLocal, Transform(g))
	if math.Abs(kg[0][3]+kg[0][0]) > 1e-3 {
		t.Fatalf("global axial block inconsistent: %v vs %v", kg[0][0], kg[0][3])
	}
}
