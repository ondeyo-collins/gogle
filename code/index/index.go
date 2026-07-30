package index

import (
	"github.com/collins994/gogle/code/text"
)

type posting struct {
	documentID uint32
	freq       uint32
}

type document struct {
	documentID   uint32
	documentName string
}
type index struct {
	dictionary map[string][]posting
	documents  []document
}

func (index *index) Add(documentName string, terms []string) {
	var (
		numberOfterms = len(terms)
		term          string
		ok            bool
	)

	// add the document to the documents list
	var newDocument = document{
		documentID:   uint32(len(index.documents)),
		documentName: documentName,
	}
	index.documents = append(index.documents, newDocument)

	for count := 0; count < numberOfterms; count++ {
		term = terms[count]
		// get the posting list for the term,
		// if the term is not in the dictionary, create a new (term: posting list) entry
		_, ok = index.dictionary[term]
		if !ok {
			index.dictionary[term] = []posting{}
		}

		var postingsCount = len(index.dictionary[term])
		for count := 0; count < postingsCount; count++ {
			if index.dictionary[term][count].documentID == newDocument.documentID {
				index.dictionary[term][count].freq++
				goto endOfLoop
			}
		}

		// the current document does not have a posting yet for this term
		// we create one
		index.dictionary[term] = append(index.dictionary[term], posting{
			documentID: newDocument.documentID,
			freq:       1,
		})

	endOfLoop:
	}
}

func (index *index) Search(term string) {
	var termRunes = []rune(term)
	text.Normalize(&termRunes)
	text.PorterStem(&termRunes)
	// println(string(termRunes))
	_, ok := index.dictionary[string(termRunes)]
	if !ok {
		println("no results for ",term)
		return
	}

	println(len(index.dictionary[string(termRunes)]), " results for ",term, ": ");
	var postingsCount = len(index.dictionary[string(termRunes)])
	for count := 0; count < postingsCount; count++ { // looping through postings
		var posting = index.dictionary[string(termRunes)][count]
		println(index.documents[posting.documentID].documentName, posting.freq);
	}
}

func NewIndex() *index {
	return &index{
		dictionary: map[string][]posting{},
		documents:  []document{},
	}
}
