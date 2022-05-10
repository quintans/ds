package collections

import (
	"strings"
)

const default_hashmap_capacity = 16

type linkedEntry[K, V any] struct {
	key   K
	value V
	next  *linkedEntry[K, V]
}

// check if it implements Map interface
var _ Map[string, any] = (*HashMap[string, any])(nil)

type HashMap[K, V any] struct {
	maxThreshold    float32
	minThreshold    float32
	maxSize         int
	minSize         int
	tableSize       int
	size            int
	table           []*linkedEntry[K, V]
	initialCapacity int
	equals          func(a, b K) bool
	hashCode        func(a K) int
}

func NewHashMap[K, V any](cmp func(a, b K) bool, hash func(a K) int) *HashMap[K, V] {
	return NewHashMapWithCapacity[K, V](cmp, hash, default_hashmap_capacity)
}

func NewHashMapWithCapacity[K, V any](cmp func(a, b K) bool, hash func(a K) int, capacity int) *HashMap[K, V] {
	hm := &HashMap[K, V]{
		maxThreshold:    0.75,
		minThreshold:    0.25,
		tableSize:       capacity,
		initialCapacity: capacity,
		equals:          cmp,
		hashCode:        hash,
	}
	hm.Clear()
	return hm
}

func (h *HashMap[K, V]) Clear() {
	h.maxSize = int(float32(h.tableSize) * h.maxThreshold)
	h.minSize = int(float32(h.tableSize) * h.minThreshold)
	h.size = 0
	h.table = make([]*linkedEntry[K, V], h.tableSize)
}

func (h *HashMap[K, V]) resize(newSize int) {
	oldTableSize := h.tableSize
	h.tableSize = newSize
	h.maxSize = int(float32(h.tableSize) * h.maxThreshold)
	h.minSize = int(float32(h.tableSize) * h.minThreshold)
	oldTable := h.table
	h.table = make([]*linkedEntry[K, V], h.tableSize)
	h.size = 0
	for hash := 0; hash < oldTableSize; hash++ {
		if oldTable[hash] != nil {
			entry := oldTable[hash]
			for entry != nil {
				h.Put(entry.key, entry.value)
				entry = entry.next
			}
			// discard
			oldTable[hash] = nil
		}
	}
}

func (h *HashMap[K, V]) index(key K) int {
	return ((h.hashCode(key) & 0x7FFFFFFF) % h.tableSize)
}

func (h *HashMap[K, V]) Get(key K) (V, bool) {
	hash := h.index(key)
	if h.table[hash] != nil {
		entry := h.table[hash]
		for entry != nil && !h.equals(entry.key, key) {
			entry = entry.next
		}
		if entry != nil {
			return entry.value, true
		}
	}
	return Zero[V](), false
}

func (h *HashMap[K, V]) Put(key K, value V) (V, bool) {
	hash := h.index(key)
	if h.table[hash] == nil {
		h.table[hash] = &linkedEntry[K, V]{key, value, nil}
		h.size++
	} else {
		entry := h.table[hash]
		var prevEntry *linkedEntry[K, V]
		for entry != nil && !h.equals(entry.key, key) {
			prevEntry = entry
			entry = entry.next
		}
		if entry == nil {
			prevEntry.next = &linkedEntry[K, V]{key, value, nil}
			h.size++
		} else {
			old := entry.value
			entry.value = value
			return old, true
		}
	}
	if h.size >= h.maxSize {
		h.resize(h.tableSize * 2)
	}

	var zero V
	return zero, false
}

func (h *HashMap[K, V]) ContainsKey(key K) bool {
	_, ok := h.Get(key)
	return ok
}

func (h *HashMap[K, V]) Delete(key K) (V, bool) {
	hash := h.index(key)
	if entry := h.table[hash]; entry != nil {
		var prevEntry *linkedEntry[K, V]
		for entry != nil && !h.equals(entry.key, key) {
			prevEntry = entry
			entry = entry.next
		}
		if entry != nil {
			old := entry.value
			if prevEntry == nil {
				h.table[hash] = entry.next
			} else {
				prevEntry.next = entry.next
			}
			h.size--
			return old, true
		}
	}

	if h.size <= h.minSize && h.tableSize > h.initialCapacity {
		h.resize(h.tableSize / 2)
	}

	var zero V
	return zero, false
}

func (this *HashMap[K, V]) Size() int {
	return this.size
}

func (h *HashMap[K, V]) Entries() []KeyValue[K, V] {
	data := make([]KeyValue[K, V], h.size)
	i := 0
	for it := h.Iterator(); it.HasNext(); {
		data[i] = it.Next()
		i++
	}
	return data
}

func (h *HashMap[K, V]) Values() []V {
	data := make([]V, h.size)
	i := 0
	for it := h.Iterator(); it.HasNext(); {
		data[i] = it.Next().Value
		i++
	}
	return data
}

func (h *HashMap[K, V]) ForEach(fn func(K, V)) {
	for it := h.Iterator(); it.HasNext(); {
		entry := it.Next()
		fn(entry.Key, entry.Value)
	}
}

func (h *HashMap[K, V]) ReplaceAll(fn func(K, V) V) {
	it := &HashMapIterator[K, V]{hashmap: h}
	for entry := it.next(); entry != nil; {
		entry.value = fn(entry.key, entry.value)
	}
}

func (h *HashMap[K, V]) String() string {
	var s strings.Builder
	s.WriteString("[")
	counter := 0
	for it := h.Iterator(); it.HasNext(); {
		if counter > 1 {
			s.WriteString(", ")
		}
		s.WriteString(it.Next().String())
		counter++
	}
	s.WriteString("]")

	return s.String()
}

func (h *HashMap[K, V]) Clone() *HashMap[K, V] {
	m := NewHashMapWithCapacity[K, V](h.equals, h.hashCode, h.initialCapacity)
	for it := h.Iterator(); it.HasNext(); {
		kv := it.Next()
		m.Put(kv.Key, kv.Value)
	}

	return m
}

// returns a function that in every call return the next value
// if key is null, no value was retrieved
func (h *HashMap[K, V]) Iterator() Iterator[KeyValue[K, V]] {
	it := &HashMapIterator[K, V]{hashmap: h}
	// initiates
	it.next()
	return it
}

type HashMapIterator[K, V any] struct {
	hashmap   *HashMap[K, V]
	hash      int
	prevEntry *linkedEntry[K, V]
	entry     *linkedEntry[K, V]
}

func (it *HashMapIterator[K, V]) HasNext() bool {
	return it.entry != nil
}

func (it *HashMapIterator[K, V]) Next() KeyValue[K, V] {
	if it.entry != nil {
		kv := KeyValue[K, V]{it.entry.key, it.entry.value}
		it.next()
		return kv
	}
	var zero KeyValue[K, V]
	return zero
}

func (it *HashMapIterator[K, V]) next() *linkedEntry[K, V] {
	max := len(it.hashmap.table)
	var aEntry *linkedEntry[K, V]
	for idx := it.hash; aEntry == nil && idx < max; idx++ {
		if it.entry == nil {
			it.prevEntry = nil
			it.entry = it.hashmap.table[idx]
		} else {
			it.prevEntry = it.entry
			it.entry = it.entry.next
		}
		aEntry = it.entry
		if it.entry != nil {
			it.hash = idx
		}
	}
	return aEntry
}

func (it *HashMapIterator[K, V]) Remove() {
	if it.entry != nil {
		if it.prevEntry == nil {
			it.hashmap.table[it.hash] = it.entry.next
		} else {
			it.prevEntry.next = it.entry.next
		}
		it.hashmap.size--
	}
}
