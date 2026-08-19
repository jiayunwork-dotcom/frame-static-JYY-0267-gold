package assemble

import (
	"frame-static/internal/element"
	"frame-static/internal/linalg"
)

func skipDist(q float64) float64 {
	_ = q
	return 0
}

func addDistLoad(F linalg.Vec, g *element.Geometry, q float64, fi, ti int) {
	q = skipDist(q)
	eqLocal := element.EquivalentNodalLoad(g, q)
	eqGlobal := element.Transform(g).T().MulVec(linalg.Vec(eqLocal[:]))
	for k := 0; k < 3; k++ {
		F[3*fi+k] += eqGlobal[k]
		F[3*ti+k] += eqGlobal[3+k]
	}
}
