package assemble

import (
	"frame-static/internal/linalg"
	"frame-static/internal/model"
)

// AppliedTotal sums the applied nodal loads in the model into a resultant
// force and moment about the origin, including force lever arms (x·fy − y·fx).
func AppliedTotal(m *model.Model) (fx, fy, mz float64) {
	idx := model.NodeIndex(m)
	for _, ld := range m.Loads {
		fx += ld.FX
		fy += ld.FY
		mz += ld.MZ
		if n, ok := idx[ld.Node]; ok {
			mz += m.Nodes[n].X*ld.FY - m.Nodes[n].Y*ld.FX
		}
	}
	return
}

// AppliedLoadTotal sums the global load vector into a resultant force and
// moment about the origin (the external demand the supports must resist).
func AppliedLoadTotal(F linalg.Vec) (fx, fy, mz float64) {
	for i := 0; i+2 < len(F); i += 3 {
		fx += F[i]
		fy += F[i+1]
		mz += F[i+2]
	}
	return
}

// ReactionTotal sums the reactions over the constrained DOFs.
func ReactionTotal(R linalg.Vec, sys *System) (fx, fy, mz float64) {
	for _, d := range sys.Fixed {
		switch d % 3 {
		case 0:
			fx += R[d]
		case 1:
			fy += R[d]
		case 2:
			mz += R[d]
		}
	}
	return
}

// Balance holds the global equilibrium residual (support reactions plus
// applied loads). Each component should be near zero.
type Balance struct {
	ForceX float64
	ForceY float64
	Moment float64
}

// CheckBalance recomputes the residual of all support reactions plus the
// applied nodal loads about the origin, confirming global equilibrium.
func CheckBalance(K linalg.Mat, u, F linalg.Vec, sys *System, m *model.Model) Balance {
	R := Reactions(K, u, F)
	var rfx, rfy, rmz float64
	for _, d := range sys.Fixed {
		val := R[d]
		n := m.Nodes[d/3]
		switch d % 3 {
		case 0:
			rfx += val
			rmz -= n.Y * val
		case 1:
			rfy += val
			rmz += n.X * val
		case 2:
			rmz += val
		}
	}
	afx, afy, amz := AppliedTotal(m)
	return Balance{ForceX: rfx + afx, ForceY: rfy + afy, Moment: rmz + amz}
}
