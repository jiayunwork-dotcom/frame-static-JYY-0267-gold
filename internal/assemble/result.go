package assemble

import "sort"

// NodeResult holds the solved displacements of one joint.
type NodeResult struct {
	ID    string
	UX    float64
	UY    float64
	Theta float64
}

// MemberEndForce holds the local end forces of one member, ordered by end.
// Axial N is positive in tension; V and M follow the local y and theta axes.
type MemberEndForce struct {
	From   string
	To     string
	Ni, Vi, Mi float64
	Nj, Vj, Mj float64
	Length    float64
}

// Reaction holds one support reaction component at a constrained DOF.
type Reaction struct {
	Node  string
	DOF   string // "ux" | "uy" | "theta"
	Force float64
}

// Result is the full solution returned to the caller.
type Result struct {
	Nodes     []NodeResult
	Members   []MemberEndForce
	Reactions []Reaction
}

// MaxAbsMoment returns the largest absolute member end moment.
func (r *Result) MaxAbsMoment() float64 {
	m := 0.0
	for _, e := range r.Members {
		m = maxAbs(m, e.Mi)
		m = maxAbs(m, e.Mj)
	}
	return m
}

// MaxAbsDeflection returns the largest absolute nodal displacement.
func (r *Result) MaxAbsDeflection() float64 {
	m := 0.0
	for _, n := range r.Nodes {
		m = maxAbs(m, n.UX)
		m = maxAbs(m, n.UY)
		m = maxAbs(m, n.Theta)
	}
	return m
}

// MaxAbsReaction returns the largest absolute support reaction component.
func (r *Result) MaxAbsReaction() float64 {
	m := 0.0
	for _, x := range r.Reactions {
		m = maxAbs(m, x.Force)
	}
	return m
}

func maxAbs(a, b float64) float64 {
	if b < 0 {
		b = -b
	}
	if b > a {
		return b
	}
	return a
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// LargestReaction returns the reaction with the largest absolute force.
func (r *Result) LargestReaction() (Reaction, bool) {
	var best Reaction
	found := false
	for _, x := range r.Reactions {
		if !found || abs(x.Force) > abs(best.Force) {
			best = x
			found = true
		}
	}
	return best, found
}

// SortedMembers returns member end forces ordered by (from, to).
func (r *Result) SortedMembers() []MemberEndForce {
	out := append([]MemberEndForce(nil), r.Members...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].From != out[j].From {
			return out[i].From < out[j].From
		}
		return out[i].To < out[j].To
	})
	return out
}
