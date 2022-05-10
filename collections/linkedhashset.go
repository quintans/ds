package collections

const default_linkedhashset_capacity = 16

type LinkedHashSet[K any] struct {
	HashSet[K]
}

// check if it implements Collection interface
var _ Collection[string] = (*LinkedHashSet[string])(nil)

func NewLinkedHashSet[K any](cmp func(a, b K) bool, hash func(a K) int) *LinkedHashSet[K] {
	return NewLinkedHashSetWithCapacity(cmp, hash, default_linkedhashset_capacity)
}

func NewLinkedHashSetWithCapacity[K any](cmp func(a, b K) bool, hash func(a K) int, capacity int) *LinkedHashSet[K] {
	s := &LinkedHashSet[K]{
		HashSet: HashSet[K]{
			equals:          cmp,
			hashCode:        hash,
			initialCapacity: capacity,
		},
	}
	s.Clear()
	return s
}

func (s *LinkedHashSet[K]) Clear() {
	s.entries = NewLinkedHashMapWithCapacity[K, struct{}](s.equals, s.hashCode, s.initialCapacity)
}
