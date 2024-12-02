package queue

import "github.com/quintans/dstruct/collections"

var _ collections.Queuer[int] = (*Queue[int])(nil)

type Option[K any] func(*Queue[K])

func WithCapacity[K, V any](capacity int) Option[K] {
	return func(l *Queue[K]) {
		l.capacity = capacity
	}
}

type Queue[T any] struct {
	head     *item[T]
	tail     *item[T]
	size     int
	capacity int
}

type item[T any] struct {
	next  *item[T]
	value T
}

// New creates a circular queue.
// The Idea is to have a FIFO with a windowing (circular) feature.
// If the max size is reached, the oldest element will be removed.
// If capacity is 0 it will add until memory is exhausted
func New[T any](options ...Option[T]) *Queue[T] {
	q := &Queue[T]{
		capacity: 0,
	}

	for _, opt := range options {
		opt(q)
	}

	return q
}

func (f *Queue[T]) Size() int {
	return f.size
}

// Clear resets the queue.
func (f *Queue[T]) Clear() {
	f.head = nil
	f.tail = nil
	f.size = 0
}

// Offer adds an element to the head of the fifo.
// If the capacity was exceeded returns the element that had to be pushed out, otherwise returns the zero value.
// The boolean value indicates if a value was pushed out.
func (f *Queue[T]) Offer(value T) (T, bool) {
	var old T
	var hasOld bool
	// if capacity == 0 it will add until memory is exhausted
	if f.capacity > 0 && f.size == f.capacity {
		old, hasOld = f.pop()
	}
	// adds new element
	e := &item[T]{value: value}
	if f.head != nil {
		f.head.next = e
	}
	f.head = e

	if f.tail == nil {
		f.tail = e
	}

	f.size++

	return old, hasOld
}

func (f *Queue[T]) pop() (T, bool) {
	var value T
	if f.tail == nil {
		return value, false
	}

	value = f.tail.value
	f.tail = f.tail.next
	f.size--
	return value, true
}

// Poll returns the tail element removing it.
// The boolean value indicates if a value was found
func (f *Queue[T]) Poll() (T, bool) {
	return f.pop()
}

// Peek returns the tail element without removing it.
// The boolean value indicates if a value was found
func (f *Queue[T]) Peek() (T, bool) {
	if f.tail != nil {
		return f.tail.value, true
	}
	var zero T
	return zero, false
}
