package arraylist

import (
	"fmt"

	"github.com/quintans/ds/collections"
)

const defaultCapacity = 16

type Option[T any] func(*List[T])

func WithCapacity[T any](capacity int) Option[T] {
	return func(l *List[T]) {
		l.initialCapacity = capacity
	}
}

type List[T any] struct {
	elements        []T
	initialCapacity int
	equals          func(a, b T) bool
}

// check if it implements List interface
var _ collections.List[string] = (*List[string])(nil)

func New[T comparable](options ...Option[T]) *List[T] {
	return NewCmp(collections.Equals[T], options...)
}

func NewCmp[T any](cmp func(a, b T) bool, options ...Option[T]) *List[T] {
	l := &List[T]{
		elements:        make([]T, 0, defaultCapacity),
		initialCapacity: defaultCapacity,
		equals:          cmp,
	}

	for _, opt := range options {
		opt(l)
	}

	return l
}

func (a *List[T]) Clear() {
	a.elements = make([]T, 0, a.initialCapacity)
}

func (a *List[T]) ToSlice() []T {
	data := make([]T, len(a.elements), cap(a.elements))
	copy(data, a.elements)
	return data
}

func (a List[T]) Size() int {
	return len(a.elements)
}

func (a *List[T]) Get(pos int) (T, error) {
	var zero T
	if pos < 0 || pos >= a.Size() {
		return zero, fmt.Errorf("index out of bounds [0 - %d): %d", len(a.elements), pos)
	}
	return a.elements[pos], nil
}

func (a *List[T]) Set(pos int, value T) error {
	if pos < 0 || pos >= a.Size() {
		return fmt.Errorf("index out of bounds [0 - %d): %d", len(a.elements), pos)
	}
	a.elements[pos] = value
	return nil
}

func (a *List[T]) AddAll(c collections.Collection[T]) {
	c.ForEach(func(_ int, t T) {
		a.elements = append(a.elements, t)
	})
}

func (this *List[T]) Add(data ...T) {
	for _, v := range data {
		this.elements = append(this.elements, v)
	}
}

func (a *List[T]) IndexOf(value T) int {
	for k, v := range a.elements {
		if a.equals(v, value) {
			return k
		}
	}
	return -1
}

func (a *List[T]) Contains(value T) bool {
	return a.IndexOf(value) != -1
}

func (a *List[T]) Delete(value T) bool {
	k := a.IndexOf(value)
	if k == -1 {
		return false
	}
	_, err := a.DeleteAt(k)
	return err == nil
}

func (a *List[T]) DeleteAt(pos int) (T, error) {
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

func (a *List[T]) ForEach(fn func(int, T)) {
	for k, v := range a.elements {
		fn(k, v)
	}
}

func (a *List[T]) ReplaceAll(fn func(int, T) T) {
	for k, v := range a.elements {
		a.elements[k] = fn(k, v)
	}
}

func (a *List[T]) String() string {
	return fmt.Sprint(a.elements)
}

func (a *List[T]) Clone() collections.Collection[T] {
	return &List[T]{
		elements:        a.ToSlice(),
		initialCapacity: a.initialCapacity,
	}
}

// returns a function that in every call return the next value
// and a flag to see if a value was retrived, even if it was nil
func (a *List[T]) Iterator() collections.Iterator[T] {
	return &Iterator[T]{list: a, pos: -1}
}

type Iterator[T any] struct {
	list *List[T]
	pos  int
}

func (a *Iterator[T]) HasNext() bool {
	return a.pos+1 < a.list.Size()
}

func (a *Iterator[T]) Next() T {
	if a.pos < a.list.Size() {
		a.pos++
		return a.list.elements[a.pos]
	}

	var zero T
	return zero
}

func (a *Iterator[T]) Remove() {
	if a.pos >= 0 && a.pos <= a.list.Size() {
		a.list.DeleteAt(a.pos - 1)
		// since this position was removed steps to the previous position
		a.pos--
	}
}
