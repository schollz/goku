package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
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

type haikuCluster struct {
	haikus    []string
	isHaikus  []bool
	numHaikus int
}

func getHaikus(words []string) (haikus []string, isHaikus []bool, numHaikus int) {
	numHaikus = 0
	for {
		haikuString, i, gotHaiku := isHaiku(words)
		haikus = append(haikus, haikuString)
		isHaikus = append(isHaikus, gotHaiku)
		if gotHaiku == false {
			return
		} else {
			numHaikus = numHaikus + 1
		}

		words = words[i+1:]
	}
}

func isHaiku(words []string) (string, int, bool) {
	checks := []int{5, 7, 5}
	curCheck := 0
	runningTotal := 0
	currentHaiku := ""
	for i, word := range words {
		slbles := cmudict[word]
		runningTotal = runningTotal + slbles
		currentHaiku = currentHaiku + word + " "
		if runningTotal == checks[curCheck] {
			curCheck += 1
			if curCheck == 3 {
				return currentHaiku, i, true
			} else {
				currentHaiku = currentHaiku + "\n"

			}
			runningTotal = 0
		}
	}
	return currentHaiku, -1, false
}

func randChoices(limits []int) (choices []int) {
	for _, l := range limits {
		choices = append(choices, rand.Intn(l))
	}
	return
}

func listAlternates(input []string) (output [50][]string) {
	totals := make([]int, len(input))
	for i, w := range input {
		totals[i] = len(getSynonyms(w))
	}
	for i := 0; i < len(output); i++ {
		choices := randChoices(totals)
		for j, w := range input {
			output[i] = append(output[i], getSynonyms(w)[choices[j]])
		}
	}
	return
}

func init() {
	// initialize the thesaurs and the syllable dictionary

	thesaurus, err = loadThesaurus("./resources/th_en_US_new.dat")
	if err != nil {
		panic(err)
	}
	cmudict, err = loadCmudict("./resources/cmudict.0.7a")
	if err != nil {
		panic(err)
	}
	skipwords = []string{"a", "about", "above", "above", "across", "after", "afterwards", "again", "against", "all", "almost", "alone", "along", "already", "also", "although", "always", "am", "among", "amongst", "amoungst", "amount", "an", "and", "another", "any", "anyhow", "anyone", "anything", "anyway", "anywhere", "are", "around", "as", "at", "back", "be", "became", "because", "become", "becomes", "becoming", "been", "before", "beforehand", "behind", "being", "below", "beside", "besides", "between", "beyond", "bill", "both", "bottom", "but", "by", "call", "can", "cannot", "cant", "co", "con", "could", "couldnt", "cry", "de", "describe", "detail", "do", "done", "down", "due", "during", "each", "eg", "eight", "either", "eleven", "else", "elsewhere", "empty", "enough", "etc", "even", "ever", "every", "everyone", "everything", "everywhere", "except", "few", "fifteen", "fify", "fill", "find", "fire", "first", "five", "for", "former", "formerly", "forty", "found", "four", "from", "front", "full", "further", "get", "give", "go", "had", "has", "hasnt", "have", "he", "hence", "her", "here", "hereafter", "hereby", "herein", "hereupon", "hers", "herself", "him", "himself", "his", "how", "however", "hundred", "ie", "if", "in", "inc", "indeed", "interest", "into", "is", "it", "its", "itself", "keep", "last", "latter", "latterly", "least", "less", "ltd", "made", "many", "may", "me", "meanwhile", "might", "mill", "mine", "more", "moreover", "most", "mostly", "move", "much", "must", "my", "myself", "name", "namely", "neither", "never", "nevertheless", "next", "nine", "no", "nobody", "none", "noone", "nor", "not", "nothing", "now", "nowhere", "of", "off", "often", "on", "once", "one", "only", "onto", "or", "other", "others", "otherwise", "our", "ours", "ourselves", "out", "over", "own", "part", "per", "perhaps", "please", "put", "rather", "re", "same", "see", "seem", "seemed", "seeming", "seems", "serious", "several", "she", "should", "show", "side", "since", "sincere", "six", "sixty", "so", "some", "somehow", "someone", "something", "sometime", "sometimes", "somewhere", "still", "such", "system", "take", "ten", "than", "that", "the", "their", "them", "themselves", "then", "thence", "there", "thereafter", "thereby", "therefore", "therein", "thereupon", "these", "they", "thickv", "thin", "third", "this", "those", "though", "three", "through", "throughout", "thru", "thus", "to", "together", "too", "top", "toward", "towards", "twelve", "twenty", "two", "un", "under", "until", "up", "upon", "us", "very", "via", "was", "we", "well", "were", "what", "whatever", "when", "whence", "whenever", "where", "whereafter", "whereas", "whereby", "wherein", "whereupon", "wherever", "whether", "which", "while", "whither", "who", "whoever", "whole", "whom", "whose", "why", "will", "with", "within", "without", "would", "yet", "you", "your", "yours", "yourself", "yourselves", "the"}

}

func main() {
	// words := strings.Split(`want to play a game with seventeen syllables we write some poem want to play a game with seventeen syllables we write some poem something else`, " ")
	// fmt.Println(getHaikus(words))
	//input := []string{"cat", "is", "nice"}
	sentence := `This is the sixth time we have had the pleasure of writing a birthday blog post for Go, and we would not be doing so if not for the wonderful and passionate people in our community. The Go team would like to thank everyone who has contributed code, written an open source library, authored a blog post, helped a new gopher, or just given Go a try. `
	fmt.Println("ORIGINAL:")
	fmt.Println(sentence)
	sentence = strings.ToLower(sentence)
	sentence = strings.Replace(sentence, "don't", "do not", -1)
	sentence = strings.Replace(sentence, "'", "", -1)
	words := regexp.MustCompile(`(\w+)`).FindAllString(sentence, -1)
	alternatives := listAlternates(words)
	bestNum := 0
	var bestHaikus []string
	var bestIsHaikus []bool

	for i := 0; i < len(alternatives); i++ {
		for j := 0; j < 3; j++ {
			haikuString, isHaikus, numHaikus := getHaikus(alternatives[i][j:])
			if numHaikus > bestNum {
				if j > 0 {
					haikuString = append([]string{strings.Join(alternatives[i][:j], " ")}, haikuString...)
					isHaikus = append([]bool{false}, isHaikus...)
				}
				bestHaikus = haikuString
				bestIsHaikus = isHaikus
			}
		}
	}
	fmt.Println("\n\nBEST HAIKU:")
	for i, bestHaiku := range bestHaikus {
		if bestIsHaikus[i] == true {
			fmt.Printf("\n")
		}
		fmt.Println(bestHaiku)
		if bestIsHaikus[i] == true {
			fmt.Printf("\n")
		}
	}

	// Todo: replace each word with the puncuation near the word in the original

}
