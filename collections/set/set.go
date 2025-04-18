package set

import (
	"fmt"
	"iter"
	"maps"
	"strings"
)

const defaultCapacity = 16

type Option[K comparable, V any] func(*Set[K, V])

func WithCapacity[K comparable, V any](capacity int) Option[K, V] {
	return func(l *Set[K, V]) {
		l.initialCapacity = capacity
	}
}

type Set[K comparable, V any] struct {
	entries         map[K]V
	initialCapacity int
	keyFunc         func(V) K
}

func New[K comparable](options ...Option[K, K]) *Set[K, K] {
	return NewFunc(func(k K) K { return k }, options...)
}

func NewFunc[K comparable, V any](keyFunc func(V) K, options ...Option[K, V]) *Set[K, V] {
	s := &Set[K, V]{
		entries:         make(map[K]V, defaultCapacity),
		initialCapacity: defaultCapacity,
		keyFunc:         keyFunc,
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}

func (h *Set[K, V]) Clear() {
	h.entries = make(map[K]V, h.initialCapacity)
}

func (h *Set[K, V]) Size() int {
	return len(h.entries)
}

func (h *Set[K, V]) Contains(value V) bool {
	key := h.keyFunc(value)
	_, ok := h.entries[key]
	return ok
}

func (h *Set[K, V]) Add(values ...V) {
	for _, v := range values {
		k := h.keyFunc(v)
		h.entries[k] = v
	}
}

func (h *Set[K, V]) Delete(value V) bool {
	key := h.keyFunc(value)
	_, ok := h.entries[key]
	if ok {
		delete(h.entries, key)
	}
	return ok
}

func (h *Set[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range h.entries {
			if !yield(v) {
				return
			}
		}
	}
}

func (h *Set[K, V]) String() string {
	var s strings.Builder
	s.WriteString("[")
	counter := 0
	for _, v := range h.entries {
		if counter > 1 {
			s.WriteString(", ")
		}
		s.WriteString(fmt.Sprintf("%+v", v))
		counter++
	}
	s.WriteString("]")

	return s.String()
}

func (h *Set[K, V]) Clone() *Set[K, V] {
	return &Set[K, V]{
		entries:         maps.Clone(h.entries),
		initialCapacity: h.initialCapacity,
		keyFunc:         h.keyFunc,
	}
}
