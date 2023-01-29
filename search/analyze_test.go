package search

import (
	"reflect"
	"testing"
)

func Test_tokenize(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string
	}{
		{
			name: "nothin to split",
			text: "hello",
			want: []string{"hello"},
		},
		{
			name: "split on space",
			text: "hello world",
			want: []string{"hello", "world"},
		},
		{
			name: "split on punctuation",
			text: "hello, world!",
			want: []string{"hello", "world"},
		},
		{
			name: "split on punctuation and space",
			text: "Did I hear it right? Did the quick brown fox jump over the lazy dog?",
			want: []string{"Did", "I", "hear", "it", "right", "Did", "the", "quick", "brown", "fox", "jump", "over", "the", "lazy", "dog"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tokenize(tt.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenize() = %v, wantIndex %v", got, tt.want)
			}
		})
	}
}

func Test_stopwordFilter(t *testing.T) {

	tests := []struct {
		name   string
		tokens []string
		want   []string
	}{
		{
			name:   "no stopWords",
			tokens: []string{"hello", "world"},
			want:   []string{"hello", "world"},
		},
		{
			name:   "stopWords",
			tokens: []string{"hello", "world", "the", "and", "a"},
			want:   []string{"hello", "world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := englishStopWordFilter(tt.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("englishStopWordFilter() = %v, wantIndex %v", got, tt.want)
			}
		})
	}
}

func TestAnalyzer_Analyze(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		want     []string
		Analyzer *Analyzer
	}{
		{
			name:     "default analyzer",
			Analyzer: &Analyzer{},
			text:     "Did I hear it right? Did the quick brown fox jump over the lazy dog?",
			want:     []string{"Did", "I", "hear", "it", "right", "Did", "the", "quick", "brown", "fox", "jump", "over", "the", "lazy", "dog"},
		},
		{
			name:     "english analyzer",
			Analyzer: NewEnglishAnalyzer(),
			text:     "Did I hear it right? Did the quick brown fox jump over the lazy dog?",
			want:     []string{"did", "hear", "it", "right", "did", "quick", "brown", "fox", "jump", "over", "lazi", "dog"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.Analyzer.Analyze(tt.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Analyzer.Analyze() = \n%v, wantIndex \n%v", got, tt.want)
			}
		})
	}
}
