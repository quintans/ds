package linkedlist_test

import (
	"testing"

	"github.com/quintans/dstruct/collections/linkedlist"
	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	expect := []int{10, 20, 30, 40, 50}
	ll := linkedlist.New[int]()
	ll.Add(expect...)

	actual := []int{}
	for it := ll.Iterator(); it.HasNext(); {
		v := it.Next()
		actual = append(actual, v)
		if v%20 == 0 {
			it.Remove()
		}
	}
	assert.Equal(t, expect, actual)

	expect = []int{10, 30, 50}
	actual = []int{}
	for it := ll.Iterator(); it.HasNext(); {
		actual = append(actual, it.Next())
	}
	assert.Equal(t, expect, actual)
}
