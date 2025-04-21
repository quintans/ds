package linkedmap

import (
	"fmt"
	"iter"
	"maps"
	"strings"

	"github.com/quintans/ds/collections/linkedlist"
)

const defaultCapacity = 16

type Option[K comparable, V any] func(*Map[K, V])

func WithCapacity[K comparable, V any](capacity int) Option[K, V] {
	return func(l *Map[K, V]) {
		l.initialCapacity = capacity
	}
}

type Map[K comparable, V any] struct {
	keyOrder        *linkedlist.List[K]
	entries         map[K]*entry[V]
	initialCapacity int
}

type entry[V any] struct {
	value           V
	keyOrderRemover func()
}

func New[K comparable, V any](options ...Option[K, V]) *Map[K, V] {
	m := &Map[K, V]{
		initialCapacity: defaultCapacity,
	}

	for _, opt := range options {
		opt(m)
	}

	m.Clear()
	return m
}

func (m *Map[K, V]) Clear() {
	m.keyOrder = linkedlist.New[K]()
	m.entries = make(map[K]*entry[V], defaultCapacity)
}

func (m *Map[K, V]) Size() int {
	return len(m.entries)
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	v, ok := m.entries[key]
	if !ok {
		var zero V
		return zero, false
	}
	return v.value, ok
}

func (m *Map[K, V]) Put(key K, value V) (V, bool) {
	var old V
	e, ok := m.entries[key]
	if ok {
		old = e.value
		e.value = value
	} else {
		m.keyOrder.Add(key)
		c := m.keyOrder.Tail()
		e = &entry[V]{
			value:           value,
			keyOrderRemover: c.Remove,
		}
		m.entries[key] = e
	}

	return old, ok
}

func (m *Map[K, V]) ContainsKey(key K) bool {
	_, ok := m.entries[key]
	return ok
}

func (l *Map[K, V]) Delete(key K) (V, bool) {
	old, ok := l.entries[key]
	if !ok {
		var zero V
		return zero, false
	}

	old.keyOrderRemover()
	delete(l.entries, key)
	return old.value, ok
}

func (l *Map[K, V]) Entries() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k := range l.keyOrder.Values() {
			v, ok := l.entries[k]
			if !ok {
				continue
			}
			if !yield(k, v.value) {
				return
			}
		}
	}
}

func (l *Map[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range l.keyOrder.Values() {
			if !yield(k) {
				return
			}
		}
	}
}

func (l *Map[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for k := range l.keyOrder.Values() {
			v, ok := l.entries[k]
			if !ok {
				continue
			}
			if !yield(v.value) {
				return
			}
		}
	}
}

func (l *Map[K, V]) Clone() *Map[K, V] {
	return &Map[K, V]{
		keyOrder:        l.keyOrder.Clone(),
		entries:         maps.Clone(l.entries),
		initialCapacity: l.initialCapacity,
	}
}

// String returns a string representation of container
func (m *Map[K, V]) String() string {
	sb := strings.Builder{}

	sb.WriteString("LinkedHashMap\nmap[")
	c := 0
	for k := range m.keyOrder.Values() {
		v, ok := m.entries[k]
		if !ok {
			continue
		}
		if c > 0 {
			sb.WriteString(" ")
		}
		c++
		sb.WriteString(fmt.Sprintf("%v:%v", k, v))
	}
	sb.WriteString("]")
	return sb.String()
}
