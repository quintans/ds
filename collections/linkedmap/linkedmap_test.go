package linkedmap_test

import (
	"testing"

	"github.com/quintans/ds/collections/linkedmap"
	"github.com/stretchr/testify/require"
)

func TestMapSet(t *testing.T) {
	m := linkedmap.New[int, string]()
	m.Set(5, "e")
	m.Set(6, "f")
	m.Set(7, "g")
	m.Set(3, "c")
	m.Set(4, "d")
	m.Set(1, "x")
	m.Set(2, "b")
	m.Set(1, "a") // overwrite

	require.Equal(t, 7, m.Size())
	require.EqualValues(t, []int{5, 6, 7, 3, 4, 1, 2}, m.Keys())
	require.EqualValues(t, []string{"e", "f", "g", "c", "d", "a", "b"}, m.Values())

	// key,expectedValue,expectedFound
	tests1 := [][]interface{}{
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
	m.Set(5, "e")
	m.Set(6, "f")
	m.Set(7, "g")
	m.Set(3, "c")
	m.Set(4, "d")
	m.Set(1, "x")
	m.Set(2, "b")
	m.Set(1, "a") // overwrite

	m.Delete(5)
	m.Delete(6)
	m.Delete(7)
	m.Delete(8)
	m.Delete(5)

	require.EqualValues(t, []int{3, 4, 1, 2}, m.Keys())
	require.EqualValues(t, []string{"c", "d", "a", "b"}, m.Values())
	require.Equal(t, 4, m.Size())

	tests2 := [][]interface{}{
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

	require.Len(t, m.Keys(), 0)
	require.Len(t, m.Values(), 0)
	require.Equal(t, 0, m.Size())
}

func TestMapRange(t *testing.T) {
	expected := []string{"c", "b", "a"}

	m := linkedmap.New[string, int]()
	for k, v := range expected {
		m.Set(v, k+1)
	}

	count := 0
	m.Range(func(key string, value int, idx int) bool {
		count++
		require.Equal(t, idx+1, value)
		require.Equal(t, expected[idx], key)
		return true
	})
	require.Equal(t, 3, count)
}

func TestMapEach(t *testing.T) {
	expected := []string{"c", "b", "a"}

	m := linkedmap.New[string, int]()
	for k, v := range expected {
		m.Set(v, k+1)
	}

	count := 0
	m.Each(func(key string, value int) {
		count++
		require.Equal(t, count, value)
		require.Equal(t, expected[count-1], key)
	})
	require.Equal(t, 3, count)
}

func TestMapEntries(t *testing.T) {
	m := linkedmap.New[string, int]()
	m.Set("c", 3)
	m.Set("a", 1)
	m.Set("b", 2)

	expectedKeys := []string{"c", "a", "b"}
	expectedValues := []int{3, 1, 2}
	entries := m.Entries()
	require.Len(t, entries, 3)
	for k, v := range entries {
		require.Equal(t, expectedKeys[k], v.Key())
		require.Equal(t, expectedValues[k], v.Value())
	}
}
