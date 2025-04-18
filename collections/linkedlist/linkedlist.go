package linkedlist

import (
	"container/list"
	"errors"
	"fmt"

	"github.com/quintans/ds/collections"
)

type Element[T any] struct {
	Value      T
	next, prev *Element[T]
	list       *List[T]
	x          list.List
}

func (e *Element[T]) Next() *Element[T] {
	if e.next == &e.list.root {
		return nil
	}
	return e.next
}

func (e *Element[T]) Previous() *Element[T] {
	if e.prev == &e.list.root {
		return nil
	}
	return e.prev
}

func (e *Element[T]) Remove() {
	e.list.cut(e)
}

var _ collections.List[string] = (*List[string])(nil)

type List[T any] struct {
	size int
	// sentinel
	root   Element[T]
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
	if l.size == 0 {
		return nil
	}
	return l.root.next
}

func (l *List[T]) Tail() *Element[T] {
	if l.size == 0 {
		return nil
	}
	return l.root.prev
}

// Clear empty this linked list, O(1)
func (l *List[T]) Clear() {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.size = 0
}

// Size returns the size of this linked list
func (l *List[T]) Size() int {
	return l.size
}

// Add adds element to the tail of the linked list, O(1)
func (l *List[T]) AddAll(c collections.Collection[T]) {
	c.ForEach(func(_ int, b T) {
		l.Add(b)
	})
}

// Add adds element to the tail of the linked list, O(1)
func (l *List[T]) Add(elems ...T) {
	for _, elem := range elems {
		l.AddLast(elem)
	}
}

// AddLast adds elements to the tail of the linked list, O(1)
func (l *List[T]) AddLast(data T) *Element[T] {
	elem := &Element[T]{Value: data, list: l}
	l.insert(l.root.prev, elem)
	return elem
}

// AddFirst adds elements to the beginning (head) of this linked list, O(1)
func (l *List[T]) AddFirst(data T) *Element[T] {
	elem := &Element[T]{Value: data, list: l}
	l.insert(&l.root, elem)
	return elem
}

func (l *List[T]) insert(at, e *Element[T]) {
	e.next = at.next
	e.prev = at

	e.prev.next = e
	e.next.prev = e

	e.list = l

	l.size++
}

// MoveToLast moves element to the tail of the linked list, O(1)
func (l *List[T]) MoveToLast(e *Element[T]) {
	if l.root.prev == e {
		return
	}

	l.move(l.root.prev, e)
}

// MoveToFirst moves element to the beginning (head) of this linked list, O(1)
func (l *List[T]) MoveToFirst(e *Element[T]) {
	l.move(&l.root, e)
}

func (l *List[T]) move(at, e *Element[T]) {
	l.cut(e)
	l.insert(at, e)
}

func (l *List[T]) cut(e *Element[T]) T {
	if e.list == l {
		// removes element
		e.prev.next = e.next
		e.next.prev = e.prev
		e.next = nil
		e.prev = nil
		e.list = nil
		l.size--
	}
	return e.Value
}

// AddAt adds an element at a specified index
func (l *List[T]) AddAt(index int, data T) error {
	at, err := l.findElementByIndex(index)
	if err != nil {
		return err
	}

	l.insert(at, &Element[T]{Value: data, list: l})

	return nil
}

// AddAt adds an element at a specified index
func (l *List[T]) Set(index int, data T) error {
	elem, err := l.findElementByIndex(index)
	if err != nil {
		return err
	}
	elem.Value = data

	return nil
}

func (l *List[T]) Get(index int) (T, error) {
	var zero T
	elem, err := l.findElementByIndex(index)
	if err != nil {
		return zero, err
	}
	return elem.Value, nil
}

func (l *List[T]) findElementByIndex(index int) (*Element[T], error) {
	if index < 0 || index >= l.size {
		return nil, fmt.Errorf("index out of bounds [0 - %d): %d", l.size, index)
	}

	var temp *Element[T]
	// Search from the front of the list
	if index < l.size/2 {
		temp = l.root.next
		for i := 0; i != index; i++ {
			temp = temp.next
		}
		// Search from the back of the list
	} else {
		temp = l.root.prev
		for i := l.size - 1; i != index; i-- {
			temp = temp.prev
		}
	}
	return temp, nil
}

// PeekFirst checks the value of the first element (head) if it exists, O(1)
func (l *List[T]) PeekFirst() (T, error) {
	var zero T
	if l.size == 0 {
		return zero, errors.New("empty list")
	}
	return l.root.next.Value, nil
}

// PeekLast checks the value of the last element (tail) if it exists, O(1)
func (l *List[T]) PeekLast() (T, error) {
	var zero T
	if l.size == 0 {
		return zero, errors.New("empty list")
	}
	return l.root.prev.Value, nil
}

// RemoveFirst removes the first value at the head of the linked list, O(1)
func (l *List[T]) RemoveFirst() (T, error) {
	if l.size == 0 {
		var zero T
		return zero, errors.New("empty list")
	}

	return l.cut(l.root.next), nil
}

// RemoveLast removes the last value at the tail of the linked list, O(1)
func (l *List[T]) RemoveLast() (T, error) {
	// Can't remove data from an empty list
	if l.size == 0 {
		var zero T
		return zero, errors.New("empty list")
	}

	return l.cut(l.root.prev), nil
}

// DeleteAt removes a element at a particular index, O(n)
func (l *List[T]) DeleteAt(index int) (T, error) {
	elem, err := l.findElementByIndex(index)
	if err != nil {
		var zero T
		return zero, err
	}

	return l.cut(elem), nil
}

// Delete removes a particular value in the linked list, O(n)
func (l *List[T]) Delete(obj T) bool {
	temp := &l.root
	for i := 0; i < l.size; i++ {
		if l.equals(temp.Value, obj) {
			l.cut(temp)
			return true
		}
		temp = temp.next
	}
	return false
}

// IndexOf finds the index of a particular value in the linked list, O(n)
func (l *List[T]) IndexOf(obj T) int {
	temp := &l.root
	for i := 0; i < l.size; i++ {
		temp = temp.next
		if l.equals(temp.Value, obj) {
			return i
		}
	}
	return -1
}

// Check is a value is contained within the linked list
func (l *List[T]) Contains(obj T) bool {
	return l.IndexOf(obj) != -1
}

func (l *List[T]) ForEach(fn func(int, T)) {
	temp := &l.root
	for i := 0; i < l.size; i++ {
		temp = temp.next
		fn(i, temp.Value)
	}
}

func (l *List[T]) ReplaceAll(fn func(int, T) T) {
	temp := &l.root
	for i := 0; i < l.size; i++ {
		temp = temp.next
		temp.Value = fn(i, temp.Value)
	}
}

func (l *List[T]) ToSlice() []T {
	elems := make([]T, 0, l.size)
	temp := &l.root
	for i := 0; i < l.size; i++ {
		temp = temp.next
		elems = append(elems, temp.Value)
	}
	return elems
}

func (l *List[T]) Clone() *List[T] {
	d := NewCmp(l.equals)
	d.AddAll(d)
	return d
}

func (l *List[T]) Iterator() collections.Iterator[T] {
	return &Iterator[T]{
		next: l.Head(),
		cut:  l.cut,
	}
}

type Iterator[T any] struct {
	current *Element[T]
	next    *Element[T]
	cut     func(*Element[T]) T
}

func (i *Iterator[T]) HasNext() bool {
	return i.next != nil
}

func (i *Iterator[T]) Next() T {
	i.current = i.next
	i.next = i.next.Next()
	return i.current.Value
}

func (i *Iterator[T]) Remove() {
	if i.current != nil {
		current := i.current
		i.current = i.current.prev
		i.cut(current)
	}
}
