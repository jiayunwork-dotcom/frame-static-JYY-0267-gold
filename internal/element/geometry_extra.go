package element

import (
	"math"

	"frame-static/internal/model"
)

// Midpoint returns the global coordinates of the member's midpoint.
func (g *Geometry) Midpoint(n1, n2 model.Node) (mx, my float64) {
	return (n1.X + n2.X) / 2, (n1.Y + n2.Y) / 2
}

// AngleDeg returns the member orientation angle in degrees, measured
// counterclockwise from the global x axis.
func (g *Geometry) AngleDeg() float64 {
	return math.Atan2(g.S, g.C) * 180 / math.Pi
}

// PerpDirection returns the direction cosines of the local y axis (rotated
// 90° counterclockwise from the member axis) in global coordinates.
func (g *Geometry) PerpDirection() (cx, cy float64) {
	return -g.S, g.C
}

// LocalProjection resolves a global displacement vector into its local axial
// and transverse components for the member.
func (g *Geometry) LocalProjection(dx, dy float64) (axial, transverse float64) {
	axial = g.C*dx + g.S*dy
	transverse = -g.S*dx + g.C*dy
	return
}

// GlobalPoint maps a local station x (0..Length) along the member axis to
// global coordinates given the from-node position.
func (g *Geometry) GlobalPoint(x0, y0, x float64) (gx, gy float64) {
	return x0 + g.C*x, y0 + g.S*x
}
