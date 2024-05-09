// Copyright 2024 The IC Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/pointlander/ic"
)

func main() {
	data, err := os.ReadFile("10.txt.utf-8")
	if err != nil {
		panic(err)
	}
	tree := ic.BuildSuffixTree(data)
	sep := "God"
	fmt.Println(tree.Index(sep))
	node := tree.Edges[uint(sep[0])].EndNode
	fmt.Println(tree.Nodes[node].Count)
}
