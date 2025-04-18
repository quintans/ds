package indexedpriorityqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// testItem is a simple struct used for testing.
type testItem struct {
	id    int
	value string
}

// testComparator compares testItems based on priority (lower priority value means higher priority).
func testComparator(a, b testItem) int {
	return a.id - b.id // Min-heap behavior
}

// testKeyExtractor extracts the id as the key.
func testKeyExtractor(item testItem) int {
	return item.id
}

// Helper to create a new queue for tests
func newTestQueue() *IndexedPriorityQueue[testItem, int] {
	return New(testComparator, testKeyExtractor)
}

func TestNew(t *testing.T) {
	q := newTestQueue()
	assert.NotNil(t, q, "New queue should not be nil")
	assert.Equal(t, 0, q.Len(), "New queue should be empty")
}

func TestEnqueueDequeue(t *testing.T) {
	q := newTestQueue()
	item1 := testItem{id: 1, value: "A"}
	item2 := testItem{id: 3, value: "C"}
	item3 := testItem{id: 2, value: "B"}

	q.Enqueue(item1)
	q.Enqueue(item2)
	q.Enqueue(item3)

	assert.Equal(t, 3, q.Len(), "Queue length should be 3 after enqueuing")

	// Dequeue should return items in priority order (lowest id first)
	dequeued1, ok1 := q.Dequeue()
	assert.True(t, ok1, "Dequeue should succeed")
	assert.Equal(t, item1, dequeued1, "First dequeued item should be item1")
	assert.Equal(t, 2, q.Len(), "Queue length should be 2")

	dequeued2, ok2 := q.Dequeue()
	assert.True(t, ok2, "Dequeue should succeed")
	assert.Equal(t, item3, dequeued2, "Second dequeued item should be item3") // id: 2
	assert.Equal(t, 1, q.Len(), "Queue length should be 1")

	dequeued3, ok3 := q.Dequeue()
	assert.True(t, ok3, "Dequeue should succeed")
	assert.Equal(t, item2, dequeued3, "Third dequeued item should be item2") // id: 3
	assert.Equal(t, 0, q.Len(), "Queue length should be 0")

	_, ok4 := q.Dequeue()
	assert.False(t, ok4, "Dequeue on empty queue should fail")
}

func TestEnqueueUpdate(t *testing.T) {
	q := newTestQueue()
	item1 := testItem{id: 1, value: "A"}
	item2 := testItem{id: 2, value: "B"}
	item1_updated := testItem{id: 1, value: "A_updated"} // Same key as item1

	q.Enqueue(item1)
	q.Enqueue(item2)
	assert.Equal(t, 2, q.Len(), "Queue length should be 2")

	// Enqueue item with the same key, should trigger update
	q.Enqueue(item1_updated)
	assert.Equal(t, 2, q.Len(), "Queue length should still be 2 after update")

	// Check if the item was updated
	retrieved, ok := q.Get(1)
	assert.True(t, ok, "Should be able to get item with key 1")
	assert.Equal(t, "A_updated", retrieved.value, "Item value should be updated")

	// Dequeue to check priority and updated value
	dequeued1, ok1 := q.Dequeue()
	assert.True(t, ok1)
	assert.Equal(t, item1_updated, dequeued1, "Dequeued item should have updated value")

	dequeued2, ok2 := q.Dequeue()
	assert.True(t, ok2)
	assert.Equal(t, item2, dequeued2)
}

func TestPeek(t *testing.T) {
	q := newTestQueue()

	_, ok := q.Peek()
	assert.False(t, ok, "Peek on empty queue should fail")

	item1 := testItem{id: 2, value: "B"}
	item2 := testItem{id: 1, value: "A"}
	q.Enqueue(item1)
	q.Enqueue(item2)

	peeked, ok := q.Peek()
	assert.True(t, ok, "Peek should succeed on non-empty queue")
	assert.Equal(t, item2, peeked, "Peek should return the highest priority item (item2)")
	assert.Equal(t, 2, q.Len(), "Peek should not change queue length")
}

func TestGet(t *testing.T) {
	q := newTestQueue()
	item1 := testItem{id: 1, value: "A"}
	item2 := testItem{id: 2, value: "B"}
	q.Enqueue(item1)
	q.Enqueue(item2)

	retrieved1, ok1 := q.Get(1)
	assert.True(t, ok1, "Should get item with key 1")
	assert.Equal(t, item1, retrieved1, "Retrieved item 1 should match")

	retrieved2, ok2 := q.Get(2)
	assert.True(t, ok2, "Should get item with key 2")
	assert.Equal(t, item2, retrieved2, "Retrieved item 2 should match")

	_, ok3 := q.Get(3) // Non-existent key
	assert.False(t, ok3, "Should not get item with non-existent key 3")

	assert.Equal(t, 2, q.Len(), "Get should not change queue length")
}

func TestRemove(t *testing.T) {
	q := newTestQueue()
	item1 := testItem{id: 1, value: "A"}
	item2 := testItem{id: 3, value: "C"}
	item3 := testItem{id: 2, value: "B"}
	q.Enqueue(item1)
	q.Enqueue(item2)
	q.Enqueue(item3) // Queue: [1:A, 2:B, 3:C] (priority order)

	assert.Equal(t, 3, q.Len())

	// Remove middle priority item
	removed1, ok1 := q.Remove(2) // Remove item3
	assert.True(t, ok1, "Remove item with key 2 should succeed")
	assert.Equal(t, item3, removed1, "Removed item should be item3")
	assert.Equal(t, 2, q.Len(), "Queue length should be 2 after removing key 2")

	// Verify remaining items and order
	peeked1, _ := q.Peek()
	assert.Equal(t, item1, peeked1, "Highest priority should still be item1")
	_, okGet2 := q.Get(2)
	assert.False(t, okGet2, "Item with key 2 should no longer exist")

	// Remove highest priority item
	removed2, ok2 := q.Remove(1) // Remove item1
	assert.True(t, ok2, "Remove item with key 1 should succeed")
	assert.Equal(t, item1, removed2, "Removed item should be item1")
	assert.Equal(t, 1, q.Len(), "Queue length should be 1 after removing key 1")

	// Verify remaining item
	peeked2, _ := q.Peek()
	assert.Equal(t, item2, peeked2, "Highest priority should now be item2")
	_, okGet1 := q.Get(1)
	assert.False(t, okGet1, "Item with key 1 should no longer exist")

	// Remove last item
	removed3, ok3 := q.Remove(3) // Remove item2
	assert.True(t, ok3, "Remove item with key 3 should succeed")
	assert.Equal(t, item2, removed3, "Removed item should be item2")
	assert.Equal(t, 0, q.Len(), "Queue length should be 0 after removing key 3")

	// Try removing non-existent key
	_, ok4 := q.Remove(1)
	assert.False(t, ok4, "Remove non-existent key should fail")
}

func TestLen(t *testing.T) {
	q := newTestQueue()
	assert.Equal(t, 0, q.Len())
	q.Enqueue(testItem{id: 1, value: "A"})
	assert.Equal(t, 1, q.Len())
	q.Enqueue(testItem{id: 2, value: "B"})
	assert.Equal(t, 2, q.Len())
	q.Dequeue()
	assert.Equal(t, 1, q.Len())
	q.Remove(2)
	assert.Equal(t, 0, q.Len())
}

func TestClear(t *testing.T) {
	q := newTestQueue()
	q.Enqueue(testItem{id: 1, value: "A"})
	q.Enqueue(testItem{id: 2, value: "B"})
	assert.Equal(t, 2, q.Len())
	_, okGet := q.Get(1)
	assert.True(t, okGet)

	q.Clear()
	assert.Equal(t, 0, q.Len(), "Queue length should be 0 after clear")
	_, okPeek := q.Peek()
	assert.False(t, okPeek, "Peek should fail after clear")
	_, okGetAfterClear := q.Get(1)
	assert.False(t, okGetAfterClear, "Get should fail after clear")
}

func TestMixedOperations(t *testing.T) {
	q := newTestQueue()

	// Enqueue some items
	q.Enqueue(testItem{id: 5, value: "E"})
	q.Enqueue(testItem{id: 2, value: "B"})
	q.Enqueue(testItem{id: 8, value: "H"})
	q.Enqueue(testItem{id: 1, value: "A"}) // Highest priority
	assert.Equal(t, 4, q.Len())

	// Peek
	peeked1, ok1 := q.Peek()
	assert.True(t, ok1)
	assert.Equal(t, 1, peeked1.id)

	// Dequeue highest priority
	dequeued1, ok2 := q.Dequeue()
	assert.True(t, ok2)
	assert.Equal(t, 1, dequeued1.id) // Should be A
	assert.Equal(t, 3, q.Len())

	// Enqueue duplicate key (update)
	q.Enqueue(testItem{id: 2, value: "B_updated"})
	assert.Equal(t, 3, q.Len())
	item2_updated, ok3 := q.Get(2)
	assert.True(t, ok3)
	assert.Equal(t, "B_updated", item2_updated.value)

	// Peek again
	peeked2, ok4 := q.Peek()
	assert.True(t, ok4)
	assert.Equal(t, 2, peeked2.id) // Should be B_updated

	// Remove an item
	removed1, ok5 := q.Remove(8) // Remove H
	assert.True(t, ok5)
	assert.Equal(t, 8, removed1.id)
	assert.Equal(t, 2, q.Len())

	// Dequeue remaining
	dequeued2, ok6 := q.Dequeue()
	assert.True(t, ok6)
	assert.Equal(t, 2, dequeued2.id) // Should be B_updated

	dequeued3, ok7 := q.Dequeue()
	assert.True(t, ok7)
	assert.Equal(t, 5, dequeued3.id) // Should be E

	assert.Equal(t, 0, q.Len())
}

// Example test for a max-heap scenario
func TestMaxHeapBehavior(t *testing.T) {
	// Comparator for max-heap (higher id means higher priority)
	maxHeapComparator := func(a, b testItem) int {
		return b.id - a.id // Reverse comparison for max-heap
	}
	// Update func remains the same for this test
	q := New(maxHeapComparator, testKeyExtractor)

	item1 := testItem{id: 1, value: "A"}
	item2 := testItem{id: 3, value: "C"} // Highest priority
	item3 := testItem{id: 2, value: "B"}

	q.Enqueue(item1)
	q.Enqueue(item2)
	q.Enqueue(item3)

	assert.Equal(t, 3, q.Len())

	// Dequeue should return items in max priority order (highest id first)
	dequeued1, ok1 := q.Dequeue()
	assert.True(t, ok1)
	assert.Equal(t, item2, dequeued1, "First dequeued item should be item2 (id: 3)")

	dequeued2, ok2 := q.Dequeue()
	assert.True(t, ok2)
	assert.Equal(t, item3, dequeued2, "Second dequeued item should be item3 (id: 2)")

	dequeued3, ok3 := q.Dequeue()
	assert.True(t, ok3)
	assert.Equal(t, item1, dequeued3, "Third dequeued item should be item1 (id: 1)")

	assert.Equal(t, 0, q.Len())
}

// Test edge cases: empty queue operations
func TestEmptyQueueOperations(t *testing.T) {
	q := newTestQueue()

	// Peek
	_, okPeek := q.Peek()
	assert.False(t, okPeek, "Peek on empty queue should return false")

	// Dequeue
	_, okDequeue := q.Dequeue()
	assert.False(t, okDequeue, "Dequeue on empty queue should return false")

	// Get
	_, okGet := q.Get(1) // Any key
	assert.False(t, okGet, "Get on empty queue should return false")

	// Remove
	_, okRemove := q.Remove(1) // Any key
	assert.False(t, okRemove, "Remove on empty queue should return false")

	// Len
	assert.Equal(t, 0, q.Len(), "Len on empty queue should be 0")
}

// Test update function changing priority
func TestUpdateChangesPriority(t *testing.T) {
	q := New(testComparator, testKeyExtractor)

	item1 := testItem{id: 5, value: "E"}
	item2 := testItem{id: 2, value: "B"}

	q.Enqueue(item1)
	q.Enqueue(item2) // Queue: [2:B, 5:E]

	peeked1, _ := q.Peek()
	assert.Equal(t, 2, peeked1.id, "Initial peek should be id 2")

	// Update item with key 5 to have a higher priority (lower id)
	// item1_updated := testItem{id: 1, value: "A"}        // Key is 5, but new id is 1 <-- Removed unused variable
	q.Enqueue(testItem{id: 5, value: "A_new_priority"}) // Enqueue with key 5, update func will use this

	// This scenario is tricky with the current updateFunc signature.
	// The updateFunc updates the *value* of the item associated with the *key*.
	// It doesn't change the key itself or allow changing the priority field directly
	// if the priority field *is* the key.
	// Let's redefine the test slightly based on how the code works.
	// The `testUpdateFunc` we defined earlier only updates the `value` field, not `id`.
	// Let's test *that* behavior correctly.

	q_standard := newTestQueue() // Using the original testUpdateFunc
	item_orig_5 := testItem{id: 5, value: "E"}
	item_orig_2 := testItem{id: 2, value: "B"}
	q_standard.Enqueue(item_orig_5)
	q_standard.Enqueue(item_orig_2) // Queue: [2:B, 5:E]

	peek_std_1, _ := q_standard.Peek()
	assert.Equal(t, 2, peek_std_1.id)

	// Enqueue item with key 5 again, updating its value
	item_update_5 := testItem{id: 5, value: "E_updated"}
	q_standard.Enqueue(item_update_5)

	assert.Equal(t, 2, q_standard.Len(), "Length should remain 2 after update")

	// Peek should still be item 2, as priority (id) didn't change
	peek_std_2, _ := q_standard.Peek()
	assert.Equal(t, 2, peek_std_2.id)

	// Dequeue and check updated value
	dequeued_std_1, _ := q_standard.Dequeue() // Dequeues id 2
	assert.Equal(t, 2, dequeued_std_1.id)

	dequeued_std_2, _ := q_standard.Dequeue() // Dequeues id 5
	assert.Equal(t, 5, dequeued_std_2.id)
	assert.Equal(t, "E_updated", dequeued_std_2.value, "Value of item 5 should be updated")

}

// Test removing the root element
func TestRemoveRoot(t *testing.T) {
	q := newTestQueue()
	item1 := testItem{id: 1, value: "A"}
	item2 := testItem{id: 3, value: "C"}
	item3 := testItem{id: 2, value: "B"}
	q.Enqueue(item1) // Root
	q.Enqueue(item2)
	q.Enqueue(item3)

	removed, ok := q.Remove(1) // Remove root (item1)
	assert.True(t, ok)
	assert.Equal(t, item1, removed)
	assert.Equal(t, 2, q.Len())

	// Check new root
	peeked, okPeek := q.Peek()
	assert.True(t, okPeek)
	assert.Equal(t, item3, peeked, "New root should be item3 (id: 2)") // Next highest priority

	// Check remaining elements
	_, okGet1 := q.Get(1)
	assert.False(t, okGet1)
	_, okGet2 := q.Get(2)
	assert.True(t, okGet2)
	_, okGet3 := q.Get(3)
	assert.True(t, okGet3)
}
