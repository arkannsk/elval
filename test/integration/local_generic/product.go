package local_generic

import (
	"fmt"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/arkannsk/elval/test/integration/local_generic/model"
)

type Review struct {
	// @evl:validate required, min:3, max:500
	Comment string
	// @evl:validate optional, min:1, max:5
	Rating model.Option[int]
}

type Product struct {
	// @evl:validate required
	Name string

	// @evl:validate optional, not-empty
	// каждый присутствующий Review должен быть валидным
	Reviews []model.Option[Review]
}

func (p Product) Validate() error {
	// Name: required
	if p.Name == "" {
		return &errs.ValidationError{
			Field:   "Name",
			Rule:    errs.ErrRequired.Rule,
			Message: errs.ErrRequired.Message,
		}
	}

	// Reviews: валидируем только присутствующие элементы
	for i, optReview := range p.Reviews {
		if review, ok := optReview.Value(); ok {
			if err := review.Validate(); err != nil {
				field := "Reviews"
				if ve, ok := err.(*errs.ValidationError); ok && ve.Field != "" {
					field = fmt.Sprintf("Reviews[%d].%s", i, ve.Field)
				} else {
					field = fmt.Sprintf("Reviews[%d]", i)
				}

				// Сохраняем оригинальный Rule и сообщение
				if ve, ok := err.(*errs.ValidationError); ok {
					return &errs.ValidationError{
						Field:   field,
						Rule:    ve.Rule,
						Message: ve.Message,
					}
				}
				// Fallback, если ошибка не нашего типа
				return &errs.ValidationError{
					Field:   field,
					Rule:    "nested",
					Message: err.Error(),
				}
			}
		}
	}

	return nil
}

func (r Review) Validate() error {
	// Comment: required + length
	if r.Comment == "" {
		return &errs.ValidationError{
			Field:   "Comment",
			Rule:    errs.ErrRequired.Rule,
			Message: errs.ErrRequired.Message,
		}
	}
	if len(r.Comment) < 3 {
		return errs.NewValidationError(
			"min:3",
			"comment must be at least %d characters, got: %d",
			3, len(r.Comment),
		)
	}
	if len(r.Comment) > 500 {
		return errs.NewValidationError(
			"max:500",
			"comment must not exceed %d characters, got: %d",
			500, len(r.Comment),
		)
	}

	// Rating: optional + range
	if rating, ok := r.Rating.Value(); ok {
		if rating < 1 || rating > 5 {
			return errs.NewValidationError(
				"min:1,max:5",
				"rating must be between %d and %d, got: %d",
				1, 5, rating,
			)
		}
	}

	return nil
}
