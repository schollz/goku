package main

import (
	"fmt"
	"math/rand"
)

func synonyms(a string) []string {
	if a == "x" {
		return []string{"a", "b", "c"}
	}
	return []string{"d", "e"}
}

func randChoices(limits []int) (choices []int) {
	for _, l := range limits {
		choices = append(choices, rand.Intn(l))
	}
	return
}

func listAlternates(input []string) (output [10][]string) {
	totals := make([]int, len(input))
	for i, w := range input {
		totals[i] = len(synonyms(w))
	}
	for i := 0; i < len(output); i++ {
		choices := randChoices(totals)
		for j, w := range input {
			output[i] = append(output[i], synonyms(w)[choices[j]])
		}
	}
	return
}

func main() {
	input := []string{"x", "y"}
	fmt.Println(listAlternates(input))
}
