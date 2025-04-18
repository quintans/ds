package hashset

import (
	"fmt"
	"strings"

	"github.com/quintans/ds/collections"
	"github.com/quintans/ds/collections/hashmap"
)

const defaultCapacity = 16

type Option[K any] func(*Set[K])

func WithCapacity[K any](capacity int) Option[K] {
	return func(l *Set[K]) {
		l.initialCapacity = capacity
	}
}

// check if it implements Collection interface
var _ collections.Collection[string] = (*Set[string])(nil)

type Set[K any] struct {
	entries         collections.Map[K, struct{}]
	initialCapacity int
	equals          func(a, b K) bool
	hashCode        func(a K) int
}

func New[K comparable](options ...Option[K]) *Set[K] {
	return NewFunc[K](collections.Equals[K], collections.HashCode[K], options...)
}

func NewFunc[K any](equals func(a, b K) bool, hashCode func(a K) int, options ...Option[K]) *Set[K] {
	s := &Set[K]{
		initialCapacity: defaultCapacity,
		equals:          equals,
		hashCode:        hashCode,
	}

	for _, opt := range options {
		opt(s)
	}

	s.Clear()
	return s
}

func (h *Set[K]) Clear() {
	h.entries = hashmap.NewFunc(h.equals, h.hashCode, hashmap.WithCapacity[K, struct{}](h.initialCapacity))
}

func (h *Set[K]) Size() int {
	return h.entries.Size()
}

func (h *Set[K]) Contains(key K) bool {
	return h.entries.ContainsKey(key)
}

func (h *Set[K]) Add(keys ...K) {
	for _, key := range keys {
		_, ok := h.entries.Get(key)
		if !ok {
			h.entries.Put(key, struct{}{})
		}
	}
}

func (h *Set[K]) AddAll(c collections.Collection[K]) {
	c.ForEach(func(_ int, k K) {
		h.Add(k)
	})
}

func (h *Set[K]) Delete(key K) bool {
	_, ok := h.entries.Delete(key)
	return ok
}

func (h *Set[K]) ForEach(fn func(int, K)) {
	index := 0
	h.entries.ForEach(func(k K, s struct{}) {
		fn(index, k)
	})
}

func (h *Set[K]) ToSlice() []K {
	data := make([]K, 0, h.entries.Size())
	h.entries.ForEach(func(k K, s struct{}) {
		data = append(data, k)
	})
	return data
}

func (h *Set[K]) String() string {
	var s strings.Builder
	s.WriteString("[")
	counter := 0
	for it := h.Iterator(); it.HasNext(); {
		if counter > 1 {
			s.WriteString(", ")
		}
		s.WriteString(fmt.Sprintf("%+v", it.Next()))
		counter++
	}
	s.WriteString("]")

	return s.String()
}

func (h *Set[K]) Clone() *Set[K] {
	m := hashmap.NewFunc(h.equals, h.hashCode, hashmap.WithCapacity[K, struct{}](h.initialCapacity))
	h.entries.ForEach(func(k K, s struct{}) {
		m.Put(k, s)
	})
	return &Set[K]{
		entries:         m,
		initialCapacity: h.initialCapacity,
	}
}

// returns a function that in every call return the next value
// if no value was retrived
func (this *Set[K]) Iterator() collections.Iterator[K] {
	return &Iterator[K]{this.entries.Iterator()}
}

type Iterator[K any] struct {
	iterator collections.Iterator[collections.KV[K, struct{}]]
}

func (it *Iterator[K]) HasNext() bool {
	return it.iterator.HasNext()
}

func (it *Iterator[K]) Next() K {
	kv := it.iterator.Next()
	return kv.Key
}

func (this *Iterator[K]) Remove() {
	this.iterator.Remove()
}
