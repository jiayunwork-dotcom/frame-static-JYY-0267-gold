package assemble

import (
	"math"
	"testing"

	"frame-static/internal/model"
)

func approx(a, b, tol float64) bool { return math.Abs(a-b) <= tol }

func reactionAt(res *Result, node, dof string) (float64, bool) {
	for _, r := range res.Reactions {
		if r.Node == node && r.DOF == dof {
			return r.Force, true
		}
	}
	return 0, false
}

func TestBuildSystemPartition(t *testing.T) {
	m := &model.Model{
		Nodes: []model.Node{
			{ID: "A", X: 0, Y: 0, Support: "fixed"},
			{ID: "B", X: 3, Y: 0},
		},
		Elements: []model.Element{{From: "A", To: "B", E: 1, A: 1, I: 1}},
	}
	sys := BuildSystem(m)
	if sys.N != 6 {
		t.Fatalf("N = %d, want 6", sys.N)
	}
	if len(sys.Free) != 3 || len(sys.Fixed) != 3 {
		t.Fatalf("free=%d fixed=%d, want 3/3", len(sys.Free), len(sys.Fixed))
	}
}

func TestGlobalStiffnessSymmetric(t *testing.T) {
	m := &model.Model{
		Nodes: []model.Node{
			{ID: "A", X: 0, Y: 0, Support: "pin"},
			{ID: "B", X: 4, Y: 3},
		},
		Elements: []model.Element{{From: "A", To: "B", E: 2e11, A: 0.01, I: 4e-5}},
	}
	sys := BuildSystem(m)
	K := GlobalStiffness(m, sys)
	for i := 0; i < sys.N; i++ {
		for j := 0; j < sys.N; j++ {
			if math.Abs(K[i][j]-K[j][i]) > 1e-9 {
				t.Fatalf("global K not symmetric at [%d][%d]", i, j)
			}
		}
	}
}

func TestSolveSimplySupportedCentralLoad(t *testing.T) {
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
	ra, _ := reactionAt(res, "A", "uy")
	rb, _ := reactionAt(res, "B", "uy")
	if !approx(ra, 500, 1e-3) {
		t.Fatalf("R_Ay = %v, want 500", ra)
	}
	if !approx(rb, 500, 1e-3) {
		t.Fatalf("R_By = %v, want 500", rb)
	}
}

func TestSolveCantileverTipLoad(t *testing.T) {
	m := &model.Model{
		Nodes: []model.Node{
			{ID: "A", X: 0, Y: 0, Support: "fixed"},
			{ID: "B", X: 3, Y: 0},
		},
		Elements: []model.Element{{From: "A", To: "B", E: 2.1e11, A: 0.01, I: 4e-5}},
		Loads:    []model.NodeLoad{{Node: "B", FY: -500}},
	}
	res, err := Solve(m)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	ray, _ := reactionAt(res, "A", "uy")
	ram, _ := reactionAt(res, "A", "theta")
	if !approx(ray, 500, 1e-3) {
		t.Fatalf("R_Ay = %v, want 500", ray)
	}
	if !approx(ram, 500*3, 1e-2) {
		t.Fatalf("R_A_theta = %v, want 1500", ram)
	}
}

func TestSolveUDLSimplySupported(t *testing.T) {
	m := &model.Model{
		Nodes: []model.Node{
			{ID: "A", X: 0, Y: 0, Support: "pin"},
			{ID: "B", X: 4, Y: 0, Restraint: [3]bool{false, true, false}},
		},
		Elements: []model.Element{
			{From: "A", To: "B", E: 2.1e11, A: 0.01, I: 4e-5, Dist: &model.DistLoad{Q: 200}},
		},
	}
	res, err := Solve(m)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	ra, _ := reactionAt(res, "A", "uy")
	rb, _ := reactionAt(res, "B", "uy")
	// qL/2 = 200*4/2 = 400
	if !approx(ra, 400, 1e-3) {
		t.Fatalf("R_Ay = %v, want 400", ra)
	}
	if !approx(rb, 400, 1e-3) {
		t.Fatalf("R_By = %v, want 400", rb)
	}
}

func TestSolveMechanismSingular(t *testing.T) {
	m := &model.Model{
		Nodes: []model.Node{
			{ID: "A", X: 0, Y: 0},
			{ID: "B", X: 4, Y: 0},
		},
		Elements: []model.Element{{From: "A", To: "B", E: 2.1e11, A: 0.01, I: 4e-5}},
	}
	if _, err := Solve(m); err == nil {
		t.Fatalf("expected singular-structure error")
	}
}
