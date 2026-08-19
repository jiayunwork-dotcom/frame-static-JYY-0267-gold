package element

import "frame-static/internal/linalg"

// GlobalEndForces transforms a local 6-vector of member end forces into global
// coordinates: F_global = T^T * F_local (T maps global DOFs to local).
func GlobalEndForces(g *Geometry, local linalg.Vec) linalg.Vec {
	return Transform(g).T().MulVec(local)
}

// GlobalEndDisplacements transforms a local 6-vector of member end
// displacements into global coordinates: d_global = T^T * d_local.
func GlobalEndDisplacements(g *Geometry, local linalg.Vec) linalg.Vec {
	return Transform(g).T().MulVec(local)
}

// LocalEndDisplacements transforms a global 6-vector of member end
// displacements into local coordinates: d_local = T * d_global.
func LocalEndDisplacements(g *Geometry, global linalg.Vec) linalg.Vec {
	return Transform(g).MulVec(global)
}
