package assemble

import (
	"math"
	"testing"

	"frame-static/internal/model"
)

func memberEnd(res *Result, from, to string) (MemberEndForce, bool) {
	for _, mf := range res.Members {
		if mf.From == from && mf.To == to {
			return mf, true
		}
	}
	return MemberEndForce{}, false
}

func appliedTotal(m *model.Model) (fx, fy, mz float64) {
	for _, ld := range m.Loads {
		fx += ld.FX
		fy += ld.FY
		mz += ld.MZ
	}
	return
}

func TestEquilibriumJointForces(t *testing.T) {
	m := &model.Model{
		Nodes: []model.Node{
			{ID: "A", X: 0, Y: 0, Support: "pin"},
			{ID: "C", X: 2, Y: 0},
			{ID: "B", X: 4, Y: 0, Restraint: [3]bool{false, true, false}},
		},
		Elements: []model.Element{
			{From: "A", To: "C", E: 2.1e11, A: 0.01, I: 4e-5},
			{From: "C", To: "B", E: 2.1e11, A: 0.01, I: 4e-5},
		},
		Loads: []model.NodeLoad{{Node: "C", FY: -1000}},
	}
	res, err := Solve(m)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	ac, _ := memberEnd(res, "A", "C")
	cb, _ := memberEnd(res, "C", "B")
	// Member end forces are the force the node applies to the element, so at a
	// joint their sum equals the applied nodal load (joint equilibrium).
	if !approx(ac.Vj+cb.Vi, -1000, 1e-3) {
		t.Fatalf("joint C vertical not balanced: Vj(AC)=%v Vi(CB)=%v", ac.Vj, cb.Vi)
	}
	if !approx(ac.Mj+cb.Mi, 0, 1e-3) {
		t.Fatalf("joint C moment not balanced: Mj(AC)=%v Mi(CB)=%v", ac.Mj, cb.Mi)
	}
	// Each member: end axial forces are equal and opposite.
	if !approx(ac.Ni+ac.Nj, 0, 1e-3) {
		t.Fatalf("member A-C axial not balanced: Ni=%v Nj=%v", ac.Ni, ac.Nj)
	}
	if !approx(cb.Ni+cb.Nj, 0, 1e-3) {
		t.Fatalf("member C-B axial not balanced: Ni=%v Nj=%v", cb.Ni, cb.Nj)
	}
}

func TestEquilibriumGlobalBalance(t *testing.T) {
	m := &model.Model{
		Nodes: []model.Node{
			{ID: "A", X: 0, Y: 0, Support: "pin"},
			{ID: "C", X: 2, Y: 0},
			{ID: "B", X: 4, Y: 0, Restraint: [3]bool{false, true, false}},
		},
		Elements: []model.Element{
			{From: "A", To: "C", E: 2.1e11, A: 0.01, I: 4e-5},
			{From: "C", To: "B", E: 2.1e11, A: 0.01, I: 4e-5},
		},
		Loads: []model.NodeLoad{{Node: "C", FY: -1000}},
	}
	sys := BuildSystem(m)
	K := GlobalStiffness(m, sys)
	F := LoadVector(m, sys)
	u, err := solveReduced(K, F, sys)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	bal := CheckBalance(K, u, F, sys, m)
	if math.Abs(bal.ForceX) > 1e-6 || math.Abs(bal.ForceY) > 1e-6 || math.Abs(bal.Moment) > 1e-6 {
		t.Fatalf("global balance off: %+v", bal)
	}
}

func TestConsistencyInclinedAxial(t *testing.T) {
	// Inclined cantilever under a purely axial tip load: reactions must balance
	// the load and the local member axial force must equal the applied axial
	// magnitude. If the member transform were wrong, BOTH checks break.
	m := &model.Model{
		Nodes: []model.Node{
			{ID: "A", X: 0, Y: 0, Support: "fixed"},
			{ID: "B", X: 3, Y: 4},
		},
		Elements: []model.Element{{From: "A", To: "B", E: 2.1e11, A: 0.01, I: 4e-5}},
		Loads:    []model.NodeLoad{{Node: "B", FX: -300, FY: -400}}, // 500 along -AB
	}
	res, err := Solve(m)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	ab, _ := memberEnd(res, "A", "B")
	if !approx(ab.Ni, 500, 1e-2) {
		t.Fatalf("local axial Ni = %v, want 500 (transform-dependent)", ab.Ni)
	}
	rax, _ := reactionAt(res, "A", "ux")
	ray, _ := reactionAt(res, "A", "uy")
	if !approx(rax-300, 0, 1e-2) || !approx(ray-400, 0, 1e-2) {
		t.Fatalf("reaction balance off: R_A=(%v,%v)", rax, ray)
	}
}
