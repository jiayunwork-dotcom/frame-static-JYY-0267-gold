package api

import "frame-static/internal/assemble"

type nodeOut struct {
	ID    string  `json:"id"`
	UX    float64 `json:"ux"`
	UY    float64 `json:"uy"`
	Theta float64 `json:"theta"`
}

type memberOut struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Ni     float64 `json:"Ni"`
	Vi     float64 `json:"Vi"`
	Mi     float64 `json:"Mi"`
	Nj     float64 `json:"Nj"`
	Vj     float64 `json:"Vj"`
	Mj     float64 `json:"Mj"`
	Length float64 `json:"length"`
}

type reactionOut struct {
	Node  string  `json:"node"`
	DOF   string  `json:"dof"`
	Force float64 `json:"force"`
}

type solveResponse struct {
	OK        bool          `json:"ok"`
	Nodes     []nodeOut     `json:"nodes"`
	Members   []memberOut   `json:"members"`
	Reactions []reactionOut `json:"reactions"`
}

func buildResponse(res *assemble.Result) solveResponse {
	out := solveResponse{OK: true}
	for _, n := range res.Nodes {
		out.Nodes = append(out.Nodes, nodeOut{ID: n.ID, UX: n.UX, UY: n.UY, Theta: n.Theta})
	}
	for _, m := range res.Members {
		out.Members = append(out.Members, memberOut{
			From: m.From, To: m.To,
			Ni: m.Ni, Vi: m.Vi, Mi: m.Mi,
			Nj: m.Nj, Vj: m.Vj, Mj: m.Mj,
			Length: m.Length,
		})
	}
	for _, r := range res.Reactions {
		out.Reactions = append(out.Reactions, reactionOut{Node: r.Node, DOF: r.DOF, Force: r.Force})
	}
	return out
}
