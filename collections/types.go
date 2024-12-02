package collections

func Equals[T comparable](a, b T) bool {
	return a == b
}

func HashCode[T any](a T) int {
	return HASH_SEED.HashAny(a).Int()
}
