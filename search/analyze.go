package search

import (
	snowballeng "github.com/kljensen/snowball/english"
	"strings"
	"unicode"
)

// Analyzer receives a text and returns a slice of tokens(words) that's used for searching.
// Input text can be in any language, so the client should choose/provide an appropriate filters.
type Analyzer struct {
	Name    string
	Filters []Filter
}

func (a *Analyzer) Analyze(text string) []string {
	tokens := tokenize(text)
	for _, filter := range a.Filters {
		tokens = filter(tokens)
	}
	return tokens
}

// MarshalJSON for simplicity saves only the name of the analyzer.
func (a *Analyzer) MarshalJSON() ([]byte, error) {
	if a.Name == "" {
		a.Name = "default_analyzer"
	}
	return []byte(a.Name), nil
}

// UnmarshalJSON for simplicity rebuilds analyzer from a name.
func (a *Analyzer) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case "english_analyzer":
		a = NewEnglishAnalyzer()
	default:
		a = &Analyzer{}
	}
	return nil
}

// Filter is used in Analyzer to filter tokens(words).
type Filter func(tokens []string) []string

// NewEnglishAnalyzer creates an english analyzer to analyze English text.
func NewEnglishAnalyzer() *Analyzer {
	return &Analyzer{
		Name: "english_analyzer",
		Filters: []Filter{
			lowercaseFilter,
			englishStopWordFilter,
			englishStemmerFilter,
		},
	}
}

// tokenize splits a string into a slice of tokens(words)
func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// Split on any character that is not a letter or a number.
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

var stopWords = map[string]struct{}{
	"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
	"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
}

// englishStopWordFilter removes stopWords from a slice of tokens.
func englishStopWordFilter(tokens []string) []string {
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := stopWords[token]; !ok {
			r = append(r, token)
		}
	}
	return r
}

// englishStemmerFilter stems english words. Ex: dogs -> dog, jumping -> jump
func englishStemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false)
	}
	return r
}

func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}
