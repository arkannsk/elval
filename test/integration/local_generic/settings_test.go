package local_generic

import (
	"testing"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/arkannsk/elval/test/integration/local_generic/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserSettings_Validate(t *testing.T) {
	t.Run("custom theme without color fails", func(t *testing.T) {
		s := UserSettings{
			Theme:        ThemeCustom,
			PrimaryColor: model.None[string](), // отсутствует
		}
		err := s.Validate()
		require.Error(t, err)
		var ve *errs.ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "PrimaryColor", ve.Field)
		assert.Equal(t, errs.ErrRequired.Rule, ve.Rule)
	})

	t.Run("custom theme with invalid color", func(t *testing.T) {
		s := UserSettings{
			Theme:        ThemeCustom,
			PrimaryColor: model.Some("not-a-color"),
		}
		err := s.Validate()
		require.Error(t, err)
		var ve *errs.ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "pattern:hex", ve.Rule)
	})

	t.Run("custom theme with valid color", func(t *testing.T) {
		s := UserSettings{
			Theme:        ThemeCustom,
			PrimaryColor: model.Some("#FF5733"),
		}
		assert.NoError(t, s.Validate())
	})

	t.Run("light theme ignores PrimaryColor", func(t *testing.T) {
		s := UserSettings{
			Theme:        ThemeLight,
			PrimaryColor: model.Some("ignored"), // не валидируется
		}
		assert.NoError(t, s.Validate())
	})

	t.Run("optional email validation", func(t *testing.T) {
		s := UserSettings{
			Theme:             ThemeDark,
			NotificationEmail: model.Some("not-an-email"),
		}
		err := s.Validate()
		require.Error(t, err)
		var ve *errs.ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "pattern:email", ve.Rule)
	})
}
