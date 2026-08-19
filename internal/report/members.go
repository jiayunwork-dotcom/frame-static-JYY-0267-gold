package report

import (
	"fmt"
	"strings"

	"frame-static/internal/format"
)

func (rep *Report) membersSection() string {
	var b strings.Builder
	b.WriteString("\n杆端力（局部坐标：轴力 N / 剪力 V / 弯矩 M）\n")
	b.WriteString("  杆            Ni       Vi       Mi       Nj       Vj       Mj   长度\n")
	for _, m := range rep.r.Members {
		fmt.Fprintf(&b, "  %-6s-%s %8.3f %8.3f %8.3f %8.3f %8.3f %8.3f %7.3f\n",
			m.From, m.To, m.Ni, m.Vi, m.Mi, m.Nj, m.Vj, m.Mj, m.Length)
	}
	return b.String()
}

func (rep *Report) reactionsSection() string {
	var b strings.Builder
	b.WriteString("\n支座反力\n")
	b.WriteString("  节点   自由度    反力\n")
	for _, rct := range rep.r.Reactions {
		fmt.Fprintf(&b, "  %-6s %-6s %s\n", rct.Node, rct.DOF, format.Signed(rct.Force, 4))
	}
	return b.String()
}

// Summary returns a one-line digest of the extreme result magnitudes.
func (rep *Report) Summary() string {
	maxM := rep.r.MaxAbsMoment()
	maxD := rep.r.MaxAbsDeflection()
	maxR := rep.r.MaxAbsReaction()
	return fmt.Sprintf("最大弯矩 %.4f / 最大位移 %.4f / 最大反力 %.4f", maxM, maxD, maxR)
}
