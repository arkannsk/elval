package types

import "time"

// AllTypes содержит поля разных типов для проверки парсинга
type AllTypes struct {
	// Примитивы
	// @evl:validate required
	// @evl:validate min:3
	// @evl:validate max:50
	StringField string

	// @evl:validate required
	// @evl:validate min:18
	// @evl:validate max:99
	IntField int

	// @evl:validate min:0.1
	// @evl:validate max:999.99
	FloatField float64

	// @evl:validate required
	BoolField bool

	// Слайсы примитивов
	// @evl:validate required
	// @evl:validate min:1
	// @evl:validate max:10
	StringSlice []string

	// @evl:validate len:3
	IntSlice []int

	// @evl:validate not-zero
	FloatSlice []float64

	// Указатели на примитивы
	// @evl:validate required
	// @evl:validate min:18
	AgePtr *int

	// @evl:validate optional
	NamePtr *string

	// Слайс указателей
	// @evl:validate min:1
	UserPtrs []*User

	// time.Time (не указатель)
	// @evl:validate required
	// @evl:validate after:2020-01-01
	// @evl:validate before:2030-12-31
	CreatedAt time.Time

	// Указатель на time.Time - валидируем как time.Time
	// @evl:validate optional
	// @evl:validate after:2020-01-01
	UpdatedAt *time.Time

	// time.Duration (не указатель)
	// @evl:validate required
	// @evl:validate min:1s
	// @evl:validate max:24h
	Timeout time.Duration

	// Указатель на time.Duration - валидируем как time.Duration
	// @evl:validate optional
	// @evl:validate min:100ms
	Interval *time.Duration
}

// User вспомогательная структура
type User struct {
	// @evl:validate required
	// @evl:validate min:3
	Name string

	// @evl:validate required
	// @evl:validate pattern:email
	Email string
}
