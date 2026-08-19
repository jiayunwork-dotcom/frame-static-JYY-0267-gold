package report

import (
	"strings"
	"testing"

	"frame-static/internal/assemble"
	"frame-static/internal/model"
)

func TestReportRenders(t *testing.T) {
	m := &model.Model{
		Nodes: []model.Node{
			{ID: "A", X: 0, Y: 0, Support: "pin"},
			{ID: "B", X: 4, Y: 0, Restraint: [3]bool{false, true, false}},
		},
		Elements: []model.Element{{From: "A", To: "B", E: 2.1e11, A: 0.01, I: 4e-5}},
		Loads:    []model.NodeLoad{{Node: "B", FY: -500}},
	}
	res, err := assemble.Solve(m)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	out := Build(m, res).String()
	if !strings.Contains(out, "节点位移") || !strings.Contains(out, "支座反力") {
		t.Fatalf("report missing sections:\n%s", out)
	}
	if !strings.Contains(Build(m, res).Summary(), "最大弯矩") {
		t.Fatalf("summary missing")
	}
}
