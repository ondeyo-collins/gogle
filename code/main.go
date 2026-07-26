package main

import (
	// "fmt"
	"bufio"
	"github.com/collins994/gogle/code/index"
	"github.com/collins994/gogle/code/parser"
	"github.com/collins994/gogle/code/text"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	var (
		file          *os.File
		index         = index.NewIndex()
		nextEventFunc func(*parser.ParserEvent)
		terms         = make([]string, 0)
		event         = parser.ParserEvent{
			Buffer: make([]rune, 256),
			Error:  nil,
		}
		query string
		scanner = bufio.NewScanner(os.Stdin)
	)

	filepath.WalkDir("files", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			println("[ERROR]: ", err.Error())
			return err
		}
		if d.IsDir() {
			return nil
		}
		file, err = os.Open(path)
		if err != nil {
			println("[ERROR]: ", err.Error())
			return nil
		}

		terms = terms[:0]
		nextEventFunc = parser.Parse(file)
		for {
			nextEventFunc(&event)
			if event.Type == parser.EventTypeEndDocument {
				break
			}
			if event.Type == parser.EventTypeTextNode {
				text.Normalize(&event.Buffer)
				text.PorterStem(&event.Buffer)
				terms = append(terms, string(event.Buffer))
			}
		}
		sort.Strings(terms)
		index.Add(file.Name(), terms)

		file.Close()
		return nil
	})
	// index.Search("GL_PIXEL_UNPACK_BUFFER")
	// index.Search("GL_PIXEL_PACK_BUFFER")
	// index.Search("GL_PIXEL_PACK_BUFFER_BINDING")
	// index.Search("GL_ELEMENT_ARRAY_BUFFER_BINDING")
	// index.Search("hyperbolic")
	// index.Search("arc")
	// index.Search("imageloaD")
	// index.Search("and")
	// index.Search("COLLIN")

	println("gogle; a small search engine")
	println("type \"exit\" to exit the program")
	print("> ")
	for scanner.Scan() {
		query = (scanner.Text())
		if query == "exit" {
			break;
		}
		index.Search(query);
		print("> ")
	}
}
