// Copyright 2024 The IC Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package books

import (
	"bytes"
	"compress/bzip2"
	"embed"
	"io"
)

//go:embed *.bz2
var f embed.FS

func GetBible() []byte {
	bible := []byte{}
	files := []string{
		"10.txt.utf-8.bz2",
		"130.txt.utf-8.bz2",
		"131.txt.utf-8.bz2",
		"398.txt.utf-8.bz2",
		"470.txt.utf-8.bz2",
		"1653.txt.utf-8.bz2",
		"3296.txt.utf-8.bz2",
		"36402.txt.utf-8.bz2",
		"5657.txt.utf-8.bz2",
		"57121.txt.utf-8.bz2",
		"59041.txt.utf-8.bz2",
	}
	for _, file := range files {
		data, err := f.ReadFile(file)
		if err != nil {
			panic(err)
		}
		reader := bzip2.NewReader(bytes.NewReader(data))
		data, err = io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		bible = append(bible, data...)
	}
	return bible
}
