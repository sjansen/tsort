package tsort

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

var _ GraphOfIDs[string] = Graph{}

type Graph map[string][]string

func (g Graph) NodeIDs() []string {
	nodes := make([]string, 0, len(g))
	for node := range g {
		nodes = append(nodes, node)
	}

	sort.Strings(nodes) // Sorting is required for consistent result.
	return nodes
}

func (g Graph) EdgeIDs(node string) []string {
	edges := g[node]

	sort.Strings(edges) // Sorting is required for consistent result.
	return edges
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

			actual, err := SortIDs[string](tc.graph)
			require.NoError(err)
			require.Equal(tc.expected, actual)
		})
	}
}
