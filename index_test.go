package gogle

import "testing"
import "slices"

type expectedPosition struct {
	documentName string
	positions    []uint32
}

func TestIndexAdd(t *testing.T) {
	var input = []struct {
		documentName string
		terms        []string
	}{
		{documentName: "one.html", terms: []string{"hello", "hello", "my", "nigga"}},
		{documentName: "two.html", terms: []string{"hello", "nigga", "my"}},
	}

	var index = NewIndex()
	index.Add(input[0].documentName, input[0].terms)
	index.Add(input[1].documentName, input[1].terms)

	var read = func(term string, documentName string, expectedPositions []uint32) {
		var matchingPosting posting
		if postings, ok := index.dictionary[term]; !ok {
			t.Fatalf("term not in index dictionary: (%s)\n", term)
		} else {
			for i := 0; i < len(postings); i++ {
				if index.documents[postings[i].documentID].documentName == documentName {
					matchingPosting = postings[i]
					goto endOfElse
				}
			}
			t.Fatalf("no \"%s\" postings for the document \"%s\"", term, documentName)
		endOfElse:
		}

		// confirm the expectedPositions
		for i := 0; i < len(expectedPositions); i++ {
			if !slices.Contains(matchingPosting.positions, expectedPositions[i]) {
				t.Fatalf("the term \"%s\" does not appear at position %v in %v", term, expectedPositions[i], documentName)
			}
		}
	}

	read("hello", "one.html", []uint32{0, 1})
	read("my", "one.html", []uint32{2})
	read("my", "one.html", []uint32{2})
	read("hello", "two.html", []uint32{0})
}
