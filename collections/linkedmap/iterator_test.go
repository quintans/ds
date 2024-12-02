package linkedmap_test

import (
	"testing"

	"github.com/quintans/dstruct/collections/linkedmap"
	"github.com/stretchr/testify/require"
)

func TestMapIteratorNextOnEmpty(t *testing.T) {
	m := linkedmap.New[string, int]()
	it := m.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty map")
	}
}

func TestMapIteratorPrevOnEmpty(t *testing.T) {
	m := linkedmap.New[string, int]()
	it := m.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty map")
	}
}

func TestMapIteratorNext(t *testing.T) {
	m := linkedmap.New[string, int]()
	m.Set("c", 1)
	m.Set("a", 2)
	m.Set("b", 3)

	expectedKeys := []string{"c", "a", "b"}
	expectedValues := []int{1, 2, 3}

	it := m.Iterator()
	count := 0
	for it.Next() {
		require.Equal(t, expectedKeys[count], it.Key())
		require.Equal(t, expectedValues[count], it.Value())
		count++
	}
	require.Equal(t, 3, count)
}

func TestMapIteratorPrev(t *testing.T) {
	m := linkedmap.New[string, int]()
	m.Set("c", 1)
	m.Set("a", 2)
	m.Set("b", 3)

	expectedKeys := []string{"c", "a", "b"}
	expectedValues := []int{1, 2, 3}

	it := m.Iterator()
	countDown := m.Size()
	for it.Prev() {
		countDown--
		require.Equal(t, expectedKeys[countDown], it.Key())
		require.Equal(t, expectedValues[countDown], it.Value())
	}
	require.Equal(t, 0, countDown)
}

func TestMapIteratorBegin(t *testing.T) {
	m := linkedmap.New[int, string]()
	it := m.Iterator()
	m.Set(3, "c")
	m.Set(1, "a")
	m.Set(2, "b")
	for it.Next() {
	}
	require.True(t, it.IsOut())
	it.Reset()
	it.Next()
	require.Equal(t, 3, it.Key())
	require.Equal(t, "c", it.Value())
}

func TestMapIteratorEnd(t *testing.T) {
	m := linkedmap.New[int, string]()
	it := m.Iterator()
	m.Set(3, "c")
	m.Set(1, "a")
	m.Set(2, "b")
	it.Next()
	it.Reset()
	it.Prev()
	require.Equal(t, 2, it.Key())
	require.Equal(t, "b", it.Value())
}

func TestMapIteratorFirst(t *testing.T) {
	m := linkedmap.New[int, string]()
	m.Set(3, "c")
	m.Set(1, "a")
	m.Set(2, "b")
	it := m.Iterator()
	require.Equal(t, true, it.First())
	require.Equal(t, 3, it.Key())
	require.Equal(t, "c", it.Value())
}

func TestMapIteratorLast(t *testing.T) {
	m := linkedmap.New[int, string]()
	m.Set(3, "c")
	m.Set(1, "a")
	m.Set(2, "b")
	it := m.Iterator()
	require.Equal(t, true, it.Last())
	require.Equal(t, 2, it.Key())
	require.Equal(t, "b", it.Value())
}

func TestMapIteratorRemoveNone(t *testing.T) {
	m := linkedmap.New[int, string]()
	m.Set(3, "c")
	m.Set(1, "a")
	m.Set(2, "b")
	it := m.Iterator()
	it.Remove()
	require.EqualValues(t, []int{3, 1, 2}, m.Keys())
}

func TestMapIteratorNextRemove(t *testing.T) {
	m := linkedmap.New[int, string]()
	m.Set(3, "c")
	m.Set(1, "a")
	m.Set(2, "b")
	it := m.Iterator()
	it.Next()
	it.Next()
	it.Remove()
	require.EqualValues(t, []int{3, 2}, m.Keys())
	require.Equal(t, 3, it.Key())
	require.Equal(t, "c", it.Value())
}

func TestMapIteratorPrevRemove(t *testing.T) {
	m := linkedmap.New[int, string]()
	m.Set(3, "c")
	m.Set(1, "a")
	m.Set(2, "b")
	it := m.Iterator()
	it.Prev()
	it.Prev()
	it.Remove()
	require.EqualValues(t, []int{3, 2}, m.Keys())
	require.Equal(t, 2, it.Key())
	require.Equal(t, "b", it.Value())
}

func TestMapIteratorRemoveFirst(t *testing.T) {
	m := linkedmap.New[int, string]()
	m.Set(3, "c")
	m.Set(1, "a")
	m.Set(2, "b")
	it := m.Iterator()
	it.Next()
	it.Remove()
	require.EqualValues(t, []int{1, 2}, m.Keys())
	require.Equal(t, 0, it.Key())
	require.Equal(t, "", it.Value())
	require.True(t, it.IsOut())
	it.Next()
	require.Equal(t, 1, it.Key())
	require.Equal(t, "a", it.Value())
}

func TestMapIteratorRemoveLast(t *testing.T) {
	m := linkedmap.New[int, string]()
	m.Set(3, "c")
	m.Set(1, "a")
	m.Set(2, "b")
	it := m.Iterator()
	it.Prev()
	it.Remove()
	require.EqualValues(t, []int{3, 1}, m.Keys())
	require.Equal(t, 0, it.Key())
	require.Equal(t, "", it.Value())
	require.True(t, it.IsOut())
	it.Prev()
	require.Equal(t, 1, it.Key())
	require.Equal(t, "a", it.Value())
}
