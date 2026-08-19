package element

import "frame-static/internal/linalg"

// Transform returns the 6x6 transformation matrix T that maps global DOFs to
// local DOFs: d_local = T * d_global. It is block-diagonal with three identical
// 2D rotation blocks plus the identity on the rotation DOF.
func Transform(g *Geometry) linalg.Mat {
	c, s := g.C, g.S
	t := linalg.NewMat(6, 6)
	blocks := [][]float64{{c, s, 0}, {-s, c, 0}, {0, 0, 1}}
	for b := 0; b < 2; b++ {
		base := b * 3
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				t[base+i][base+j] = blocks[i][j]
			}
		}
	}
	return t
}

// GlobalStiffness returns the member stiffness in global coordinates:
// K_global = T^T * k_local * T.
func GlobalStiffness(kLocal, t linalg.Mat) linalg.Mat {
	return t.T().Mul(kLocal).Mul(t)
}
