package linkedlist

import (
	"errors"
	"fmt"
	"iter"

	"github.com/quintans/ds/collections"
)

type Element[T any] struct {
	value      T
	next, prev *Element[T]
	list       *List[T]
}

func (e *Element[T]) Value() T {
	return e.value
}

func (e *Element[T]) Next() *Element[T] {
	return e.next
}

func (e *Element[T]) Previous() *Element[T] {
	return e.prev
}

func (e *Element[T]) Remove() {
	e.list.cut(e)
}

type List[T any] struct {
	size   int
	head   *Element[T]
	tail   *Element[T]
	equals func(a, b T) bool
}

func New[T comparable]() *List[T] {
	return NewCmp(collections.Equals[T])
}

func NewCmp[T any](cmp func(a, b T) bool) *List[T] {
	l := &List[T]{
		equals: cmp,
	}
	l.Clear()
	return l
}

func (l *List[T]) Head() *Element[T] {
	return l.head
}

func (l *List[T]) Tail() *Element[T] {
	return l.tail
}

// Clear empty this linked list, O(1)
func (l *List[T]) Clear() {
	l.head = nil
	l.tail = nil
	l.size = 0
}

// Size returns the size of this linked list
func (l *List[T]) Size() int {
	return l.size
}

// AddLast adds elements to the tail of the linked list, O(1)
func (l *List[T]) Add(data T) *Element[T] {
	elem := &Element[T]{value: data, list: l}

	if l.tail == nil {
		l.head = elem
		l.tail = elem
		l.size++
	} else {
		l.insertAfter(l.tail, elem)
	}

	return elem
}

// AddFirst adds elements to the beginning (head) of this linked list, O(1)
func (l *List[T]) AddFirst(data T) *Element[T] {
	elem := &Element[T]{value: data, list: l}
	if l.head == nil {
		l.head = elem
		l.tail = elem
		l.size++
	} else {
		l.insertBefore(l.head, elem)
	}
	return elem
}

func (l *List[T]) insertAfter(at, e *Element[T]) {
	e.prev = at
	e.next = at.next

	if e.next != nil {
		e.next.prev = e
	} else {
		l.tail = e
	}
	at.next = e

	e.list = l

	l.size++
}

func (l *List[T]) insertBefore(at, e *Element[T]) {
	e.next = at
	e.prev = at.prev

	if e.prev != nil {
		e.prev.next = e
	} else {
		l.head = e
	}
	at.prev = e

	e.list = l

	l.size++
}

// MoveToLast moves element to the tail of the linked list, O(1)
func (l *List[T]) MoveToLast(e *Element[T]) {
	l.moveAfter(l.tail, e)
}

// MoveToFirst moves element to the beginning (head) of this linked list, O(1)
func (l *List[T]) MoveToFirst(e *Element[T]) {
	l.moveBefore(l.head, e)
}

func (l *List[T]) moveAfter(at, e *Element[T]) {
	l.cut(e)
	l.insertAfter(at, e)
}

func (l *List[T]) moveBefore(at, e *Element[T]) {
	l.cut(e)
	l.insertBefore(at, e)
}

func (l *List[T]) cut(e *Element[T]) {
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		l.head = e.next
	}

	if e.next != nil {
		e.next.prev = e.prev
	} else {
		l.tail = e.prev
	}

	e.next = nil
	e.prev = nil

	e.list = nil

	l.size--
}

// AddAt adds an element at a specified index
func (l *List[T]) AddAt(index int, data T) error {
	at, err := l.findElementByIndex(index)
	if err != nil {
		return err
	}

	l.insertBefore(at, &Element[T]{value: data, list: l})

	return nil
}

// AddAt adds an element at a specified index
func (l *List[T]) Set(index int, data T) error {
	elem, err := l.findElementByIndex(index)
	if err != nil {
		return err
	}
	elem.value = data

	return nil
}

func (l *List[T]) Get(index int) (T, error) {
	var zero T
	elem, err := l.findElementByIndex(index)
	if err != nil {
		return zero, err
	}
	return elem.value, nil
}

func (l *List[T]) findElementByIndex(index int) (*Element[T], error) {
	if index < 0 || index >= l.size {
		return nil, fmt.Errorf("index out of bounds [0 - %d): %d", l.size, index)
	}

	var temp *Element[T]
	// Search from the front of the list
	if index < l.size/2 {
		temp = l.head
		for i := 0; i != index; i++ {
			temp = temp.next
		}
		// Search from the back of the list
	} else {
		temp = l.tail
		for i := l.size - 1; i != index; i-- {
			temp = temp.prev
		}
	}
	return temp, nil
}

// PeekFirst checks the value of the first element (head) if it exists, O(1)
func (l *List[T]) PeekFirst() (T, error) {
	if l.size == 0 {
		var zero T
		return zero, errors.New("empty list")
	}
	return l.head.value, nil
}

// PeekLast checks the value of the last element (tail) if it exists, O(1)
func (l *List[T]) PeekLast() (T, error) {
	var zero T
	if l.size == 0 {
		return zero, errors.New("empty list")
	}
	return l.tail.value, nil
}

// RemoveFirst removes the first value at the head of the linked list, O(1)
func (l *List[T]) RemoveFirst() (T, error) {
	if l.size == 0 {
		var zero T
		return zero, errors.New("empty list")
	}

	value := l.head.value
	l.cut(l.head)
	return value, nil
}

// RemoveLast removes the last value at the tail of the linked list, O(1)
func (l *List[T]) RemoveLast() (T, error) {
	// Can't remove data from an empty list
	if l.size == 0 {
		var zero T
		return zero, errors.New("empty list")
	}

	value := l.tail.value
	l.cut(l.tail)

	return value, nil
}

// DeleteAt removes a element at a particular index, O(n)
func (l *List[T]) DeleteAt(index int) (T, error) {
	elem, err := l.findElementByIndex(index)
	if err != nil {
		var zero T
		return zero, err
	}

	value := elem.value
	l.cut(elem)

	return value, nil
}

// Delete removes a particular value in the linked list, O(n)
func (l *List[T]) Delete(obj T) bool {
	temp := l.head
	for i := 0; i < l.size; i++ {
		if l.equals(temp.value, obj) {
			l.cut(temp)
			return true
		}
		temp = temp.next
	}
	return false
}

// IndexOf finds the index of a particular value in the linked list, O(n)
func (l *List[T]) IndexOf(obj T) int {
	temp := l.head
	for i := 0; i < l.size; i++ {
		if l.equals(temp.value, obj) {
			return i
		}
		temp = temp.next
	}
	return -1
}

// Check is a value is contained within the linked list
func (l *List[T]) Contains(obj T) bool {
	return l.IndexOf(obj) != -1
}

func (l *List[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		temp := l.head
		for range l.size {
			if !yield(temp.value) {
				return
			}
			temp = temp.next
		}
	}
}

func (l *List[T]) ReplaceAll(fn func(int, T) T) {
	temp := l.head
	for i := 0; i < l.size; i++ {
		temp.value = fn(i, temp.value)
		temp = temp.next
	}
}

func (l *List[T]) Clone() *List[T] {
	d := NewCmp(l.equals)
	for v := range l.Values() {
		d.Add(v)
	}
	return d
}
