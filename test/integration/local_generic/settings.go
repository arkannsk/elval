// test/integration/local_generic/settings.go
package local_generic

import (
	"strings"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/arkannsk/elval/test/integration/local_generic/model"
)

type Theme string

const (
	ThemeLight  Theme = "light"
	ThemeDark   Theme = "dark"
	ThemeCustom Theme = "custom"
)

type UserSettings struct {
	// @evl:validate required, oneof:light,dark,custom
	Theme Theme

	// @evl:validate required_if:Theme,custom, pattern:^#[0-9A-Fa-f]{6}$
	// обязателен только если Theme == "custom", и должен быть hex-цвет
	PrimaryColor model.Option[string]

	// @evl:validate optional, email
	NotificationEmail model.Option[string]
}

func (s UserSettings) Validate() error {
	// Theme: required + oneof
	switch s.Theme {
	case ThemeLight, ThemeDark, ThemeCustom:
		// ok
	default:
		return errs.NewValidationError(
			"oneof:light,dark,custom",
			"theme must be one of [light, dark, custom], got: %q",
			s.Theme,
		)
	}

	// PrimaryColor: required_if Theme == custom
	if s.Theme == ThemeCustom {
		if !s.PrimaryColor.IsPresent() {
			return &errs.ValidationError{
				Field:   "PrimaryColor",
				Rule:    errs.ErrRequired.Rule,
				Message: errs.ErrRequired.Message,
			}
		}
		if color, ok := s.PrimaryColor.Value(); ok {
			// Простая проверка hex: #RRGGBB
			if len(color) != 7 || color[0] != '#' {
				return errs.NewValidationError(
					"pattern:hex",
					"color must be in hex format #RRGGBB, got: %q",
					color,
				)
			}
			// Можно добавить проверку символов, но для примера достаточно
		}
	}
	// Если Theme != custom — PrimaryColor игнорируем (даже если указан)

	// NotificationEmail: optional + email pattern
	if email, ok := s.NotificationEmail.Value(); ok {
		if !isValidEmail(email) {
			return errs.NewValidationError(
				"pattern:email",
				"invalid email: %q",
				email,
			)
		}
	}

	return nil
}

// Вспомогательная функция (можно вынести в pkg/errs или internal)
func isValidEmail(s string) bool {
	// Упрощённая проверка для примера
	// В реальном коде используйте regexp или библиотеку
	return strings.Contains(s, "@") && strings.Contains(s, ".")
}
