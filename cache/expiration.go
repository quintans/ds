package cache

import (
	"runtime"
	"sync"
	"time"
	"weak"
)

type Expiration[K comparable, V any] struct {
	items   *LRU[K, *item[V]]
	timeout time.Duration
	quit    chan struct{}
	mu      sync.Mutex
	locks   map[K]*sync.Mutex
}

type item[V any] struct {
	value      V
	expiration time.Time
}

// Returns true if the item has expired.
func (i *item[V]) expired() bool {
	return i.expiration.Before(time.Now())
}

func NewExpiration[K comparable, V any](capacity int, timeout time.Duration, interval time.Duration, onEvict func(key K, value V)) *Expiration[K, V] {
	quit := make(chan struct{})
	cache := &Expiration[K, V]{
		items: NewLRU[K, *item[V]](capacity, func(key K, value *item[V]) {
			onEvict(key, value.value)
		}),
		timeout: timeout,
		quit:    quit,
		locks:   make(map[K]*sync.Mutex, capacity),
	}

	runtime.AddCleanup(cache, func(quit chan struct{}) {
		close(quit)
	}, quit)

	go cleanup(weak.Make(cache), interval, quit)

	return cache
}

func (c *Expiration[K, V]) Dispose() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.quit != nil {
		close(c.quit)
		c.quit = nil
		c.items.Clear()
		c.locks = nil
	}
}

func (c *Expiration[K, V]) GetIfPresent(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.items.Get(key)
	if ok {
		v.expiration = time.Now().Add(c.timeout)
		return v.value, true
	}
	var zero V
	return zero, false
}

func (c *Expiration[K, V]) Get(key K, callback func() (V, error)) (V, error) {
	c.mu.Lock()
	it, ok := c.items.Get(key)
	c.mu.Unlock()

	if !ok {
		// per key lock to avoid multiple goroutines creating the same item
		c.mu.Lock()
		l, ok := c.locks[key]
		if !ok {
			l = &sync.Mutex{}
			c.locks[key] = l
		}
		c.mu.Unlock()

		l.Lock()
		defer l.Unlock()

		// recheck if the item was created while waiting for the lock
		c.mu.Lock()
		it, ok = c.items.Get(key)
		c.mu.Unlock()
		if ok {
			it.expiration = time.Now().Add(c.timeout)
			return it.value, nil
		}

		// create the item
		v, err := callback()
		if err != nil {
			return *new(V), err
		}
		it = &item[V]{v, time.Now().Add(c.timeout)}
		c.mu.Lock()
		c.items.Put(key, it)
		c.mu.Unlock()
	}
	return it.value, nil
}

func (c *Expiration[K, V]) Put(key K, value V) {
	c.mu.Lock()
	// defer now since I do not know what will happen in a out of memory error
	defer c.mu.Unlock()
	c.items.Put(key, &item[V]{value, time.Now().Add(c.timeout)})
}

func (c *Expiration[K, V]) Extend(key K) {
	c.mu.Lock()
	v, ok := c.items.Get(key)
	if ok {
		v.expiration = time.Now().Add(c.timeout)
	}
	c.mu.Unlock()
}

func (c *Expiration[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items.Delete(key)
}

func cleanup[K comparable, V any](wp weak.Pointer[Expiration[K, V]], interval time.Duration, quit chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c := wp.Value()
			if c == nil {
				return
			}
			c.deleteExpired()
		case <-quit:
			return
		}
	}
}

// Delete all expired items from the cache.
func (c *Expiration[K, V]) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.items.ReverseIterator() {
		if !v.expired() {
			// since the items are ordered by last access time, we can stop here
			break
		}

		c.items.Delete(k)
	}
}
