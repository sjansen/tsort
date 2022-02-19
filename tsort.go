package tsort

type Node interface {
	comparable
}

type Cycle[N Node] struct {
	Path []N
}

func (c *Cycle[Node]) Error() string {
	return "cycle detected"
}

type DiGraph[N Node] interface {
	Nodes(func(N) error) error
	Edges(N, func(N) error) error
}

type state[N Node] struct {
	graph  DiGraph[N]
	frozen map[N]bool
	result []N
}

func TSort[G DiGraph[N], N Node](g G) ([]N, error) {
	s := &state[N]{
		graph:  g,
		frozen: make(map[N]bool),
		result: make([]N, 0),
	}

	err := g.Nodes(func(node N) error {
		cycle, err := s.dfs(node)
		if cycle != nil {
			return cycle
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return s.result, nil
}

func (s *state[Node]) dfs(node Node) (*Cycle[Node], error) {
	if frozen, ok := s.frozen[node]; ok {
		if frozen {
			return nil, nil
		} else {
			return &Cycle[Node]{
				Path: []Node{node},
			}, nil
		}
	}
	s.frozen[node] = false

	err := s.graph.Edges(node, func(edge Node) error {
		cycle, err := s.dfs(edge)
		if cycle != nil {
			cycle.Path = append(cycle.Path, edge)
			return cycle
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	s.result = append(s.result, node)
	s.frozen[node] = true

	return nil, nil
}
