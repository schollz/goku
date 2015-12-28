package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

var thesaurus map[string][]string
var cmudict map[string][]int
var skipwords []string
var err error

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

func getSynonyms(w string) (possibilities []string) {
	// returns a list of synonyms that have a unique number of syllables
	possibilities = append(possibilities, w)

	for _, skipword := range skipwords {
		if w == skipword {
			return
		}
	}

	syllablesAccountedFor := make(map[int]string)
	for _, synonym := range thesaurus[w] {
		sylbls := cmudict[synonym][0]
		if sylbls > 0 {
			if _, ok := syllablesAccountedFor[sylbls]; ok {
				// pass
			} else {
				syllablesAccountedFor[sylbls] = synonym
			}
		}
	}

	for _, word := range syllablesAccountedFor {
		possibilities = append(possibilities, word)
	}
	return
}

func formatPoem(poem []string, syllables []int) (string, error) {

	l := 0
	c := syllables[l]

	h := &bytes.Buffer{}

	processed := 0

	for _, p := range poem {
		processed++
		word := p

		word = strings.ToUpper(word)

		word = strings.Map(func(r rune) rune {
			if 'A' <= r && r <= 'Z' {
				return r
			}
			return -1
		}, word)

		a := cmudict[word]
		if a == nil || len(a) == 0 {

			prefixes := []struct {
				prefix    string
				syllables int
			}{
				{"ANTI", 2},
				{"DE", 1},
				{"DIS", 1},
				{"EX", 1},
				{"MEGA", 2},
				{"MINI", 2},
				{"MIS", 1},
				{"MULTI", 2},
				{"NON", 1},
				{"POST", 1},
				{"PRE", 1},
				{"PRO", 1},
				{"PROTO", 2},
				{"QUASI", 2},
				{"RE", 1},
				{"SEMI", 2},
				{"UN", 1},
				{"VICE", 1},
			}

			for _, p := range prefixes {
				if strings.HasPrefix(word, p.prefix) {
					w := strings.TrimPrefix(word, p.prefix)
					a = cmudict[w]
					if a != nil && len(a) != 0 {
						a[0] += p.syllables
						break
					}
				}
			}

			if a == nil || len(a) == 0 {
				return "", errors.New("unknown word: " + word)
			}
		}

		if len(a) > 1 {
			return "", errors.New("don't yet handle words with multiple syllable counts")
		}

		c -= a[0]

		if c < 0 {
			break
		}

		fmt.Fprint(h, p)

		if c > 0 {
			fmt.Fprint(h, " ")
		} else {
			fmt.Fprint(h, "\n")
			l++
			if l >= len(syllables) {
				break
			}
			c = syllables[l]
		}
	}

	if processed != len(poem) || c != 0 {
		return "", errors.New("not a haiku")
	}

	return h.String(), nil
}

func init() {
	// initialize the thesaurs and the syllable dictionary

	thesaurus, err = loadThesaurus("./resources/th_en_US_new.dat")
	if err != nil {
		panic(err)
	}
	fmt.Println(thesaurus["chocolate"])
	cmudict, err = loadCmudict("./resources/cmudict.0.7a")
	if err != nil {
		panic(err)
	}
	fmt.Println(cmudict["chocolate"])

	skipwords = []string{"a", "about", "above", "above", "across", "after", "afterwards", "again", "against", "all", "almost", "alone", "along", "already", "also", "although", "always", "am", "among", "amongst", "amoungst", "amount", "an", "and", "another", "any", "anyhow", "anyone", "anything", "anyway", "anywhere", "are", "around", "as", "at", "back", "be", "became", "because", "become", "becomes", "becoming", "been", "before", "beforehand", "behind", "being", "below", "beside", "besides", "between", "beyond", "bill", "both", "bottom", "but", "by", "call", "can", "cannot", "cant", "co", "con", "could", "couldnt", "cry", "de", "describe", "detail", "do", "done", "down", "due", "during", "each", "eg", "eight", "either", "eleven", "else", "elsewhere", "empty", "enough", "etc", "even", "ever", "every", "everyone", "everything", "everywhere", "except", "few", "fifteen", "fify", "fill", "find", "fire", "first", "five", "for", "former", "formerly", "forty", "found", "four", "from", "front", "full", "further", "get", "give", "go", "had", "has", "hasnt", "have", "he", "hence", "her", "here", "hereafter", "hereby", "herein", "hereupon", "hers", "herself", "him", "himself", "his", "how", "however", "hundred", "ie", "if", "in", "inc", "indeed", "interest", "into", "is", "it", "its", "itself", "keep", "last", "latter", "latterly", "least", "less", "ltd", "made", "many", "may", "me", "meanwhile", "might", "mill", "mine", "more", "moreover", "most", "mostly", "move", "much", "must", "my", "myself", "name", "namely", "neither", "never", "nevertheless", "next", "nine", "no", "nobody", "none", "noone", "nor", "not", "nothing", "now", "nowhere", "of", "off", "often", "on", "once", "one", "only", "onto", "or", "other", "others", "otherwise", "our", "ours", "ourselves", "out", "over", "own", "part", "per", "perhaps", "please", "put", "rather", "re", "same", "see", "seem", "seemed", "seeming", "seems", "serious", "several", "she", "should", "show", "side", "since", "sincere", "six", "sixty", "so", "some", "somehow", "someone", "something", "sometime", "sometimes", "somewhere", "still", "such", "system", "take", "ten", "than", "that", "the", "their", "them", "themselves", "then", "thence", "there", "thereafter", "thereby", "therefore", "therein", "thereupon", "these", "they", "thickv", "thin", "third", "this", "those", "though", "three", "through", "throughout", "thru", "thus", "to", "together", "too", "top", "toward", "towards", "twelve", "twenty", "two", "un", "under", "until", "up", "upon", "us", "very", "via", "was", "we", "well", "were", "what", "whatever", "when", "whence", "whenever", "where", "whereafter", "whereas", "whereby", "wherein", "whereupon", "wherever", "whether", "which", "while", "whither", "who", "whoever", "whole", "whom", "whose", "why", "will", "with", "within", "without", "would", "yet", "you", "your", "yours", "yourself", "yourselves", "the"}

}

func main() {

	fmt.Println(getSynonyms("the"))

	scan := `In my old home
which I forsook, the cherries
are in bloom.`

	var poem []string

	var sentenceFinished bool

	for _, word := range strings.Split(scan, " ") {
		t := word

		poem = append(poem, t)
		lastRune := t[len(t)-1]

		if lastRune == '.' || lastRune == '?' || lastRune == '!' {
			sentenceFinished = true
		}

		if sentenceFinished {

			switch t {
			case "Mr.":
				fallthrough
			case "Dr.":
				fallthrough
			case "Ms.":
				fallthrough
			case "Mrs.":
				fallthrough
			case "Sr.":
				sentenceFinished = false
			}

			if sentenceFinished {

				syllables := []int{5, 7, 5}

				p, err := formatPoem(poem, syllables)
				if err == nil {
					fmt.Println(p)
				}
				sentenceFinished = false
				poem = poem[:0]
			}
		}
	}

	fmt.Println(poem)

}
