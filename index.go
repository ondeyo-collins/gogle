package gogle

type posting struct {
	documentID uint32
	freq       uint32
	positions  []uint32
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

	for termPosition := 0; termPosition < numberOfterms; termPosition++ {
		term = terms[termPosition]
		// get the posting list for the term,
		// if the term is not in the dictionary, create a new (term: posting list) entry
		_, ok = index.dictionary[term]
		if !ok {
			index.dictionary[term] = []posting{}
		}

		// get the term's posting for the document
		var postingsCount = len(index.dictionary[term])
		var currentPosting *posting
		for i := 0; i <= postingsCount; i++ {
			if i == postingsCount { // we reached the end,, the term's posting list does not include a posting for the current document
				index.dictionary[term] = append(index.dictionary[term], posting{documentID: newDocument.documentID, positions: make([]uint32, 0), freq: 0})
				currentPosting = &index.dictionary[term][postingsCount]
				break
			}
			if index.dictionary[term][i].documentID == newDocument.documentID {
				currentPosting = &index.dictionary[term][i]
				break
			}
		}

		currentPosting.freq++
		currentPosting.positions = append(currentPosting.positions, uint32(termPosition))
	}
}

func (index *index) Search(term string) {
	var termRunes = []rune(term)
	Normalize(&termRunes)
	PorterStem(&termRunes)
	_, ok := index.dictionary[string(termRunes)]
	if !ok {
		println("zero results for ", term)
		return
	}

	println(len(index.dictionary[string(termRunes)]), " results for ", term, ": ")
	var postingsCount = len(index.dictionary[string(termRunes)])
	for count := 0; count < postingsCount; count++ { // looping through postings
		var posting = index.dictionary[string(termRunes)][count]
		println(index.documents[posting.documentID].documentName, posting.freq)
	}
}

func NewIndex() *index {
	return &index{
		dictionary: map[string][]posting{},
		documents:  []document{},
	}
}
