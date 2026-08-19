package model

// Node is a joint of the plane frame. Coordinates are in a global 2D system;
// the third DOF is the in-plane rotation theta about the z axis.
type Node struct {
	ID       string  `json:"id"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Restraint [3]bool `json:"restraint"` // [ux, uy, theta]; true => constrained
	Support   string  `json:"support"`   // optional: pin | fixed | roller-x | roller-y
}

// EffectiveRestraint resolves the three DOF constraints, letting a symbolic
// support type expand into the [ux, uy, theta] booleans when no explicit
// boolean array was given.
func (n Node) EffectiveRestraint() [3]bool {
	if n.Restraint != [3]bool{false, false, false} {
		return n.Restraint
	}
	switch n.Support {
	case "pin", "pinned", "hinge", "铰支":
		return [3]bool{true, true, false}
	case "fixed", "固支":
		return [3]bool{true, true, true}
	case "roller-x", "滑动-x", "roller":
		return [3]bool{false, true, false}
	case "roller-y", "滑动-y":
		return [3]bool{true, false, false}
	default:
		return n.Restraint
	}
}

// Element is a prismatic frame member between two nodes.
type Element struct {
	From string    `json:"from"`
	To   string    `json:"to"`
	E    float64   `json:"E"`  // Young's modulus
	A    float64   `json:"A"`  // cross-section area
	I    float64   `json:"I"`  // second moment of area
	Dist *DistLoad `json:"dist"` // optional local transverse UDL (downward +)
}

// DistLoad is a uniformly distributed transverse load along the element,
// expressed in the element's local coordinate system (perpendicular to the
// member axis, positive downward / -y_local).
type DistLoad struct {
	Q float64 `json:"q"`
}

// NodeLoad is a concentrated force and/or moment applied at a node.
type NodeLoad struct {
	Node string  `json:"node"`
	FX   float64 `json:"fx"`
	FY   float64 `json:"fy"`
	MZ   float64 `json:"mz"`
}

// Model is the full structural description submitted by the user.
type Model struct {
	Nodes    []Node      `json:"nodes"`
	Elements []Element   `json:"elements"`
	Loads    []NodeLoad  `json:"loads"`
}
