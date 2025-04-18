package linkedset

import (
	"fmt"
	"iter"
	"strings"

	"github.com/quintans/ds/collections/linkedmap"
)

const defaultCapacity = 16

type Options struct {
	Capacity int
}

type Option func(*Options)

func WithCapacity(capacity int) Option {
	return func(l *Options) {
		l.Capacity = capacity
	}
}

type Set[K comparable, V any] struct {
	lhm     *linkedmap.Map[K, V]
	keyFunc func(V) K
}

func New[K comparable](options ...Option) *Set[K, K] {
	return NewFunc(func(k K) K { return k }, options...)
}

func NewFunc[K comparable, V any](keyFunc func(V) K, options ...Option) *Set[K, V] {
	opts := &Options{
		Capacity: defaultCapacity,
	}

	for _, opt := range options {
		opt(opts)
	}

	return &Set[K, V]{
		lhm: linkedmap.New(
			linkedmap.WithCapacity[K, V](opts.Capacity),
		),
		keyFunc: keyFunc,
	}
}

func (h *Set[K, V]) Clear() {
	h.lhm.Clear()
}

func (h *Set[K, V]) Size() int {
	return h.lhm.Size()
}

func (h *Set[K, V]) Contains(value V) bool {
	key := h.keyFunc(value)
	return h.lhm.ContainsKey(key)
}

func (h *Set[K, V]) Add(values ...V) {
	for _, v := range values {
		k := h.keyFunc(v)
		h.lhm.Put(k, v)
	}
}

func (h *Set[K, V]) Delete(value V) bool {
	key := h.keyFunc(value)
	_, ok := h.lhm.Delete(key)

	return ok
}

func (h *Set[K, V]) Values() iter.Seq[V] {
	return h.lhm.Values()
}

func (h *Set[K, V]) String() string {
	var s strings.Builder
	s.WriteString("[")
	counter := 0
	for v := range h.lhm.Values() {
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
		lhm:     h.lhm.Clone(),
		keyFunc: h.keyFunc,
	}
}
