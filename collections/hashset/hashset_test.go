package hashset_test

import (
	"testing"

	"github.com/quintans/dstruct/collections/hashset"
	"github.com/stretchr/testify/require"
)

var unsortedArray = []int{10, 2, 6, 71, 3}

func TestHashSetSame(t *testing.T) {
	r := require.New(t)

	list := hashset.New[int]()
	list.Add(unsortedArray...)
	list.Add(2)
	r.Equal(5, list.Size())
}

func TestHashSetContains(t *testing.T) {
	r := require.New(t)

	list := hashset.New[int]()
	list.Add(unsortedArray...)

	r.False(list.Contains(25))
	r.True(list.Contains(2))
}

func TestHashSetInEnumerator(t *testing.T) {
	r := require.New(t)

	hs := hashset.New[int]()
	hs.Add(unsortedArray...)

	tmp := []int{}
	for it := hs.Iterator(); it.HasNext(); {
		tmp = append(tmp, it.Next())
	}
	r.ElementsMatch(unsortedArray, tmp)

	tmp = []int{}
	hs.ForEach(func(i int, h int) {
		tmp = append(tmp, h)
	})
	r.ElementsMatch(unsortedArray, tmp)

	hs.Delete(unsortedArray[4])

	trimmed := unsortedArray[:4]
	tmp = []int{}
	for it := hs.Iterator(); it.HasNext(); {
		tmp = append(tmp, it.Next())
	}
	r.ElementsMatch(trimmed, tmp)

	tmp = []int{}
	hs.ForEach(func(i int, h int) {
		tmp = append(tmp, h)
	})
	r.ElementsMatch(trimmed, tmp)
}
