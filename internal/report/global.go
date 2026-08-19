package report

import (
	"fmt"
	"strings"

	"frame-static/internal/assemble"
	"frame-static/internal/element"
	"frame-static/internal/linalg"
	"frame-static/internal/model"
)

// globalSection renders member end forces transformed into global
// coordinates. It exercises the same transform used for assembly, so a broken
// rotation would show up here as well as in the reactions.
func (rep *Report) globalSection() string {
	var b strings.Builder
	b.WriteString("\n杆端力（整体坐标）\n")
	b.WriteString("  杆             Fx_i      Fy_i      Mz_i      Fx_j      Fy_j      Mz_j\n")
	idx := model.NodeIndex(rep.m)
	for _, mf := range rep.r.Members {
		fi, ok1 := idx[mf.From]
		ti, ok2 := idx[mf.To]
		if !ok1 || !ok2 {
			continue
		}
		g, err := element.GeometryOf(rep.m.Nodes[fi], rep.m.Nodes[ti])
		if err != nil {
			continue
		}
		local := linalg.Vec{mf.Ni, mf.Vi, mf.Mi, mf.Nj, mf.Vj, mf.Mj}
		glob := element.GlobalEndForces(g, local)
		fmt.Fprintf(&b, "  %-6s-%s %9.3f %9.3f %9.3f %9.3f %9.3f %9.3f\n",
			mf.From, mf.To, glob[0], glob[1], glob[2], glob[3], glob[4], glob[5])
	}
	return b.String()
}

// BalanceSection renders the global equilibrium residuals.
func (rep *Report) BalanceSection(bal assemble.Balance) string {
	return fmt.Sprintf("\n整体平衡残差（应≈0）: Fx %.3e  Fy %.3e  M %.3e\n",
		bal.ForceX, bal.ForceY, bal.Moment)
}
