package main

import (
	"regexp"
	"strings"
)

var aprostrophes = map[string]string{
	"let's": "us",
	"'m":    " am",
	"'re":   " are",
	"'ve":   " have",
	"'s":    " is",
	"'d":    " would",
	"'ll":   "will",
	"n't":   " not",
}

func sentanceToWords(sentence string) (words []string) {
	// lower case everything
	sentence = strings.ToLower(sentence)
	// remove all aprostrophes
	for k, v := range aprostrophes {
		sentence = strings.Replace(sentence, k, v, -1)
	}
	return regexp.MustCompile(`(\w+)`).FindAllString(sentence, -1)
}
