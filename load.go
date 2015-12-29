package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

type haikuCluster struct {
	haikus    []string
	isHaikus  []bool
	numHaikus int
}

type node struct {
	haikus  []string
	syns    map[int][]string
	numSyns []int
	start   int
	end     int
}

var thesaurus map[string][]string
var cmudict map[string]int
var skipwords []string
var err error
var nodes []node
var sentenceWords []string

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
		for j, _ := range words {
			for i, w := range words {
				if j > 0 && i > 0 && i != j {
					if stringInSlice(w, m[words[j]]) == false {
						m[words[j]] = append(m[words[j]], w)
					}
				}
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

func listAlternates(input []string) (output [1000000][]string) {
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

// Iterator is a struct that holds a []int array
// containing the maximum values that it should
// iterate up to. So Iterator{Limit:{2,2,3}} will
// iterate over all non-negative integer vectors
// of length 3 with values of the form:
// [0<=...<2, 0<=2...<3, 0<=...<3].
// Calling next until the returned value is nil
// will iterate over all these vectors in some order.
// Should be lazy and take very little memory. Typical
// use would be:
// ```
// itr := Iterator{{2,2,3}}
// for arr := itr.Next(); arr != nil; arr = itr.Next() {
//   ... // do something with `v`
// }
// ```
type Iterator struct {
	Limit []int
	arr   []int
}

// Next returns the next []int in the sequence.
// So something like {0,0} -> {0,1} -> {1,1} -> ...
// When you get to the end calling Next will return nil.
func (i *Iterator) Next() []int {
	if i.arr == nil && len(i.Limit) > 0 {
		i.arr = make([]int, len(i.Limit))
	} else {
		if !itrNext(i.arr, i.Limit) {
			return nil
		}
	}
	return i.arr
}

// true if can be incremented
func itrNext(arr, max []int) bool {
	if len(max) == 0 || len(arr) == 0 {
		panic("must have non-zero lengths")
	}
	if len(arr) == 1 {
		if arr[0] < max[0]-1 {
			arr[0]++
			return true
		}
		return false
	}
	if itrNext(arr[1:], max[1:]) {
		return true
	}
	if arr[0] < max[0]-1 {
		arr[0]++
		for i := 1; i < len(arr); i++ {
			arr[i] = 0
		}
		return true
	}
	return false
}

func init() {
	// initialize the thesaurs and the syllable dictionary

	thesaurus, err = loadThesaurus("./resources/th_en_US_new.dat")
	if err != nil {
		panic(err)
	}
	fmt.Println(thesaurus["define"])
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
	//sentence := `This is the sixth time we have had the pleasure of writing a birthday blog post for Go, and we would not be doing so if not for the wonderful and passionate people in our community. The Go team would like to thank everyone who has contributed code, written an open source library, authored a blog post, helped a new gopher, or just given Go a try. `
	sentence := `Classical thinkers employed classification as a way to define and assess the quality of poetry. Notably, Aristotle's Poetics describes the three genres of poetry: the epic, comic, and tragic, and develops rules to distinguish the highest-quality poetry of each genre, based on the underlying purposes of that genre`
	fmt.Println("ORIGINAL:")
	fmt.Println(sentence)
	sentence = strings.ToLower(sentence)
	sentence = strings.Replace(sentence, "don't", "do not", -1)
	sentence = strings.Replace(sentence, "'", "", -1)
	sentenceWords = regexp.MustCompile(`(\w+)`).FindAllString(sentence, -1)
	fmt.Println(sentenceWords)
	var goodNodes []node
	syllableTarget := 5
	for i := 1; i <= syllableTarget; i++ {
		fmt.Println(sentenceWords[0:i])
		nodes = append(nodes, node{start: 0, end: i})
	}
	for i := 0; i < len(nodes); i++ {
		nodes[i].syns = make(map[int][]string)
		for j := nodes[i].start; j < nodes[i].end; j++ {
			fmt.Println(sentenceWords[j])
			nodes[i].syns[j-nodes[i].start] = getSynonyms(sentenceWords[j])
			nodes[i].numSyns = append(nodes[i].numSyns, len(nodes[i].syns[j-nodes[i].start]))
		}
		fmt.Println(nodes[i])
		itr := &Iterator{Limit: nodes[i].numSyns}
		for arr := itr.Next(); arr != nil; arr = itr.Next() {
			totalSyllables := 0
			testSentence := ""
			for k := 0; k < len(arr); k++ {
				totalSyllables = totalSyllables + cmudict[nodes[i].syns[k][arr[k]]]
				testSentence = testSentence + nodes[i].syns[k][arr[k]] + " "
			}
			// fmt.Println(testSentence)
			// fmt.Println(totalSyllables)
			if totalSyllables == syllableTarget {
				fmt.Println("GOOD: " + testSentence)
				goodNodes = append(goodNodes, node{haikus: []string{testSentence}, start: nodes[i].start, end: nodes[i].end})
			}
		}

	}

	nodes = nil
	syllableTarget := 5
	for n := 0; n < len(goodNodes); n++ {
		for i := 1; i <= syllableTarget; i++ {
			fmt.Println(goodNodes.haikus[0])
			fmt.Println(sentenceWords[goodNodes[n].end : goodNodes[n].start+i])
			nodes = append(nodes, node{start: goodNodes[n].end, end: goodNodes[n].start + i})
		}

	}

	// alternatives := listAlternates(words)
	// bestNum := 0
	// var bestHaikus []string
	// var bestIsHaikus []bool

	// for i := 0; i < len(alternatives); i++ {
	// 	for j := 0; j < 1; j++ {
	// 		haikuString, isHaikus, numHaikus := getHaikus(alternatives[i][j:])
	// 		if numHaikus > bestNum {
	// 			if j > 0 {
	// 				haikuString = append([]string{strings.Join(alternatives[i][:j], " ")}, haikuString...)
	// 				isHaikus = append([]bool{false}, isHaikus...)
	// 			}
	// 			bestHaikus = haikuString
	// 			bestIsHaikus = isHaikus
	// 		}
	// 	}
	// }
	// fmt.Println("\n\nBEST HAIKU:")
	// for i, bestHaiku := range bestHaikus {
	// 	if bestIsHaikus[i] == true {
	// 		fmt.Printf("\n")
	// 	}
	// 	fmt.Println(bestHaiku)
	// 	if bestIsHaikus[i] == true {
	// 		fmt.Printf("\n")
	// 	}
	// }

	// Todo: replace each word with the puncuation near the word in the original

}
