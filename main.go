package main

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main2() {
	var terms = strings.Split("hello    world", " ")
	for count := 0; count < len(terms); count++{
		println(terms[count])
	}
}

func main() {
	if len(os.Args) < 2 {
		println("[USAGE]: gogle <root folder>")
		return
	}
	if folderStat, err := os.Stat(os.Args[1]); err != nil {
		println(err.Error())
		return
	} else if !folderStat.IsDir() {
		println("[ERROR] ", os.Args[1], " is not a folder")
		return
	}

	var (
		rootFolder    string = os.Args[1]
		file          *os.File
		index         = index.NewIndex()
		nextEventFunc func(*parser.ParserEvent)
		terms         = make([]string, 0)
		event         = parser.ParserEvent{
			Buffer: make([]rune, 256),
			Error:  nil,
		}
		query    string
		scanner  = bufio.NewScanner(os.Stdin)
		spinners = []string{"/", "\\"}
	)

	var count = 0
	print("Indexing: ", rootFolder, "  ")
	filepath.WalkDir(rootFolder, func(path string, d fs.DirEntry, err error) error {
		print("\b", spinners[count%2])
		count++
		if err != nil {
			println("[ERROR]: ", err.Error())
			return err
		}
		if d.IsDir() {
			return nil
		}

		if matches, _ := (filepath.Match("*html", filepath.Ext(path))); !matches {
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
	print("\b", " ")
	println("\ngogle; a small search engine")
	println("type \".exit\" to exit the program")
	println("[INFO]: the current Version of gogle can only search for a single term at a time")
	println("[INFO]: gogle does not have any typo tolerance yet")
	print("gogle> ")
	for scanner.Scan() {
		query = (strings.TrimSpace(scanner.Text()))
		if query == ".exit" {
			break
		}
		index.Search(query)
		print("> ")
	}
}


/*
	thoughts(collins):
	Maybe the index methods should not do any Normalization or stemming 
	maybe we should have a wrapper that converts the user's query into a searchable format (normalized, stemmed terms)
	eg
	Search(query string) 	
		-> get the first term, use it to find the postings in the index, 
		-> then get the second term, and use it on the already retrieved postings
*/
