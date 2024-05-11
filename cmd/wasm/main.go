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

var tree *ic.SuffixTree

func main() {
	input := books.GetBible()
	tree = ic.BuildSuffixTree(input)
	js.Global().Set("inference", inferenceWrapper())
	<-make(chan struct{})
}

func Inference(prefix string, seed int64, size, count int) string {
	return tree.Inference(prefix, seed, size, count)
}

func inferenceWrapper() js.Func {
	inferenceFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 4 {
			return "Invalid no of arguments passed"
		}
		return tree.Inference(args[0].String(), int64(args[1].Int()), args[2].Int(), args[3].Int())
	})
	return inferenceFunc
}