package linkedlist_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/quintans/ds/collections/linkedlist"
)

func TestNew(t *testing.T) {
	l := linkedlist.New[int]()
	assert.NotNil(t, l)
	assert.Equal(t, 0, l.Size())
	_, err := l.PeekFirst()
	assert.Error(t, err)
	_, err = l.PeekLast()
	assert.Error(t, err)
}

func TestAdd(t *testing.T) {
	l := linkedlist.New[int]()
	e1 := l.Add(1)
	assert.Equal(t, 1, l.Size())
	v, err := l.PeekFirst()
	require.NoError(t, err)
	assert.Equal(t, 1, v)
	v, err = l.PeekLast()
	require.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, l.Head(), e1)
	assert.Equal(t, l.Tail(), e1)

	e2 := l.Add(2)
	assert.Equal(t, 2, l.Size())
	v, err = l.PeekFirst()
	require.NoError(t, err)
	assert.Equal(t, 1, v)
	v, err = l.PeekLast()
	require.NoError(t, err)
	assert.Equal(t, 2, v)
	assert.Equal(t, l.Head(), e1)
	assert.Equal(t, l.Tail(), e2)
	assert.Equal(t, e1.Next(), e2)
	assert.Equal(t, e2.Previous(), e1)

	e3 := l.Add(3)
	assert.Equal(t, 3, l.Size())
	v, err = l.PeekLast()
	require.NoError(t, err)
	assert.Equal(t, 3, v)
	assert.Equal(t, l.Tail(), e3)
	assert.Equal(t, e2.Next(), e3)
	assert.Equal(t, e3.Previous(), e2)
}

func TestAddFirst(t *testing.T) {
	l := linkedlist.New[string]()
	e1 := l.AddFirst("a")
	assert.Equal(t, 1, l.Size())
	v, err := l.PeekFirst()
	require.NoError(t, err)
	assert.Equal(t, "a", v)
	v, err = l.PeekLast()
	require.NoError(t, err)
	assert.Equal(t, "a", v)
	assert.Equal(t, l.Head(), e1)
	assert.Equal(t, l.Tail(), e1)

	e2 := l.AddFirst("b")
	assert.Equal(t, 2, l.Size())
	v, err = l.PeekFirst()
	require.NoError(t, err)
	assert.Equal(t, "b", v)
	v, err = l.PeekLast()
	require.NoError(t, err)
	assert.Equal(t, "a", v)
	assert.Equal(t, l.Head(), e2)
	assert.Equal(t, l.Tail(), e1)
	assert.Equal(t, e2.Next(), e1)
	assert.Equal(t, e1.Previous(), e2)

	e3 := l.AddFirst("c")
	assert.Equal(t, 3, l.Size())
	v, err = l.PeekFirst()
	require.NoError(t, err)
	assert.Equal(t, "c", v)
	assert.Equal(t, l.Head(), e3)
	assert.Equal(t, e3.Next(), e2)
	assert.Equal(t, e2.Previous(), e3)
}

func TestClear(t *testing.T) {
	l := linkedlist.New[int]()
	l.Add(1)
	l.Add(2)
	assert.Equal(t, 2, l.Size())
	l.Clear()
	assert.Equal(t, 0, l.Size())
	assert.Nil(t, l.Head())
	assert.Nil(t, l.Tail())
	_, err := l.PeekFirst()
	assert.Error(t, err)
}

func TestGetSetAddAt(t *testing.T) {
	l := linkedlist.New[int]()
	l.Add(10)
	l.Add(20)
	l.Add(30)

	// Get
	v, err := l.Get(0)
	require.NoError(t, err)
	assert.Equal(t, 10, v)
	v, err = l.Get(1)
	require.NoError(t, err)
	assert.Equal(t, 20, v)
	v, err = l.Get(2)
	require.NoError(t, err)
	assert.Equal(t, 30, v)
	_, err = l.Get(3)
	assert.Error(t, err)
	_, err = l.Get(-1)
	assert.Error(t, err)

	// Set
	err = l.Set(1, 25)
	require.NoError(t, err)
	v, err = l.Get(1)
	require.NoError(t, err)
	assert.Equal(t, 25, v)
	err = l.Set(3, 35)
	assert.Error(t, err)
	err = l.Set(-1, 5)
	assert.Error(t, err)

	// AddAt
	err = l.AddAt(1, 15) // Insert 15 between 10 and 25
	require.NoError(t, err)
	assert.Equal(t, 4, l.Size())
	v, err = l.Get(0)
	require.NoError(t, err)
	assert.Equal(t, 10, v)
	v, err = l.Get(1)
	require.NoError(t, err)
	assert.Equal(t, 15, v)
	v, err = l.Get(2)
	require.NoError(t, err)
	assert.Equal(t, 25, v)
	v, err = l.Get(3)
	require.NoError(t, err)
	assert.Equal(t, 30, v)

	err = l.AddAt(0, 5) // Insert 5 at the beginning
	require.NoError(t, err)
	assert.Equal(t, 5, l.Size())
	v, err = l.Get(0)
	require.NoError(t, err)
	assert.Equal(t, 5, v)
	v, err = l.PeekFirst()
	require.NoError(t, err)
	assert.Equal(t, 5, v)

	// AddAt index == size is not allowed, use Add()
	err = l.AddAt(5, 35)
	assert.Error(t, err)

	err = l.AddAt(-1, 0)
	assert.Error(t, err)
}

func TestRemoveFirst(t *testing.T) {
	l := linkedlist.New[string]()
	_, err := l.RemoveFirst()
	assert.Error(t, err)

	l.Add("one")
	l.Add("two")
	l.Add("three")

	v, err := l.RemoveFirst()
	require.NoError(t, err)
	assert.Equal(t, "one", v)
	assert.Equal(t, 2, l.Size())
	p, err := l.PeekFirst()
	require.NoError(t, err)
	assert.Equal(t, "two", p)

	v, err = l.RemoveFirst()
	require.NoError(t, err)
	assert.Equal(t, "two", v)
	assert.Equal(t, 1, l.Size())
	p, err = l.PeekFirst()
	require.NoError(t, err)
	assert.Equal(t, "three", p)
	p, err = l.PeekLast()
	require.NoError(t, err)
	assert.Equal(t, "three", p)

	v, err = l.RemoveFirst()
	require.NoError(t, err)
	assert.Equal(t, "three", v)
	assert.Equal(t, 0, l.Size())
	_, err = l.PeekFirst()
	assert.Error(t, err)
	_, err = l.PeekLast()
	assert.Error(t, err)

	_, err = l.RemoveFirst()
	assert.Error(t, err)
}

func TestRemoveLast(t *testing.T) {
	l := linkedlist.New[string]()
	_, err := l.RemoveLast()
	assert.Error(t, err)

	l.Add("one")
	l.Add("two")
	l.Add("three")

	v, err := l.RemoveLast()
	require.NoError(t, err)
	assert.Equal(t, "three", v)
	assert.Equal(t, 2, l.Size())
	p, err := l.PeekLast()
	require.NoError(t, err)
	assert.Equal(t, "two", p)

	v, err = l.RemoveLast()
	require.NoError(t, err)
	assert.Equal(t, "two", v)
	assert.Equal(t, 1, l.Size())
	p, err = l.PeekLast()
	require.NoError(t, err)
	assert.Equal(t, "one", p)
	p, err = l.PeekFirst()
	require.NoError(t, err)
	assert.Equal(t, "one", p)

	v, err = l.RemoveLast()
	require.NoError(t, err)
	assert.Equal(t, "one", v)
	assert.Equal(t, 0, l.Size())
	_, err = l.PeekFirst()
	assert.Error(t, err)
	_, err = l.PeekLast()
	assert.Error(t, err)

	_, err = l.RemoveLast()
	assert.Error(t, err)
}

func TestDeleteAt(t *testing.T) {
	l := linkedlist.New[int]()
	l.Add(10)
	l.Add(20)
	l.Add(30)
	l.Add(40)

	_, err := l.DeleteAt(4)
	assert.Error(t, err)
	_, err = l.DeleteAt(-1)
	assert.Error(t, err)

	// Delete middle
	v, err := l.DeleteAt(1)
	require.NoError(t, err)
	assert.Equal(t, 20, v)
	assert.Equal(t, 3, l.Size())
	vals := collectValues(l)
	assert.Equal(t, []int{10, 30, 40}, vals)

	// Delete head
	v, err = l.DeleteAt(0)
	require.NoError(t, err)
	assert.Equal(t, 10, v)
	assert.Equal(t, 2, l.Size())
	vals = collectValues(l)
	assert.Equal(t, []int{30, 40}, vals)
	p, err := l.PeekFirst()
	require.NoError(t, err)
	assert.Equal(t, 30, p)

	// Delete tail
	v, err = l.DeleteAt(1)
	require.NoError(t, err)
	assert.Equal(t, 40, v)
	assert.Equal(t, 1, l.Size())
	vals = collectValues(l)
	assert.Equal(t, []int{30}, vals)
	p, err = l.PeekLast()
	require.NoError(t, err)
	assert.Equal(t, 30, p)

	// Delete last remaining
	v, err = l.DeleteAt(0)
	require.NoError(t, err)
	assert.Equal(t, 30, v)
	assert.Equal(t, 0, l.Size())
	vals = collectValues(l)
	assert.Empty(t, vals)
	_, err = l.PeekFirst()
	assert.Error(t, err)
}

func TestDelete(t *testing.T) {
	l := linkedlist.New[string]()
	l.Add("a")
	l.Add("b")
	l.Add("c")
	l.Add("b") // duplicate
	l.Add("d")

	assert.False(t, l.Delete(equalFunc("e")))
	assert.Equal(t, 5, l.Size())

	// Delete middle (first occurrence)
	assert.True(t, l.Delete(equalFunc("b")))
	assert.Equal(t, 4, l.Size())
	vals := collectValues(l)
	assert.Equal(t, []string{"a", "c", "b", "d"}, vals)

	// Delete head
	assert.True(t, l.Delete(equalFunc("a")))
	assert.Equal(t, 3, l.Size())
	vals = collectValues(l)
	assert.Equal(t, []string{"c", "b", "d"}, vals)
	p, err := l.PeekFirst()
	require.NoError(t, err)
	assert.Equal(t, "c", p)

	// Delete tail
	assert.True(t, l.Delete(equalFunc("d")))
	assert.Equal(t, 2, l.Size())
	vals = collectValues(l)
	assert.Equal(t, []string{"c", "b"}, vals)
	p, err = l.PeekLast()
	require.NoError(t, err)
	assert.Equal(t, "b", p)

	// Delete remaining
	assert.True(t, l.Delete(equalFunc("c")))
	assert.True(t, l.Delete(equalFunc("b")))
	assert.Equal(t, 0, l.Size())
	assert.False(t, l.Delete(equalFunc("c"))) // Already deleted
}

func equalFunc(s string) func(string) bool {
	return func(v string) bool {
		return v == s
	}
}

func TestMoveToFirstLast(t *testing.T) {
	l := linkedlist.New[int]()
	e1 := l.Add(1)
	e2 := l.Add(2)
	e3 := l.Add(3)
	e4 := l.Add(4)

	// Move middle to first
	l.MoveToFirst(e3) // 1, 2, 4 -> 3, 1, 2, 4
	assert.Equal(t, 4, l.Size())
	vals := collectValues(l)
	assert.Equal(t, []int{3, 1, 2, 4}, vals)
	assert.Equal(t, e3, l.Head())
	assert.Equal(t, e4, l.Tail())
	assert.Equal(t, e1, e3.Next())
	assert.Equal(t, e3, e1.Previous())
	assert.Equal(t, e2, e1.Next())
	assert.Equal(t, e1, e2.Previous())
	assert.Equal(t, e4, e2.Next())
	assert.Equal(t, e2, e4.Previous())

	// Move head to last
	l.MoveToLast(e3) // 3, 1, 2, 4 -> 1, 2, 4, 3
	assert.Equal(t, 4, l.Size())
	vals = collectValues(l)
	assert.Equal(t, []int{1, 2, 4, 3}, vals)
	assert.Equal(t, e1, l.Head())
	assert.Equal(t, e3, l.Tail())
	assert.Equal(t, e3, e4.Next())
	assert.Equal(t, e4, e3.Previous())
	assert.Nil(t, e3.Next())

	// Move tail to first
	l.MoveToFirst(e3) // 1, 2, 4, 3 -> 3, 1, 2, 4
	assert.Equal(t, 4, l.Size())
	vals = collectValues(l)
	assert.Equal(t, []int{3, 1, 2, 4}, vals)
	assert.Equal(t, e3, l.Head())
	assert.Equal(t, e4, l.Tail())
	assert.Nil(t, e4.Next())
	assert.Equal(t, e2, e4.Previous())

	// Move first to last
	l.MoveToLast(e3) // 3, 1, 2, 4 -> 1, 2, 4, 3
	assert.Equal(t, 4, l.Size())
	vals = collectValues(l)
	assert.Equal(t, []int{1, 2, 4, 3}, vals)
	assert.Equal(t, e1, l.Head())
	assert.Equal(t, e3, l.Tail())

	// Move middle to last
	l.MoveToLast(e2) // 1, 2, 4, 3 -> 1, 4, 3, 2
	assert.Equal(t, 4, l.Size())
	vals = collectValues(l)
	assert.Equal(t, []int{1, 4, 3, 2}, vals)
	assert.Equal(t, e1, l.Head())
	assert.Equal(t, e2, l.Tail())
	assert.Equal(t, e4, e1.Next())
	assert.Equal(t, e1, e4.Previous())
	assert.Equal(t, e3, e4.Next())
	assert.Equal(t, e4, e3.Previous())
	assert.Equal(t, e2, e3.Next())
	assert.Equal(t, e3, e2.Previous())
	assert.Nil(t, e2.Next())
}

func TestElementRemove(t *testing.T) {
	l := linkedlist.New[int]()
	e1 := l.Add(1)
	e2 := l.Add(2)
	e3 := l.Add(3)

	// Remove middle
	e2.Remove()
	assert.Equal(t, 2, l.Size())
	vals := collectValues(l)
	assert.Equal(t, []int{1, 3}, vals)
	assert.Equal(t, e1, l.Head())
	assert.Equal(t, e3, l.Tail())
	assert.Equal(t, e3, e1.Next())
	assert.Equal(t, e1, e3.Previous())
	assert.Nil(t, e2.Next()) // Ensure removed element is detached
	assert.Nil(t, e2.Previous())

	// Remove head
	e1.Remove()
	assert.Equal(t, 1, l.Size())
	vals = collectValues(l)
	assert.Equal(t, []int{3}, vals)
	assert.Equal(t, e3, l.Head())
	assert.Equal(t, e3, l.Tail())
	assert.Nil(t, e1.Next())
	assert.Nil(t, e1.Previous())

	// Remove tail (last element)
	e3.Remove()
	assert.Equal(t, 0, l.Size())
	vals = collectValues(l)
	assert.Empty(t, vals)
	assert.Nil(t, l.Head())
	assert.Nil(t, l.Tail())
	assert.Nil(t, e3.Next())
	assert.Nil(t, e3.Previous())
}

func TestValues(t *testing.T) {
	l := linkedlist.New[string]()
	l.Add("x")
	l.Add("y")
	l.Add("z")

	var result []string
	for v := range l.Values() {
		result = append(result, v)
	}
	assert.Equal(t, []string{"x", "y", "z"}, result)

	// Test empty list
	l.Clear()
	result = nil
	for v := range l.Values() {
		result = append(result, v) // Should not run
	}
	assert.Empty(t, result)

	// Test early exit
	l.Add("a")
	l.Add("b")
	l.Add("c")
	result = nil
	for v := range l.Values() {
		result = append(result, v)
		if v == "b" {
			break // Simulate early exit
		}
	}
	assert.Equal(t, []string{"a", "b"}, result)
}

func TestReplaceAll(t *testing.T) {
	l := linkedlist.New[int]()
	l.Add(1)
	l.Add(2)
	l.Add(3)

	l.ReplaceAll(func(idx int, val int) int {
		return val*10 + idx
	})

	vals := collectValues(l)
	assert.Equal(t, []int{10, 21, 32}, vals)

	// Test empty list
	l.Clear()
	l.ReplaceAll(func(idx int, val int) int {
		t.Fail() // Should not be called
		return 0
	})
	assert.Equal(t, 0, l.Size())
}

func TestClone(t *testing.T) {
	l1 := linkedlist.New[string]()
	l1.Add("one")
	l1.Add("two")

	l2 := l1.Clone()
	assert.NotSame(t, l1, l2)
	assert.Equal(t, l1.Size(), l2.Size())
	assert.Equal(t, collectValues(l1), collectValues(l2))

	// Modify original, clone should be unaffected
	l1.Add("three")
	l1.RemoveFirst()
	assert.Equal(t, 2, l1.Size())
	assert.Equal(t, 2, l2.Size())
	assert.Equal(t, []string{"two", "three"}, collectValues(l1))
	assert.Equal(t, []string{"one", "two"}, collectValues(l2))

	// Ensure elements are different instances
	e1_1, _ := l1.Get(0)
	e2_1, _ := l2.Get(0)
	assert.NotSame(t, e1_1, e2_1)

	// Test cloning empty list
	l3 := linkedlist.New[int]()
	l4 := l3.Clone()
	assert.NotSame(t, l3, l4)
	assert.Equal(t, 0, l4.Size())
	assert.Nil(t, l4.Head())
	assert.Nil(t, l4.Tail())
}

// Helper function to collect values from the list for easier assertion
func collectValues[T any](l *linkedlist.List[T]) []T {
	var values []T
	for v := range l.Values() {
		values = append(values, v)
	}
	return values
}
