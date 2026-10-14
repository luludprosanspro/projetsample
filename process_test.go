package main

import (
	"os"
	"testing"
)

func TestProcessExamples(t *testing.T) {
	tests := []struct{ in, want string }{
		{"1E (hex) files were added", "30 files were added"},
		{"It has been 10 (bin) years", "It has been 2 years"},
		{"Ready, set, go (up) !", "Ready, set, GO!"},
		{"I should stop SHOUTING (low)", "I should stop shouting"},
		{"Welcome to the Brooklyn bridge (cap)", "Welcome to the Brooklyn Bridge"},
		{"This is so exciting (up, 2)", "This is SO EXCITING"},
		{"I was sitting over there ,and then BAMM !!", "I was sitting over there, and then BAMM!!"},
		{"I was thinking ... You were right", "I was thinking... You were right"},
		{"I am exactly how they describe me: ' awesome '", "I am exactly how they describe me: 'awesome'"},
		{"As Elton John said: ' I am the most well-known homosexual in the world '", "As Elton John said: 'I am the most well-known homosexual in the world'"},
		{"There it was. A amazing rock!", "There it was. An amazing rock!"},
	}
	for _, tc := range tests {
		if got := Process(tc.in); got != tc.want {
			t.Errorf("Process(%q)\n got: %q\nwant: %q", tc.in, got, tc.want)
		}
	}
}

func TestProcessFile(t *testing.T) {
	in, err := os.ReadFile("testdata/sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/expected.txt")
	if err != nil {
		t.Fatal(err)
	}
	if got := Process(string(in)); got != string(want) {
		t.Errorf("sortie différente\n got:\n%s\nwant:\n%s", got, want)
	}
}
