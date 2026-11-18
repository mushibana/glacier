package domain

type Optional[T any] struct {
	Set   bool
	Value T
}

func SetOptional[T any](v T) Optional[T] {
	return Optional[T]{Set: true, Value: v}
}
