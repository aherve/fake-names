package graph

import (
	"math"
	"math/rand"
)

const (
	startChar = "^" // fake char that is the entry node of the graph
	endChar   = "$" //  fake char that represents a string termination
)

type (
	// Graph holds the structure that represent the Markov chain, and provide high-level methods
	Graph struct {
		nodes map[string]*node
		names []string
	}

	node struct {
		key       string
		edges     map[string]int
		weightSum int
	}
)

// InitializeGraph initializes a Graph structure, and build the connections from a list of names
func InitializeGraph(names []string) *Graph {
	g := newGraph(names)
	for _, name := range g.names {
		if len(name) < 3 {
			continue
		}

		// Level 1: monograms
		prev := startChar
		for _, r := range name + endChar {
			g.nodes[prev].addEdge(g.findOrCreateNode(string(r)).key)
			prev = string(r)
		}

		// level 2: bigrams
		prev = startChar + string([]rune(name)[0])
		for i, r := range name + endChar {
			// skip first
			if i == 0 {
				continue
			}
			g.findOrCreateNode(prev).addEdge(g.findOrCreateNode(string(r)).key)
			prev = string([]rune(prev)[1]) + string(r)
		}

		// level 3: trigrams
		prev = startChar + string([]rune(name)[0]) + string([]rune(name)[1])
		for i, r := range name + endChar {
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

// GenerateName generates one or many fake names
func (g Graph) GenerateName() string {
	c1 := startChar
	c2 := g.nextProbable(c1)
	c3 := g.nextProbable(c1 + c2)
	res := c2 + c3

	for c3 != endChar {
		c1, c2, c3 = c2, c3, g.nextProbable(c1+c2+c3)
		res += c3
	}

	res = res[0 : len(res)-1]
	if g.doesExist(res) {
		return g.GenerateName()
	}
	return res
}

func newGraph(names []string) *Graph {
	return &Graph{
		names: names,
		nodes: map[string]*node{
			startChar: {
				key:   startChar,
				edges: map[string]int{},
			},
		},
	}
}

func (n *node) addEdge(r string) {
	n.edges[r] = n.edges[r] + 1
	n.weightSum++
}

// either find a node, or create and return it
func (g *Graph) findOrCreateNode(r string) *node {
	if found, ok := g.nodes[r]; ok {
		return found
	}

	g.nodes[r] = &node{
		key:   r,
		edges: map[string]int{},
	}
	return g.nodes[r]
}

func (g *Graph) link(a, b string) {
	na, _ := g.findOrCreateNode(a), g.findOrCreateNode(b)
	na.addEdge(b)
}

// Weighted reservoir sampling for k=1
// https://en.wikipedia.org/wiki/Reservoir_sampling#Algorithm_A-Res
func (g *Graph) nextProbable(input string) string {
	start := g.findOrCreateNode(input)
	if len(start.edges) == 0 {
		return endChar
	}
	maxWeight := 0.0
	currentRes := endChar

	for r, w := range start.edges {
		rand := math.Pow(rand.Float64(), 1.0/float64(w))
		if rand > maxWeight {
			maxWeight = rand
			currentRes = r
		}
	}
	return currentRes
}

// Checks wether a string already exist or not
func (g Graph) doesExist(input string) bool {
	for _, name := range g.names {
		if name == input {
			return true
		}
	}
	return false
}
