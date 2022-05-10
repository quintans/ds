package collections

import (
	"fmt"
	"strings"
)

const default_hashset_capacity = 16

// check if it implements Collection interface
var _ Collection[string] = (*HashSet[string])(nil)

type HashSet[K any] struct {
	entries         Map[K, struct{}]
	initialCapacity int
	equals          func(a, b K) bool
	hashCode        func(a K) int
}

func NewHashSet[K any](equals func(a, b K) bool, hashCode func(a K) int) *HashSet[K] {
	return NewHashSetWithCapacity[K](equals, hashCode, default_hashset_capacity)
}

func NewHashSetWithCapacity[K any](equals func(a, b K) bool, hashCode func(a K) int, capacity int) *HashSet[K] {
	s := &HashSet[K]{
		initialCapacity: capacity,
		equals:          equals,
		hashCode:        hashCode,
	}
	s.Clear()
	return s
}

func (h *HashSet[K]) Clear() {
	h.entries = NewHashMapWithCapacity[K, struct{}](
		h.equals,
		h.hashCode,
		h.initialCapacity,
	)
}

func (h *HashSet[K]) Size() int {
	return h.entries.Size()
}

func (h *HashSet[K]) Contains(key K) bool {
	return h.entries.ContainsKey(key)
}

func (h *HashSet[K]) Add(keys ...K) {
	for _, key := range keys {
		_, ok := h.entries.Get(key)
		if !ok {
			h.entries.Put(key, struct{}{})
		}
	}
}

func (h *HashSet[K]) AddAll(c Collection[K]) {
	c.ForEach(func(_ int, k K) {
		h.Add(k)
	})
}

func (h *HashSet[K]) Delete(key K) bool {
	_, ok := h.entries.Delete(key)
	return ok
}

func (h *HashSet[K]) ForEach(fn func(int, K)) {
	index := 0
	h.entries.ForEach(func(k K, s struct{}) {
		fn(index, k)
	})
}

func (h *HashSet[K]) ToSlice() []K {
	data := make([]K, 0, h.entries.Size())
	h.entries.ForEach(func(k K, s struct{}) {
		data = append(data, k)
	})
	return data
}

func (h *HashSet[K]) String() string {
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

func (h *HashSet[K]) Clone() *HashSet[K] {
	m := NewHashMapWithCapacity[K, struct{}](
		h.equals,
		h.hashCode,
		h.initialCapacity,
	)
	h.entries.ForEach(func(k K, s struct{}) {
		m.Put(k, s)
	})
	return &HashSet[K]{
		entries:         m,
		initialCapacity: h.initialCapacity,
	}
}

// returns a function that in every call return the next value
// if no value was retrived
func (this *HashSet[K]) Iterator() Iterator[K] {
	return &HashSetIterator[K]{this.entries.Iterator()}
}

type HashSetIterator[K any] struct {
	iterator Iterator[KeyValue[K, struct{}]]
}

func (it *HashSetIterator[K]) HasNext() bool {
	return it.iterator.HasNext()
}

func (it *HashSetIterator[K]) Next() K {
	kv := it.iterator.Next()
	return kv.Key
}

func (this *HashSetIterator[K]) Remove() {
	this.iterator.Remove()
}
