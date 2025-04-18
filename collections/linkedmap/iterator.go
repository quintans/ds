package linkedmap

import "github.com/quintans/ds/collections/linkedlist"

// Iterator gives an iterator over all entries.
// Call Next() to fetch the first element if any.
// Call Prev() to fetch the last element if any.
//
//	it := Iterator()
//	for it.Next() {
//		fmt.Println("Key:", it.Key(), "Value:", it.Value())
//		if it.Key() == "BAD" {
//			it.Remove()
//		}
//	}
func (m *Map[K, V]) Iterator() *Iterator[K, V] {
	return &Iterator[K, V]{
		forward: true,
		lm:      m,
		cursor:  m.list.Head(),
		out:     true,
	}
}

type Iterator[K comparable, V any] struct {
	forward bool
	lm      *Map[K, V]
	cursor  *linkedlist.Element[*Entry[K, V]]
	out     bool
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (i *Iterator[K, V]) Next() bool {
	i.forward = true

	if i.out {
		i.cursor = i.lm.list.Head()
	} else {
		i.cursor = i.cursor.Next()
	}

	i.out = i.cursor == nil
	return !i.out
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
// If Prev() was called for the first time, then it will point the iterator to the last element if it exists.
// Modifies the state of the iterator.
func (i *Iterator[K, V]) Prev() bool {
	i.forward = false

	if i.out {
		i.cursor = i.lm.list.Tail()
	} else {
		i.cursor = i.cursor.Previous()
	}

	i.out = i.cursor == nil
	return !i.out
}

// Reset resets the iterator to its initial state.
// Call Next() to fetch the first element if any.
// Call Prev() to fetch the last element if any.
func (i *Iterator[K, V]) Reset() {
	i.out = true
	i.cursor = nil
}

// First moves the iterator to the first element and returns true if there was a first element in the container.
// If First() returns true, then first element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator
func (i *Iterator[K, V]) First() bool {
	i.Reset()
	return i.Next()
}

// Last moves the iterator to the last element and returns true if there was a last element in the container.
// If Last() returns true, then last element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (i *Iterator[K, V]) Last() bool {
	i.Reset()
	return i.Prev()
}

// Remove remove the entry from the map and moves the iterator
// to the previous position if it is a forward iterator or
// to the next position if it is reverse iterator.
func (i *Iterator[K, V]) Remove() {
	if i.out {
		return
	}

	cursor := i.cursor

	if i.forward {
		i.Prev()
	} else {
		i.Next()
	}
	i.lm.Delete(cursor.Value.key)
}

// IsOut returns if the iterator is pointing to an entry.
// If true, when calling Key() and Value() you will the zero values.
func (i *Iterator[K, V]) IsOut() bool {
	return i.out
}

func (i *Iterator[K, V]) Key() K {
	if i.cursor == nil {
		var zero K
		return zero
	}
	return i.cursor.Value.key
}

func (i *Iterator[K, V]) Value() V {
	if i.cursor == nil {
		var zero V
		return zero
	}
	return i.cursor.Value.value
}
