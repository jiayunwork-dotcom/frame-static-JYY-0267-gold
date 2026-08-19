package model

import "testing"

func TestNodeByIDAndQuery(t *testing.T) {
	m := sample()
	n, i, ok := m.NodeByID("B")
	if !ok || n.X != 3 || i != 1 {
		t.Fatalf("NodeByID B = (%v, %d, %v)", n, i, ok)
	}
	if _, _, ok := m.NodeByID("Z"); ok {
		t.Fatalf("expected missing node")
	}
	if len(m.ElementsIncident("A")) != 1 {
		t.Fatalf("incident count wrong")
	}
	if m.RestrainedCount() != 4 { // fixed A (3) + roller B (uy)
		t.Fatalf("restrained = %d", m.RestrainedCount())
	}
	if m.FreeCount() != TotalDOF(m)-m.RestrainedCount() {
		t.Fatalf("free count inconsistent")
	}
	if m.MaxSpan() != 3 {
		t.Fatalf("max span = %v", m.MaxSpan())
	}
}

func TestCloneIndependent(t *testing.T) {
	m := sample()
	c := m.Clone()
	c.Nodes[0].X = 999
	c.Elements[0].I = 1
	if m.Nodes[0].X == 999 || m.Elements[0].I == 1 {
		t.Fatalf("clone shares memory")
	}
}

func TestLoadsQuery(t *testing.T) {
	m := sample()
	m.Loads = []NodeLoad{{Node: "B", FY: -500}, {Node: "A", FX: 10}}
	fx, fy, _ := m.TotalAppliedForce()
	if fx != 10 || fy != -500 {
		t.Fatalf("applied = (%v, %v)", fx, fy)
	}
	if len(m.LoadedNodes()) != 2 {
		t.Fatalf("loaded nodes = %v", m.LoadedNodes())
	}
	if m.HasDistLoads() {
		t.Fatalf("no dist loads expected")
	}
}
