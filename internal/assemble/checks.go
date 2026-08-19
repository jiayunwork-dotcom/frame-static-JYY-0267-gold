package assemble

import "frame-static/internal/model"

// SolveChecked runs the full solve and also returns the global equilibrium
// residual, so callers can verify balance in a single call.
func SolveChecked(m *model.Model) (*Result, *Balance, error) {
	sys, K, F, u, err := prepare(m)
	if err != nil {
		return nil, nil, err
	}
	res := buildResult(m, sys, u, K, F)
	bal := CheckBalance(K, u, F, sys, m)
	return res, &bal, nil
}

// NodeDOFs returns the three global DOF indices of a joint.
func (sys *System) NodeDOFs(node int) [3]int {
	return [3]int{3 * node, 3*node + 1, 3*node + 2}
}

// DOFName returns "ux", "uy", or "theta" for a DOF within a joint.
func DOFName(d int) string {
	switch d % 3 {
	case 0:
		return "ux"
	case 1:
		return "uy"
	default:
		return "theta"
	}
}

// FreeCount returns how many DOFs participate in the reduced system.
func (sys *System) FreeCount() int { return len(sys.Free) }

// FixedCount returns how many DOFs are constrained.
func (sys *System) FixedCount() int { return len(sys.Fixed) }

// IsRestrained reports whether the given global DOF index is constrained.
func (sys *System) IsRestrained(d int) bool {
	for _, f := range sys.Fixed {
		if f == d {
			return true
		}
	}
	return false
}
