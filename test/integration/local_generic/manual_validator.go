package local_generic

import (
	"regexp"

	"github.com/arkannsk/elval/pkg/errs"
)

var _emailRegex = regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

func (u UserProfile) Validate() error {
	// ── Email: required + pattern ─────────────────────────────
	if !u.Email.IsPresent() {
		// Вариант А: использовать готовую константу + указать поле
		return &errs.ValidationError{
			Field:   "Email",
			Rule:    errs.ErrRequired.Rule,
			Message: errs.ErrRequired.Message,
		}

		// Вариант Б (если добавите хелпер):
		// return errs.WithField(errs.ErrRequired, "Email")
	}

	if emailVal, ok := u.Email.Value(); ok {
		if !_emailRegex.MatchString(emailVal) {
			// NewValidationError теперь возвращает *ValidationError напрямую
			return errs.NewValidationError(
				"pattern:email",
				"invalid email: %q",
				emailVal,
			)
		}
	}

	// ── Age: optional + min ───────────────────────────────────
	if ageVal, ok := u.Age.Value(); ok {
		if ageVal < 18 {
			return errs.NewValidationError(
				"min:18",
				"age must be at least %d, got: %d",
				18, ageVal,
			)
		}
	}

	// ── Meta: required + nested ───────────────────────────────
	if !u.Metadata.IsPresent() {
		return &errs.ValidationError{
			Field:   "Metadata",
			Rule:    errs.ErrRequired.Rule,
			Message: errs.ErrRequired.Message,
		}
	}

	if metaVal, ok := u.Metadata.Value(); ok {
		if err := metaVal.Validate(); err != nil {
			// Оборачиваем с контекстом поля
			return errs.NewValidationError(
				"nested",
				"validation failed for field %q: %v",
				"Metadata", err,
			)
		}
	}
	return nil
}
