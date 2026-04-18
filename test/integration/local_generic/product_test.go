package local_generic

import (
	"testing"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/arkannsk/elval/test/integration/local_generic/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProduct_Validate(t *testing.T) {
	t.Run("valid product with mixed reviews", func(t *testing.T) {
		p := Product{
			Name: "Gadget",
			Reviews: []model.Option[Review]{
				model.Some(Review{Comment: "Great product!", Rating: model.Some(5)}),
				model.None[Review](), // пропущенный ревью — ок
				model.Some(Review{Comment: "Good", Rating: model.None[int]()}), // рейтинг не указан — ок
			},
		}
		assert.NoError(t, p.Validate())
	})

	t.Run("review with invalid comment length", func(t *testing.T) {
		p := Product{
			Name: "Gadget",
			Reviews: []model.Option[Review]{
				model.Some(Review{Comment: "OK", Rating: model.Some(4)}), // too short
			},
		}
		err := p.Validate()
		require.Error(t, err)
		var ve *errs.ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "Reviews[0]", ve.Field) // проверяем индекс в сообщении
		assert.Contains(t, ve.Message, "comment must be at least 3")
	})

	t.Run("rating out of range", func(t *testing.T) {
		p := Product{
			Name: "Gadget",
			Reviews: []model.Option[Review]{
				model.Some(Review{Comment: "Nice", Rating: model.Some(10)}),
			},
		}
		err := p.Validate()
		require.Error(t, err)
		var ve *errs.ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "min:1,max:5", ve.Rule)
	})

	t.Run("empty reviews slice is ok (optional)", func(t *testing.T) {
		p := Product{Name: "Gadget", Reviews: []model.Option[Review]{}}
		assert.NoError(t, p.Validate())
	})
}
