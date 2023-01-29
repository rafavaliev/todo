package search

import (
	"golang.org/x/exp/slices"
	"time"
)

// Document is a searchable document
type Document struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

// Index is an inverted index of token -> list of document IDs which contain the token
type Index map[string][]string

// UserIndex contains an inverted index of a user's documents and analyzer used to tokenize documents
type UserIndex struct {
	UserID    uint      `json:"user_id"`
	Index     Index     `json:"index"`
	Analyzer  *Analyzer `json:"analyzer"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Insert adds a document to the user index
func (idx *UserIndex) Insert(document Document) {
	// Split document content into tokens
	if idx.Analyzer == nil {
		idx.Analyzer = NewEnglishAnalyzer()
	}
	tokens := idx.Analyzer.Analyze(document.Content)

	// Insert document ID to each token's
	for _, token := range tokens {
		ids, ok := idx.Index[token]
		if !ok {
			ids = []string{}
		}
		if slices.Contains(ids, document.ID) {
			continue
		}
		idx.Index[token] = append(ids, document.ID)
	}
}

// Search returns a list of document IDs which contain all tokens in the query
func (idx *UserIndex) Search(text string) []string {
	// Split query into tokens
	if idx.Analyzer == nil {
		idx.Analyzer = NewEnglishAnalyzer()
	}
	// analyzer search query
	tokens := idx.Analyzer.Analyze(text)

	// Get deduplicated document IDs for each token
	set := map[string]struct{}{}
	for _, token := range tokens {
		for _, id := range idx.Index[token] {
			set[id] = struct{}{}
		}
	}

	// Convert set to slice
	var ids []string
	for k := range set {
		ids = append(ids, k)
	}

	return ids
}

// Delete searches for a document occurrences in the index and removes it
func (idx *UserIndex) Delete(document Document) {
	if idx.Analyzer == nil {
		idx.Analyzer = NewEnglishAnalyzer()
	}
	// Split document content into tokens
	tokens := idx.Analyzer.Analyze(document.Content)

	// Search for all tokens that document contain
	for _, token := range tokens {
		ids, ok := idx.Index[token]
		if !ok {
			continue
		}
		// And remove document ID from token's list
		for i, id := range ids {
			if id == document.ID {
				ids = append(ids[:i], ids[i+1:]...)
				break
			}
		}

		if len(ids) == 0 {
			// if token has no more documents associated, remove it from index
			delete(idx.Index, token)
		} else {
			idx.Index[token] = ids
		}

	}
}
