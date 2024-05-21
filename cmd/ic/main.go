// Copyright 2024 The IC Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"math"
	"strings"

	"github.com/pointlander/ic"
	"github.com/pointlander/ic/books"
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

	/*books := []string{
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
	}*/
	input, ranges := books.GetBible()
	tree := ic.BuildSuffixTree(input, ranges)
	pair := tree.Recursive(ic.Pair{Str: []byte(*FlagPrefix)}, 9)
	index := pair.Idx - 1024
	if index < 0 {
		index = 0
	}
	line := ""
	for i := 0; i < 80; i++ {
		line += "_"
	}
	end := pair.Idx
	idx := strings.LastIndex(string(tree.Buffer[index:end]), string(pair.Str))
	if idx > 0 {
		end = index + idx + len(*FlagPrefix)
	}
	pair.Str = append(pair.Str, tree.Buffer[pair.Idx+1:pair.Idx+1024]...)
	tree.GetBooks(&pair)
	prefix := string(tree.Buffer[index:end]) + "\n" + line + "\n"
	for _, set := range pair.Bok {
		for _, value := range set {
			fmt.Printf("'%s' ", ranges[value].Title)
		}
		fmt.Println()
	}
	fmt.Println(prefix + string(pair.Str))
	fmt.Println(len(string(pair.Str)), len(pair.Bok))
}
