package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var thesaurus map[string][]string

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

func main() {

	thesaurus, err := loadThesaurus("./resources/th_en_US_new.dat")
	if err != nil {
		panic(err)
	}
	fmt.Println(thesaurus["chocolate"])

}
