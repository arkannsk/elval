package local_generic

import (
	"testing"
	"time"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/arkannsk/elval/test/integration/local_generic/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvent_Validate(t *testing.T) {
	now := time.Now()
	past := now.AddDate(-1, 0, 0)
	future := now.AddDate(0, 0, 30)

	t.Run("valid event with dates", func(t *testing.T) {
		e := Event{
			Name:      "Conference",
			StartDate: model.Some(future),
			EndDate:   model.Some(future.AddDate(0, 0, 3)),
		}
		assert.NoError(t, e.Validate())
	})

	t.Run("start date in past", func(t *testing.T) {
		e := Event{
			Name:      "Old Event",
			StartDate: model.Some(past),
		}
		err := e.Validate()
		require.Error(t, err)
		var ve *errs.ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "StartDate", ve.Field)
		assert.Equal(t, "min_time:now", ve.Rule)
	})

	t.Run("end before start", func(t *testing.T) {
		e := Event{
			Name:      "Bad Order",
			StartDate: model.Some(future.AddDate(0, 0, 10)),
			EndDate:   model.Some(future), // раньше start
		}
		err := e.Validate()
		require.Error(t, err)
		var ve *errs.ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "EndDate", ve.Field)
		assert.Equal(t, "after_field:StartDate", ve.Rule)
	})

	t.Run("optional dates absent", func(t *testing.T) {
		e := Event{Name: "Flexible Event"} // даты не указаны
		assert.NoError(t, e.Validate())
	})
}
