package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCapacity(t *testing.T) {
	delCalls := []string{}
	onEvict := func(key string, value string) {
		delCalls = append(delCalls, key)
	}
	lru := NewLRU(2, onEvict)
	lru.Put("one", "um")
	lru.Put("two", "dois")
	lru.Put("three", "tres") // will make "one" to be dropped

	_, found := lru.Get("one")
	assert.False(t, found)

	v, found := lru.Get("two")
	assert.True(t, found)
	assert.Equal(t, "dois", v)

	v, found = lru.Get("three")
	assert.True(t, found)
	assert.Equal(t, "tres", v)

	lru.Put("three", "yyy") // will move it to the front

	lru.Delete("three")
	require.Equal(t, 1, lru.Size())

	assert.Len(t, delCalls, 2)
	assert.Equal(t, []string{"one", "three"}, delCalls)

	lru.Clear()
	assert.Equal(t, 0, lru.Size())
	assert.Len(t, delCalls, 3) // "two" is also evicted
	assert.Equal(t, []string{"one", "three", "two"}, delCalls)
}
