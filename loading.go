package main

import (
	"bufio"
	"os"
	"strings"
)

var (
	thesaurus map[string][]string
	cmudict   map[string]int
)

func loadCmudict(path string) (map[string]int, error) {

	m := make(map[string]int, 140000)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scan := bufio.NewScanner(f)

	for scan.Scan() {
		s := scan.Text()
		s = strings.ToLower(s)
		if s[0] == ';' {
			// skip comments
			continue
		}

		// find first word
		idx := strings.Index(s, " ")
		w := s[0:idx]

		if w[idx-1] == ')' {
			w = w[:idx-3]
		}

		c := 0
		// count syllables == digits in remaining string
		for _, r := range s[idx:] {
			if r >= '0' && r <= '9' {
				c++
			}
		}

		m[w] = c
	}

	if err := scan.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func loadThesaurus(path string) (map[string][]string, error) {

	m := make(map[string][]string)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scan := bufio.NewScanner(f)

	for scan.Scan() {
		s := scan.Text()
		words := strings.Split(s, "|")
		if len(words) <= 2 {
			// skip if doesn't have syllables
			continue
		}
		for j := range words {
			for i, w := range words {
				if j > 0 && i > 0 && i != j {
					if stringInSlice(w, m[words[j]]) == false {
						m[words[j]] = append(m[words[j]], w)
					}
				}
			}
		}
	}

	return m, scan.Err()
}
