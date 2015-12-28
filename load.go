package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var thesaurus map[string][]string
var syllables map[string][]int

func loadCmudict(path string) (map[string][]int, error) {

	m := make(map[string][]int, 140000)

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

		syl := m[w]
		syl = appendIfUnique(w, syl, c)
		m[w] = syl
	}

	if err := scan.Err(); err != nil {
		return nil, err
	}

	return m, nil
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

		for i, w := range words {
			if i > 1 {
				m[words[1]] = append(m[words[1]], w)
			}
		}
	}

	if err := scan.Err(); err != nil {
		return nil, err
	}

	return m, nil

}

func appendIfUnique(w string, l []int, n int) []int {

	if l == nil {
		return append(l, n)
	}

	// ugly, but len(l) == 0 almost always, and 1 or 2 very rarely
	for _, v := range l {
		if v == n {
			return l
		}
	}

	return append(l, n)
}

func main() {

	thesaurus, err := loadThesaurus("./resources/th_en_US_new.dat")
	if err != nil {
		panic(err)
	}
	fmt.Println(thesaurus["chocolate"])
	syllables, err := loadCmudict("./resources/cmudict.0.7a")
	if err != nil {
		panic(err)
	}
	fmt.Println(syllables["chocolate"])

}
