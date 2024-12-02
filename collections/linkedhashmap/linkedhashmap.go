package linkedhashmap

import (
	"github.com/quintans/dstruct/collections"
	"github.com/quintans/dstruct/collections/hashmap"
	"github.com/quintans/dstruct/collections/linkedlist"
)

const defaultCapacity = 16

type Option[K, V any] func(*Map[K, V])

func WithCapacity[K, V any](capacity int) Option[K, V] {
	return func(l *Map[K, V]) {
		l.initialCapacity = capacity
	}
}

type Map[K, V any] struct {
	keyOrder        *linkedlist.List[K]
	entries         *hashmap.Map[K, entry[V]]
	initialCapacity int
	equals          func(a, b K) bool
	hashCode        func(a K) int
}

type entry[V any] struct {
	value           V
	keyOrderRemover func()
}

// check if it implements Map interface
var _ collections.Map[string, any] = (*Map[string, any])(nil)

func New[K comparable, V any]() *Map[K, V] {
	return NewFunc[K, V](collections.Equals[K], collections.HashCode[K])
}

func NewFunc[K, V any](cmp func(a, b K) bool, hash func(a K) int, options ...Option[K, V]) *Map[K, V] {
	m := &Map[K, V]{
		initialCapacity: defaultCapacity,
		equals:          cmp,
		hashCode:        hash,
	}

	for _, opt := range options {
		opt(m)
	}

	m.Clear()
	return m
}

func (l *Map[K, V]) Clear() {
	l.keyOrder = linkedlist.NewCmp[K](l.equals)
	l.entries = hashmap.NewFunc[K, entry[V]](l.equals, l.hashCode, hashmap.WithCapacity[K, entry[V]](l.initialCapacity))
}

func (l *Map[K, V]) Size() int {
	return l.entries.Size()
}

func (l *Map[K, V]) Get(key K) (V, bool) {
	e, ok := l.entries.Get(key)
	return e.value, ok
}

func (l *Map[K, V]) Put(key K, value V) (V, bool) {
	var old V
	e, ok := l.entries.Get(key)
	if ok {
		old = e.value
		e.value = value
	} else {
		l.keyOrder.Add(key)
		c := l.keyOrder.Tail()
		e = entry[V]{
			value:           value,
			keyOrderRemover: c.Remove,
		}
		l.entries.Put(key, e)
	}

	return old, ok
}

func (h *Map[K, V]) ContainsKey(key K) bool {
	_, ok := h.Get(key)
	return ok
}

func (l *Map[K, V]) Delete(key K) (V, bool) {
	old, deleted := l.entries.Delete(key)
	if deleted {
		old.keyOrderRemover()
	}
	return old.value, deleted
}

func (l *Map[K, V]) Entries() []collections.KV[K, V] {
	data := make([]collections.KV[K, V], 0, l.entries.Size())
	for it := l.Iterator(); it.HasNext(); {
		data = append(data, it.Next())
	}
	return data
}

func (l *Map[K, V]) Values() []V {
	data := make([]V, 0, l.entries.Size())
	for it := l.Iterator(); it.HasNext(); {
		data = append(data, it.Next().Value)
	}
	return data
}

func (l *Map[K, V]) ForEach(fn func(K, V)) {
	for it := l.Iterator(); it.HasNext(); {
		entry := it.Next()
		fn(entry.Key, entry.Value)
	}
}

func (l *Map[K, V]) ReplaceAll(fn func(K, V) V) {
	l.entries.ReplaceAll(func(k K, e entry[V]) entry[V] {
		e.value = fn(k, e.value)
		return e
	})
}

func (l *Map[K, V]) Clone() *Map[K, V] {
	return &Map[K, V]{
		keyOrder:        l.keyOrder.Clone(),
		entries:         l.entries.Clone(),
		initialCapacity: l.initialCapacity,
	}
}

// returns a function that in every call return the next value
// if key is null, no value was retrieved
func (l *Map[K, V]) Iterator() collections.Iterator[collections.KV[K, V]] {
	return &Iterator[K, V]{
		hashmap:  l.entries,
		iterator: l.keyOrder.Iterator(),
	}
}

type Iterator[K, V any] struct {
	hashmap  *hashmap.Map[K, entry[V]]
	iterator collections.Iterator[K]
	lastKey  K
}

func (l *Iterator[K, V]) HasNext() bool {
	return l.iterator.HasNext()
}

func (l *Iterator[K, V]) Next() collections.KV[K, V] {
	k := l.iterator.Next()
	l.lastKey = k
	return l.getEntry(k)
}

func (l *Iterator[K, V]) getEntry(k K) collections.KV[K, V] {
	e, ok := l.hashmap.Get(k)
	if ok {
		return collections.KV[K, V]{Key: k, Value: e.value}
	}

	return collections.KV[K, V]{}
}

func (l *Iterator[K, V]) Remove() {
	l.hashmap.Delete(l.lastKey)
}
