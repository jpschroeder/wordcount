package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

func printHelp() {
	fmt.Println("stats: print statistics")
	fmt.Println("help: print help")
	fmt.Println("exit: quit")
}

type stats map[string]int

func read(rd io.Reader) (stats, error) {
	reader := bufio.NewReader(rd)
	stats := make(stats)

	// read until 'exit' or io.EOF
	for {
		word, err := readWord(reader)
		if len(word) > 0 {
			switch word {
			case "stats":
				fmt.Println(stats)
			case "help":
				printHelp()
			case "exit":
				return stats, nil
			default:
				if _, exists := stats[word]; !exists {
					stats[word] = 0
				}
				stats[word]++
			}
		}
		if err == io.EOF {
			return stats, nil
		} else if err != nil {
			return nil, err
		}
	}
}

func readWord(r *bufio.Reader) (string, error) {
	var sb strings.Builder
	ch, _, err := r.ReadRune()

	// Skip non-word characters
	for err == nil && !isWordChar(ch) {
		ch, _, err = r.ReadRune()
	}

	// Read while we have word characters
	for err == nil && isWordChar(ch) {
		sb.WriteRune(unicode.ToLower(ch))
		ch, _, err = r.ReadRune()
	}

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
	// For example: "john's book" or "run-of-the-mill"
	return isRegexWordChar || ch == '-' || ch == '\''
}

func main() {
	/* todo
	DONE: unit tests
	DONE: check ispunct

	optional ignore case
	code cleanup
	prompt
	handle redirected input (no prompt) - just print stats
	better keywords :exit or \stats
	maybe don't need bufio, if no UnreadRune()
	reset stats?
	usage documentation
	*/
	printHelp()
	read(os.Stdin)
}
