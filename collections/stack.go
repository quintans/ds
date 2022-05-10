package collections

var _ Stacker[int] = (*Stack[int])(nil)

type element[T any] struct {
	data T
	next *element[T]
}

type Stack[T any] struct {
	head *element[T]
	size int
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (l *Stack[T]) Push(data T) {
	element := &element[T]{
		data: data,
		next: l.head,
	}
	l.head = element
	l.size++
}

func (l *Stack[T]) Pop() (T, bool) {
	if l.head == nil {
		var zero T
		return zero, false
	}

	r := l.head.data
	l.head = l.head.next
	l.size--

	return r, true
}

func (l *Stack[T]) Peek() (T, bool) {
	if l.head == nil {
		var zero T
		return zero, false
	}

	return l.head.data, true
}

func (l *Stack[T]) Size() int {
	return l.size
}

func (l *Stack[T]) Clear() {
	l.head = nil
	l.size = 0
}
