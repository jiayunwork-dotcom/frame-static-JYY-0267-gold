package assemble

import "frame-static/internal/linalg"

// Reactions returns R = K*u - F, the out-of-balance forces at every DOF. At
// free DOFs this is ~0 by construction; at fixed DOFs it is the support reaction.
func Reactions(K linalg.Mat, u, F linalg.Vec) linalg.Vec {
	return K.MulVec(u).Sub(F)
}
