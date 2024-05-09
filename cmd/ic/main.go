// Copyright 2024 The IC Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/pointlander/ic"
)

const (
	// S is the scaling factor for the softmax
	S = 1.0 - 1e-300
)

func softmax(values []float64) {
	max := 0.0
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	s := max * S
	sum := 0.0
	for j, value := range values {
		values[j] = math.Exp(value - s)
		sum += values[j]
	}
	for j, value := range values {
		values[j] = value / sum
	}
}

var (
	// FlagPrefix is the prefix
	FlagPrefix = flag.String("prefix", "God", "the prefix string")
)

func main() {
	flag.Parse()

	rng := rand.New(rand.NewSource(1))
	books := []string{
		"books/10.txt.utf-8",
		"books/84.txt.utf-8",
		"books/145.txt.utf-8",
		"books/1342.txt.utf-8",
		"books/1513.txt.utf-8",
		"books/2641.txt.utf-8",
		"books/2701.txt.utf-8",
		"books/37106.txt.utf-8",
	}
	var input []byte
	for _, book := range books {
		data, err := os.ReadFile(book)
		if err != nil {
			panic(err)
		}
		input = append(input, data...)
	}
	tree := ic.BuildSuffixTree(input)
	sep := *FlagPrefix
	for s := 0; s < 256; s++ {
		dist, sum := []float64{}, 0.0
		for i := 0; i < 256; i++ {
			node := tree.Index(fmt.Sprintf("%s%c", sep, i))
			if node < 0 {
				dist = append(dist, 0)
				continue
			}
			value := float64(tree.Nodes[node].Count)
			sum += value
			dist = append(dist, value)
		}
		for key, value := range dist {
			dist[key] = value / sum
		}
		//softmax(dist)
		selected, sum := rng.Float64(), 0
		for i, value := range dist {
			sum += value
			if sum > selected {
				sep = fmt.Sprintf("%s%c", sep, i)
				break
			}
		}
		fmt.Println(sep)
		fmt.Println("-----------------------------------")
	}
}
