package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCapacity(t *testing.T) {
	lru := NewLRU[string](2)
	lru.Put("one", "um")
	lru.Put("two", "dois")
	lru.Put("three", "tres") // will make "one" to be dropped

	_, found := lru.GetIfPresent("one")
	require.False(t, found)
	_, found = lru.GetIfPresent("two")
	require.True(t, found)
	_, found = lru.GetIfPresent("three")
	require.True(t, found)
	v, found := lru.Get("one", func() string {
		return "xxx"
	})
	require.Equal(t, "xxx", v)
	require.False(t, found)

	lru.Put("three", "yyy") // will move it to the front

	lru.Delete("three")
	require.Equal(t, 1, lru.Size())
}
