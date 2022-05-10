package cache

import (
	"sync"
	"time"
)

var _ Cache[string] = (*ExpirationCache[string])(nil)

type ExpirationCache[V any] struct {
	items    map[string]*item[V]
	timeout  time.Duration
	interval time.Duration
	quit     chan struct{}
	mu       sync.Mutex
}

type item[V any] struct {
	value      V
	expiration time.Time
}

// Returns true if the item has expired.
func (i *item[V]) expired() bool {
	return i.expiration.Before(time.Now())
}

func NewExpirationCache[V any](timeout time.Duration, interval time.Duration) *ExpirationCache[V] {
	cache := &ExpirationCache[V]{
		items:    map[string]*item[V]{},
		timeout:  timeout,
		interval: interval,
		quit:     make(chan struct{}),
	}
	go cache.cleanup(cache.quit)
	return cache
}

func (c *ExpirationCache[V]) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.quit != nil {
		close(c.quit)
		c.quit = nil
	}
}

func (c *ExpirationCache[V]) cleanup(quit chan struct{}) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-quit:
			return
		}
	}
}

// Delete all expired items from the cache.
func (c *ExpirationCache[V]) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.items {
		if v.expired() {
			delete(c.items, k)
		}
	}
}

func (c *ExpirationCache[V]) GetIfPresentAndTouch(key string) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.items[key]
	if ok {
		v.expiration = time.Now().Add(c.timeout)
		return v.value, true
	}
	var zero V
	return zero, false
}

func (c *ExpirationCache[V]) GetIfPresent(key string) (V, bool) {
	c.mu.Lock()
	v, ok := c.items[key]
	c.mu.Unlock()
	if ok {
		return v.value, true
	}
	var zero V
	return zero, false
}

func (c *ExpirationCache[V]) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

func (c *ExpirationCache[V]) Get(key string, callback func() V) (V, bool) {
	return c.GetWithDuration(key, callback, c.timeout)
}

func (c *ExpirationCache[V]) GetWithDuration(key string, callback func() V, duration time.Duration) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.items[key]
	if !ok {
		v = &item[V]{callback(), time.Now().Add(c.timeout)}
		c.items[key] = v
	}
	return v.value, ok
}

func (c *ExpirationCache[V]) Put(key string, value V) {
	c.PutWithDuration(key, value, c.timeout)
}

// put a value in the cache, overwriting any previous value for that key
func (c *ExpirationCache[V]) PutWithDuration(key string, value V, duration time.Duration) {
	c.mu.Lock()
	// defer now sice I do not know what will happen in a out of memory error
	defer c.mu.Unlock()
	c.items[key] = &item[V]{value, time.Now().Add(duration)}
}

func (c *ExpirationCache[V]) Touch(key string) {
	c.TouchWithDuration(key, c.timeout)
}

func (c *ExpirationCache[V]) TouchWithDuration(key string, duration time.Duration) {
	c.mu.Lock()
	v, ok := c.items[key]
	if ok {
		v.expiration = time.Now().Add(duration)
	}
	c.mu.Unlock()
}
