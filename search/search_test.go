package search

import (
	"reflect"
	"testing"
)

func TestUserIndex_Add(t *testing.T) {
	tests := []struct {
		name      string
		userIndex UserIndex
		inputDocs []Document
		wantIndex Index
	}{
		{
			name:      "add 1 document",
			userIndex: UserIndex{Index: Index{}},
			inputDocs: []Document{
				{ID: "1", Content: "Did I hear it right? Did the quick brown fox jump over the lazy dog?"},
			},
			wantIndex: Index{
				"did":   []string{"1"},
				"brown": []string{"1"},
				"dog":   []string{"1"},
				"fox":   []string{"1"},
				"hear":  []string{"1"},
				"it":    []string{"1"},
				"jump":  []string{"1"},
				"lazi":  []string{"1"},
				"over":  []string{"1"},
				"quick": []string{"1"},
				"right": []string{"1"},
			},
		},
		{

			name:      "add a few documents",
			userIndex: UserIndex{Index: Index{"some_random_existing_token": []string{"15"}}},
			inputDocs: []Document{
				{ID: "1", Content: "Did I hear it right? Did the quick brown fox jump over the lazy dog?"},
				{ID: "2", Content: "Did you hear that fox?"},
				{ID: "3", Content: "I heard something, I think it was a fox jumping over my dog!"},
			},
			wantIndex: Index{
				"did":                        []string{"1", "2"},
				"brown":                      []string{"1"},
				"dog":                        []string{"1", "3"},
				"fox":                        []string{"1", "2", "3"},
				"hear":                       []string{"1", "2"},
				"heard":                      []string{"3"},
				"it":                         []string{"1", "3"},
				"my":                         []string{"3"},
				"jump":                       []string{"1", "3"},
				"lazi":                       []string{"1"},
				"over":                       []string{"1", "3"},
				"quick":                      []string{"1"},
				"right":                      []string{"1"},
				"some_random_existing_token": []string{"15"},
				"someth":                     []string{"3"},
				"think":                      []string{"3"},
				"was":                        []string{"3"},
				"you":                        []string{"2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, document := range tt.inputDocs {
				tt.userIndex.Insert(document)
			}
			if !reflect.DeepEqual(tt.userIndex.Index, tt.wantIndex) {
				t.Errorf("UserIndex.Insert() = \n%v, wantIndex \n%v", tt.userIndex.Index, tt.wantIndex)
			}
		})
	}
}

func TestUserIndex_Remove(t *testing.T) {
	defaultIndex := func() Index {
		return Index{
			"did":                        []string{"1", "2"},
			"brown":                      []string{"1"},
			"dog":                        []string{"1", "3"},
			"fox":                        []string{"1", "2", "3"},
			"hear":                       []string{"1", "2"},
			"heard":                      []string{"3"},
			"it":                         []string{"1", "3"},
			"my":                         []string{"3"},
			"jump":                       []string{"1", "3"},
			"lazi":                       []string{"1"},
			"over":                       []string{"1", "3"},
			"quick":                      []string{"1"},
			"right":                      []string{"1"},
			"some_random_existing_token": []string{"15"},
			"someth":                     []string{"3"},
			"think":                      []string{"3"},
			"was":                        []string{"3"},
			"you":                        []string{"2"},
		}
	}
	tests := []struct {
		name      string
		userIndex UserIndex
		deleteDoc Document
		wantIndex Index
	}{
		{
			name:      "remove non-existing document of index",
			userIndex: UserIndex{Index: defaultIndex()},
			deleteDoc: Document{ID: "345", Content: "some_random_existing_token"},
			wantIndex: defaultIndex(),
		},
		{
			name:      "remove existing document of index",
			userIndex: UserIndex{Index: defaultIndex()},
			deleteDoc: Document{ID: "3", Content: "I heard something, I think it was a fox jumping over my dog!"},
			wantIndex: Index{
				"did":                        []string{"1", "2"},
				"brown":                      []string{"1"},
				"dog":                        []string{"1"},
				"fox":                        []string{"1", "2"},
				"hear":                       []string{"1", "2"},
				"it":                         []string{"1"},
				"jump":                       []string{"1"},
				"lazi":                       []string{"1"},
				"over":                       []string{"1"},
				"quick":                      []string{"1"},
				"right":                      []string{"1"},
				"you":                        []string{"2"},
				"some_random_existing_token": []string{"15"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.userIndex.Delete(tt.deleteDoc)
			if !reflect.DeepEqual(tt.userIndex.Index, tt.wantIndex) {
				t.Errorf("UserIndex.Delete() = \n%v\n, wantIndex \n%v", tt.userIndex.Index, tt.wantIndex)
			}
		})
	}
}
