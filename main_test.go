package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prod(a []int) int {
	p := 1
	for _, i := range a {
		p = p * i
	}
	return p
}

func TestIterator(t *testing.T) {
	cases := [][]int{
		[]int{1, 2, 3},
		[]int{2, 2},
		[]int{1, 1, 1, 1},
	}
	for _, c := range cases {
		var i int
		itr := &Iterator{Limit: c}
		for a := itr.Next(); a != nil; a = itr.Next() {
			i++
		}
		assert.Equal(t, prod(c), i, fmt.Sprintf("%+v", c))
	}
}

func TestToWords(t *testing.T) {
	cases := []struct {
		in  string
		out []string
	}{
		{"this is a sentance", []string{"this", "is", "a", "sentance"}},
		{"don't listen", []string{"do", "not", "listen"}},
		{"a tree falls", []string{"a", "tree", "falls"}},
	}
	for _, c := range cases {
		assert.True(t, reflect.DeepEqual(sentanceToWords(c.in), c.out), fmt.Sprintf("%+v\n%+v", c, sentanceToWords(c.in)))
	}
}
