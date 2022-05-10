package collections

import (
	"errors"
	"fmt"
)

type Element[T any] struct {
	Value      T
	next, prev *Element[T]
	list       *DoublyLinkedList[T]
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

var _ List[string] = (*DoublyLinkedList[string])(nil)

type DoublyLinkedList[T any] struct {
	size int
	// sentinel
	root   Element[T]
	equals func(a, b T) bool
}

func NewDoublyLinkedList[T any](cmp func(a, b T) bool) *DoublyLinkedList[T] {
	l := &DoublyLinkedList[T]{
		equals: cmp,
	}
	l.Clear()
	return l
}

func (l *DoublyLinkedList[T]) Head() *Element[T] {
	if l.size == 0 {
		return nil
	}
	return l.root.next
}

func (l *DoublyLinkedList[T]) Tail() *Element[T] {
	if l.size == 0 {
		return nil
	}
	return l.root.prev
}

// Clear empty this linked list, O(1)
func (l *DoublyLinkedList[T]) Clear() {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.size = 0
}

// Size returns the size of this linked list
func (l *DoublyLinkedList[T]) Size() int {
	return l.size
}

// Add adds element to the tail of the linked list, O(1)
func (l *DoublyLinkedList[T]) AddAll(c Collection[T]) {
	c.ForEach(func(_ int, b T) {
		l.Add(b)
	})
}

// Add adds element to the tail of the linked list, O(1)
func (l *DoublyLinkedList[T]) Add(elems ...T) {
	for _, elem := range elems {
		l.AddLast(elem)
	}
}

// AddLast adds elements to the tail of the linked list, O(1)
func (l *DoublyLinkedList[T]) AddLast(data T) {
	l.insert(l.root.prev, &Element[T]{Value: data, list: l})
}

// AddFirst adds elements to the beginning (head) of this linked list, O(1)
func (l *DoublyLinkedList[T]) AddFirst(data T) {
	l.insert(&l.root, &Element[T]{Value: data, list: l})
}

func (l *DoublyLinkedList[T]) insert(at, e *Element[T]) {
	e.next = at.next
	e.prev = at

	e.prev.next = e
	e.next.prev = e

	e.list = l

	l.size++
}

// MoveToLast moves element to the tail of the linked list, O(1)
func (l *DoublyLinkedList[T]) MoveToLast(e *Element[T]) {
	if l.root.prev == e {
		return
	}

	l.move(l.root.prev, e)
}

// MoveToFirst moves element to the beginning (head) of this linked list, O(1)
func (l *DoublyLinkedList[T]) MoveToFirst(e *Element[T]) {
	l.move(&l.root, e)
}

func (l *DoublyLinkedList[T]) move(at, e *Element[T]) {
	l.cut(e)
	l.insert(at, e)
}

func (l *DoublyLinkedList[T]) cut(e *Element[T]) T {
	// removes element
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil
	e.prev = nil
	e.list = nil
	l.size--
	return e.Value
}

// AddAt adds an element at a specified index
func (l *DoublyLinkedList[T]) AddAt(index int, data T) error {
	at, err := l.findElementByIndex(index)
	if err != nil {
		return err
	}

	l.insert(at, &Element[T]{Value: data, list: l})

	return nil
}

// AddAt adds an element at a specified index
func (l *DoublyLinkedList[T]) Set(index int, data T) error {
	elem, err := l.findElementByIndex(index)
	if err != nil {
		return err
	}
	elem.Value = data

	return nil
}

func (l *DoublyLinkedList[T]) Get(index int) (T, error) {
	var zero T
	elem, err := l.findElementByIndex(index)
	if err != nil {
		return zero, err
	}
	return elem.Value, nil
}

func (l *DoublyLinkedList[T]) findElementByIndex(index int) (*Element[T], error) {
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
func (l *DoublyLinkedList[T]) PeekFirst() (T, error) {
	var zero T
	if l.size == 0 {
		return zero, errors.New("empty list")
	}
	return l.root.next.Value, nil
}

// PeekLast checks the value of the last element (tail) if it exists, O(1)
func (l *DoublyLinkedList[T]) PeekLast() (T, error) {
	var zero T
	if l.size == 0 {
		return zero, errors.New("empty list")
	}
	return l.root.prev.Value, nil
}

// RemoveFirst removes the first value at the head of the linked list, O(1)
func (l *DoublyLinkedList[T]) RemoveFirst() (T, error) {
	if l.size == 0 {
		var zero T
		return zero, errors.New("empty list")
	}

	return l.cut(l.root.next), nil
}

// RemoveLast removes the last value at the tail of the linked list, O(1)
func (l *DoublyLinkedList[T]) RemoveLast() (T, error) {
	// Can't remove data from an empty list
	if l.size == 0 {
		var zero T
		return zero, errors.New("empty list")
	}

	return l.cut(l.root.prev), nil
}

// DeleteAt removes a element at a particular index, O(n)
func (l *DoublyLinkedList[T]) DeleteAt(index int) (T, error) {
	elem, err := l.findElementByIndex(index)
	if err != nil {
		var zero T
		return zero, err
	}

	return l.cut(elem), nil
}

// Delete removes a particular value in the linked list, O(n)
func (l *DoublyLinkedList[T]) Delete(obj T) bool {
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
func (l *DoublyLinkedList[T]) IndexOf(obj T) int {
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
func (l *DoublyLinkedList[T]) Contains(obj T) bool {
	return l.IndexOf(obj) != -1
}

func (l *DoublyLinkedList[T]) ForEach(fn func(int, T)) {
	temp := &l.root
	for i := 0; i < l.size; i++ {
		temp = temp.next
		fn(i, temp.Value)
	}
}

func (l *DoublyLinkedList[T]) ReplaceAll(fn func(int, T) T) {
	temp := &l.root
	for i := 0; i < l.size; i++ {
		temp = temp.next
		temp.Value = fn(i, temp.Value)
	}
}

func (l *DoublyLinkedList[T]) ToSlice() []T {
	elems := make([]T, 0, l.size)
	temp := &l.root
	for i := 0; i < l.size; i++ {
		temp = temp.next
		elems = append(elems, temp.Value)
	}
	return elems
}

func (l *DoublyLinkedList[T]) Clone() *DoublyLinkedList[T] {
	d := NewDoublyLinkedList(l.equals)
	d.AddAll(d)
	return d
}

func (l *DoublyLinkedList[T]) Iterator() Iterator[T] {
	return &DoublyLinkedListIterator[T]{
		list: l,
		trav: &l.root,
	}
}

type DoublyLinkedListIterator[T any] struct {
	list *DoublyLinkedList[T]
	trav *Element[T]
}

func (i *DoublyLinkedListIterator[T]) HasNext() bool {
	return i.trav.next != &i.list.root
}

func (i *DoublyLinkedListIterator[T]) Next() T {
	data := i.trav.Value
	i.trav = i.trav.next
	return data
}

func (i *DoublyLinkedListIterator[T]) Remove() {
	if i.trav.prev != nil {
		i.list.cut(i.trav.prev)
	}
}
