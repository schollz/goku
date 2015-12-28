package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var thesaurus map[string][]string
var cmudict map[string]int
var skipwords []string
var err error

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
		sylbls := cmudict[synonym]
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

func getHaikus(words []string) (haikus []string) {
	gotHaiku := true
	i := 0
	for {
		i, gotHaiku = isHaiku(words)
		if gotHaiku == false {
			return
		}
		fmt.Println(i)
		fmt.Println(gotHaiku)
		haikus = append(haikus, strings.Join(words[:i], " "))
		words = words[i:]
	}
}

func isHaiku(words []string) (int, bool) {
	checks := []int{5, 7, 5}
	curCheck := 0
	runningTotal := 0

	for i, word := range words {
		slbles := cmudict[word]
		runningTotal = runningTotal + slbles
		if runningTotal == checks[curCheck] {
			curCheck += 1
			if curCheck == 3 {
				return i + 1, true
			}
			runningTotal = 0
		}
		if runningTotal > checks[curCheck] {
			return -1, false
		}
	}
	return -1, false
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

	words := strings.Split(`want to play a game with seventeen syllables we write some poem want to play a game with seventeen syllables we write some poem something else`, " ")
	fmt.Println(getHaikus(words))

}
