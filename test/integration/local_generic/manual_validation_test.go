package local_generic

import (
	"testing"

	"github.com/arkannsk/elval"
	"github.com/arkannsk/elval/pkg/errs"
	"github.com/arkannsk/elval/test/integration/local_generic/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserProfile_Validate_Manual(t *testing.T) {
	t.Run("valid all present", func(t *testing.T) {
		u := UserProfile{
			Email:    model.Some("user@example.com"),
			Age:      model.Some(25),
			Metadata: model.Some(UserMeta{DisplayName: "Alice"}),
		}
		assert.NoError(t, u.Validate())
	})

	t.Run("email missing", func(t *testing.T) {
		u := UserProfile{
			Age:      model.Some(25),
			Metadata: model.Some(UserMeta{DisplayName: "Alice"}),
		}
		err := u.Validate()
		require.Error(t, err)
		var ve *elval.ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "Email", ve.Field)
		assert.Equal(t, errs.ErrRequired.Rule, ve.Rule)
	})

	t.Run("nested validation fails", func(t *testing.T) {
		u := UserProfile{
			Email:    model.Some("user@example.com"),
			Age:      model.Some(25),
			Metadata: model.Some(UserMeta{DisplayName: ""}), // нарушает required
		}
		err := u.Validate()
		require.Error(t, err)
		var ve *elval.ValidationError
		require.ErrorAs(t, err, &ve)
	})
}

func TestUserProfile_Validate(t *testing.T) {
	t.Run("email required", func(t *testing.T) {
		u := UserProfile{ /* Email absent */ }
		err := u.Validate()

		var ve *errs.ValidationError
		require.ErrorAs(t, err, &ve)

		assert.Equal(t, "Email", ve.Field)
		assert.Equal(t, errs.ErrRequired.Rule, ve.Rule)
		assert.Equal(t, errs.ErrRequired.Message, ve.Message)
	})
}
