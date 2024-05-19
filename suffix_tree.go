// Copyright 2024 The IC Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate ./generate.sh

package ic

import (
	"fmt"
	"math/rand"

	"github.com/pointlander/ic/books"
)

type Edge struct {
	FirstIndex, LastIndex, StartNode, EndNode int
}

type Node struct {
	Node  int
	Count int
	Books []int
}

type SuffixTree struct {
	Edges  map[uint]Edge
	Nodes  []Node
	Buffer []byte
	Ranges []books.Range
}

const SYMBOL_SIZE = 9

func BuildSuffixTree(input []byte, ranges []books.Range) *SuffixTree {
	length := len(input)
	size := 2 * length
	edges, nodes := make(map[uint]Edge, size), make([]Node, size)
	for i := range nodes {
		nodes[i].Node = -1
	}

	putEdge := func(edge Edge) {
		symbol := uint(input[edge.FirstIndex])
		edges[(uint(edge.StartNode)<<SYMBOL_SIZE)|symbol] = edge
	}

	getEdge := func(node, index int) Edge {
		symbol := uint(input[index])
		return edges[(uint(node)<<SYMBOL_SIZE)|symbol]
	}

	hasEdge := func(node, index int) bool {
		symbol := uint(input[index])
		_, has := edges[(uint(node)<<SYMBOL_SIZE)|symbol]
		return has
	}

	putNode := func(node, parent int) {
		nodes[node].Node = parent
		nodes[node].Count++
		for i, r := range ranges {
			if parent >= r.Begin && parent <= r.End {
				for _, n := range nodes[node].Books {
					if n == i {
						return
					}
				}
				nodes[node].Books = append(nodes[node].Books, i)
				break
			}
		}
	}

	getNode := func(node int) int {
		return nodes[node].Node
	}

	node_count, parent_node, origin, first_index, last_index := 1, 0, 0, 0, -1

	findEdge := func(i int, v byte) bool {
		if first_index > last_index {
			if hasEdge(origin, i) {
				return true
			}
		} else {
			edge, last_edge := getEdge(origin, first_index), last_index-first_index
			last_edge += edge.FirstIndex
			next_edge := last_edge + 1
			if v == input[next_edge] {
				return true
			}
			putEdge(Edge{FirstIndex: edge.FirstIndex, LastIndex: last_edge, StartNode: origin, EndNode: node_count})
			putNode(node_count, origin)
			edge.FirstIndex, edge.StartNode = next_edge, node_count
			putEdge(edge)
			parent_node = node_count
			node_count++
		}
		return false
	}

	canonize := func() {
		if first_index > last_index {
			return
		}
		edge := getEdge(origin, first_index)
		span := edge.LastIndex - edge.FirstIndex
		for span <= (last_index - first_index) {
			first_index += span + 1
			origin = edge.EndNode
			if first_index <= last_index {
				edge = getEdge(edge.EndNode, first_index)
				span = edge.LastIndex - edge.FirstIndex
			}
		}
	}

	addEdge := func(i int) {
		putEdge(Edge{FirstIndex: i, LastIndex: length, StartNode: parent_node, EndNode: node_count})
		node_count++
		if origin == 0 {
			first_index++
		} else {
			origin = getNode(origin)
		}
		canonize()
	}

	for i, v := range input {
		parent_node = origin
		if findEdge(i, v) {
			last_index++
			canonize()
			continue
		}
		addEdge(i)
		last_parent_node := parent_node
		parent_node = origin
		for !findEdge(i, v) {
			addEdge(i)
			putNode(last_parent_node, parent_node)
			last_parent_node, parent_node = parent_node, origin
		}
		putNode(last_parent_node, parent_node)
		last_index++
		canonize()
	}

	/*add the last end nodes*/
	putEdge = func(edge Edge) {
		symbol := uint(256)
		if int(edge.FirstIndex) < length {
			symbol = uint(input[edge.FirstIndex])
		}

		edges[(uint(edge.StartNode)<<SYMBOL_SIZE)|symbol] = edge
	}

	getEdge = func(node, index int) Edge {
		symbol := uint(256)
		if index < length {
			symbol = uint(input[index])
		}

		return edges[(uint(node)<<SYMBOL_SIZE)|symbol]
	}

	hasEdge = func(node, index int) bool {
		symbol := uint(256)
		if index < length {
			symbol = uint(input[index])
		}

		_, has := edges[(uint(node)<<SYMBOL_SIZE)|symbol]
		return has
	}

	findEdge = func(i int, v byte) bool {
		if first_index > last_index {
			if hasEdge(origin, i) {
				return true
			}
		} else {
			edge, last_edge := getEdge(origin, first_index), last_index-first_index
			last_edge += edge.FirstIndex
			next_edge := last_edge + 1
			if next_edge == length {
				return true
			}
			putEdge(Edge{FirstIndex: edge.FirstIndex, LastIndex: last_edge, StartNode: origin, EndNode: node_count})
			putNode(node_count, origin)
			edge.FirstIndex, edge.StartNode = next_edge, node_count
			putEdge(edge)
			parent_node = node_count
			node_count++
		}
		return false
	}

	tree := &SuffixTree{
		Edges:  edges,
		Nodes:  nodes,
		Buffer: input,
		Ranges: ranges,
	}
	parent_node = origin
	if findEdge(length, 0) {
		return tree
	}
	addEdge(length)
	last_parent_node := parent_node
	parent_node = origin
	for !findEdge(length, 0) {
		addEdge(length)
		putNode(last_parent_node, parent_node)
		last_parent_node, parent_node = parent_node, origin
	}
	putNode(last_parent_node, parent_node)
	return tree
}

func (tree *SuffixTree) Index(sep string) (Edge, int) {
	i, node, last_i := 0, 0, 0
	var last_edge Edge
search:
	for i < len(sep) {
		edge, has := tree.Edges[(uint(node)<<SYMBOL_SIZE)+uint(sep[i])]
		if !has {
			return edge, -1
		}
		/*fmt.Printf("at node %v %v %v %v\n", edge.first_index, edge.last_index, edge.start_node, edge.end_node)
		  fmt.Printf("found '%c'\n", sep[i])*/
		node, last_edge, last_i, i = int(edge.EndNode), edge, i, i+1
		if edge.FirstIndex >= edge.LastIndex {
			continue search
		}
		for index := edge.FirstIndex + 1; index <= edge.LastIndex && i < len(sep); index++ {
			if sep[i] != tree.Buffer[index] {
				return edge, -1
			}
			/*fmt.Printf("found '%c'\n", sep[i])*/
			i++
		}
	}
	/*fmt.Printf("%v\n", string(tree.buffer[int(last_edge.first_index) - last_i:int(last_edge.first_index) - last_i + len(sep)]))*/
	return last_edge, int(last_edge.FirstIndex) - last_i
}

func (tree *SuffixTree) MultiIndex(sep string) []Pair {
	i, node, last_i := 0, 0, 0
	var last_edge Edge
	pairs := make([]Pair, 256)
search:
	for i < len(sep) {
		edge, has := tree.Edges[(uint(node)<<SYMBOL_SIZE)+uint(sep[i])]
		if !has {
			return pairs
		}
		/*fmt.Printf("at node %v %v %v %v\n", edge.first_index, edge.last_index, edge.start_node, edge.end_node)
		  fmt.Printf("found '%c'\n", sep[i])*/
		node, last_edge, last_i, i = int(edge.EndNode), edge, i, i+1
		if edge.FirstIndex >= edge.LastIndex {
			continue search
		}
		for index := edge.FirstIndex + 1; index <= edge.LastIndex && i < len(sep); index++ {
			if sep[i] != tree.Buffer[index] {
				return pairs
			}
			/*fmt.Printf("found '%c'\n", sep[i])*/
			i++
			if i == len(sep) {
				index++
				if index <= edge.LastIndex {
					j := tree.Buffer[index]
					//fmt.Println(tree.Nodes[edge.StartNode].Books)
					pairs[j].Int = tree.Nodes[edge.StartNode].Count
					pairs[j].Str = fmt.Sprintf("%s%c", sep, j)
					pairs[j].Idx = index
					return pairs
				} else {
					node, last_edge = int(edge.EndNode), edge
					break search
				}
			}
		}
	}
	_, _ = last_i, last_edge
	for i := range pairs {
		edge, has := tree.Edges[(uint(node)<<SYMBOL_SIZE)+uint(i)]
		if !has {
			continue
		}
		j := tree.Buffer[edge.FirstIndex+1]
		pairs[j].Int = tree.Nodes[edge.StartNode].Count
		pairs[j].Str = fmt.Sprintf("%s%c", sep, i)
		pairs[j].Idx = edge.FirstIndex + 1
	}
	/*fmt.Printf("%v\n", string(tree.buffer[int(last_edge.first_index) - last_i:int(last_edge.first_index) - last_i + len(sep)]))*/
	return pairs
}

func (tree *SuffixTree) Brute(prefix string, seed int64, size, count int) []string {
	results := make([]string, 256)
	for i := range results {
		found, result := tree.Inference(fmt.Sprintf("%s%c", prefix, i), seed, size, count)
		if found {
			results[i] = result
		}
	}
	return results
}

func (tree *SuffixTree) Inference(prefix string, seed int64, size, count int) (bool, string) {
	rng := rand.New(rand.NewSource(seed))
	found := false
	for s := 0; s < count; s++ {
		dist, sum := []float64{}, 0.0
		for i := 0; i < 256; i++ {
			edge, has := tree.Index(fmt.Sprintf("%s%c", prefix, i))
			node := tree.Nodes[edge.StartNode]
			if has < 0 {
				dist = append(dist, 0)
				continue
			}
			value := float64(node.Count)
			sum += value
			dist = append(dist, value)
		}
		for key, value := range dist {
			dist[key] = value / sum
		}
		selected, sum := rng.Float64(), 0
		for i, value := range dist {
			sum += value
			if sum > selected {
				prefix = fmt.Sprintf("%s%c", prefix, i)
				found = true
				break
			}
		}
	}
	return found, prefix
}

// Pair is a pair of values
type Pair struct {
	Int int
	Str string
	Idx int
}

func (tree *SuffixTree) Recursive(prefix Pair, count int) Pair {
	if count == 0 {
		return prefix
	}

	sum := prefix.Int
	entries := tree.MultiIndex(prefix.Str)
	found := false
	for i, pair := range entries {
		if pair.Str != "" {
			found = true
			pair = tree.Recursive(pair, count-1)
			entries[i] = pair
			sum += pair.Int
		}
	}

	if found {
		max, pair := prefix.Int, Pair{
			Int: sum,
			Str: prefix.Str,
			Idx: prefix.Idx,
		}
		for _, v := range entries {
			if v.Int > max {
				max, pair.Str, pair.Idx = v.Int, v.Str, v.Idx
			}
		}
		return pair
	}
	return prefix
}
