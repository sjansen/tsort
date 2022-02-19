package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/sjansen/tsort"
)

type Graph map[string][]string

var _ tsort.DiGraph[string] = Graph{}

func (g Graph) Nodes(fn func(node string) error) error {
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

func main() {
	const brew = "brew install librsvg"
	const convert = "rsvg-convert -h 64 logo.svg > logo.png"
	const curl = "curl -v -o logo.svg https://go.dev/images/go-logo-white.svg"
	const open = "open /tmp/logo.png"
	const popd = "popd"
	const pushd = "pushd /tmp"

	g := Graph{
		open:    []string{convert, popd},
		brew:    []string{},
		convert: []string{brew, curl, pushd},
		curl:    []string{pushd},
		popd:    []string{convert, curl},
		pushd:   []string{},
	}

	if script, err := tsort.TSort[Graph, string](g); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		fmt.Println("#!/bin/sh")
		for _, command := range script {
			fmt.Println(command)
		}
		fmt.Println()
	}
}
