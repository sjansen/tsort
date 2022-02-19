package tsort

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

type Graph map[string][]string

var _ DiGraph[string] = Graph{}

func (g Graph) Nodes(fn func(node string) error) error {
	// Sorting is required for stable result order.
	nodes := make([]string, 0, len(g))
	for node, edges := range g {
		nodes = append(nodes, node)
		sort.Strings(edges)
	}
	sort.Strings(nodes)

	for _, node := range nodes {
		if err := fn(node); err != nil {
			return err
		}
	}
	return nil
}

func (g Graph) Edges(node string, fn func(edge string) error) error {
	for _, edge := range g[node] {
		if err := fn(edge); err != nil {
			return err
		}
	}
	return nil
}

func TestHappyPath(t *testing.T) {
	for name, tc := range map[string]struct {
		expected []string
		graph    Graph
	}{
		"simple": {
			expected: []string{"b", "d", "c", "a"},
			graph: Graph{
				"a": []string{"b", "c", "d"},
				"b": []string{},
				"c": []string{"d"},
				"d": []string{},
			},
		},
		"mixed": {
			expected: []string{"b", "d", "c", "a", "f", "g", "h", "e"},
			graph: Graph{
				"h": []string{},
				"g": []string{},
				"f": []string{},
				"e": []string{"h", "g", "f"},
				"d": []string{},
				"c": []string{"d"},
				"b": []string{},
				"a": []string{"c", "b", "d"},
			},
		},
		"bad wolf": {
			expected: []string{"b", "a", "d", " ", "w", "o", "l", "f"},
			graph: Graph{
				" ": []string{"a", "d"},
				"a": []string{"b"},
				"b": []string{},
				"d": []string{"a"},
				"f": []string{" ", "a", "l"},
				"l": []string{" ", "b", "o"},
				"o": []string{"d", "w"},
				"w": []string{" ", "d"},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)

			actual, err := TSort[Graph, string](tc.graph)
			require.NoError(err)
			require.Equal(tc.expected, actual)
		})
	}
}
