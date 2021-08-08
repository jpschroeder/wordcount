package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestStatistics(t *testing.T) {
	testWordCount(t, "", Stats{})
	testWordCount(t, "hello there", Stats{"hello": 1, "there": 1})
	testWordCount(t, "hello there hello there hello there", Stats{"hello": 3, "there": 3})
	testWordCount(t, "hello hello hello there there there", Stats{"hello": 3, "there": 3})
}

func TestWhiteSpace(t *testing.T) {
	testWordCount(t,
		"    hello\t\tthere    hello\n\nthere  hello \n\t there   \n  \t  \r\n",
		Stats{"hello": 3, "there": 3})
}

func TestCaseInsensitivity(t *testing.T) {
	testWordCount(t, "hello Hello HELLO heLLo", Stats{"hello": 4})
}

func TestPunctuation(t *testing.T) {
	testWordCount(t,
		"Hi, my name is john!  This; is my program... Is it cool?",
		Stats{"hi": 1, "my": 2, "name": 1, "is": 3, "john": 1, "this": 1, "program": 1, "it": 1, "cool": 1})
}

func TestWordCharacters(t *testing.T) {
	testWordCount(t, "Words with under_scores under_scores.", Stats{"words": 1, "with": 1, "under_scores": 2})
	testWordCount(t, "Words with number5 number5 1 2 thr33", Stats{"words": 1, "with": 1, "number5": 2, "1": 1, "2": 1, "thr33": 1})
	testWordCount(t, "John's book. John's book.", Stats{"john's": 2, "book": 2})
	testWordCount(t, "This is a run-of-the-mill high-school.", Stats{"this": 1, "is": 1, "a": 1, "run-of-the-mill": 1, "high-school": 1})
}

func TestKeywords(t *testing.T) {
	testWordCount(t, "hello -stats world", Stats{"hello": 1, "world": 1})
	testWordCount(t, "hello -blah world", Stats{"hello": 1, "world": 1, "-blah": 1})
	testWordCount(t, "-word1 word2-stats word3", Stats{"-word1": 1, "word2-stats": 1, "word3": 1})
	testWordCount(t, "hello -stats stats stats world", Stats{"hello": 1, "world": 1, "stats": 2})
}

func TestReset(t *testing.T) {
	testWordCount(t, "hello hello -reset hello hello hello", Stats{"hello": 3})
}

func TestStatString(t *testing.T) {
	s := Stats{"how": 10, "you": 5, "are": 5, "world": 50, "hello": 500}
	expected := `hello: 500
world:  50
how:    10
are:     5
you:     5
`
	if s.String() != expected {
		t.Errorf("\nExpected:\n%s\nActual:\n%s", expected, s.String())
	}
}

func testWordCount(t *testing.T, input string, expected Stats) {
	quiet = true
	actual, err := countWords(strings.NewReader(input))
	if err != nil {
		t.Errorf("\nInput: %s\nExpected: %v\nError: %v", input, expected, err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nInput: %s\nExpected: %v\nActual: %v", input, expected, actual)
	}
}
