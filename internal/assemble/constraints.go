package assemble

import "frame-static/internal/model"

// System captures the free/fixed DOF partition derived from node restraints.
// Constrained DOFs are removed ("划行划列") before solving.
type System struct {
	N       int
	Free    []int
	Fixed   []int
	freePos map[int]int
}

// BuildSystem partitions the 3*nodeCount global DOFs into free and fixed sets.
func BuildSystem(m *model.Model) *System {
	n := model.TotalDOF(m)
	free := make([]int, 0, n)
	fixed := make([]int, 0, n)
	pos := make(map[int]int, n)
	for i, node := range m.Nodes {
		r := node.EffectiveRestraint()
		for d := 0; d < 3; d++ {
			dof := 3*i + d
			if r[d] {
				fixed = append(fixed, dof)
			} else {
				pos[dof] = len(free)
				free = append(free, dof)
			}
		}
	}
	return &System{N: n, Free: free, Fixed: fixed, freePos: pos}
}
