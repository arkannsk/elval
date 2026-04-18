package local_generic

import (
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
