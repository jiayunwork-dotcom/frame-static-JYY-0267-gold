package model

// TotalAppliedForce sums all nodal loads into a resultant force and moment.
func (m *Model) TotalAppliedForce() (fx, fy, mz float64) {
	for _, ld := range m.Loads {
		fx += ld.FX
		fy += ld.FY
		mz += ld.MZ
	}
	return
}

// HasDistLoads reports whether any member carries a distributed load.
func (m *Model) HasDistLoads() bool {
	for _, e := range m.Elements {
		if e.Dist != nil {
			return true
		}
	}
	return false
}

// TotalDistLoad sums the magnitudes of all member distributed loads.
func (m *Model) TotalDistLoad() float64 {
	var s float64
	for _, e := range m.Elements {
		if e.Dist != nil {
			s += e.Dist.Q
		}
	}
	return s
}

// LoadedNodes returns the ids of nodes that carry at least one applied load.
func (m *Model) LoadedNodes() []string {
	var out []string
	for _, ld := range m.Loads {
		out = append(out, ld.Node)
	}
	return out
}

// LoadCount returns how many nodal loads are defined.
func (m *Model) LoadCount() int { return len(m.Loads) }
