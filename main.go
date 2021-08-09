package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//==============================================================================

func countWords(rd io.Reader) (*Stats, error) {
	reader := bufio.NewReader(rd)
	stats := NewStats()

	// read until 'exit' or io.EOF
	for {
		word, err := readWord(reader)
		if len(word) > 0 {
			switch word {
			case "-stats":
				fmt.Println(stats)
			case "-help":
				help()
			case "-reset":
				stats = NewStats()
			case "-exit":
				return stats, nil
			default:
				stats.Add(word)
			}
		}
		if err == io.EOF {
			return stats, nil
		} else if err != nil {
			return nil, err
		}
	}
}

func readWord(reader *bufio.Reader) (string, error) {
	var sb strings.Builder
	ch, _, err := reader.ReadRune()

	// Skip non-word characters
	for err == nil && !isWordChar(ch) {
		if ch == '\n' {
			prompt()
		}
		ch, _, err = reader.ReadRune()
	}

	// Read while we have word characters
	for err == nil && isWordChar(ch) {
		sb.WriteRune(unicode.ToLower(ch))
		ch, _, err = reader.ReadRune()
	}

	reader.UnreadRune()
	return sb.String(), err
}

func isWordChar(ch rune) bool {
	// https://docs.microsoft.com/en-us/dotnet/standard/base-types/character-classes-in-regular-expressions?redirectedfrom=MSDN#WordCharacter
	isRegexWordChar := unicode.In(ch,
		unicode.Ll, // Letter, Lowercase
		unicode.Lu, // Letter, Uppercase
		unicode.Lt, // Letter, Titlecase
		unicode.Lo, // Letter, Other
		unicode.Lm, // Letter, Modifier
		unicode.Mn, // Mark, Nonspacing
		unicode.Nd, // Number, Decimal Digit
		unicode.Pc) // Punctuation, Connector
	// Also including appostrophes and hyphens.
	// For example: "john's book" or "high-school"
	return isRegexWordChar || ch == '-' || ch == '\''
}

//==============================================================================

type Stats struct {
	entries []timeEntry
	counts  map[string]int
}

type timeEntry struct {
	word  string
	added time.Time
}

func NewStats() *Stats {
	return &Stats{
		entries: make([]timeEntry, 0),
		counts:  make(map[string]int),
	}
}

// Add a single count of 'word' to the statistics
func (s *Stats) Add(word string) {
	s.AddTime(word, time.Now())
}

func (s *Stats) AddTime(word string, added time.Time) {
	s.entries = append(s.entries, timeEntry{word, added})

	if _, exists := s.counts[word]; !exists {
		s.counts[word] = 0
	}
	s.counts[word]++
}

// Remove entries older than a minute
func (s *Stats) RemoveOldEntries() {
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)
	var i int
	for i = 0; i < len(s.entries); i++ {
		if s.entries[i].added.After(oneMinuteAgo) {
			break
		}
		s.counts[s.entries[i].word]--
		if s.counts[s.entries[i].word] == 0 {
			delete(s.counts, s.entries[i].word)
		}
	}
	s.entries = s.entries[i:]
}

func (s *Stats) Map() map[string]int {
	s.RemoveOldEntries()
	return s.counts
}

func (s *Stats) String() string {
	return PrettyPrint(s.Map())
}

// Pretty print the word counts, sorted and padded. Example:
// hello: 500
// world:  50
// how:    10
// are:     5
// you:     5
func PrettyPrint(s map[string]int) string {
	if len(s) == 0 {
		return "-- no entries --"
	}

	// extract entries for sorting
	type countEntry struct {
		word  string
		count int
	}
	words := make([]countEntry, len(s))

	// determine max length for padding
	var maxWord, maxDigits, i int

	for word, count := range s {
		words[i] = countEntry{word, count}
		if len(word) > maxWord {
			maxWord = len(word)
		}
		digits := len(strconv.Itoa(count))
		if digits > maxDigits {
			maxDigits = digits
		}
		i++
	}

	sort.Slice(words, func(i, j int) bool {
		// sort by count descending
		if words[i].count != words[j].count {
			return words[i].count > words[j].count
		}
		// then by word ascending
		return words[i].word < words[j].word
	})

	var sb strings.Builder
	for _, entry := range words {
		padding := maxDigits + maxWord - len(entry.word)
		line := fmt.Sprintf("%s: %*d\n", entry.word, padding, entry.count)
		sb.WriteString(line)
	}
	ret := sb.String()
	return ret[:len(ret)-1] // trim the last newline
}

//==============================================================================

// Normally I would not use globals, but this case seemed simple enough
// If things got any more complicated, I would refactor
var quiet = false

func help() {
	if quiet {
		return
	}
	fmt.Printf(`Enter text to capture word counts.
Words will be removed from the count after 1 minute.
The following keywords will not be counted (beginning with '-')
-stats: print statistics
-reset: reset statistics
-help: print help
-exit: quit
`)
}

func prompt() {
	if quiet {
		return
	}
	fmt.Printf("> ")
}

func usage() {
	fmt.Printf(`wordcount - print the number of unique words in a stream

Usage: wordcount [file]

Example:

interactive mode: wordcount
read from file:   wordcount test.txt
read from pipe:   echo lorem ipsum lorem | wordcount

Output:

hello: 500
world:  50
how:    10
are:     5
you:     5
`)
	os.Exit(1)
}

func isInputRedirected() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}

func main() {
	if len(os.Args) > 2 {
		usage()
	}
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		usage()
	}

	input := os.Stdin
	quiet = isInputRedirected()

	if len(os.Args) == 2 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		input = file
		quiet = true
	}

	help()
	prompt()
	s, err := countWords(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}
