package main

import (
	"unicode"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/arkannsk/elval/pkg/validator"
)

func init() {
	validator.RegisterCustom("x-strong-password", func(value any, params string) *errs.ValidationError {
		str, ok := value.(string)
		if !ok {
			return errs.NewValidationError("", "type", "expected string for x-strong-password")
		}
		if len(str) < 8 {
			return errs.NewValidationError("", "x-strong-password", "password must be at least 8 characters")
		}
		var hasUpper, hasLower, hasDigit bool
		for _, r := range str {
			if unicode.IsUpper(r) {
				hasUpper = true
			}
			if unicode.IsLower(r) {
				hasLower = true
			}
			if unicode.IsDigit(r) {
				hasDigit = true
			}
		}
		if !hasUpper || !hasLower || !hasDigit {
			return errs.NewValidationError("", "x-strong-password", "password must contain upper, lower, and digit")
		}
		return nil
	})
}
