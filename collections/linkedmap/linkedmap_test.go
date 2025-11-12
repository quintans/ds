package linkedmap_test

import (
	"slices"
	"testing"

	"github.com/quintans/ds/collections/linkedmap"
	"github.com/stretchr/testify/require"
)

func TestMapSet(t *testing.T) {
	m := linkedmap.New[int, string]()
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")
	m.Put(1, "a") // overwrite

	require.Equal(t, 7, m.Size())
	require.EqualValues(t, []int{5, 6, 7, 3, 4, 1, 2}, slices.Collect(m.Keys()))
	require.EqualValues(t, []string{"e", "f", "g", "c", "d", "a", "b"}, slices.Collect(m.Values()))

	// key,expectedValue,expectedFound
	tests1 := [][]any{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "", false},
	}

	for _, test := range tests1 {
		// retrievals
		actualValue, actualFound := m.Get(test[0].(int))
		require.Equal(t, test[1], actualValue)
		require.Equal(t, test[2], actualFound)
	}
}

func TestMapDelete(t *testing.T) {
	m := linkedmap.New[int, string]()
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")
	m.Put(1, "a") // overwrite

	m.Delete(5)
	m.Delete(6)
	m.Delete(7)
	m.Delete(8)
	m.Delete(5)

	require.EqualValues(t, []int{3, 4, 1, 2}, slices.Collect(m.Keys()))
	require.EqualValues(t, []string{"c", "d", "a", "b"}, slices.Collect(m.Values()))
	require.Equal(t, 4, m.Size())

	tests2 := [][]any{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "", false},
		{6, "", false},
		{7, "", false},
		{8, "", false},
	}

	for _, test := range tests2 {
		actualValue, actualFound := m.Get(test[0].(int))
		require.Equal(t, test[1], actualValue)
		require.Equal(t, test[2], actualFound)
	}

	m.Delete(1)
	m.Delete(4)
	m.Delete(2)
	m.Delete(3)
	m.Delete(2)
	m.Delete(2)

	require.Len(t, slices.Collect(m.Keys()), 0)
	require.Len(t, slices.Collect(m.Values()), 0)
	require.Equal(t, 0, m.Size())
}

func TestMapValues(t *testing.T) {
	expected := []string{"c", "b", "a"}

	m := linkedmap.New[string, int]()
	for k, v := range expected {
		m.Put(v, k+1)
	}

	count := 0
	for v := range m.Values() {
		require.Equal(t, count+1, v)
		count++
	}

	require.Equal(t, 3, count)
}

func TestMapEntries(t *testing.T) {
	m := linkedmap.New[string, int]()
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)

	expectedKeys := []string{"c", "a", "b"}
	expectedValues := []int{3, 1, 2}
	require.Equal(t, 3, m.Size())
	cnt := 0
	for k, v := range m.Entries() {
		require.Equal(t, expectedKeys[cnt], k)
		require.Equal(t, expectedValues[cnt], v)
		cnt++
	}
	require.Equal(t, 3, cnt)
}
