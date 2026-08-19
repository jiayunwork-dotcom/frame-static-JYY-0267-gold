package element

import (
	"errors"
	"math"

	"frame-static/internal/model"
)

// Geometry holds a member's length and direction cosines in global space.
type Geometry struct {
	Length float64
	C      float64 // cos(angle from "from" node to "to" node)
	S      float64 // sin(angle)
}

// GeometryOf computes member geometry between two nodes.
func GeometryOf(n1, n2 model.Node) (*Geometry, error) {
	dx, dy := n2.X-n1.X, n2.Y-n1.Y
	L := math.Hypot(dx, dy)
	if L == 0 {
		return nil, errors.New("zero-length element")
	}
	return &Geometry{Length: L, C: dx / L, S: dy / L}, nil
}
