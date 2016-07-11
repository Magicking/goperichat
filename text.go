package main

import (
	"log"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func isMn(r rune) bool {
	    return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func TextReduce(word string) ([]string, int) {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

	ret, _, err := transform.String(t, word)
	if err != nil {
		log.Printf("Failed transform for %s: %s\n", word, err)
	}
	words := strings.Split(ret, " ")
	word_count := 0
	for i, word := range words {
		not_word := utf8.RuneCountInString(word) < 3 || strings.HasPrefix(word, "@")
		if not_word {
			words[i] = ""
			continue
		}
		word_count++
	}
	return words, word_count
}
