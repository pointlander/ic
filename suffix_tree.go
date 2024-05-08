// Copyright 2024 The IC Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ic

type Edge struct {
	first_index, last_index, start_node, end_node int
}

type SuffixTree struct {
	Edges  map[uint]Edge
	Nodes  []int
	Buffer []byte
}

const SYMBOL_SIZE = 9

func BuildSuffixTree(input []byte) *SuffixTree {
	length := len(input)
	size := 2 * length
	edges, nodes := make(map[uint]Edge, size), make([]int, size)
	for i := range nodes {
		nodes[i] = -1
	}

	putEdge := func(edge Edge) {
		symbol := uint(input[edge.first_index])
		edges[(uint(edge.start_node)<<SYMBOL_SIZE)|symbol] = edge
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
		nodes[node] = parent
	}

	getNode := func(node int) int {
		return nodes[node]
	}

	node_count, parent_node, origin, first_index, last_index := 1, 0, 0, 0, -1

	findEdge := func(i int, v byte) bool {
		if first_index > last_index {
			if hasEdge(origin, i) {
				return true
			}
		} else {
			edge, last_edge := getEdge(origin, first_index), last_index-first_index
			last_edge += edge.first_index
			next_edge := last_edge + 1
			if v == input[next_edge] {
				return true
			}
			putEdge(Edge{first_index: edge.first_index, last_index: last_edge, start_node: origin, end_node: node_count})
			putNode(node_count, origin)
			edge.first_index, edge.start_node = next_edge, node_count
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
		span := edge.last_index - edge.first_index
		for span <= (last_index - first_index) {
			first_index += span + 1
			origin = edge.end_node
			if first_index <= last_index {
				edge = getEdge(edge.end_node, first_index)
				span = edge.last_index - edge.first_index
			}
		}
	}

	addEdge := func(i int) {
		putEdge(Edge{first_index: i, last_index: length, start_node: parent_node, end_node: node_count})
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
		if int(edge.first_index) < length {
			symbol = uint(input[edge.first_index])
		}

		edges[(uint(edge.start_node)<<SYMBOL_SIZE)|symbol] = edge
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
			last_edge += edge.first_index
			next_edge := last_edge + 1
			if next_edge == length {
				return true
			}
			putEdge(Edge{first_index: edge.first_index, last_index: last_edge, start_node: origin, end_node: node_count})
			putNode(node_count, origin)
			edge.first_index, edge.start_node = next_edge, node_count
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

func (tree *SuffixTree) Index(sep string) int {
	i, node, last_i := 0, 0, 0
	var last_edge Edge
search:
	for i < len(sep) {
		edge, has := tree.Edges[(uint(node)<<SYMBOL_SIZE)+uint(sep[i])]
		if !has {
			return -1
		}
		/*fmt.Printf("at node %v %v %v %v\n", edge.first_index, edge.last_index, edge.start_node, edge.end_node)
		  fmt.Printf("found '%c'\n", sep[i])*/
		node, last_edge, last_i, i = int(edge.end_node), edge, i, i+1
		if edge.first_index >= edge.last_index {
			continue search
		}
		for index := edge.first_index + 1; index <= edge.last_index && i < len(sep); index++ {
			if sep[i] != tree.Buffer[index] {
				return -1
			}
			/*fmt.Printf("found '%c'\n", sep[i])*/
			i++
		}
	}
	/*fmt.Printf("%v\n", string(tree.buffer[int(last_edge.first_index) - last_i:int(last_edge.first_index) - last_i + len(sep)]))*/
	return int(last_edge.first_index) - last_i
}
