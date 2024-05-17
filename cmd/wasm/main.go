// Copyright 2024 The IC Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && wasm
// +build js,wasm

package main

import (
	"syscall/js"

	"github.com/pointlander/ic"
	"github.com/pointlander/ic/books"
)

var (
	tree    *ic.SuffixTree
	loading = true
)

func main() {
	js.Global().Set("load", loadWrapper())
	js.Global().Set("inference", inferenceWrapper())
	<-make(chan struct{})
}

func Load() {
	input, ranges := books.GetBible()
	tree = ic.BuildSuffixTree(input, ranges)
	return
}

func Inference(prefix string, seed int64, size, count int) string {
	pair := tree.Recursive(ic.Pair{Str: prefix}, 8)
	_, result := tree.Inference(pair.Str, seed, size, count)
	index := pair.Idx - count
	if index < 0 {
		index = 0
	}
	line := ""
	for i := 0; i < 80; i++ {
		line += "_"
	}
	end := pair.Idx - 8
	if end < 0 {
		end = 0
	}
	fix := string(tree.Buffer[index:end]) + "\n" + line + "\n"
	return fix + result
}

func loadWrapper() js.Func {
	loadFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 0 {
			return "Invalid no of arguments passed"
		}
		Load()
		return true
	})
	return loadFunc
}

func inferenceWrapper() js.Func {
	inferenceFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 4 {
			return "Invalid no of arguments passed"
		}
		return Inference(args[0].String(), int64(args[1].Int()), args[2].Int(), args[3].Int())
	})
	return inferenceFunc
}
