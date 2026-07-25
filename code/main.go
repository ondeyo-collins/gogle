package main

import (
	"fmt"
	"github.com/collins994/gogle/code/index"
	"github.com/collins994/gogle/code/parser"
	"github.com/collins994/gogle/code/text"
	"os"
	"sort"
)

func main() {
	var (
		event = parser.ParserEvent{
			Buffer: make([]rune, 0, 256), // one kilobyte
			Error:  nil,
		}
		terms = make([]string, 0)
		index = index.NewIndex()
	)

	for count := 1; count < 352; count++ {
		file, err := os.Open(fmt.Sprintf("gl2/%d.xhtml", count))
		if err != nil {
			fmt.Printf("[ERROR]: %v\n", err)
			continue
		}
		terms = terms[:0]
		var next = parser.Parse(file)
		for {
			next(&event)
			if event.Type == parser.EventTypeEndDocument {
				break
			}
			if event.Type == parser.EventTypeTextNode {
				// print(string(event.Buffer), " ==> ")
				text.Normalize(&event.Buffer)
				text.PorterStem(&event.Buffer)
				// println(string(event.Buffer))
				terms = append(terms, string(event.Buffer))
			}
		}
		sort.Strings(terms)
		index.Add(file.Name(), terms)
		file.Close()
	}
	index.Search("GL_PIXEL_UNPACK_BUFFER")
	index.Search("GL_PIXEL_PACK_BUFFER")
	index.Search("GL_PIXEL_PACK_BUFFER_BINDING")
	index.Search("COLLIN")
}
