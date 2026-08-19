package assemble

import (
	"math"
	"testing"

	"frame-static/internal/model"
)

func TestSolveCheckedBalanced(t *testing.T) {
	m := &model.Model{
		Nodes: []model.Node{
			{ID: "A", X: 0, Y: 0, Support: "fixed"},
			{ID: "B", X: 3, Y: 0},
		},
		Elements: []model.Element{{From: "A", To: "B", E: 2.1e11, A: 0.01, I: 4e-5}},
		Loads:    []model.NodeLoad{{Node: "B", FY: -600}},
	}
	res, bal, err := SolveChecked(m)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if res == nil || len(res.Nodes) != 2 {
		t.Fatalf("unexpected result")
	}
	if math.Abs(bal.ForceY) > 1e-6 || math.Abs(bal.Moment) > 1e-6 {
		t.Fatalf("balance off: %+v", bal)
	}
}

func TestSystemHelpers(t *testing.T) {
	m := &model.Model{
		Nodes: []model.Node{
			{ID: "A", X: 0, Y: 0, Support: "pin"},
			{ID: "B", X: 3, Y: 0},
		},
		Elements: []model.Element{{From: "A", To: "B", E: 1, A: 1, I: 1}},
	}
	sys := BuildSystem(m)
	if sys.FreeCount() != 4 || sys.FixedCount() != 2 {
		t.Fatalf("free=%d fixed=%d", sys.FreeCount(), sys.FixedCount())
	}
	if !sys.IsRestrained(0) || sys.IsRestrained(2) {
		t.Fatalf("restraint partition wrong")
	}
	dofs := sys.NodeDOFs(1)
	if dofs != [3]int{3, 4, 5} {
		t.Fatalf("node DOFs = %v", dofs)
	}
	if DOFName(4) != "uy" {
		t.Fatalf("DOFName(4) = %q", DOFName(4))
	}
}
