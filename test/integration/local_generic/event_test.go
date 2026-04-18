package local_generic

import (
	"testing"
	"time"

	"github.com/arkannsk/elval/test/integration/local_generic/model"
	"github.com/stretchr/testify/assert"
)

func TestEvent_Validate(t *testing.T) {
	now := time.Now()
	future := now.AddDate(0, 0, 30)

	t.Run("valid event with dates", func(t *testing.T) {
		e := Event{
			Name:      "Conference",
			StartDate: model.Some(future),
			EndDate:   model.Some(future.AddDate(0, 0, 3)),
		}
		assert.NoError(t, e.Validate())
	})
}
