package model

// NodeIndex returns a map from node id to its position in m.Nodes.
func NodeIndex(m *Model) map[string]int {
	idx := make(map[string]int, len(m.Nodes))
	for i, n := range m.Nodes {
		idx[n.ID] = i
	}
	return idx
}

// TotalDOF returns 3 * nodeCount: every joint contributes two translations and
// one in-plane rotation.
func TotalDOF(m *Model) int { return 3 * len(m.Nodes) }
