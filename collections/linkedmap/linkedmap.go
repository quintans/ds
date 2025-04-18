package linkedmap

import (
	"fmt"
	"strings"

	"github.com/quintans/ds/collections/linkedlist"
)

type Map[K comparable, V any] struct {
	entries map[K]*linkedlist.Element[*Entry[K, V]]
	list    *linkedlist.List[*Entry[K, V]]
}

type Entry[K comparable, V any] struct {
	key   K
	value V
}

func (e Entry[K, V]) Key() K {
	return e.key
}

func (e Entry[K, V]) Value() V {
	return e.value
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		entries: map[K]*linkedlist.Element[*Entry[K, V]]{},
		list: linkedlist.NewCmp(func(a, b *Entry[K, V]) bool {
			return a.key == b.key
		}),
	}
}

func (m *Map[K, V]) Size() int {
	return len(m.entries)
}

func (m *Map[K, V]) Set(key K, value V) (val V, existed bool) {
	if entry, existed := m.entries[key]; existed { // If the key exists, it is updated
		oldValue := entry.Value.value
		entry.Value.value = value
		return oldValue, true
	}
	entry := &Entry[K, V]{
		key:   key,
		value: value,
	}
	element := m.list.AddLast(entry) // Add to linked list
	m.entries[key] = element         // Add to map
	return value, false
}

func (m *Map[K, V]) Delete(key K) (val V, existed bool) {
	if entry, exists := m.entries[key]; exists { // If present
		entry.Remove()         // Remove from linked list
		delete(m.entries, key) // Remove from map
		return entry.Value.value, true
	}
	return
}

func (m *Map[K, V]) Get(key K) (val V, existed bool) {
	if entry, existed := m.entries[key]; existed {
		return entry.Value.value, true
	}
	return
}

func (m *Map[K, V]) Range(f func(key K, value V, idx int) bool) {
	idx := 0
	for e := m.list.Head(); e != nil; e = e.Next() {
		if e.Value != nil {
			if ok := f(e.Value.key, e.Value.value, idx); !ok {
				return
			}
			idx++
		}
	}
}

func (m *Map[K, V]) Each(f func(key K, value V)) {
	for e := m.list.Head(); e != nil; e = e.Next() {
		if e.Value != nil {
			f(e.Value.key, e.Value.value)
		}
	}
}

func (m *Map[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.entries))
	m.Each(func(key K, _ V) {
		keys = append(keys, key)
	})
	return keys
}

func (m *Map[K, V]) Values() []V {
	values := make([]V, 0, len(m.entries))
	m.Each(func(_ K, value V) {
		values = append(values, value)
	})
	return values
}

func (m *Map[K, V]) Entries() []*Entry[K, V] {
	entries := make([]*Entry[K, V], 0, len(m.entries))
	m.Each(func(key K, value V) {
		entries = append(entries, &Entry[K, V]{key: key, value: value})
	})
	return entries
}

// String returns a string representation of container
func (m *Map[K, V]) String() string {
	str := "LinkedMap\nmap["
	it := m.Iterator()
	for it.Next() {
		str += fmt.Sprintf("%v:%v ", it.Key(), it.Value())
	}
	return strings.TrimRight(str, " ") + "]"
}
