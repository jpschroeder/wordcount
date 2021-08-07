package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestStatistics(t *testing.T) {
	testRead(t, "", stats{})
	testRead(t, "hello there", stats{"hello": 1, "there": 1})
	testRead(t, "hello there hello there hello there", stats{"hello": 3, "there": 3})
	testRead(t, "hello hello hello there there there", stats{"hello": 3, "there": 3})
}

func TestWhiteSpace(t *testing.T) {
	testRead(t,
		"    hello\t\tthere    hello\n\nthere  hello \n\t there   \n  \t  \r\n",
		stats{"hello": 3, "there": 3})
}

func TestCaseInsensitivity(t *testing.T) {
	testRead(t, "hello Hello HELLO heLLo", stats{"hello": 4})
}

func TestPunctuation(t *testing.T) {
	testRead(t,
		"Hi, my name is john!  This; is my program... Is it cool?",
		stats{"hi": 1, "my": 2, "name": 1, "is": 3, "john": 1, "this": 1, "program": 1, "it": 1, "cool": 1})
}

func TestWordCharacters(t *testing.T) {
	testRead(t, "Words with under_scores under_scores.", stats{"words": 1, "with": 1, "under_scores": 2})
	testRead(t, "Words with number5 number5 1 2 thr33", stats{"words": 1, "with": 1, "number5": 2, "1": 1, "2": 1, "thr33": 1})
	testRead(t, "John's book. John's book.", stats{"john's": 2, "book": 2})
	testRead(t, "This is a run-of-the-mill high-school.", stats{"this": 1, "is": 1, "a": 1, "run-of-the-mill": 1, "high-school": 1})
}

func TestKeywords(t *testing.T) {
	testRead(t, "hello \\stats world", stats{"hello": 1, "world": 1})
	testRead(t, "hello \\blah world", stats{"hello": 1, "world": 1, "blah": 1})
	testRead(t, "\\word1\\word2\\help\\word3", stats{"word1": 1, "word2": 1, "word3": 1})
	testRead(t, "hello \\stats stats stats world", stats{"hello": 1, "world": 1, "stats": 2})
}

func testRead(t *testing.T, input string, expected map[string]int) {
	actual, err := readString(input)
	if err != nil {
		t.Errorf("\nInput: %s\nExpected: %v\nError: %v", input, expected, err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nInput: %s\nExpected: %v\nActual: %v", input, expected, actual)
	}
}

func readString(input string) (map[string]int, error) {
	return read(strings.NewReader(input))
}
