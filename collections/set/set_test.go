package set_test

import (
	"testing"

	"github.com/quintans/ds/collections/set"
	"github.com/stretchr/testify/require"
)

var unsortedArray = []int{10, 2, 6, 71, 3}

func TestHashSetSame(t *testing.T) {
	r := require.New(t)

	list := set.New[int]()
	list.Add(unsortedArray...)
	list.Add(2)
	r.Equal(5, list.Size())
}

func TestHashSetContains(t *testing.T) {
	r := require.New(t)

	list := set.New[int]()
	list.Add(unsortedArray...)

	r.False(list.Contains(25))
	r.True(list.Contains(2))
}

func TestHashSetValues(t *testing.T) {
	r := require.New(t)

	hs := set.New[int]()
	hs.Add(unsortedArray...)

	tmp := []int{}
	for v := range hs.Values() {
		tmp = append(tmp, v)
	}
	r.ElementsMatch(unsortedArray, tmp)

	hs.Delete(unsortedArray[4])

	trimmed := unsortedArray[:4]
	tmp = []int{}
	for v := range hs.Values() {
		tmp = append(tmp, v)
	}
	r.ElementsMatch(trimmed, tmp)
}
