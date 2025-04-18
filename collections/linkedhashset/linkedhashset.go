package linkedhashset

import (
	"github.com/quintans/ds/collections"
	"github.com/quintans/ds/collections/hashset"
)

const defaultCapacity = 16

type Option[K any] func(*Set[K])

func WithCapacity[K any](capacity int) Option[K] {
	return func(l *Set[K]) {
		l.initialCapacity = capacity
	}
}

type Set[K any] struct {
	*hashset.Set[K]
	initialCapacity int
	equals          func(a, b K) bool
	hashCode        func(a K) int
}

// check if it implements Collection interface
var _ collections.Collection[string] = (*Set[string])(nil)

func New[K comparable](options ...Option[K]) *Set[K] {
	return NewFunc[K](collections.Equals[K], collections.HashCode[K], options...)
}

func NewFunc[K any](cmp func(a, b K) bool, hash func(a K) int, options ...Option[K]) *Set[K] {
	s := &Set[K]{
		initialCapacity: defaultCapacity,
		equals:          cmp,
		hashCode:        hash,
	}

	for _, opt := range options {
		opt(s)
	}

	s.Set = hashset.NewFunc[K](cmp, hash, hashset.WithCapacity[K](s.initialCapacity))
	s.Clear()
	return s
}
