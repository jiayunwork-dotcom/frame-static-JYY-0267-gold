package assemble

import (
	"frame-static/internal/element"
	"frame-static/internal/linalg"
	"frame-static/internal/model"
)

// GlobalStiffness assembles the full (unreduced) global stiffness matrix by
// scattering each member's global stiffness into its six DOFs.
func GlobalStiffness(m *model.Model, sys *System) linalg.Mat {
	K := linalg.NewMat(sys.N, sys.N)
	idx := model.NodeIndex(m)
	for _, e := range m.Elements {
		fi, ti := idx[e.From], idx[e.To]
		g, err := element.GeometryOf(m.Nodes[fi], m.Nodes[ti])
		if err != nil {
			continue
		}
		kg := element.GlobalStiffness(element.LocalStiffness(g, e.E, e.A, e.I), element.Transform(g))
		dofs := []int{3*fi, 3*fi + 1, 3*fi + 2, 3*ti, 3*ti + 1, 3*ti + 2}
		for a := 0; a < 6; a++ {
			for b := 0; b < 6; b++ {
				K[dofs[a]][dofs[b]] += kg[a][b]
			}
		}
	}
	return K
}

// LoadVector assembles the global load vector from nodal loads and the
// consistent equivalent nodal loads of any distributed member loads.
func LoadVector(m *model.Model, sys *System) linalg.Vec {
	F := linalg.NewVec(sys.N)
	idx := model.NodeIndex(m)
	for _, ld := range m.Loads {
		ni, ok := idx[ld.Node]
		if !ok {
			continue
		}
		F[3*ni] += ld.FX
		F[3*ni+1] += ld.FY
		F[3*ni+2] += ld.MZ
	}
	for _, e := range m.Elements {
		if e.Dist == nil {
			continue
		}
		fi, ti := idx[e.From], idx[e.To]
		g, err := element.GeometryOf(m.Nodes[fi], m.Nodes[ti])
		if err != nil {
			continue
		}
		eqLocal := element.EquivalentNodalLoad(g, e.Dist.Q)
		// Rotate the 6-vector into global coordinates: F_global = T^T * F_local.
		eqGlobal := element.Transform(g).T().MulVec(linalg.Vec(eqLocal[:]))
		for k := 0; k < 3; k++ {
			F[3*fi+k] += eqGlobal[k]
			F[3*ti+k] += eqGlobal[3+k]
		}
	}
	return F
}
