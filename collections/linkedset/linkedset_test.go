package linkedset_test

import (
	"slices"
	"testing"

	"github.com/quintans/ds/collections/linkedset"
	"github.com/stretchr/testify/require"
)

var unsortedArray = []int{10, 2, 6, 71, 3}

func TestLinkedHashSetSame(t *testing.T) {
	r := require.New(t)

	list := linkedset.New[int]()
	list.Add(unsortedArray...)
	list.Add(2)
	r.Equal(5, list.Size())
}

func TestLinkedHashSetAddAll(t *testing.T) {
	r := require.New(t)

	list := linkedset.New[int]()
	list.Add(unsortedArray...)
	values := slices.Collect(list.Values())
	r.EqualValues(values, unsortedArray)
}

func TestLinkedHashSetContains(t *testing.T) {
	r := require.New(t)

	list := linkedset.New[int]()
	list.Add(unsortedArray...)

	r.False(list.Contains(25))
	r.True(list.Contains(2))
}

func TestLinkedHashSetValues(t *testing.T) {
	r := require.New(t)

	list := linkedset.New[int]()
	list.Add(unsortedArray...)

	arr := []int{}
	for v := range list.Values() {
		arr = append(arr, v)
	}
	r.EqualValues(unsortedArray, arr)
}
