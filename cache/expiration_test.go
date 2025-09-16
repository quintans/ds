package cache_test

import (
	"sync"
	"testing"
	"time"

	"github.com/quintans/ds/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpiration(t *testing.T) {
	type call struct {
		key   string
		value string
	}
	calls := []call{}
	mu := sync.Mutex{}

	exp := cache.NewExpiration(100, 300*time.Millisecond, 100*time.Millisecond, func(key string, value string) {
		mu.Lock()
		calls = append(calls, call{key, value})
		mu.Unlock()
	})
	t.Cleanup(func() {
		exp.Dispose()
	})

	exp.Put("a", "A")
	v, ok := exp.GetIfPresent("a")
	require.True(t, ok)
	assert.Equal(t, "A", v)

	exp.Put("b", "B")

	v, ok = exp.GetIfPresent("b")
	require.True(t, ok)
	assert.Equal(t, "B", v)

	time.Sleep(500 * time.Millisecond)

	_, ok = exp.GetIfPresent("a")
	require.False(t, ok)

	_, ok = exp.GetIfPresent("b")
	require.False(t, ok)

	require.Len(t, calls, 2)
	assert.Equal(t, call{"a", "A"}, calls[0])
	assert.Equal(t, call{"b", "B"}, calls[1])
}
