package linkedhashset_test

import (
	"testing"

	"github.com/quintans/dstruct/collections/linkedhashset"
	"github.com/stretchr/testify/require"
)

var unsortedArray = []int{10, 2, 6, 71, 3}

func TestLinkedHashSetSame(t *testing.T) {
	r := require.New(t)

	list := linkedhashset.New[int]()
	list.Add(unsortedArray...)
	list.Add(2)
	r.Equal(5, list.Size())
}

func TestLinkedHashSetAddAll(t *testing.T) {
	r := require.New(t)

	list := linkedhashset.New[int]()
	list.Add(unsortedArray...)
	r.EqualValues(list.ToSlice(), unsortedArray)
}

func TestLinkedHashSetContains(t *testing.T) {
	r := require.New(t)

	list := linkedhashset.New[int]()
	list.Add(unsortedArray...)

	r.False(list.Contains(25))
	r.True(list.Contains(2))
}

func TestLinkedHashSetIterator(t *testing.T) {
	r := require.New(t)

	list := linkedhashset.New[int]()
	list.Add(unsortedArray...)

	arr := []int{}
	for e := list.Iterator(); e.HasNext(); {
		arr = append(arr, e.Next())
	}
	r.EqualValues(unsortedArray, arr)
}
