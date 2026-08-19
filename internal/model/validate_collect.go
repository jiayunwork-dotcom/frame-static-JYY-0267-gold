package model

import (
	"fmt"
	"math"
)

// ValidateAll returns every structural defect found, rather than stopping at
// the first. Duplicate ids, unknown endpoints, zero-length members, non-positive
// sections, disconnected components, and missing constraints are all reported.
func ValidateAll(m *Model) []error {
	var errs []error
	if len(m.Nodes) == 0 {
		return append(errs, &InvalidModelError{Field: "nodes", Msg: "no nodes"})
	}
	if len(m.Elements) == 0 {
		return append(errs, &InvalidModelError{Field: "elements", Msg: "no elements"})
	}
	idx := make(map[string]int, len(m.Nodes))
	for i, n := range m.Nodes {
		if n.ID == "" {
			errs = append(errs, &InvalidModelError{Field: "nodes", Msg: fmt.Sprintf("node #%d has empty id", i)})
			continue
		}
		if _, dup := idx[n.ID]; dup {
			errs = append(errs, &InvalidModelError{Field: "nodes", Msg: fmt.Sprintf("duplicate node id %q", n.ID)})
		}
		idx[n.ID] = i
	}
	adj := make([][]int, len(m.Nodes))
	for ei, e := range m.Elements {
		fi, ok1 := idx[e.From]
		ti, ok2 := idx[e.To]
		if !ok1 {
			errs = append(errs, &InvalidModelError{Field: "elements", Msg: fmt.Sprintf("element #%d references unknown node %q", ei, e.From)})
			continue
		}
		if !ok2 {
			errs = append(errs, &InvalidModelError{Field: "elements", Msg: fmt.Sprintf("element #%d references unknown node %q", ei, e.To)})
			continue
		}
		if fi == ti {
			errs = append(errs, &InvalidModelError{Field: "elements", Msg: fmt.Sprintf("element #%d is a loop on node %q", ei, e.From)})
			continue
		}
		nf, nt := m.Nodes[fi], m.Nodes[ti]
		if math.Hypot(nt.X-nf.X, nt.Y-nf.Y) == 0 {
			errs = append(errs, &InvalidModelError{Field: "elements", Msg: fmt.Sprintf("element #%d (%s-%s) has zero length", ei, e.From, e.To)})
		}
		if e.E <= 0 {
			errs = append(errs, &InvalidModelError{Field: "elements", Msg: fmt.Sprintf("element #%d E=%v must be > 0", ei, e.E)})
		}
		if e.A <= 0 {
			errs = append(errs, &InvalidModelError{Field: "elements", Msg: fmt.Sprintf("element #%d A=%v must be > 0", ei, e.A)})
		}
		if e.I <= 0 {
			errs = append(errs, &InvalidModelError{Field: "elements", Msg: fmt.Sprintf("element #%d I=%v must be > 0", ei, e.I)})
		}
		adj[fi] = append(adj[fi], ti)
		adj[ti] = append(adj[ti], fi)
	}
	// Connectivity: every node reachable from node 0.
	if len(errs) == 0 {
		seen := make([]bool, len(m.Nodes))
		stack := []int{0}
		seen[0] = true
		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, nx := range adj[cur] {
				if !seen[nx] {
					seen[nx] = true
					stack = append(stack, nx)
				}
			}
		}
		for i, s := range seen {
			if !s {
				errs = append(errs, &InvalidModelError{Field: "nodes", Msg: fmt.Sprintf("node %q is not in the same connected component as the rest", m.Nodes[i].ID)})
			}
		}
	}
	any := false
	for _, n := range m.Nodes {
		r := n.EffectiveRestraint()
		if r[0] || r[1] || r[2] {
			any = true
			break
		}
	}
	if !any {
		errs = append(errs, &InvalidModelError{Field: "nodes", Msg: "no constrained DOF: structure is a rigid body"})
	}
	return errs
}
