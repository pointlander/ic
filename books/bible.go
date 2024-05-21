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

type Range struct {
	Title string
	File  string
	Begin int
	End   int
}

func GetBible() ([]byte, []Range) {
	bible := []byte{}
	ranges := []Range{
		{
			Title: "The King James Version of the Bible",
			File:  "10.txt.utf-8.bz2",
		},
		{
			Title: "Orthodoxy",
			File:  "130.txt.utf-8.bz2",
		},
		{
			Title: "The Pilgrim's Progress from this world to that which is to come",
			File:  "131.txt.utf-8.bz2",
		},
		{
			Title: "The First Book of Adam and Eve",
			File:  "398.txt.utf-8.bz2",
		},
		{
			Title: "Heretics",
			File:  "470.txt.utf-8.bz2",
		},
		{
			Title: "The Imitation of Christ",
			File:  "1653.txt.utf-8.bz2",
		},
		{
			Title: "The Confessions of St. Augustine",
			File:  "3296.txt.utf-8.bz2",
		},
		{
			Title: "On Union with God",
			File:  "36402.txt.utf-8.bz2",
		},
		{
			Title: "The Practice of the Presence of God",
			File:  "5657.txt.utf-8.bz2",
		},
		{
			Title: "Humility: The Beauty of Holiness",
			File:  "57121.txt.utf-8.bz2",
		},
		{
			Title: "Sermons Preached at the Church of St. Paul the Apostle, New York, During the Year 1861",
			File:  "59041.txt.utf-8.bz2",
		},
	}
	for i, file := range ranges {
		data, err := f.ReadFile(file.File)
		if err != nil {
			panic(err)
		}
		reader := bzip2.NewReader(bytes.NewReader(data))
		data, err = io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		str := []byte{}
		for _, r := range string(data) {
			if r < 256 {
				str = append(str, byte(r))
			}
		}
		ranges[i].Begin = len(bible)
		bible = append(bible, str...)
		ranges[i].End = len(bible)
	}
	return bible, ranges
}
