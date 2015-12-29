package main

import "fmt"

type node struct {
	haikus  []string
	syns    map[int][]string
	numSyns []int
	start   int
	end     int
}

var (
	nodes         []node
	sentenceWords []string
)

var skipwords = []string{"a", "about", "above", "above", "across", "after", "afterwards", "again", "against", "all", "almost", "alone", "along", "already", "also", "although", "always", "am", "among", "amongst", "amoungst", "amount", "an", "and", "another", "any", "anyhow", "anyone", "anything", "anyway", "anywhere", "are", "around", "as", "at", "back", "be", "became", "because", "become", "becomes", "becoming", "been", "before", "beforehand", "behind", "being", "below", "beside", "besides", "between", "beyond", "bill", "both", "bottom", "but", "by", "call", "can", "cannot", "cant", "co", "con", "could", "couldnt", "cry", "de", "describe", "detail", "do", "done", "down", "due", "during", "each", "eg", "eight", "either", "eleven", "else", "elsewhere", "empty", "enough", "etc", "even", "ever", "every", "everyone", "everything", "everywhere", "except", "few", "fifteen", "fify", "fill", "find", "fire", "first", "five", "for", "former", "formerly", "forty", "found", "four", "from", "front", "full", "further", "get", "give", "go", "had", "has", "hasnt", "have", "he", "hence", "her", "here", "hereafter", "hereby", "herein", "hereupon", "hers", "herself", "him", "himself", "his", "how", "however", "hundred", "ie", "if", "in", "inc", "indeed", "interest", "into", "is", "it", "its", "itself", "keep", "last", "latter", "latterly", "least", "less", "ltd", "made", "many", "may", "me", "meanwhile", "might", "mill", "mine", "more", "moreover", "most", "mostly", "move", "much", "must", "my", "myself", "name", "namely", "neither", "never", "nevertheless", "next", "nine", "no", "nobody", "none", "noone", "nor", "not", "nothing", "now", "nowhere", "of", "off", "often", "on", "once", "one", "only", "onto", "or", "other", "others", "otherwise", "our", "ours", "ourselves", "out", "over", "own", "part", "per", "perhaps", "please", "put", "rather", "re", "same", "see", "seem", "seemed", "seeming", "seems", "serious", "several", "she", "should", "show", "side", "since", "sincere", "six", "sixty", "so", "some", "somehow", "someone", "something", "sometime", "sometimes", "somewhere", "still", "such", "system", "take", "ten", "than", "that", "the", "their", "them", "themselves", "then", "thence", "there", "thereafter", "thereby", "therefore", "therein", "thereupon", "these", "they", "thickv", "thin", "third", "this", "those", "though", "three", "through", "throughout", "thru", "thus", "to", "together", "too", "top", "toward", "towards", "twelve", "twenty", "two", "un", "under", "until", "up", "upon", "us", "very", "via", "was", "we", "well", "were", "what", "whatever", "when", "whence", "whenever", "where", "whereafter", "whereas", "whereby", "wherein", "whereupon", "wherever", "whether", "which", "while", "whither", "who", "whoever", "whole", "whom", "whose", "why", "will", "with", "within", "without", "would", "yet", "you", "your", "yours", "yourself", "yourselves", "the"}

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

func init() {
	// initialize the thesaurs and the syllable dictionary
	var err error
	thesaurus, err = loadThesaurus("./resources/th_en_US_new.dat")
	if err != nil {
		panic(err)
	}
	fmt.Println(thesaurus["define"])
	cmudict, err = loadCmudict("./resources/cmudict.0.7a")
	if err != nil {
		panic(err)
	}
}

func main() {
	sentence := `Classical thinkers employed classification as a way to define and assess the quality of poetry. Notably, Aristotle's Poetics describes the three genres of poetry: the epic, comic, and tragic, and develops rules to distinguish the highest-quality poetry of each genre, based on the underlying purposes of that genre`
	sentenceWords = sentanceToWords(sentence)
	fmt.Println("ORIGINAL:")
	fmt.Println(sentence)
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

	fmt.Println("7")
	nodes = nil
	syllableTarget = 7
	for n := 0; n < len(goodNodes); n++ {
		for i := 1; i <= syllableTarget; i++ {
			// fmt.Println(goodNodes[n].haikus[len(goodNodes[n].haikus)-1])
			// fmt.Println(sentenceWords[goodNodes[n].end : goodNodes[n].end+i])
			nodes = append(nodes, node{haikus: goodNodes[n].haikus, start: goodNodes[n].end, end: goodNodes[n].end + i})
		}
	}
	goodNodes = nil
	for i := 0; i < len(nodes); i++ {
		nodes[i].syns = make(map[int][]string)
		for j := nodes[i].start; j < nodes[i].end; j++ {
			fmt.Println(sentenceWords[j])
			nodes[i].syns[j-nodes[i].start] = getSynonyms(sentenceWords[j])
			nodes[i].numSyns = append(nodes[i].numSyns, len(nodes[i].syns[j-nodes[i].start]))
		}
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
				goodNodes = append(goodNodes, node{haikus: append(nodes[i].haikus, testSentence), start: nodes[i].start, end: nodes[i].end})
			}
		}

	}
	fmt.Println(goodNodes)

	// Todo: replace each word with the puncuation near the word in the original

}
