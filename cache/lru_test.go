package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCapacity(t *testing.T) {
	lru := NewLRU[string, string](2)
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
}
