package cache

type Cache[V any] interface {
	GetIfPresent(key string) (V, bool)
	Delete(key string)
	// Get returns the value under the key
	// true is returned if the value was found in the cache
	// false is return if the value was not found in the cache and was created by the callback
	Get(key string, callback func() V) (V, bool)
	Put(key string, value V)
}
