package invalid

// InvalidStruct содержит невалидные аннотации для проверки ошибок
type InvalidStruct struct {
	// @evl:validate required
	// @evl:validate min:10
	// @evl:validate max:5
	// Ошибка: min (10) > max (5) - но это не ловится на этапе парсинга,
	// потому что min и max - это отдельные директивы
	// Для проверки min > max нужно использовать min-max
	Name string

	// @evl:validate required
	// @evl:validate min-max:10,5
	// Ошибка: min (10) > max (5)
	Age int

	// @evl:validate pattern:email
	// Ошибка: pattern не поддерживается для int
	Count int

	// @evl:validate unknown
	// Ошибка: неизвестная директива
	Field string

	// @evl:validate min:invalid
	// Ошибка: параметр не число
	InvalidNumber int

	// @evl:validate required
	// @evl:validate min:1
	// @evl:validate max:10
	// @evl:validate pattern:email
	// Ошибка: pattern не поддерживается для слайса
	Tags []string
}
