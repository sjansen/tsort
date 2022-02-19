package tsort

type GraphOfIDs[ID NodeID] interface {
	NodeIDs() []ID
	EdgeIDs(ID) []ID
}

type NodeID comparable

type Cycle[ID NodeID] struct {
	Path []ID
}

func (c *Cycle[NodeID]) Error() string {
	return "cycle detected"
}

func SortIDs[ID NodeID](g GraphOfIDs[ID]) ([]ID, error) {
	ids := g.NodeIDs()
	s := &state[ID]{
		graph:  g,
		frozen: make(map[ID]bool, len(ids)),
		result: make([]ID, 0),
	}

	for _, id := range ids {
		if cycle := s.dfs(id); cycle != nil {
			return nil, cycle
		}
	}

	return s.result, nil
}

type state[ID NodeID] struct {
	graph  GraphOfIDs[ID]
	frozen map[ID]bool
	result []ID
}

func (s *state[NodeID]) dfs(id NodeID) *Cycle[NodeID] {
	if frozen, ok := s.frozen[id]; ok {
		if frozen {
			return nil
		} else {
			return &Cycle[NodeID]{
				Path: []NodeID{id},
			}
		}
	}
	s.frozen[id] = false

	edges := s.graph.EdgeIDs(id)
	for _, edge := range edges {
		if cycle := s.dfs(edge); cycle != nil {
			cycle.Path = append(cycle.Path, edge)
			return cycle
		}
	}

	s.result = append(s.result, id)
	s.frozen[id] = true

	return nil
}
