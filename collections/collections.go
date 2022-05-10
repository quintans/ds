package collections

import (
	"fmt"
)

type Collection[T any] interface {
	Size() int
	Clear()
	Contains(value T) bool
	AddAll(Collection[T])
	Add(...T)
	Delete(key T) bool

	Iterator() Iterator[T]
	ForEach(func(int, T))
	ToSlice() []T
}

type List[T any] interface {
	Collection[T]

	Get(pos int) (T, error)
	Set(pos int, value T) error
	IndexOf(value T) int
	DeleteAt(pos int) (T, error)
	ReplaceAll(func(int, T) T)
}

type Iterator[T any] interface {
	HasNext() bool
	Next() T
	Remove()
}

type KeyValue[K, V any] struct {
	Key   K
	Value V
}

func (kv KeyValue[K, V]) String() string {
	return fmt.Sprintf("{%v: %v}", kv.Key, kv.Value)
}

type Map[K, V any] interface {
	Size() int
	Clear()
	Delete(K) (V, bool)

	// Put adds or replace a value.
	// The old value is returned and true if it existed
	Put(K, V) (V, bool)
	Get(K) (V, bool)
	ContainsKey(K) bool

	Entries() []KeyValue[K, V]
	Values() []V
	ForEach(func(K, V))
	ReplaceAll(func(K, V) V)
	Iterator() Iterator[KeyValue[K, V]]
}

type Queuer[T any] interface {
	Size() int
	Clear()
	Offer(value T) (T, bool)
	Poll() (T, bool)
	Peek() (T, bool)
}

type Stacker[T any] interface {
	Size() int
	Clear()
	Push(value T)
	Pop() (T, bool)
	Peek() (T, bool)
}
