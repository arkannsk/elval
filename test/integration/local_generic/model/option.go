package model

// @oa:ignore
type Option[T any] struct {
	value T
	ok    bool
}

func Some[T any](v T) Option[T] { return Option[T]{value: v, ok: true} }
func None[T any]() Option[T]    { return Option[T]{} }

func (o Option[T]) IsPresent() bool  { return o.ok }
func (o Option[T]) Value() (T, bool) { return o.value, o.ok }
func (o Option[T]) GetOr(fallback T) T {
	if o.ok {
		return o.value
	}
	return fallback
}
