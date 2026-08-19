package element

import "frame-static/internal/linalg"

// LocalStiffness returns the 6x6 member stiffness matrix in local coordinates.
// DOF order is [u_i, v_i, theta_i, u_j, v_j, theta_j], with u along the member
// axis, v perpendicular, theta counterclockwise.
func LocalStiffness(g *Geometry, E, A, I float64) linalg.Mat {
	L := g.Length
	k := linalg.NewMat(6, 6)
	ea := E * A / L
	ei := E * I
	l3 := L * L * L
	l2 := L * L
	// Axial terms.
	k[0][0] = ea
	k[0][3] = -ea
	k[3][0] = -ea
	k[3][3] = ea
	// Bending terms.
	k[1][1] = 12 * ei / l3
	k[1][2] = 6 * ei / l2
	k[1][4] = -12 * ei / l3
	k[1][5] = 6 * ei / l2
	k[2][1] = 6 * ei / l2
	k[2][2] = 4 * ei / L
	k[2][4] = -6 * ei / l2
	k[2][5] = 2 * ei / L
	k[4][1] = -12 * ei / l3
	k[4][2] = -6 * ei / l2
	k[4][4] = 12 * ei / l3
	k[4][5] = -6 * ei / l2
	k[5][1] = 6 * ei / l2
	k[5][2] = 2 * ei / L
	k[5][4] = -6 * ei / l2
	k[5][5] = 4 * ei / L
	return k
}

// AxialStiffnessPart returns the 6x6 matrix holding only the axial (EA/L)
// terms; all bending entries are zero.
func AxialStiffnessPart(g *Geometry, E, A float64) linalg.Mat {
	k := linalg.NewMat(6, 6)
	ea := E * A / g.Length
	k[0][0], k[0][3] = ea, -ea
	k[3][0], k[3][3] = -ea, ea
	return k
}

// BendingStiffnessPart returns the 6x6 matrix holding only the bending (EI)
// terms; all axial entries are zero.
func BendingStiffnessPart(g *Geometry, E, I float64) linalg.Mat {
	L := g.Length
	ei := E * I
	l3 := L * L * L
	l2 := L * L
	k := linalg.NewMat(6, 6)
	k[1][1] = 12 * ei / l3
	k[1][2] = 6 * ei / l2
	k[1][4] = -12 * ei / l3
	k[1][5] = 6 * ei / l2
	k[2][1] = 6 * ei / l2
	k[2][2] = 4 * ei / L
	k[2][4] = -6 * ei / l2
	k[2][5] = 2 * ei / L
	k[4][1] = -12 * ei / l3
	k[4][2] = -6 * ei / l2
	k[4][4] = 12 * ei / l3
	k[4][5] = -6 * ei / l2
	k[5][1] = 6 * ei / l2
	k[5][2] = 2 * ei / L
	k[5][4] = -6 * ei / l2
	k[5][5] = 4 * ei / L
	return k
}

// LocalStiffnessFromParts composes the local stiffness from its axial and
// bending parts; the sum equals the full matrix from LocalStiffness.
func LocalStiffnessFromParts(g *Geometry, E, A, I float64) linalg.Mat {
	return AxialStiffnessPart(g, E, A).Add(BendingStiffnessPart(g, E, I))
}
