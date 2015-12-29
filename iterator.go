package main

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
