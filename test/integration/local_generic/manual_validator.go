package local_generic

import (
	"regexp"

	"github.com/arkannsk/elval"
	"github.com/arkannsk/elval/pkg/errs"
)

var _emailRegex = regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

func (u UserProfile) Validate() error {
	// ── Email: required + pattern ─────────────────────────────
	if !u.Email.IsPresent() {
		return &elval.ValidationError{
			Field:   "Email",
			Rule:    "required",
			Message: errs.ErrRequired.Error(),
		}
	}
	if emailVal, ok := u.Email.Value(); ok {
		if !_emailRegex.MatchString(emailVal) {
			return errs.NewValidationError(
				"pattern:email",
				"невалидный email адрес: %q",
				emailVal,
			)
		}
	}

	// ── Age: optional + min ───────────────────────────────────
	if ageVal, ok := u.Age.Value(); ok {
		if ageVal < 18 {
			return errs.NewValidationError(
				"min:18",
				"возраст должен быть не менее %d, получено: %d",
				18, ageVal,
			)
		}
	}

	// ── Meta: required + nested ───────────────────────────────
	if !u.Metadata.IsPresent() {
		return &elval.ValidationError{
			Field:   "Metadata",
			Rule:    "required",
			Message: errs.ErrRequired.Error(),
		}
	}
	if metaVal, ok := u.Metadata.Value(); ok {
		// Прямо вызываем валидацию вложенной структуры
		if err := metaVal.Validate(); err != nil {
			// Вариант 1: просто пробросить ошибку вверх
			// return err

			// Вариант 2: обернуть с указанием поля (рекомендую)
			return errs.NewValidationError(
				"nested",
				"ошибка валидации поля %q: %v",
				"Metadata", err,
			)
		}
	}

	return nil
}
