package collections

import (
	"fmt"
)

type ArrayList[T any] struct {
	elements []T
	capacity int
	equals   func(a, b T) bool
}

// check if it implements List interface
var _ List[string] = (*ArrayList[string])(nil)

func NewArrayList[T any](cmp func(a, b T) bool) *ArrayList[T] {
	return NewArrayListWithCapacity(cmp, 16)
}

func NewArrayListWithCapacity[T any](cmp func(a, b T) bool, capacity int) *ArrayList[T] {
	return &ArrayList[T]{
		elements: make([]T, 0, capacity),
		capacity: capacity,
		equals:   cmp,
	}
}

func (a *ArrayList[T]) Clear() {
	a.elements = make([]T, 0, a.capacity)
}

func (a *ArrayList[T]) ToSlice() []T {
	data := make([]T, len(a.elements), cap(a.elements))
	copy(data, a.elements)
	return data
}

func (a ArrayList[T]) Size() int {
	return len(a.elements)
}

func (a *ArrayList[T]) Get(pos int) (T, error) {
	var zero T
	if pos < 0 || pos >= a.Size() {
		return zero, fmt.Errorf("index out of bounds [0 - %d): %d", len(a.elements), pos)
	}
	return a.elements[pos], nil
}

func (a *ArrayList[T]) Set(pos int, value T) error {
	if pos < 0 || pos >= a.Size() {
		return fmt.Errorf("index out of bounds [0 - %d): %d", len(a.elements), pos)
	}
	a.elements[pos] = value
	return nil
}

func (a *ArrayList[T]) AddAll(c Collection[T]) {
	c.ForEach(func(_ int, t T) {
		a.elements = append(a.elements, t)
	})
}

func (this *ArrayList[T]) Add(data ...T) {
	for _, v := range data {
		this.elements = append(this.elements, v)
	}
}

func (a *ArrayList[T]) IndexOf(value T) int {
	for k, v := range a.elements {
		if a.equals(v, value) {
			return k
		}
	}
	return -1
}

func (a *ArrayList[T]) Contains(value T) bool {
	return a.IndexOf(value) != -1
}

func (a *ArrayList[T]) Delete(value T) bool {
	k := a.IndexOf(value)
	if k == -1 {
		return false
	}
	_, err := a.DeleteAt(k)
	return err == nil
}

func (a *ArrayList[T]) DeleteAt(pos int) (T, error) {
	var zero T
	if pos < 0 || pos >= a.Size() {
		return zero, fmt.Errorf("index out of bounds [0 - %d): %d", len(a.elements), pos)
	}

	data := a.elements[pos]
	// since the slice might have a non-primitive, we have to zero it
	copy(a.elements[pos:], a.elements[pos+1:])
	a.elements[len(a.elements)-1] = zero // zero it
	a.elements = a.elements[:len(a.elements)-1]

	return data, nil
}

func (a *ArrayList[T]) ForEach(fn func(int, T)) {
	for k, v := range a.elements {
		fn(k, v)
	}
}

func (a *ArrayList[T]) ReplaceAll(fn func(int, T) T) {
	for k, v := range a.elements {
		a.elements[k] = fn(k, v)
	}
}

func (a *ArrayList[T]) String() string {
	return fmt.Sprint(a.elements)
}

func (a *ArrayList[T]) Clone() Collection[T] {
	return &ArrayList[T]{
		elements: a.ToSlice(),
		capacity: a.capacity,
	}
}

// returns a function that in every call return the next value
// and a flag to see if a value was retrived, even if it was nil
func (a *ArrayList[T]) Iterator() Iterator[T] {
	return &ArrayListIterator[T]{list: a, pos: -1}
}

type ArrayListIterator[T any] struct {
	list *ArrayList[T]
	pos  int
}

func (a *ArrayListIterator[T]) HasNext() bool {
	return a.pos+1 < a.list.Size()
}

func (a *ArrayListIterator[T]) Next() T {
	if a.pos < a.list.Size() {
		a.pos++
		return a.list.elements[a.pos]
	}

	return Zero[T]()
}

func (a *ArrayListIterator[T]) Remove() {
	if a.pos >= 0 && a.pos <= a.list.Size() {
		a.list.DeleteAt(a.pos - 1)
		// since this position was removed steps to the previous position
		a.pos--
	}
}
