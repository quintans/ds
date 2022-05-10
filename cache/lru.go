package cache

import (
	"github.com/quintans/dstruct/collections"
)

var _ Cache[string] = (*LRUCache[string])(nil)

type LRUCache[V any] struct {
	entries *collections.DoublyLinkedList[*entry[V]]
	table   map[string]*collections.Element[*entry[V]]

	capacity int
}

type entry[V any] struct {
	key   string
	value V
}

func NewLRUCache[V any](capacity int) *LRUCache[V] {
	c := &LRUCache[V]{
		capacity: capacity,
	}
	c.Clear()
	return c
}

func (c *LRUCache[V]) Clear() {
	c.entries = collections.NewDoublyLinkedList(func(a, b *entry[V]) bool {
		return a.key == b.key
	})
	c.table = map[string]*collections.Element[*entry[V]]{}
}

func (c *LRUCache[V]) Size() int {
	return len(c.table)
}

func (c *LRUCache[V]) GetIfPresent(key string) (V, bool) {
	return c.get(key)
}

func (c *LRUCache[V]) get(key string) (V, bool) {
	element := c.table[key]
	if element != nil {
		// move to front
		c.entries.MoveToFirst(element)
		return element.Value.value, true
	}
	var zero V
	return zero, false
}

// Get gets a the value under the key. If the value is not found it will use the value of the callback function, store it and return it.
// ok=true indicating that it was found in the cache
// ok=false indicating that it was not found in the cache and was created by the callback
func (c *LRUCache[V]) Get(key string, callback func() V) (V, bool) {
	value, ok := c.get(key)
	if !ok {
		value = callback()
		c.add(key, value)
	}
	return value, ok
}

func (c *LRUCache[V]) Put(key string, value V) {
	element := c.table[key]
	if element != nil {
		e := element.Value
		e.value = value
		c.entries.MoveToFirst(element)
	} else {
		c.add(key, value)
	}
}

func (c *LRUCache[V]) add(key string, value V) {
	if c.entries.Size() == c.capacity {
		// if at full capacity recycle last element
		element := c.entries.Tail()
		e := element.Value
		element.Value = &entry[V]{key, value}
		c.entries.MoveToFirst(element)

		delete(c.table, e.key)
		c.table[key] = element
	} else {
		c.entries.AddFirst(&entry[V]{key, value})
		c.table[key] = c.entries.Head()
	}
}

func (c *LRUCache[V]) Delete(key string) {
	element := c.table[key]
	if element != nil {
		// remove from list
		element.Remove()
	}
	delete(c.table, key)
}
