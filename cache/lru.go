package cache

// NodeKV represents a node in the doubly linked list
type NodeKV[K comparable, V any] struct {
	key   K
	value V
	next  *NodeKV[K, V]
	prev  *NodeKV[K, V]
}

// LRU represents a Least Recently Used cache
//
// Head -> Node -> <-prev- Node(value) -next-> <- Node <-Tail
type LRU[K comparable, V any] struct {
	capacity int
	cache    map[K]*NodeKV[K, V]
	head     *NodeKV[K, V]
	tail     *NodeKV[K, V]
}

func NewLRU[K comparable, V any](capacity int) *LRU[K, V] {
	return &LRU[K, V]{
		capacity: capacity,
		cache:    make(map[K]*NodeKV[K, V], capacity),
	}
}

func (l *LRU[K, V]) Get(key K) (V, bool) {
	if node, ok := l.cache[key]; ok {
		l.moveToFront(node)
		return node.value, true
	}
	var zero V
	return zero, false
}

func (l *LRU[K, V]) Put(key K, value V) {
	if node, ok := l.cache[key]; ok {
		node.value = value
		l.moveToFront(node)
		return
	}

	if l.Size() == l.capacity {
		delete(l.cache, l.tail.key)
		l.remove(l.tail)
	}

	newNode := &NodeKV[K, V]{key: key, value: value}
	l.cache[key] = newNode
	l.add(newNode)
}

func (l *LRU[K, V]) Delete(key K) {
	if node, ok := l.cache[key]; ok {
		l.remove(node)
		delete(l.cache, key)
	}
}

func (l *LRU[K, V]) Clear() {
	l.cache = make(map[K]*NodeKV[K, V], l.capacity)
	l.head = nil
	l.tail = nil
}

func (l *LRU[K, V]) Size() int {
	return len(l.cache)
}

func (l *LRU[K, V]) moveToFront(node *NodeKV[K, V]) {
	l.remove(node)
	l.add(node)
}

func (l *LRU[K, V]) remove(node *NodeKV[K, V]) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}

	node.next = nil
	node.prev = nil
}

func (l *LRU[K, V]) add(node *NodeKV[K, V]) {
	if l.head == nil {
		l.head = node
		l.tail = node
		return
	}

	l.head.prev = node
	node.next = l.head
	l.head = node
}
