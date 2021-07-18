package main

import (
	"math"
	"math/rand"
)

type (
	Graph struct {
		nodes map[string]*Node
		names []string
	}

	Node struct {
		key       string
		edges     map[string]int
		weightSum int
	}
)

func NewGraph(names []string) *Graph {
	return &Graph{
		names: names,
		nodes: map[string]*Node{
			"^": {
				key:   "^",
				edges: map[string]int{},
			},
		},
	}
}

func (n *Node) addEdge(r string) {
	n.edges[r] = n.edges[r] + 1
	n.weightSum++
}

func (g *Graph) findOrCreateNode(r string) *Node {
	if found, ok := g.nodes[r]; ok {
		return found
	}

	g.nodes[r] = &Node{
		key:   r,
		edges: map[string]int{},
	}
	return g.nodes[r]
}

func (g *Graph) link(a, b string) {
	na, _ := g.findOrCreateNode(a), g.findOrCreateNode(b)
	na.addEdge(b)
}

func (g *Graph) nextProbable(input string) string {
	start := g.findOrCreateNode(input)
	if len(start.edges) == 0 {
		return "."
	}
	maxWeight := 0.0
	currentRes := "."

	for r, w := range start.edges {
		rand := math.Pow(rand.Float64(), 1.0/float64(w))
		if rand > maxWeight {
			maxWeight = rand
			currentRes = r
		}
	}
	return currentRes
}

func (g Graph) generateName() string {
	c1 := "^"
	c2 := g.nextProbable(c1)
	c3 := g.nextProbable(c1 + c2)
	res := c2 + c3

	for c3 != "." {
		c1, c2, c3 = c2, c3, g.nextProbable(c1+c2+c3)
		res += c3
	}

	res = res[0 : len(res)-1]
	if g.doesExist(res) {
		return g.generateName()
	}
	return res
}

func (g Graph) doesExist(input string) bool {
	for _, name := range g.names {
		if name == input {
			return true
		}
	}
	return false
}

func initializeGraph(names []string) *Graph {
	g := NewGraph(names)
	for _, name := range g.names {
		if len(name) < 3 {
			continue
		}

		// Level 1: monograms
		prev := "^"
		for _, r := range name + "." {
			g.nodes[prev].addEdge(g.findOrCreateNode(string(r)).key)
			prev = string(r)
		}

		// level 2: bigrams
		prev = "^" + string([]rune(name)[0])
		for i, r := range name + "." {
			// skip first
			if i == 0 {
				continue
			}
			g.findOrCreateNode(prev).addEdge(g.findOrCreateNode(string(r)).key)
			prev = string([]rune(prev)[1]) + string(r)
		}

		// level 3: trigrams
		prev = "^" + string([]rune(name)[0]) + string([]rune(name)[1])
		for i, r := range name + "." {
			// skip first two
			if i < 2 {
				continue
			}
			g.findOrCreateNode(prev).addEdge(g.findOrCreateNode(string(r)).key)
			prev = string([]rune(prev)[1]) + string([]rune(prev)[2]) + string(r)
		}
	}

	return g
}
