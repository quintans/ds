package collections

const default_linkedhashmap_capacity = 16

type LinkedHashMap[K, V any] struct {
	keyOrder        *DoublyLinkedList[K]
	entries         *HashMap[K, entry[V]]
	initialCapacity int
	equals          func(a, b K) bool
	hashCode        func(a K) int
}

type entry[V any] struct {
	value           V
	keyOrderRemover func()
}

// check if it implements Map interface
var _ Map[string, any] = (*LinkedHashMap[string, any])(nil)

func NewLinkedHashMap[K, V any](cmp func(a, b K) bool, hash func(a K) int) *LinkedHashMap[K, V] {
	return NewLinkedHashMapWithCapacity[K, V](cmp, hash, default_linkedhashmap_capacity)
}

func NewLinkedHashMapWithCapacity[K, V any](cmp func(a, b K) bool, hash func(a K) int, capacity int) *LinkedHashMap[K, V] {
	m := &LinkedHashMap[K, V]{
		initialCapacity: capacity,
		equals:          cmp,
		hashCode:        hash,
	}
	m.Clear()
	return m
}

func (l *LinkedHashMap[K, V]) Clear() {
	l.keyOrder = NewDoublyLinkedList[K](l.equals)
	l.entries = NewHashMapWithCapacity[K, entry[V]](l.equals, l.hashCode, l.initialCapacity)
}

func (l *LinkedHashMap[K, V]) Size() int {
	return l.entries.Size()
}

func (l *LinkedHashMap[K, V]) Get(key K) (V, bool) {
	e, ok := l.entries.Get(key)
	return e.value, ok
}

func (l *LinkedHashMap[K, V]) Put(key K, value V) (V, bool) {
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

func (h *LinkedHashMap[K, V]) ContainsKey(key K) bool {
	_, ok := h.Get(key)
	return ok
}

func (l *LinkedHashMap[K, V]) Delete(key K) (V, bool) {
	old, deleted := l.entries.Delete(key)
	if deleted {
		old.keyOrderRemover()
	}
	return old.value, deleted
}

func (l *LinkedHashMap[K, V]) Entries() []KeyValue[K, V] {
	data := make([]KeyValue[K, V], 0, l.entries.Size())
	for it := l.Iterator(); it.HasNext(); {
		data = append(data, it.Next())
	}
	return data
}

func (l *LinkedHashMap[K, V]) Values() []V {
	data := make([]V, 0, l.entries.Size())
	for it := l.Iterator(); it.HasNext(); {
		data = append(data, it.Next().Value)
	}
	return data
}

func (l *LinkedHashMap[K, V]) ForEach(fn func(K, V)) {
	for it := l.Iterator(); it.HasNext(); {
		entry := it.Next()
		fn(entry.Key, entry.Value)
	}
}

func (l *LinkedHashMap[K, V]) ReplaceAll(fn func(K, V) V) {
	it := &HashMapIterator[K, entry[V]]{hashmap: l.entries}
	for entry := it.next(); entry != nil; {
		entry.value.value = fn(entry.key, entry.value.value)
	}
}

func (l *LinkedHashMap[K, V]) Clone() *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{
		keyOrder:        l.keyOrder.Clone(),
		entries:         l.entries.Clone(),
		initialCapacity: l.initialCapacity,
	}
}

// returns a function that in every call return the next value
// if key is null, no value was retrieved
func (l *LinkedHashMap[K, V]) Iterator() Iterator[KeyValue[K, V]] {
	return &LinkedHashMapIterator[K, V]{
		hashmap:  l.entries,
		iterator: l.keyOrder.Iterator(),
	}
}

type LinkedHashMapIterator[K, V any] struct {
	hashmap  *HashMap[K, entry[V]]
	iterator Iterator[K]
	lastKey  K
}

func (l *LinkedHashMapIterator[K, V]) HasNext() bool {
	return l.iterator.HasNext()
}

func (l *LinkedHashMapIterator[K, V]) Next() KeyValue[K, V] {
	k := l.iterator.Next()
	l.lastKey = k
	return l.getEntry(k)
}

func (l *LinkedHashMapIterator[K, V]) getEntry(k K) KeyValue[K, V] {
	e, ok := l.hashmap.Get(k)
	if ok {
		return KeyValue[K, V]{k, e.value}
	}

	return KeyValue[K, V]{}
}

func (l *LinkedHashMapIterator[K, V]) Remove() {
	l.hashmap.Delete(l.lastKey)
}
