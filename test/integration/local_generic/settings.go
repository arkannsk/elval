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
	// @evl:validate required
	// @evl:validate enum:light,dark,custom
	Theme Theme

	// @evl:validate pattern:^#[0-9A-Fa-f]{6}$
	PrimaryColor model.Option[string]

	// @evl:validate optional
	// @evl:validate pattern:email
	NotificationEmail model.Option[string]
}
