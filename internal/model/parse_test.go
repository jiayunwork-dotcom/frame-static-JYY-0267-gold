package model

import (
	"strings"
	"testing"
)

func TestParseValid(t *testing.T) {
	src := `{
		"nodes": [
			{"id":"A","x":0,"y":0,"support":"fixed"},
			{"id":"B","x":3,"y":0,"restraint":[false,true,false]}
		],
		"elements":[{"from":"A","to":"B","E":2.1e11,"A":0.01,"I":4e-5}],
		"loads":[{"node":"B","fy":-10}]
	}`
	m, err := ParseModelBytes([]byte(src))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(m.Nodes) != 2 || len(m.Elements) != 1 {
		t.Fatalf("counts: %d nodes %d elems", len(m.Nodes), len(m.Elements))
	}
	r := m.Nodes[0].EffectiveRestraint()
	if !r[0] || !r[1] || !r[2] {
		t.Fatalf("fixed support should restrain all 3 DOF, got %v", r)
	}
}

func TestParseRejectsUnknownField(t *testing.T) {
	src := `{"nodes":[{"id":"A","x":0,"y":0}],"bogus":1}`
	if _, err := ParseModelBytes([]byte(src)); err == nil {
		t.Fatalf("expected error for unknown field")
	}
}

func TestEffectiveRestraintFallback(t *testing.T) {
	n := Node{ID: "X", X: 0, Y: 0}
	if got := n.EffectiveRestraint(); got != [3]bool{false, false, false} {
		t.Fatalf("empty node should be fully free, got %v", got)
	}
	n2 := Node{ID: "P", X: 0, Y: 0, Support: "pin"}
	want := [3]bool{true, true, false}
	if n2.EffectiveRestraint() != want {
		t.Fatalf("pin -> %v, want %v", n2.EffectiveRestraint(), want)
	}
}

func TestParseErrorWrapped(t *testing.T) {
	_, err := ParseModelBytes([]byte("not json"))
	if err == nil || !strings.Contains(err.Error(), "json") {
		t.Fatalf("expected json error, got %v", err)
	}
}
