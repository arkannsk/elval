package elval

// Generic — универсальная обёртка для извлечения значения из любого Optional-типа.
type Generic[T any] struct {
	val T
	ok  bool
}

func (g Generic[T]) IsPresent() bool  { return g.ok }
func (g Generic[T]) Value() (T, bool) { return g.val, g.ok }

// OptionalExtractor интерфейс для типов, которые можно извлекать через Unwrap
type OptionalExtractor[T any] interface {
	Value() (T, bool)
}

// GetExtractor интерфейс для типов с методом Get()
type GetExtractor[T any] interface {
	Get() (T, bool)
}

// Unwrap извлекает значение из произвольной обёртки.
// Поддерживает типы, реализующие OptionalExtractor[T] интерфейс
func Unwrap[T any](src OptionalExtractor[T]) Generic[T] {
	if src == nil {
		return Generic[T]{}
	}

	val, ok := src.Value()
	return Generic[T]{val: val, ok: ok}
}

// UnwrapWithGet извлекает значение из произвольной обёртки с методом Get().
// Поддерживает типы, реализующие GetExtractor[T] интерфейс
func UnwrapWithGet[T any](src GetExtractor[T]) Generic[T] {
	if src == nil {
		return Generic[T]{}
	}

	val, ok := src.Get()
	return Generic[T]{val: val, ok: ok}
}

// UnwrapGeneric извлекает значение из произвольной обёртки, поддерживающей оба интерфейса.
// Пытается использовать OptionalExtractor сначала, затем GetExtractor.
func UnwrapGeneric[T any](src any) Generic[T] {
	// Проверяем, реализует ли тип OptionalExtractor
	if extractor, ok := src.(OptionalExtractor[T]); ok {
		return Unwrap(extractor)
	}

	// Проверяем, реализует ли тип GetExtractor
	if extractor, ok := src.(GetExtractor[T]); ok {
		return UnwrapWithGet(extractor)
	}

	// Если не реализует ни один из интерфейсов, возвращаем пустой Generic
	return Generic[T]{}
}
