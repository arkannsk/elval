package elval

import "reflect"

// Generic — универсальная обёртка для извлечения значения из любого Optional-типа.
type Generic[T any] struct {
	val T
	ok  bool
}

func (g Generic[T]) IsPresent() bool  { return g.ok }
func (g Generic[T]) Value() (T, bool) { return g.val, g.ok }

// Unwrap извлекает значение из произвольной обёртки.
// Поддерживает методы: Value(), Get(), Unwrap(), Ok(), GetOrZero()
func Unwrap[T any](src any) Generic[T] {
	if src == nil {
		return Generic[T]{}
	}

	v := reflect.ValueOf(src)
	if !v.IsValid() {
		return Generic[T]{}
	}

	// Приоритетный список методов для извлечения
	methods := []string{"Value", "Get", "Unwrap", "Ok", "GetOrZero"}

	for _, name := range methods {
		m := v.MethodByName(name)
		if !m.IsValid() {
			continue
		}

		out := m.Call(nil)
		if len(out) == 0 {
			continue
		}

		// Сигнатура: (T, bool)
		if len(out) == 2 && out[1].Kind() == reflect.Bool {
			val := out[0].Interface()
			if tVal, ok := val.(T); ok {
				return Generic[T]{val: tVal, ok: out[1].Bool()}
			}
		}
		// Сигнатура: (T)
		if len(out) == 1 {
			if tVal, ok := out[0].Interface().(T); ok {
				return Generic[T]{val: tVal, ok: true}
			}
		}
	}

	// Методы не найдены → значение отсутствует
	return Generic[T]{}
}
