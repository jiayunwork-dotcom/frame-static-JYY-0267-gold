package model

import "testing"

func sample() *Model {
	return &Model{
		Nodes: []Node{
			{ID: "A", X: 0, Y: 0, Support: "fixed"},
			{ID: "B", X: 3, Y: 0, Restraint: [3]bool{false, true, false}},
		},
		Elements: []Element{
			{From: "A", To: "B", E: 2.1e11, A: 0.01, I: 4e-5},
		},
	}
}

func TestValidateOK(t *testing.T) {
	if err := Validate(sample()); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}
}

func TestValidateDuplicateNode(t *testing.T) {
	m := sample()
	m.Nodes = append(m.Nodes, Node{ID: "A", X: 1, Y: 1})
	if err := Validate(m); err == nil {
		t.Fatalf("expected duplicate node error")
	}
}

func TestValidateUnknownEndpoint(t *testing.T) {
	m := sample()
	m.Elements[0].To = "Z"
	if err := Validate(m); err == nil {
		t.Fatalf("expected unknown endpoint error")
	}
}

func TestValidateZeroLength(t *testing.T) {
	m := sample()
	m.Nodes[1].X = 0
	m.Nodes[1].Y = 0
	if err := Validate(m); err == nil {
		t.Fatalf("expected zero-length error")
	}
}

func TestValidateNonPositiveSection(t *testing.T) {
	m := sample()
	m.Elements[0].I = 0
	if err := Validate(m); err == nil {
		t.Fatalf("expected I<=0 error")
	}
}

func TestValidateNoConstraint(t *testing.T) {
	m := sample()
	for i := range m.Nodes {
		m.Nodes[i].Restraint = [3]bool{false, false, false}
		m.Nodes[i].Support = ""
	}
	if err := Validate(m); err == nil {
		t.Fatalf("expected no-constraint error")
	}
}

func TestValidateDisconnected(t *testing.T) {
	m := sample()
	// dangling node with no element
	m.Nodes = append(m.Nodes, Node{ID: "C", X: 9, Y: 9, Support: "fixed"})
	if err := Validate(m); err == nil {
		t.Fatalf("expected disconnected-component error")
	}
}
