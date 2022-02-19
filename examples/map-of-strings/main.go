package main

import (
	"fmt"
	"os"

	"github.com/sjansen/tsort"
)

var _ tsort.GraphOfIDs[string] = Graph{}

type Graph map[string][]string

func NewGraph() Graph {
	return Graph{}
}

func (g Graph) AddNode(node string, edges ...string) Graph {
	g[node] = edges
	return g
}

func (g Graph) NodeIDs() []string {
	nodes := make([]string, 0, len(g))
	for node := range g {
		nodes = append(nodes, node)
	}
	return nodes
}

func (g Graph) EdgeIDs(node string) []string {
	edges := g[node]
	return edges
}

func main() {
	const brew = "brew install librsvg"
	const convert = "rsvg-convert -h 64 logo.svg > logo.png"
	const curl = "curl -v -o logo.svg https://go.dev/images/go-logo-white.svg"
	const open = "open /tmp/logo.png"
	const popd = "popd"
	const pushd = "pushd /tmp"

	g := NewGraph().
		AddNode(brew).
		AddNode(convert, brew, curl, pushd).
		AddNode(curl, pushd).
		AddNode(open, convert, popd).
		AddNode(popd, convert, curl).
		AddNode(pushd)

	if script, err := tsort.SortIDs[string](g); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		fmt.Println("#!/bin/sh")
		for _, command := range script {
			fmt.Println(command)
		}
		fmt.Println()
	}
}
