package indexedpriorityqueue

import (
	"cmp"
	"container/heap"
)

// Comparator defines the function signature for comparing two elements.
// It should return:
//   - negative value if a < b
//   - zero if a == b
//   - positive value if a > b
type Comparator[T any] func(a, b T) int

// KeyExtractor defines the function signature for extracting a key K from an element T.
type KeyExtractor[T any, K comparable] func(item T) K

// item represents an element in the priority queue, holding the value and its index in the heap.
type item[T any, K comparable] struct {
	value T
	key   K
	index int // The index of the item in the heap.
}

// pqInternal implements the heap.Interface and stores the items.
type pqInternal[T any, K comparable] struct {
	items      []*item[T, K]
	comparator Comparator[T]
}

// Len returns the number of items in the priority queue.
func (pq *pqInternal[T, K]) Len() int { return len(pq.items) }

// Less compares two items based on the provided comparator.
func (pq *pqInternal[T, K]) Less(i, j int) bool {
	// We want Pop to give us the highest priority item, so we use > here.
	// If comparator(a, b) < 0, it means a has higher priority.
	return pq.comparator(pq.items[i].value, pq.items[j].value) < 0
}

// Swap swaps two items in the priority queue and updates their indices.
func (pq *pqInternal[T, K]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

// Push adds an item to the priority queue.
func (pq *pqInternal[T, K]) Push(x any) {
	n := len(pq.items)
	it := x.(*item[T, K])
	it.index = n
	pq.items = append(pq.items, it)
}

// Pop removes and returns the highest priority item from the priority queue.
func (pq *pqInternal[T, K]) Pop() any {
	old := pq.items
	n := len(old)
	it := old[n-1]
	old[n-1] = nil // avoid memory leak
	it.index = -1  // for safety
	pq.items = old[0 : n-1]
	return it
}

// IndexedPriorityQueue represents a priority queue where elements can be accessed and updated by a key.
type IndexedPriorityQueue[T any, K comparable] struct {
	pq      *pqInternal[T, K]
	index   map[K]*item[T, K] // Map from key K to the item in the heap
	keyFunc KeyExtractor[T, K]
}

// New creates a new IndexedPriorityQueue.
//   - comparator: Function to determine the priority between two elements.
//   - keyFunc: Function to extract the unique key K from an element T.
func New[T any, K comparable](comparator Comparator[T], keyFunc KeyExtractor[T, K]) *IndexedPriorityQueue[T, K] {
	return &IndexedPriorityQueue[T, K]{
		pq: &pqInternal[T, K]{
			items:      make([]*item[T, K], 0),
			comparator: comparator,
		},
		index:   make(map[K]*item[T, K]),
		keyFunc: keyFunc,
	}
}

func NewOrdered[O cmp.Ordered]() *IndexedPriorityQueue[O, O] {
	return New(func(a, b O) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}, func(item O) O {
		return item
	})
}

// Len returns the number of elements in the queue.
func (ipq *IndexedPriorityQueue[T, K]) Len() int {
	return ipq.pq.Len()
}

// Clear removes all elements from the queue.
func (ipq *IndexedPriorityQueue[T, K]) Clear() {
	ipq.pq.items = make([]*item[T, K], 0)
	ipq.index = make(map[K]*item[T, K])
	// No need to explicitly call heap.Init as it's empty
}

// Peek returns the highest priority element without removing it.
// Returns the element and true if the queue is not empty, otherwise the zero value of T and false.
func (ipq *IndexedPriorityQueue[T, K]) Peek() (T, bool) {
	if ipq.pq.Len() == 0 {
		var zero T
		return zero, false
	}
	return ipq.pq.items[0].value, true
}

// Get retrieves the element associated with the given key without removing it.
// Returns the element and true if the key exists, otherwise the zero value of T and false.
func (ipq *IndexedPriorityQueue[T, K]) Get(key K) (T, bool) {
	if item, ok := ipq.index[key]; ok {
		return item.value, true
	}
	var zero T
	return zero, false
}

// Enqueue adds an element to the priority queue.
// If an element with the same key already exists it is replaced by the new one.
func (ipq *IndexedPriorityQueue[T, K]) Enqueue(value T) {
	key := ipq.keyFunc(value)
	if existingItem, ok := ipq.index[key]; ok {
		// Key exists, replace the existing item's value
		existingItem.value = value
	} else {
		// Key doesn't exist, add a new item
		newItem := &item[T, K]{
			value: value,
			key:   key,
			// index will be set by heap.Push
		}
		ipq.index[key] = newItem
		heap.Push(ipq.pq, newItem)
	}
}

// Dequeue removes and returns the highest priority element from the queue.
// Returns the element and true if the queue is not empty, otherwise the zero value of T and false.
func (ipq *IndexedPriorityQueue[T, K]) Dequeue() (T, bool) {
	if ipq.pq.Len() == 0 {
		var zero T
		return zero, false
	}
	item := heap.Pop(ipq.pq).(*item[T, K])
	delete(ipq.index, item.key)
	return item.value, true
}

// Remove removes the element associated with the given key from the queue.
// Returns the removed element and true if the key existed, otherwise the zero value of T and false.
func (ipq *IndexedPriorityQueue[T, K]) Remove(key K) (T, bool) {
	itemToRemove, ok := ipq.index[key]
	if !ok {
		var zero T
		return zero, false
	}

	removedItem := heap.Remove(ipq.pq, itemToRemove.index).(*item[T, K])
	delete(ipq.index, key)

	// Ensure the removed item's index is invalidated
	removedItem.index = -1 // Mark as removed

	return removedItem.value, true
}
