package main

import (
	"fmt"
	"os"

	"github.com/tidwall/btree"

	"github.com/sjansen/tsort"
)

type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

var _ tsort.GraphOfIDs[string] = &Graph[string]{}

type Graph[ID ordered] struct {
	nodes btree.Map[ID, *btree.Set[ID]]
}

func NewGraph[ID ordered]() *Graph[ID] {
	return &Graph[ID]{}
}

func (g *Graph[ID]) AddNode(id ID, edges ...ID) *Graph[ID] {
	tmp := &btree.Set[ID]{}
	for _, edge := range edges {
		tmp.Insert(edge)
	}
	g.nodes.Set(id, tmp)
	return g
}

func (g *Graph[ID]) NodeIDs() []ID {
	result := make([]ID, 0, g.nodes.Len())
	g.nodes.Scan(func(id ID, value *btree.Set[ID]) bool {
		result = append(result, id)
		return true
	})
	return result
}

func (g *Graph[ID]) EdgeIDs(id ID) []ID {
	edges, ok := g.nodes.Get(id)
	if !ok {
		return nil
	}
	result := make([]ID, 0, edges.Len())
	edges.Scan(func(id ID) bool {
		result = append(result, id)
		return true
	})
	return result
}

func main() {
	g := NewGraph[int]().
		AddNode(1, 2).
		AddNode(2, 3).
		AddNode(3, 4).
		AddNode(4, 5).
		AddNode(5)

	if result, err := tsort.SortIDs[int](g); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		for _, s := range result {
			fmt.Println(s)
		}
		fmt.Println()
	}
}
