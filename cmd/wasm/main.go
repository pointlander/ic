// Copyright 2024 The IC Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"strings"
	"syscall/js"
	"unicode"

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
	pair := tree.Recursive(ic.Pair{Str: []byte(prefix)}, 8)
	index := pair.Idx - count
	if index < 0 {
		index = 0
	}
	end := pair.Idx
	idx := strings.LastIndex(string(tree.Buffer[index:end]), string(pair.Str))
	if idx > 0 {
		end = index + idx + len(prefix)
	}
	fix := string(tree.Buffer[index:end]) + "<hr/>"
	pair.Str = append(pair.Str, tree.Buffer[pair.Idx+1:pair.Idx+count]...)
	tree.GetBooks(&pair)
	output := strings.TrimSpace(fix + string(pair.Str))
	word := false
	html := ""
	for _, value := range output {
		if unicode.IsSpace(value) {
			if word {
				html += fmt.Sprintf("</span>%c", value)
				word = false
			} else {
				html += fmt.Sprintf("%c", value)
			}
		} else {
			if !word {
				html += fmt.Sprintf("<span onclick=\"selectWord(event);\" class=\"fragment\">%c", value)
				word = true
			} else {
				html += fmt.Sprintf("%c", value)
			}
		}
	}
	return html
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
