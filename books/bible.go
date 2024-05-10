// Copyright 2024 The IC Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package books

import (
	"bytes"
	"compress/gzip"
	"embed"
	"io"
)

//go:embed 10.txt.utf-8.gz
var f embed.FS

func GetBible() []byte {
	data, err := f.ReadFile("10.txt.utf-8.gz")
	if err != nil {
		panic(err)
	}
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	data, err = io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return data
}
