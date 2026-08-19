package report

import (
	"fmt"
	"strings"

	"frame-static/internal/assemble"
	"frame-static/internal/model"
)

// Report is a human-readable rendering of a solved frame.
type Report struct {
	m *model.Model
	r *assemble.Result
}

// Build assembles a report from the model and its solution.
func Build(m *model.Model, res *assemble.Result) *Report {
	return &Report{m: m, r: res}
}

// String renders the full report.
func (rep *Report) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "frame-static 求解报告（%d 节点 / %d 杆件）\n", len(rep.m.Nodes), len(rep.m.Elements))
	b.WriteString(rep.nodesSection())
	b.WriteString(rep.membersSection())
	b.WriteString(rep.globalSection())
	b.WriteString(rep.reactionsSection())
	return b.String()
}

func (rep *Report) nodesSection() string {
	var b strings.Builder
	b.WriteString("\n节点位移\n")
	b.WriteString("  id             ux          uy        theta\n")
	for _, n := range rep.r.Nodes {
		fmt.Fprintf(&b, "  %-12s %11.6f %11.6f %9.7f\n", n.ID, n.UX, n.UY, n.Theta)
	}
	return b.String()
}
